import axios from 'axios'

export default defineNuxtPlugin({
  name: 'axios',
  setup() {
    const config = useRuntimeConfig()

    const axiosInstance = axios.create({
      baseURL: config.public.apiUrl,
      withCredentials: true,
    })

    return {
      provide: { axios: axiosInstance },
    }
  },
})
