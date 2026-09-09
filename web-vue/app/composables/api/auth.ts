export interface CurrentUser {
  id: string
  firstName: string
  lastName: string
  email: string
  role: string
}

export interface CustomerItem {
  id: string
  firstName: string
  lastName: string
  email: string
  slug: string
  createdAt: string
}

export type CustomerRole = 'DIGITAL_PHOTOBOOTH' | 'SOFTWARE_PHOTOBOOTH'

export function useAuthApi() {
  const axios = useAxios()

  const login = (email: string, password: string) =>
    axios.post('/auth/login', { email, password }).then((r) => r.data.data as CurrentUser)

  const logout = () => axios.post('/auth/logout')

  const getCurrentUser = () => axios.get('/auth/me').then((r) => r.data.data as CurrentUser)

  // Onboards a new customer account (DIGITAL_PHOTOBOOTH or
  // SOFTWARE_PHOTOBOOTH), superadmin-only (enforced by the backend).
  const registerCustomer = (
    firstName: string,
    lastName: string,
    email: string,
    password: string,
    slug: string,
    role: CustomerRole = 'DIGITAL_PHOTOBOOTH',
  ) =>
    axios.post('/auth/register', { firstName, lastName, email, password, slug, role }).then((r) => r.data.data as CustomerItem)

  const getCustomers = (role: CustomerRole = 'DIGITAL_PHOTOBOOTH') =>
    axios.get(`/auth/customers?role=${role}`).then((r) => r.data.data as CustomerItem[])

  return { login, logout, getCurrentUser, registerCustomer, getCustomers }
}
