// Mirrors the old web/src/utils/serverApi.ts fetchPublic helper, with an 8s
// timeout, swallowing any error and returning null so landing page sections
// can render a fallback instead of crashing.
//
// Always fetches the relative `/api/**` path (proxied to the backend by the
// routeRule in nuxt.config.ts) rather than reading the private
// `backendOrigin` runtime config directly: this composable can re-run on the
// client when a page is reached via client-side navigation (e.g. a NuxtLink
// click) rather than a fresh SSR page load, and private runtime config keys
// are stripped from the client bundle — using them there resolves to
// `undefined`, producing requests like "undefined/api/...".
export async function useServerFetch<T>(path: string, key: string) {
  const { data } = await useAsyncData<T | null>(key, async () => {
    try {
      const json = await $fetch<{ data: T }>(`/api${path}`, {
        signal: AbortSignal.timeout(8000),
      })
      return json.data
    } catch {
      return null
    }
  })

  return data
}
