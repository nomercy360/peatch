import { verificationStatus } from '~/gen/types'

type FiltersProps = {
	filterStatus: () => string
	setFilterStatus: (status: string) => void
}

export default function Filters({ filterStatus, setFilterStatus }: FiltersProps) {
	return (
		<div class="mb-6 flex flex-col space-y-4 md:flex-row md:space-y-0 md:space-x-4">
			<div class="flex-1">
				<label class="block text-sm font-medium text-gray-700 mb-1">Status Filter</label>
				<select
					class="w-full rounded-md border border-gray-300 p-2"
					value={filterStatus()}
					onChange={(e) => setFilterStatus(e.target.value)}
				>
					<option value={verificationStatus.VerificationStatusPending}>Pending</option>
					<option value={verificationStatus.VerificationStatusVerified}>Verified</option>
					<option value={verificationStatus.VerificationStatusDenied}>Denied</option>
					<option value={verificationStatus.VerificationStatusBlocked}>Blocked</option>
					<option value={verificationStatus.VerificationStatusUnverified}>Unverified</option>
				</select>
			</div>
		</div>
	)
}
