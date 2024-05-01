export function ProgressBar(props: { screen: number; totalScreens: number }) {
  return (
    <div class="h-1.5 w-[160px] rounded-lg bg-main">
      <div
        class="bg-accent h-1.5 rounded-lg"
        style={`width: ${(props.screen / props.totalScreens) * 100}%`}
      ></div>
    </div>
  );
}
