import { defineStore } from 'pinia'

export interface AuthUser {
  id?: string
  firstName: string
  lastName: string
  email?: string
  role: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: {
      id: undefined,
      firstName: '',
      lastName: '',
      email: undefined,
      role: '',
    } as AuthUser,
    isChecked: false,
  }),
  actions: {
    setAuth(user: AuthUser) {
      this.user = user
    },
    clearAuth() {
      this.user = { firstName: '', lastName: '', role: '' }
    },
    setChecked(checked: boolean) {
      this.isChecked = checked
    },
  },
})
