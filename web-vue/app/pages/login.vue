<script setup lang="ts">
import { useForm, useField } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { loginSchema } from '~/schemas/login.schema'

const authStore = useAuthStore()
const { login } = useAuthApi()
const router = useRouter()

const showPassword = ref(false)
const loginError = ref('')
const submitting = ref(false)

const { handleSubmit } = useForm({
  validationSchema: toTypedSchema(loginSchema),
  initialValues: { email: '', password: '' },
})

const { value: email, errorMessage: emailError } = useField<string>('email')
const { value: password, errorMessage: passwordError } = useField<string>('password')

const onSubmit = handleSubmit(async (values) => {
  loginError.value = ''
  submitting.value = true
  try {
    const user = await login(values.email, values.password)
    authStore.setAuth(user)
    router.push('/dashboard')
  } catch (error: any) {
    loginError.value = error?.response?.data?.message ?? 'Gagal login, silakan coba lagi'
  } finally {
    submitting.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-[#920f0f] flex items-center justify-center px-4">
    <div class="w-full max-w-md rounded-2xl bg-white shadow-xl p-8">
      <div class="flex flex-col items-center mb-6">
        <span class="relative flex h-12 w-12 items-center justify-center overflow-hidden rounded-full shadow-md">
          <img src="/nora_logo.jpg" alt="Nora Photobooth" class="h-full w-full scale-75 object-contain" />
        </span>
        <h1 class="mt-2 font-aloja text-2xl tracking-wide text-[#1E2537]">Nora Photobooth</h1>
        <p class="text-sm text-gray-500">Admin Dashboard</p>
      </div>

      <h2 class="mb-6 text-center text-lg font-semibold text-gray-800">Login to your account</h2>

      <form @submit="onSubmit">
        <div class="mb-4">
          <label class="mb-1 block text-sm font-medium text-gray-700">Email Address</label>
          <input
            v-model="email"
            type="email"
            placeholder="name@company.com"
            class="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
          <p class="text-xs font-bold text-red-500">{{ emailError }}</p>
        </div>

        <div class="mb-4">
          <label class="mb-1 block text-sm font-medium text-gray-700">Password</label>

          <div class="relative">
            <input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Enter your password"
              class="w-full rounded-lg border border-gray-300 px-4 py-2.5 pr-10 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
            />
            <p class="text-xs font-bold text-red-500">{{ passwordError }}</p>

            <button
              type="button"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              @click="showPassword = !showPassword"
            >
              <Icon :name="showPassword ? 'heroicons:eye-slash' : 'heroicons:eye'" class="h-5 w-5" />
            </button>
          </div>
        </div>

        <div class="mb-6 flex items-center justify-between">
          <label class="flex items-center gap-2 text-sm text-gray-600">
            <input type="checkbox" class="rounded border-gray-300 text-[#920f0f] focus:ring-[#920f0f]" />
            Remember me
          </label>
        </div>

        <p v-if="loginError" class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm font-medium text-red-600">
          {{ loginError }}
        </p>

        <button
          type="submit"
          :disabled="submitting"
          class="mb-6 flex w-full items-center justify-center gap-2 rounded-xl bg-[#920f0f] py-2.5 text-sm font-semibold text-white hover:bg-[#7a0c0c] transition disabled:opacity-60"
        >
          {{ submitting ? 'Logging in...' : 'Login →' }}
        </button>
      </form>

      <div class="flex items-center justify-between text-xs text-gray-400">
        <span>Nora Photobooth Admin</span>
        <div class="flex gap-4">
          <a href="#" class="hover:text-gray-600">Support</a>
          <a href="#" class="hover:text-gray-600">Privacy</a>
        </div>
      </div>
    </div>
  </div>
</template>
