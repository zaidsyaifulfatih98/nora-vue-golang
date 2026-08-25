export interface CurrentUser {
  id: string
  firstName: string
  lastName: string
  email: string
  role: string
}

export function useAuthApi() {
  const axios = useAxios()

  const login = (email: string, password: string) =>
    axios.post('/auth/login', { email, password }).then((r) => r.data.data as CurrentUser)

  const logout = () => axios.post('/auth/logout')

  const getCurrentUser = () => axios.get('/auth/me').then((r) => r.data.data as CurrentUser)

  return { login, logout, getCurrentUser }
}
