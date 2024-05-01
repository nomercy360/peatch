import {
  createSignal,
  createContext,
  useContext,
  createEffect,
  onCleanup,
} from 'solid-js';
import { useBackButton } from './useBackButton';
import { useLocation, useNavigate } from '@solidjs/router';

interface NavigationContext {
  navigateBack: () => void;
}

const Navigation = createContext<NavigationContext>({} as NavigationContext);

export function NavigationProvider(props: { children: any }) {
  const backButton = useBackButton();

  const navigate = useNavigate();
  const location = useLocation();

  const navigateBack = () => {
    console.log('location:', location.pathname);

    const state = location.state;

    !state && navigate('/');

    const deserialize = (state: any) => {
      try {
        return JSON.parse(state);
      } catch (e) {
        return state;
      }
    };

    const stateData = deserialize(state);

    if (stateData.from && location !== stateData.from) {
      console.log('navigating back to:', stateData.from);
      navigate(stateData.from);
    } else if (stateData.back) {
      console.log('navigating back');
      navigate(-1);
    } else {
      console.log('navigating back to root');
      navigate('/');
    }
  };

  createEffect(() => {
    backButton.hide();
    if (location.pathname !== '/') {
      backButton.setVisible();
      backButton.onClick(navigateBack);
    }
  });

  onCleanup(() => {
    backButton.hide();
    backButton.offClick(navigateBack);
  });

  const value: NavigationContext = {
    navigateBack,
  };

  return (
    <Navigation.Provider value={value}>{props.children}</Navigation.Provider>
  );
}

export function useNavigation() {
  return useContext(Navigation);
}
