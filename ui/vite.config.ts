import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy:{
      '/ws': {target: 'http://localhost:3000', changeOrigin: true, ws:true},
      '/ping': {target: 'http://localhost:3000', changeOrigin: true},
    },
  },
})
