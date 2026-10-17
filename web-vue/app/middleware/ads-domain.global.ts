// Di ads.noraphotobooth.online, halaman utama (/) langsung menampilkan halaman promo.
export default defineNuxtRouteMiddleware((to) => {
  const host = useRequestURL().hostname
  if (host.startsWith('ads.') && to.path === '/') {
    return navigateTo('/promo', { redirectCode: 302 })
  }
})
