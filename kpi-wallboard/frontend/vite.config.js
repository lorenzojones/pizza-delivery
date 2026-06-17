import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// The frontend talks to the Go backend through `/api`.
// - In dev, Vite proxies `/api` to the backend (default http://localhost:8080).
// - In prod (nginx), nginx proxies `/api` to the `backend` service.
// - VITE_API_BASE can override the base at build/run time (default same-origin "").
export default defineConfig(() => {
  const proxyTarget = process.env.VITE_DEV_API_TARGET || 'http://localhost:8080';

  return {
    plugins: [svelte()],
    server: {
      host: true,
      port: 5173,
      proxy: {
        '/api': {
          target: proxyTarget,
          changeOrigin: true,
        },
      },
    },
    preview: {
      host: true,
      port: 4173,
    },
    build: {
      outDir: 'dist',
      sourcemap: false,
    },
  };
});
