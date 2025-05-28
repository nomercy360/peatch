import { splitProps, type ComponentProps } from 'solid-js'

import { cn } from '~/lib/utils'

type IconProps = ComponentProps<'svg'>

const Icon = (props: IconProps) => {
	const [, rest] = splitProps(props, ['class'])
	return (
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			class={cn('size-4', props.class)}
			{...rest}
		/>
	)
}

// ICONS

export function IconUsers(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M9 7m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0" />
			<path d="M3 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2" />
			<path d="M16 3.13a4 4 0 0 1 0 7.75" />
			<path d="M21 21v-2a4 4 0 0 0 -3 -3.85" />
		</Icon>
	)
}

export function IconUserCog(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M10 15H6a4 4 0 0 0-4 4v2" />
			<path d="m14.305 16.53.923-.382" />
			<path d="m15.228 13.852-.923-.383" />
			<path d="m16.852 12.228-.383-.923" />
			<path d="m16.852 17.772-.383.924" />
			<path d="m19.148 12.228.383-.923" />
			<path d="m19.53 18.696-.382-.924" />
			<path d="m20.772 13.852.924-.383" />
			<path d="m20.772 16.148.924.383" />
			<circle cx="18" cy="15" r="3" />
			<circle cx="9" cy="7" r="4" />
		</Icon>
	)
}

export function IconTarget(props: IconProps) {
	return (
		<Icon {...props}>
			<circle cx="12" cy="12" r="10" />
			<circle cx="12" cy="12" r="6" />
			<circle cx="12" cy="12" r="2" />
		</Icon>
	)
}

export function IconHandshake(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="m11 17 2 2a1 1 0 1 0 3-3" />
			<path
				d="m14 14 2.5 2.5a1 1 0 1 0 3-3l-3.88-3.88a3 3 0 0 0-4.24 0l-.88.88a1 1 0 1 1-3-3l2.81-2.81a5.79 5.79 0 0 1 7.06-.87l.47.28a2 2 0 0 0 1.42.25L21 4" />
			<path d="m21 3 1 11h-2" />
			<path d="M3 3 2 14l6.5 6.5a1 1 0 1 0 3-3" />
			<path d="M3 4h8" />
		</Icon>
	)
}

export function IconMapPin(props: IconProps) {
	return (
		<Icon {...props}>
			<path
				d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0" />
			<circle cx="12" cy="10" r="3" />
		</Icon>
	)
}

export function IconAward(props: IconProps) {
	return (
		<Icon {...props}>
			<path
				d="m15.477 12.89 1.515 8.526a.5.5 0 0 1-.81.47l-3.58-2.687a1 1 0 0 0-1.197 0l-3.586 2.686a.5.5 0 0 1-.81-.469l1.514-8.526" />
			<circle cx="12" cy="8" r="6" />
		</Icon>
	)
}

export function IconTrash(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M3 6h18" />
			<path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
			<path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
			<line x1="10" x2="10" y1="11" y2="17" />
			<line x1="14" x2="14" y1="11" y2="17" />
		</Icon>
	)
}
