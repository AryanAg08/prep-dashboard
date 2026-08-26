import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// The dev server proxies /api to the Go backend so the frontend can call
// relative URLs without CORS headaches during development.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  // `vite preview` serves the production build (used by PM2 in deploy).
  // It proxies /api to the backend so the app's relative API calls work.
  preview: {
    port: Number(process.env.FRONTEND_PORT) || 9001,
    host: true,
    proxy: {
      '/api': {
        target: process.env.BACKEND_URL || 'http://localhost:9000',
        changeOrigin: true,
      },
    },
  },
})
