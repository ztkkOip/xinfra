import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  base: '/auth/',
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/auth/api': 'http://127.0.0.1:8083',
      '/healthz': 'http://127.0.0.1:8083',
      '/readyz': 'http://127.0.0.1:8083',
    },
  },
})
