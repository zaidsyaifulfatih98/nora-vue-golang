import tailwindcss from '@tailwindcss/vite'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  modules: ['@pinia/nuxt', '@nuxt/icon', '@nuxtjs/i18n'],

  i18n: {
    strategy: 'no_prefix',
    defaultLocale: 'id',
    locales: [
      { code: 'id', name: 'Indonesia', language: 'id-ID', flag: '🇮🇩', file: 'id.json' },
      { code: 'en', name: 'English', language: 'en-US', flag: '🇬🇧', file: 'en.json' },
    ],
    langDir: 'locales',
    lazy: true,
    detectBrowserLanguage: false,
  },

  components: [{ path: '~/components', pathPrefix: false }],

  icon: {
    serverBundle: 'local',
    // Moved out of /api/* so the routeRules proxy below (which forwards all
    // /api/** to the Go backend) doesn't swallow Nuxt Icon's own API route.
    localApiEndpoint: '/_nuxt_icon',
  },

  imports: {
    dirs: ['composables/api'],
  },

  css: ['~/assets/css/globals.css'],
  vite: {
    plugins: [tailwindcss()],
  },

  runtimeConfig: {
    // server-only
    backendOrigin: process.env.NUXT_BACKEND_ORIGIN || 'http://localhost:8000',
    public: {
      apiUrl: process.env.NUXT_PUBLIC_API_URL || '/api',
      whatsappNumber: process.env.NUXT_PUBLIC_WHATSAPP_NUMBER || '',
    },
  },

  routeRules: {
    '/api/**': { proxy: `${process.env.NUXT_BACKEND_ORIGIN || 'http://localhost:8000'}/api/**` },
  },

  app: {
    head: {
      title: 'Nora Photobooth — Abadikan Momen Bahagiamu',
      meta: [
        { name: 'description', content: 'Sewa photobooth untuk pernikahan, ulang tahun, dan acara spesial lainnya.' },
      ],
      link: [
        { rel: 'icon', type: 'image/png', href: '/noraPhotobooth.png' },
      ],
    },
  },
})
