import { createContext, useContext } from 'solid-js';
import { HashNavigator } from '@tma.js/sdk';

const NavigatorContext = createContext();

export function NavigatorProvider(props: {
  navigator: HashNavigator;
  children: any;
}) {
  const navigator = props.navigator;

  return (
    <NavigatorContext.Provider value={navigator}>
      {props.children}
    </NavigatorContext.Provider>
  );
}

export function useNavigator() {
  return useContext(NavigatorContext);
}
