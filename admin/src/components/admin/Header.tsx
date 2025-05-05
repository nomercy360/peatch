type TabType = 'users' | 'collaborations'

type HeaderProps = {
	activeTab: () => TabType
	setActiveTab: (tab: TabType) => void
}

export default function Header({ activeTab, setActiveTab }: HeaderProps) {
	return (
		<>
			<div class="flex mb-6 border-b">
				<button
					class={`px-4 py-2 font-medium ${activeTab() === 'users' ? 'text-primary border-b-2 border-primary' : 'text-muted-foreground'}`}
					onClick={() => setActiveTab('users')}
				>
					Users
				</button>
				<button
					class={`px-4 py-2 font-medium ${activeTab() === 'collaborations' ? 'text-primary border-b-2 border-primary' : 'text-muted-foreground'}`}
					onClick={() => setActiveTab('collaborations')}
				>
					Collaborations
				</button>
			</div>
		</>
	)
}
