import type {Component} from 'solid-js';
import {createEffect} from 'solid-js';

import logo from './logo.svg';
import styles from './App.module.css';


const App: Component = () => {
    createEffect(() => {
        console.log(window.Telegram.WebApp.initData);
    })

    return (
        <div class={styles.App}>
            <header class={styles.header}>
                <img src={logo} class={styles.logo} alt="logo"/>
                <p>
                    Edit <code>src/App.tsx</code> and save to reload.
                </p>
                <a
                    class={styles.link}
                    href="https://github.com/solidjs/solid"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Learn Solid
                </a>
            </header>
        </div>
    );
};

export default App;
