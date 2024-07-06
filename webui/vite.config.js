import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  base: "spa",
  server: {
    proxy: {
      // This is the path you want to proxy
      '/api': {
        target: 'http://localhost:8085',
       changeOrigin: true,
        secure: false, // if you are using an HTTPS API, set this to true
        // Configure how cookies and headers should be handled
        cookieDomainRewrite: {
          // Change the domain of the cookies to localhost
          '*': ''
        }
      }
    }
  }
})




