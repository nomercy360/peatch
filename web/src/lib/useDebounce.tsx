import { onCleanup } from 'solid-js'

export default function useDebounce(
	signalSetter: (value: any) => void,
	delay: number,
) {
	let timerHandle: number

	function debouncedSignalSetter(value: any) {
		clearTimeout(timerHandle)
		timerHandle = setTimeout(() => signalSetter(value), delay)
	}

	onCleanup(() => clearInterval(timerHandle))
	return debouncedSignalSetter
}
