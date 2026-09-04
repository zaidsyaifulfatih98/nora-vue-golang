// Tracks page views on the public site only — dashboard/login navigation is
// the admin's own usage, not visitor traffic, so it's excluded.
//
// dependsOn: 'axios' matters — Nuxt runs client plugins alphabetically by
// default, and "analytics" sorts before "axios", so without this the $axios
// instance useAnalyticsApi() relies on isn't provided yet.
export default defineNuxtPlugin({
  name: 'analytics',
  dependsOn: ['axios'],
  setup() {
    const { trackPageView } = useAnalyticsApi()
    const router = useRouter()

    function shouldTrack(path: string) {
      return !path.startsWith('/dashboard') && !path.startsWith('/login')
    }

    function track(path: string) {
      if (shouldTrack(path)) trackPageView(path)
    }

    track(router.currentRoute.value.path)
    router.afterEach((to) => track(to.path))
  },
})
