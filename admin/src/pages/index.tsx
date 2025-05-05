import { createSignal, For, Show } from 'solid-js'
import Header from '~/components/admin/Header'
import Filters from '~/components/admin/Filters'
import { fetchCollaborations, fetchUsers, updateUserStatus, updateCollaborationStatus } from '~/lib/api'
import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query'
import { type CollaborationResponse, type UserResponse, type VerificationStatus } from '~/gen'
import UserCard from '~/components/admin/UserCard'
import Pagination from '~/components/admin/Pagination'
import CollaborationCard from '~/components/admin/CollaborationCard'

export default function AdminPanel() {
	const queryClient = useQueryClient()
	const [activeTab, setActiveTab] = createSignal<'users' | 'collaborations'>('users')
	const [page, setPage] = createSignal(1)
	const [filterStatus, setFilterStatus] = createSignal<string>('pending')

	const users = useQuery(() => ({
		queryKey: ['users', page(), filterStatus()],
		queryFn: () => fetchUsers({ pageParam: page(), status: filterStatus() }),
	}))

	const collaborations = useQuery(() => ({
		queryKey: ['collaborations', page(), filterStatus()],
		queryFn: () => fetchCollaborations({ pageParam: page(), status: filterStatus() }),
	}))

	const userStatusMutation = useMutation(() => ({
		mutationFn: ({ userId, status }: {
			userId: string,
			status: VerificationStatus
		}) => updateUserStatus(userId, status),
		onMutate: async ({ userId, status }) => {
			await queryClient.cancelQueries({ queryKey: ['users', page(), filterStatus()] })

			console.log(userId, status)

			const previousUsers = queryClient.getQueryData(['users', page(), filterStatus()])

			queryClient.setQueryData(['users', page(), filterStatus()], (old: any) => {
				if (!old) return old

				const newData = { ...old }
				newData.data = newData.data.map((user: any) => {
					if (user.id === userId) {
						return { ...user, verification_status: status }
					}
					return user
				})

				return newData
			})

			return { previousUsers }
		},
	}))

	const collaborationStatusMutation = useMutation(() => ({
		mutationFn: ({ userId, collaborationId, status }: {
			userId: string,
			collaborationId: string,
			status: VerificationStatus
		}) =>
			updateCollaborationStatus(userId, collaborationId, status),
		onMutate: async ({ userId, collaborationId, status }) => {
			// Cancel any outgoing refetches
			await queryClient.cancelQueries({ queryKey: ['collaborations', page(), filterStatus()] })

			// Snapshot the previous value
			const previousCollaborations = queryClient.getQueryData(['collaborations', page(), filterStatus()])

			// Optimistically update the collaboration status
			queryClient.setQueryData(['collaborations', page(), filterStatus()], (old: any) => {
				if (!old) return old

				const newData = { ...old }
				newData.data = newData.data.map((collab: any) => {
					if (collab.id === collaborationId && collab.user.id === userId) {
						return { ...collab, verificationStatus: status }
					}
					return collab
				})

				return newData
			})

			return { previousCollaborations }
		},
		onError: (_, __, context) => {
			// If the mutation fails, use the context returned from onMutate to roll back
			if (context?.previousCollaborations) {
				queryClient.setQueryData(['collaborations', page(), filterStatus()], context.previousCollaborations)
			}
		},
		onSuccess: () => {
			// Invalidate and refetch after successful mutation
			queryClient.invalidateQueries({ queryKey: ['collaborations', page(), filterStatus()] })
		},
	}))

	// Handle pagination
	const handlePrevPage = () => {
		if (page() > 1) {
			setPage(page() - 1)
		}
	}

	const handleNextPage = () => {
		setPage(page() + 1)
	}

	const handleUpdateUserStatus = (userId: string, status: VerificationStatus) => {
		userStatusMutation.mutate({ userId, status })
	}

	const handleUpdateCollaborationStatus = (userId: string, collaborationId: string, status: VerificationStatus) => {
		collaborationStatusMutation.mutate({ userId, collaborationId, status })
	}

	return (
		<div class="h-screen bg-secondary p-3 overflow-y-auto">
			{/* Header and Tabs */}
			<Header activeTab={activeTab} setActiveTab={setActiveTab} />

			{/* Filters */}
			<Filters
				filterStatus={filterStatus}
				setFilterStatus={setFilterStatus}
			/>
			<Show when={activeTab() === 'users'}>
				<div class="space-y-4">
					<Show when={users.data} fallback={<div>Loading users...</div>}>
						<For each={users.data?.data}>
							{(user: UserResponse) => (
								<UserCard
									user={user}
									updateUserStatus={handleUpdateUserStatus} />
							)}
						</For>

						<Pagination
							page={page}
							handlePrevPage={handlePrevPage}
							handleNextPage={handleNextPage}
						/>
					</Show>
				</div>
			</Show>

			{/* Collaborations Tab Content */}
			<Show when={activeTab() === 'collaborations'}>
				<div class="space-y-4">
					<Show when={collaborations && collaborations.data} fallback={<div>Loading collaborations...</div>}>
						<For each={collaborations.data?.data}>
							{(collab: CollaborationResponse) => (
								<CollaborationCard
									collab={collab}
									updateCollaborationStatus={updateCollaborationStatus}
								/>
							)}
						</For>

						<Pagination
							page={page}
							handlePrevPage={handlePrevPage}
							handleNextPage={handleNextPage}
						/>
					</Show>
				</div>
			</Show>
		</div>
	)
}
