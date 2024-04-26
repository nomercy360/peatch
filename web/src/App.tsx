import { Component, createEffect } from 'solid-js';


export default function App(props: any) {
  createEffect(() => {
    window.Telegram.WebApp.ready();
    window.Telegram.WebApp.expand();

    console.log('App mounted');
  });

  return (
    <>
      {props.children}
    </>
  );
};
