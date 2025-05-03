export function ProgressBar(props: { screen: number; totalScreens: number }) {
	return (
		<div class="h-1.5 w-[160px] rounded-lg bg-border">
			<div
				class="h-1.5 rounded-l-lg bg-accent"
				classList={{ 'rounded-r-lg': props.screen === props.totalScreens }}
				style={`width: ${(props.screen / props.totalScreens) * 100}%`}
			></div>
		</div>
	)
}
