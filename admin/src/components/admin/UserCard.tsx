import { For } from 'solid-js'
import { type UserResponse } from '~/gen/types'
import { type VerificationStatus, verificationStatus } from '~/gen'
import { cn } from '~/lib/utils'

type UserCardProps = {
	user: UserResponse
	updateUserStatus: (userId: string, status: VerificationStatus) => void
}

export default function UserCard({ user, updateUserStatus }: UserCardProps) {
	return (
		<div class="bg-card rounded-lg shadow p-4">
			{/* User header: avatar and name */}
			<div class="flex items-center mb-3">
				<div class="flex-shrink-0 h-12 w-12">
					{user.avatar_url ? (
						<img
							class="h-12 w-12 rounded-full"
							src={`https://assets.peatch.io/cdn-cgi/image/width=100/${user.avatar_url}`}
							alt={`${user.name}`}
						/>
					) : (
						<div class="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
							<span class="text-muted-foreground text-lg">{user.name?.[0]}</span>
						</div>
					)}
				</div>
				<div class="ml-3">
					<div class="text-base font-medium text-card-foreground">{user.name}</div>
					<div class="text-sm text-muted-foreground">@{user.username}</div>
				</div>
				<div class="ml-auto">
          <span class={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full
            ${user.verification_status === 'verified' ? 'bg-success text-success-foreground' :
						user.verification_status === 'pending' ? 'bg-warning text-warning-foreground' :
							user.verification_status === 'denied' ? 'bg-error text-error-foreground' :
								user.verification_status === 'blocked' ? 'bg-muted text-muted-foreground' :
									'bg-muted text-muted-foreground'}`}
					>
            {user.verification_status}
          </span>
				</div>
			</div>

			<div class="mb-3">
				<div class="text-sm font-medium text-card-foreground mb-1">{user.title}</div>
				<div class="text-sm text-muted-foreground mb-2">{user.description}</div>

				{user.location?.name && (
					<div class="text-sm text-foreground mb-2">
						<span class="font-medium">Location:</span> {user.location.name}
					</div>
				)}

				{user.badges && user.badges.length > 0 && (
					<div class="mb-2">
						<div class="text-xs text-muted-foreground mb-1">Badges:</div>
						<div class="flex flex-wrap gap-1">
							<For each={user.badges}>
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

				{user.opportunities && user.opportunities.length > 0 && (
					<div>
						<div class="text-xs text-muted-foreground mb-1">Opportunities:</div>
						<div class="flex flex-wrap gap-1">
							<For each={user.opportunities}>
								{(opportunity) => (
									<span
										class={cn('inline-flex items-center px-2 py-0.5 rounded text-xs font-medium text-white')}
										style={{ background: `#${opportunity.color}` }}
									>
                    {opportunity.text}
                  </span>
								)}
							</For>
						</div>
					</div>
				)}
			</div>

			{/* Actions */}
			<div class="mt-4 flex justify-end space-x-2">
				<button
					class="px-3 py-1 bg-success text-success-foreground rounded-md text-sm font-medium"
					onClick={() => updateUserStatus(user.id || '', verificationStatus.VerificationStatusVerified)}
				>
					Verify
				</button>
				<button
					class="px-3 py-1 bg-error text-error-foreground rounded-md text-sm font-medium"
					onClick={() => updateUserStatus(user.id || '', verificationStatus.VerificationStatusDenied)}
				>
					Deny
				</button>
				<button
					class="px-3 py-1 bg-muted text-muted-foreground rounded-md text-sm font-medium"
					onClick={() => updateUserStatus(user.id || '', verificationStatus.VerificationStatusBlocked)}
				>
					Block
				</button>
			</div>
		</div>
	)
}
