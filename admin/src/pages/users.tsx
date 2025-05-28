import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query'
import { For, Show, createSignal } from 'solid-js'
import {
	createSolidTable,
	getCoreRowModel,
	getPaginationRowModel,
	flexRender,
	ColumnDef,
} from '@tanstack/solid-table'
import { fetchUsers, updateUser, updateUserStatus, deleteUsers } from '~/lib/api'
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
import { UserResponse, VerificationStatus, CityResponse, UpdateUserRequest } from '~/gen/types'
import { TextField, TextFieldInput, TextFieldTextArea } from '~/components/ui/text-field'
import {
	Select,
	SelectContent,
	SelectItem, SelectLabel,
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

export default function UsersPage() {
	const [page, setPage] = createSignal(0)
	const [pageSize] = createSignal(10)
	const [searchQuery, setSearchQuery] = createSignal('')
	const [statusFilter, setStatusFilter] = createSignal('All')
	const [selectedUser, setSelectedUser] = createSignal<UserResponse | null>(null)
	const [isSheetOpen, setIsSheetOpen] = createSignal(false)
	const [editedUser, setEditedUser] = createSignal<Partial<UserResponse>>({})
	const [showDeleteDialog, setShowDeleteDialog] = createSignal(false)
	const queryClient = useQueryClient()

	const query = useQuery(() => ({
		queryKey: ['users', page(), statusFilter(), searchQuery()],
		queryFn: () => fetchUsers({
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
	} = useTableSelection<UserResponse>()

	const deleteUsersMutation = useMutation(() => ({
		mutationFn: deleteUsers,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['users'] })
			clearSelection()
			setShowDeleteDialog(false)
		},
		onError: (error) => {
			console.error('Failed to delete users:', error)
		},
	}))

	const columns: ColumnDef<UserResponse>[] = [
		{
			accessorKey: 'avatar_url',
			header: 'Avatar',
			cell: (props) => {
				const avatarUrl = props.getValue() as string
				return avatarUrl ? (
					<img
						src={`https://assets.peatch.io/cdn-cgi/image/width=100/${avatarUrl}`}
						alt="User avatar"
						class="size-8 rounded-full object-cover"
					/>
				) : (
					<div class="size-8 rounded-full bg-secondary flex items-center justify-center">
						<span class="text-secondary-foreground text-xs">
							Nil
						</span>
					</div>
				)
			},
		},
		{
			accessorKey: 'chat_id',
			header: 'Chat ID',
		},
		{
			accessorKey: 'name',
			header: 'Name',
		},
		{
			accessorKey: 'username',
			header: 'Username',
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
			accessorKey: 'location',
			header: 'City',
			cell: (props) => {
				const location = props.getValue() as CityResponse
				return <span>{location?.name || '-'}</span>
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
	]

	const columnsWithSelection = [
		createSelectionColumn(() => query.data?.users || []),
		...columns,
	]

	const updateUserMutation = useMutation(() => ({
		mutationFn: () => {
			if (!selectedUser() || !selectedUser()?.id) return Promise.reject('No user selected')
			const data: UpdateUserRequest = {
				name: editedUser().name || selectedUser()!.name,
				title: editedUser().title || selectedUser()!.title,
				description: editedUser().description || selectedUser()!.description,
			}
			return updateUser(selectedUser()?.id, data)
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['users'] })
			setIsSheetOpen(false)
		},
		onError: (error) => {
			console.error('Failed to update user:', error)
		},
	}))

	const updateStatusMutation = useMutation(() => ({
		mutationFn: () => {
			if (!selectedUser() || !selectedUser()?.id) return Promise.reject('No user selected')
			const newStatus = selectedUser()?.verification_status === 'verified' ? 'unverified' : 'verified'
			return updateUserStatus(selectedUser()?.id, newStatus)
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['users'] })
		},
		onError: (error) => {
			console.error('Failed to update user status:', error)
		},
	}))

	const table = createSolidTable({
		get data() {
			return query.data?.users || []
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
					<h1 class="text-2xl font-bold">Users</h1>
				</div>

				<div class="flex gap-4">
					<TextField
						value={searchQuery()}
						onChange={setSearchQuery}
						class="max-w-sm"
					>
						<TextFieldInput
							placeholder="Search users..."
							type="search"
						/>
					</TextField>

					<Select
						value={statusFilter()}
						onChange={(value) => {
							setStatusFilter(value || 'verified')
							setPage(0) // Reset to first page on status change
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

				<Show when={query.isLoading}>
					<div>Loading...</div>
				</Show>

				<Show when={query.error}>
					<div class="text-destructive-foreground">Error: {query.error?.message}</div>
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
													setSelectedUser(row.original)
													setEditedUser(row.original)
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
					<Show when={selectedUser()}>
						<SheetHeader>
							<SheetTitle>Edit User</SheetTitle>
							<SheetDescription>
								Update user information. Click save when you're done.
							</SheetDescription>
						</SheetHeader>

						<div class="space-y-4 mt-6">
							<div>
								<img
									src={`https://assets.peatch.io/cdn-cgi/image/width=100/${selectedUser()?.avatar_url}`}
									alt="User avatar"
									class="size-16 rounded-full object-cover mb-2"
								/>
							</div>
							<div>
								<label class="text-sm font-medium">Chat ID</label>
								<p class="text-sm text-secondary-foreground">{selectedUser()?.chat_id}</p>
							</div>
							<div>
								<label class="text-sm font-medium">Name</label>
								<TextField
									value={editedUser().name || ''}
									onChange={(value) => setEditedUser(prev => ({ ...prev, name: value }))}
								>
									<TextFieldInput placeholder="Enter name" />
								</TextField>
							</div>

							<div>
								<label class="text-sm font-medium">Title</label>
								<TextField
									value={editedUser().title || ''}
									onChange={(value) => setEditedUser(prev => ({ ...prev, title: value }))}
								>
									<TextFieldInput placeholder="Enter title" />
								</TextField>
							</div>

							<div>
								<label class="text-sm font-medium">Description</label>
								<TextField
									value={editedUser().description || ''}
									onChange={(value) => setEditedUser(prev => ({ ...prev, description: value }))}
								>
									<TextFieldTextArea
										placeholder="Enter description"
										class="min-h-[120px] resize-y"
									/>
								</TextField>
							</div>

							<div>
								<label class="text-sm font-medium">Location</label>
								<p class="text-sm text-secondary-foreground">{selectedUser()?.location?.name || 'Not set'}</p>
							</div>

							<div>
								<label class="text-sm font-medium">Verification Status</label>
								<Select
									value={selectedUser()?.verification_status || 'unverified'}
									onChange={(value) => {
										if (selectedUser()?.id && value) {
											updateStatusMutation.mutate()
										}
									}}
									disabled={updateStatusMutation.isPending}
									options={statusOptions}
									itemComponent={(props) => (
										<SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
									)}
								>
									<SelectTrigger class="w-full">
										<SelectValue<{ value: string; label: string }>>
											{(state) => state.selectedOption()?.label || 'Select status'}
										</SelectValue>
									</SelectTrigger>
									<SelectContent />
								</Select>
							</div>


							<div>
								<label class="text-sm font-medium">Created At</label>
								<p class="text-sm text-secondary-foreground">
									{formatDate(selectedUser()?.created_at)}
								</p>
							</div>

							<div>
								<label class="text-sm font-medium">Last Active</label>
								<p class="text-sm text-secondary-foreground">
									{formatDate(selectedUser()?.last_active_at)}
								</p>
							</div>

							<div class="flex gap-2 pt-4">
								<Button
									onClick={() => {
										if (selectedUser()?.id) {
											updateUserMutation.mutate()
										}
									}}
									disabled={updateUserMutation.isPending}
								>
									{updateUserMutation.isPending ? 'Saving...' : 'Save Changes'}
								</Button>
								<Button
									variant="outline"
									onClick={() => setIsSheetOpen(false)}
									disabled={updateUserMutation.isPending}
								>
									Cancel
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
							<AlertDialogTitle>Delete Users</AlertDialogTitle>
							<AlertDialogDescription>
								Are you sure you want to delete {selectedCount()} user{selectedCount() > 1 ? 's' : ''}?
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
									const selectedUsers = getSelectedItems(query.data?.users || [])
									const userIds = selectedUsers.map(u => u.id).filter(Boolean) as string[]
									deleteUsersMutation.mutate(userIds)
								}}
								disabled={deleteUsersMutation.isPending}
							>
								{deleteUsersMutation.isPending ? 'Deleting...' : 'Delete'}
							</Button>
						</div>
					</div>
				</AlertDialogContent>
			</AlertDialog>
		</>
	)
}
