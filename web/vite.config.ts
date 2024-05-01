import {defineConfig} from 'vite';
import solidPlugin from 'vite-plugin-solid';
// import devtools from 'solid-devtools/vite';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
    plugins: [
        /*
        Uncomment the following line to enable solid-devtools.
        For more info see https://github.com/thetarnav/solid-devtools/tree/main/packages/extension#readme
        */
        // devtools(),
        solidPlugin(),
        tsconfigPaths(),
    ],
    server: {
        port: 3000,
        host: '127.0.0.1',
        https: {
            key: 'certs/key.pem',
            cert: 'certs/cert.pem',
        }
    },
    build: {
        target: 'esnext',
    },
});
