// Mirrors the old web/src/utils/serverApi.ts fetchPublic helper: fetches
// directly from the backend origin (bypassing the client-side /api proxy),
// with an 8s timeout, swallowing any error and returning null so landing
// page sections can render a fallback instead of crashing.
export async function useServerFetch<T>(path: string, key: string) {
  const config = useRuntimeConfig()

  const { data } = await useAsyncData<T | null>(key, async () => {
    try {
      const json = await $fetch<{ data: T }>(`${config.backendOrigin}/api${path}`, {
        signal: AbortSignal.timeout(8000),
      })
      return json.data
    } catch {
      return null
    }
  })

  return data
}
