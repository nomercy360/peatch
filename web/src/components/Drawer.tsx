// import type { Component, ComponentProps } from 'solid-js'
// import { splitProps } from 'solid-js'
//
// import * as DrawerPrimitive from 'corvu/drawer'
//
// const Drawer = DrawerPrimitive.Root
//
// const DrawerTrigger = DrawerPrimitive.Trigger
//
// const DrawerPortal = DrawerPrimitive.Portal
//
// const DrawerClose = DrawerPrimitive.Close
//
// const DrawerOverlay: Component<DrawerPrimitive.OverlayProps> = props => {
// 	const [, rest] = splitProps(props, ['class'])
// 	const drawerContext = DrawerPrimitive.useContext()
// 	return (
// 		<DrawerPrimitive.Overlay
// 			class="fixed inset-0 z-50 data-[transitioning]:transition-colors data-[transitioning]:duration-300"
// 			style={{
// 				'background-color': `rgb(0 0 0 / ${0.8 * drawerContext.openPercentage()})`,
// 			}}
// 			{...rest}
// 		/>
// 	)
// }
//
// const DrawerContent: Component<DrawerPrimitive.ContentProps> = props => {
// 	const [, rest] = splitProps(props, ['class', 'children'])
// 	return (
// 		<DrawerPortal>
// 			<DrawerOverlay />
// 			<DrawerPrimitive.Content
// 				class="fixed inset-x-0 bottom-0 z-50 mt-24 flex h-auto flex-col rounded-t-[10px] border bg-secondary after:absolute after:inset-x-0 after:top-full after:h-1/2 after:bg-inherit data-[transitioning]:transition-transform data-[transitioning]:duration-300 md:select-none"
// 				{...rest}
// 			>
// 				<div class="bg-muted mx-auto mt-4 h-2 w-[100px] rounded-full" />
// 				{props.children}
// 			</DrawerPrimitive.Content>
// 		</DrawerPortal>
// 	)
// }
//
// const DrawerHeader: Component<ComponentProps<'div'>> = props => {
// 	const [, rest] = splitProps(props, ['class'])
// 	return (
// 		<div class="flex w-full flex-row items-center justify-between px-4">
// 			{props.children}
// 			<DrawerClose>
// 				<button>Close</button>
// 			</DrawerClose>
// 		</div>
// 	)
// }
//
// const DrawerTitle: Component<DrawerPrimitive.LabelProps> = props => {
// 	const [, rest] = splitProps(props, ['class'])
// 	return (
// 		<DrawerPrimitive.Label
// 			class="text-lg font-semibold leading-none tracking-tight text-main"
// 			{...rest}
// 		/>
// 	)
// }
//
// export {
// 	Drawer,
// 	DrawerPortal,
// 	DrawerOverlay,
// 	DrawerTrigger,
// 	DrawerClose,
// 	DrawerContent,
// 	DrawerHeader,
// 	DrawerTitle,
// }
