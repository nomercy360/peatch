import type { AnchorProps } from '@solidjs/router'
import { A } from '@solidjs/router'
import type { Component } from 'solid-js'

export const Link: Component<AnchorProps> = props => {
  const onClick = (e: MouseEvent) => {
    const targetUrl = new URL(props.href, window.location.toString())
    const currentUrl = new URL(window.location.toString())
    const isExternal =
      targetUrl.protocol !== currentUrl.protocol ||
      targetUrl.host !== currentUrl.host

    const isTelegramLink =
      targetUrl.protocol === 'https:' &&
      targetUrl.hostname === 't.me' &&
      targetUrl.pathname.startsWith('/')

    if (isTelegramLink) {
      e.preventDefault()
      window.Telegram.WebApp.openTelegramLink(targetUrl.toString())
    } else if (isExternal) {
      e.preventDefault()
      window.Telegram.WebApp.openLink(targetUrl.toString())
    }
  }

  return <A {...props} onClick={onClick} class={props.class} />
}
