export function useAxios() {
  const { $axios } = useNuxtApp()
  return $axios as import('axios').AxiosInstance
}
