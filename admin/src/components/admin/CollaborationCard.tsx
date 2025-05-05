import { For } from 'solid-js'
import { type CollaborationResponse } from '~/gen/types'
import { type VerificationStatus, verificationStatus } from '~/gen'
import { cn } from '~/lib/utils'

type CollaborationCardProps = {
	collab: CollaborationResponse
	updateCollaborationStatus: (userId: string, collabId: string, status: VerificationStatus) => void
}

export default function CollaborationCard({ collab, updateCollaborationStatus }: CollaborationCardProps) {
	return (
		<div class="bg-card rounded-lg shadow p-4">
			<div class="flex items-center justify-between mb-3">
				<div class="text-lg font-medium text-card-foreground">{collab.title}</div>
				<span class={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full
          ${collab.verification_status === 'verified' ? 'bg-success text-success-foreground' :
					collab.verification_status === 'pending' ? 'bg-warning text-warning-foreground' :
						collab.verification_status === 'denied' ? 'bg-error text-error-foreground' :
							collab.verification_status === 'blocked' ? 'bg-muted text-muted-foreground' :
								'bg-muted text-muted-foreground'}`}
				>
          {collab.verification_status}
        </span>
			</div>

			<div class="mb-3">
				<div class="text-sm text-muted-foreground mb-3">{collab.description}</div>

				<div class="flex items-center mb-3 p-2 bg-muted/50 rounded-lg">
					<div class="flex items-center">
						<div class="flex-shrink-0 h-8 w-8">
							{collab.user?.avatar_url ? (
								<img
									class="h-8 w-8 rounded-full object-cover"
									src={`https://assets.peatch.io/cdn-cgi/image/width=100/${collab.user.avatar_url}`}
									alt={`${collab.user.first_name} ${collab.user.last_name}`}
								/>
							) : (
								<div class="h-8 w-8 rounded-full bg-muted flex items-center justify-center">
                  <span
										class="text-muted-foreground">{collab.user?.first_name?.[0]}{collab.user?.last_name?.[0]}</span>
								</div>
							)}
						</div>
						<div class="ml-2">
							<div class="text-sm font-medium text-card-foreground">
								{collab.user?.first_name} {collab.user?.last_name}
							</div>
						</div>
					</div>
				</div>

				{/* Location and Payment Info */}
				<div class="flex flex-wrap gap-y-2 gap-x-4 mb-3">
					{collab.location?.name && (
						<div class="text-sm">
							<span class="font-medium text-foreground">Location:</span>
							<span class="text-muted-foreground">{collab.location.name}</span>
						</div>
					)}
					<div class="text-sm">
						<span class="font-medium text-foreground">Payment:</span>
						<span class="text-muted-foreground">{collab.is_payable ? 'Payable' : 'Non-payable'}</span>
					</div>
				</div>

				{/* Badges */}
				{collab.badges && collab.badges.length > 0 && (
					<div class="mb-2">
						<div class="text-xs text-muted-foreground mb-1">Badges:</div>
						<div class="flex flex-wrap gap-1">
							<For each={collab.badges}>
								{(badge) => (
									<span
										class={cn('inline-flex items-center px-2 py-0.5 rounded text-xs font-medium text-white')}
										style={{ background: `#${badge.color}` }}
									>
                    {badge.text}
                  </span>
								)}
							</For>
						</div>
					</div>
				)}

				{/* Opportunity */}
				{collab.opportunity && (
					<div>
						<div class="text-xs text-muted-foreground mb-1">Opportunity:</div>
						<span
							class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium text-white"
							style={{ background: `#${collab.opportunity.color}` }}
						>
              {collab.opportunity?.text}
            </span>
					</div>
				)}
			</div>

			{/* Actions */}
			<div class="mt-4 flex justify-end space-x-2">
				<button
					class="px-3 py-1 bg-success text-success-foreground rounded-md text-sm font-medium hover:bg-success/30"
					onClick={() => updateCollaborationStatus(collab.user_id || '', collab.id || '', verificationStatus.VerificationStatusVerified)}
				>
					Verify
				</button>
				<button
					class="px-3 py-1 bg-error text-error-foreground rounded-md text-sm font-medium hover:bg-error/30"
					onClick={() => updateCollaborationStatus(collab.user_id || '', collab.id || '', verificationStatus.VerificationStatusDenied)}
				>
					Deny
				</button>
				<button
					class="px-3 py-1 bg-muted text-muted-foreground rounded-md text-sm font-medium hover:bg-muted/80"
					onClick={() => updateCollaborationStatus(collab.user_id || '', collab.id || '', verificationStatus.VerificationStatusBlocked)}
				>
					Block
				</button>
			</div>
		</div>
	)
}
