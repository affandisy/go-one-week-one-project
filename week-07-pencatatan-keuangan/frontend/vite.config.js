import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
  plugins: [
    react(),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'masked-icon.svg'],
      manifest: {
        name: 'Keuangan Sederhana',
        short_name: 'Keuangan',
        description: 'Aplikasi pencatatan keuangan simpel 3-klik',
        theme_color: '#1e40af', // Warna blue-800 sesuai header kita
        background_color: '#f9fafb', // Warna gray-50
        display: 'standalone', // Membuatnya tampil full screen seperti aplikasi asli
        icons: [
          {
            src: 'pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png'
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png'
          }
        ]
      }
    })
  ]
})