export default function Badge(props: {
  icon: string
  name: string
  color: string
}) {
  return (
    <div
      class="flex h-5 flex-row items-center justify-center gap-[5px] rounded px-2.5"
      style={{ 'background-color': `#${props.color}` }}
    >
      <span class="material-symbols-rounded text-[10px] text-white">
        {String.fromCharCode(parseInt(props.icon, 16))}
      </span>
      <p class="text-xs font-semibold text-white">{props.name}</p>
    </div>
  )
}
