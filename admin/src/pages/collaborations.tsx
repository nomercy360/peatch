import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query'
import { For, Show, createSignal } from 'solid-js'
import {
	createSolidTable,
	getCoreRowModel,
	getPaginationRowModel,
	flexRender,
	ColumnDef,
} from '@tanstack/solid-table'
import { fetchCollaborations, updateCollaborationStatus, deleteCollaborations } from '~/lib/api'
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '~/components/ui/table'
import { Button } from '~/components/ui/button'
import { AdminLayout } from '~/components/AdminLayout'
import {
	CollaborationResponse,
	VerificationStatus,
	CityResponse,
	UserProfileResponse,
	OpportunityResponse,
} from '~/gen/types'
import { TextField, TextFieldInput } from '~/components/ui/text-field'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '~/components/ui/select'
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from '~/components/ui/sheet'
import {
	AlertDialog,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogTitle,
} from '~/components/ui/alert-dialog'
import { cn, formatDate } from '~/lib/utils'
import { useTableSelection } from '~/lib/useTableSelection'
import { IconTrash } from '~/components/icons'

const statusOptions = [
	'All',
	'Verified',
	'Pending',
	'Unverified',
	'Rejected',
	'Blocked',
]

export default function CollaborationsPage() {
	const [page, setPage] = createSignal(0)
	const [pageSize] = createSignal(10)
	const [searchQuery, setSearchQuery] = createSignal('')
	const [statusFilter, setStatusFilter] = createSignal('All')
	const [selectedCollaboration, setSelectedCollaboration] = createSignal<CollaborationResponse | null>(null)
	const [isSheetOpen, setIsSheetOpen] = createSignal(false)
	const [showDeleteDialog, setShowDeleteDialog] = createSignal(false)
	const queryClient = useQueryClient()

	const query = useQuery(() => ({
		queryKey: ['collaborations', page(), statusFilter(), searchQuery()],
		queryFn: () => fetchCollaborations({
			page: page(),
			limit: pageSize(),
			status: statusFilter() !== 'All' ? statusFilter().toLowerCase() : '',
		}),
	}))

	const {
		clearSelection,
		selectedCount,
		getSelectedItems,
		createSelectionColumn,
	} = useTableSelection<CollaborationResponse>()

	const deleteCollaborationsMutation = useMutation(() => ({
		mutationFn: deleteCollaborations,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['collaborations'] })
			clearSelection()
			setShowDeleteDialog(false)
		},
		onError: (error) => {
			console.error('Failed to delete collaborations:', error)
		},
	}))

	const columns: ColumnDef<CollaborationResponse>[] = [
		{
			accessorKey: 'user',
			header: 'User',
			cell: (props) => {
				const user = props.getValue() as UserProfileResponse
				return (
					<div class="flex items-center gap-2">
						{user?.avatar_url ? (
							<img
								src={`https://assets.peatch.io/cdn-cgi/image/width=100/${user.avatar_url}`}
								alt="User avatar"
								class="size-8 rounded-full object-cover"
							/>
						) : (
							<div class="size-8 rounded-full bg-secondary flex items-center justify-center">
								<span class="text-secondary-foreground text-xs">
									Nil
								</span>
							</div>
						)}
						<span class="text-sm">{user?.name || '-'}</span>
					</div>
				)
			},
		},
		{
			accessorKey: 'title',
			header: 'Title',
			cell: (props) => {
				const title = props.getValue() as string
				return <span class="max-w-xs truncate">{title || '-'}</span>
			},
		},
		{
			accessorKey: 'description',
			header: 'Description',
			cell: (props) => {
				const description = props.getValue() as string
				return <span class="max-w-xs truncate">{description || '-'}</span>
			},
		},
		{
			accessorKey: 'opportunity',
			header: 'Opportunity',
			cell: (props) => {
				const opportunity = props.getValue() as OpportunityResponse
				return <span>{opportunity?.text || '-'}</span>
			},
		},
		{
			accessorKey: 'location',
			header: 'City',
			cell: (props) => {
				const location = props.getValue() as CityResponse
				return <span>{location?.name || '-'}</span>
			},
		},
		{
			accessorKey: 'is_payable',
			header: 'Payable',
			cell: (props) => {
				const isPayable = props.getValue() as boolean
				return (
					<span class={cn('text-xs', isPayable ? 'text-green-600' : 'text-secondary-foreground')}>
						{isPayable ? 'Yes' : 'No'}
					</span>
				)
			},
		},
		{
			accessorKey: 'verification_status',
			header: 'Status',
			cell: (props) => {
				const status = props.getValue() as VerificationStatus
				return (
					<span
						class={cn('rounded-xl px-2 py-1 text-xs font-medium',
							status === 'verified'
								? 'text-success-foreground bg-success font-medium'
								: status === 'unverified'
									? 'text-muted-foreground bg-muted font-medium'
									: 'text-warning-foreground bg-warning font-medium')
						}
					>
          {status}
        </span>
				)
			},
		},
		{
			accessorKey: 'created_at',
			header: 'Created',
			cell: (props) => {
				const date = props.getValue() as string
				return <span class="text-xs text-secondary-foreground">{formatDate(date)}</span>
			},
		},
	]

	const columnsWithSelection = [
		createSelectionColumn(() => query.data?.collaborations || []),
		...columns
	]

	const updateStatusMutation = useMutation(() => ({
		mutationFn: () => {
			const collaboration = selectedCollaboration()
			if (!collaboration?.id || !collaboration?.user_id) return Promise.reject('No collaboration selected')
			const newStatus = collaboration.verification_status === 'verified' ? 'unverified' : 'verified'
			return updateCollaborationStatus(collaboration.user_id, collaboration.id, newStatus)
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['collaborations'] })
			setIsSheetOpen(false)
		},
		onError: (error) => {
			console.error('Failed to update collaboration status:', error)
		},
	}))

	const table = createSolidTable({
		get data() {
			return query.data?.collaborations || []
		},
		columns: columnsWithSelection,
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		manualPagination: true,
		pageCount: Math.ceil((query.data?.total || 0) / pageSize()),
	})

	return (
		<>
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<h1 class="text-2xl font-bold">Collaborations</h1>
					<Show when={selectedCount() > 0}>
						<div class="flex items-center gap-2">
							<span class="text-sm text-muted-foreground">
								{selectedCount()} selected
							</span>
							<Button
								variant="destructive"
								size="sm"
								onClick={() => setShowDeleteDialog(true)}
							>
								<IconTrash class="size-4 mr-2" />
								Delete
							</Button>
						</div>
					</Show>
				</div>

				<div class="flex gap-4">
					<TextField
						value={searchQuery()}
						onChange={setSearchQuery}
						class="max-w-sm"
					>
						<TextFieldInput
							placeholder="Search collaborations..."
							type="search"
						/>
					</TextField>

					<Select
						value={statusFilter()}
						onChange={(value) => {
							setStatusFilter(value || 'verified')
							setPage(0)
						}}
						options={statusOptions}
						placeholder="Filter by status"
						defaultValue={statusFilter()}
						itemComponent={(props) => (
							<SelectItem item={props.item}>
								{props.item.rawValue}
							</SelectItem>
						)}
					>
						<SelectTrigger class="w-[180px]">
							<SelectValue<string>>
								{(state) => state.selectedOption()}
							</SelectValue>
						</SelectTrigger>
						<SelectContent />
					</Select>
				</div>

				<Show when={query.isLoading}>
					<div>Loading...</div>
				</Show>

				<Show when={query.error}>
					<div class="text-red-600">Error: {query.error?.message}</div>
				</Show>

				<Show when={query.data}>
					<div class="rounded-md border">
						<Table>
							<TableHeader>
								<For each={table.getHeaderGroups()}>
									{(headerGroup) => (
										<TableRow>
											<For each={headerGroup.headers}>
												{(header) => (
													<TableHead>
														<Show when={!header.isPlaceholder}>
															{flexRender(
																header.column.columnDef.header,
																header.getContext(),
															)}
														</Show>
													</TableHead>
												)}
											</For>
										</TableRow>
									)}
								</For>
							</TableHeader>
							<TableBody>
								<Show
									when={table.getRowModel().rows?.length}
									fallback={
										<TableRow>
											<TableCell
												colSpan={columns.length}
												class="h-24 text-center"
											>
												No results.
											</TableCell>
										</TableRow>
									}
								>
									<For each={table.getRowModel().rows}>
										{(row) => (
											<TableRow
												class="cursor-pointer hover:bg-gray-50"
												onClick={(e) => {
													// Don't open sheet if clicking on checkbox
													if ((e.target as HTMLElement).closest('[role="checkbox"]')) {
														return
													}
													setSelectedCollaboration(row.original)
													setIsSheetOpen(true)
												}}
											>
												<For each={row.getVisibleCells()}>
													{(cell) => (
														<TableCell class="truncate max-w-[200px]">
															{flexRender(
																cell.column.columnDef.cell,
																cell.getContext(),
															)}
														</TableCell>
													)}
												</For>
											</TableRow>
										)}
									</For>
								</Show>
							</TableBody>
						</Table>
					</div>

					<div class="flex items-center justify-between">
						<div class="text-sm text-muted-foreground">
							Showing {page() * pageSize() + 1} to{' '}
							{Math.min((page() + 1) * pageSize(), query.data?.total)} of{' '}
							{query.data?.total} results
						</div>
						<div class="flex items-center space-x-2">
							<Button
								variant="outline"
								size="sm"
								onClick={() => setPage(p => Math.max(0, p - 1))}
								disabled={page() === 0}
							>
								Previous
							</Button>
							<Button
								variant="outline"
								size="sm"
								onClick={() => setPage(p => p + 1)}
								disabled={(page() + 1) * pageSize() >= (query.data?.total || 0)}
							>
								Next
							</Button>
						</div>
					</div>
				</Show>
			</div>

			<Sheet open={isSheetOpen()} onOpenChange={setIsSheetOpen}>
				<SheetContent class="w-[400px] sm:w-[540px] overflow-y-auto">
					<Show when={selectedCollaboration()}>
						<SheetHeader>
							<SheetTitle>Collaboration Details</SheetTitle>
							<SheetDescription>
								View and manage collaboration information.
							</SheetDescription>
						</SheetHeader>

						<div class="space-y-4 mt-6">
							<div>
								<label class="text-sm font-medium">User</label>
								<div class="flex items-center gap-2 mt-1">
									{selectedCollaboration()?.user?.avatar_url ? (
										<img
											src={`https://assets.peatch.io/cdn-cgi/image/width=100/${selectedCollaboration()?.user?.avatar_url}`}
											alt="User avatar"
											class="size-12 rounded-full object-cover"
										/>
									) : (
										<div class="size-12 rounded-full bg-secondary flex items-center justify-center">
											<span class="text-secondary-foreground text-xs">Nil</span>
										</div>
									)}
									<div>
										<p class="text-sm font-medium">{selectedCollaboration()?.user?.name || '-'}</p>
										<p class="text-sm text-secondary-foreground">@{selectedCollaboration()?.user?.username || '-'}</p>
									</div>
								</div>
							</div>

							<div>
								<label class="text-sm font-medium">Title</label>
								<p class="text-sm text-muted-foreground mt-1">{selectedCollaboration()?.title || '-'}</p>
							</div>

							<div>
								<label class="text-sm font-medium">Description</label>
								<p class="text-sm text-muted-foreground mt-1 whitespace-pre-wrap">{selectedCollaboration()?.description || '-'}</p>
							</div>

							<div>
								<label class="text-sm font-medium">Opportunity</label>
								<p class="text-sm text-muted-foreground mt-1">{selectedCollaboration()?.opportunity?.text || 'Not specified'}</p>
							</div>

							<div>
								<label class="text-sm font-medium">Location</label>
								<p class="text-sm text-muted-foreground mt-1">{selectedCollaboration()?.location?.name || 'Not specified'}</p>
							</div>

							<div>
								<label class="text-sm font-medium">Payable</label>
								<p class="text-sm text-muted-foreground mt-1">{selectedCollaboration()?.is_payable ? 'Yes' : 'No'}</p>
							</div>

							<div>
								<label class="text-sm font-medium">Has Interest</label>
								<p class="text-sm text-muted-foreground mt-1">{selectedCollaboration()?.has_interest ? 'Yes' : 'No'}</p>
							</div>

							<Show when={selectedCollaboration()?.badges && selectedCollaboration()?.badges?.length > 0}>
								<div>
									<label class="text-sm font-medium">Badges</label>
									<div class="flex flex-wrap gap-2 mt-1">
										<For each={selectedCollaboration()?.badges}>
											{(badge) => (
												<span class="text-xs bg-secondary-foreground text-muted-foreground px-2 py-1 rounded-full">
													{badge.text}
												</span>
											)}
										</For>
									</div>
								</div>
							</Show>

							<Show when={selectedCollaboration()?.links && selectedCollaboration()?.links?.length > 0}>
								<div>
									<label class="text-sm font-medium">Links</label>
									<div class="space-y-1 mt-1">
										<For each={selectedCollaboration()?.links}>
											{(link) => (
												<a href={link.url} target="_blank" rel="noopener noreferrer" class="text-sm text-blue-600 hover:underline block">
													{link.url}
												</a>
											)}
										</For>
									</div>
								</div>
							</Show>

							<div>
								<label class="text-sm font-medium">Verification Status</label>
								<div class="mt-2">
									<Button
										onClick={() => updateStatusMutation.mutate()}
										disabled={updateStatusMutation.isPending}
										variant={selectedCollaboration()?.verification_status === 'verified' ? 'destructive' : 'default'}
									>
										{updateStatusMutation.isPending
											? 'Updating...'
											: selectedCollaboration()?.verification_status === 'verified'
												? 'Unverify'
												: 'Verify'}
									</Button>
								</div>
							</div>

							<div>
								<label class="text-sm font-medium">Created At</label>
								<p class="text-sm text-secondary-foreground">
									{formatDate(selectedCollaboration()?.created_at)}
								</p>
							</div>

							<div>
								<label class="text-sm font-medium">Updated At</label>
								<p class="text-sm text-secondary-foreground">
									{formatDate(selectedCollaboration()?.updated_at)}
								</p>
							</div>

							<div class="flex gap-2 pt-4">
								<Button
									variant="outline"
									onClick={() => setIsSheetOpen(false)}
								>
									Close
								</Button>
							</div>
						</div>
					</Show>
				</SheetContent>
			</Sheet>

			<AlertDialog open={showDeleteDialog()} onOpenChange={setShowDeleteDialog}>
				<AlertDialogContent>
					<div class="space-y-4">
						<div class="space-y-2">
							<AlertDialogTitle>Delete Collaborations</AlertDialogTitle>
							<AlertDialogDescription>
								Are you sure you want to delete {selectedCount()} collaboration{selectedCount() > 1 ? 's' : ''}?
								This action cannot be undone.
							</AlertDialogDescription>
						</div>
						<div class="flex justify-end gap-2">
							<Button variant="outline" onClick={() => setShowDeleteDialog(false)}>
								Cancel
							</Button>
							<Button
								variant="destructive"
								onClick={() => {
									const selectedCollaborations = getSelectedItems(query.data?.collaborations || [])
									const collaborationIds = selectedCollaborations.map(c => c.id).filter(Boolean) as string[]
									deleteCollaborationsMutation.mutate(collaborationIds)
								}}
								disabled={deleteCollaborationsMutation.isPending}
							>
								{deleteCollaborationsMutation.isPending ? 'Deleting...' : 'Delete'}
							</Button>
						</div>
					</div>
				</AlertDialogContent>
			</AlertDialog>
		</>
	)
}
