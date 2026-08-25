import type { CurrentUser } from '~/composables/api/auth'

// Runs before every /dashboard/* page render (both SSR and client nav).
// On the server, cookies must be forwarded explicitly since $fetch doesn't
// automatically carry the incoming request's cookies.
export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()
  const config = useRuntimeConfig()

  try {
    const baseURL = import.meta.server ? `${config.backendOrigin}/api` : config.public.apiUrl
    const headers = import.meta.server ? useRequestHeaders(['cookie']) : undefined

    const res = await $fetch<{ data: CurrentUser }>('/auth/me', {
      baseURL,
      credentials: 'include',
      headers,
    })

    authStore.setAuth(res.data)
    authStore.setChecked(true)
  } catch {
    authStore.clearAuth()
    authStore.setChecked(true)
    return navigateTo('/login')
  }
})
