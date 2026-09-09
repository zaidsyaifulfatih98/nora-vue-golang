import type { CurrentUser } from '~/composables/api/auth'

// Runs before every /dashboard/* page render (both SSR and client nav).
// On the server, cookies must be forwarded explicitly since $fetch doesn't
// automatically carry the incoming request's cookies.
export default defineNuxtRouteMiddleware(async (to) => {
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

    // The main overview shows finance/analytics a customer account can't
    // access — send them straight to their own section.
    if (['DIGITAL_PHOTOBOOTH', 'SOFTWARE_PHOTOBOOTH'].includes(res.data.role) && to.path === '/dashboard') {
      return navigateTo('/dashboard/photobooth-frames')
    }
  } catch {
    authStore.clearAuth()
    authStore.setChecked(true)
    return navigateTo('/login')
  }
})
