import {
  createContext,
  createEffect,
  JSX,
  onCleanup,
  useContext,
} from 'solid-js'
import { useBackButton } from './useBackButton'
import { useLocation, useNavigate } from '@solidjs/router'

interface NavigationContext {
  navigateBack: () => void
}

const Navigation = createContext<NavigationContext>({} as NavigationContext)

export function NavigationProvider(props: { children: JSX.Element }) {
  const backButton = useBackButton()

  const navigate = useNavigate()
  const location = useLocation()

  const navigateBack = () => {
    const state = location.state

    !state && navigate('/')

    const deserialize = (state: unknown) => {
      if (typeof state === 'string') {
        try {
          return JSON.parse(state)
        } catch (e) {
          console.error('Error deserializing state:', e)
          return state
        }
      }
      return state
    }

    const stateData = deserialize(state)

    const isObject = (value: unknown) => {
      return value && typeof value === 'object' && value.constructor === Object
    }

    if (isObject(stateData) && stateData.from) {
      console.log('navigating back to:', stateData.from)
      navigate(stateData.from)
    } else if (isObject(stateData) && stateData.back) {
      console.log('navigating back')
      navigate(-1)
    } else {
      console.log('navigating back to root')
      navigate('/')
    }

    console.log('stateData:', stateData)

    if (isObject(stateData) && stateData.scroll) {
      console.log('scrolling to:', stateData.scroll)
      setTimeout(() => {
        window.scrollTo(0, stateData.scroll)
      }, 0)
    }
  }

  createEffect(() => {
    if (location.pathname !== '/') {
      backButton.setVisible()
      backButton.onClick(navigateBack)
    } else {
      backButton.hide()
    }
  })

  onCleanup(() => {
    backButton.hide()
    backButton.offClick(navigateBack)
  })

  const value: NavigationContext = {
    navigateBack,
  }

  return (
    <Navigation.Provider value={value}>{props.children}</Navigation.Provider>
  )
}

export function useNavigation() {
  return useContext(Navigation)
}
