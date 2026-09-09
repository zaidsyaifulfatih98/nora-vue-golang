<script setup lang="ts">
import type { CustomerItem } from '~/composables/api/auth'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { registerCustomer, getCustomers } = useAuthApi()
const { t } = useI18n()
const authStore = useAuthStore()
const config = useRuntimeConfig()

// Only SUPER_ADMIN can onboard customers (also enforced by the backend) —
// anyone else lands back on the dashboard overview.
if (authStore.user.role !== 'SUPER_ADMIN') {
  await navigateTo('/dashboard')
}

const customers = ref<CustomerItem[]>([])
const loading = ref(true)
const submitting = ref(false)
const formError = ref('')
const copiedId = ref('')

const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const slug = ref('')

async function loadCustomers() {
  loading.value = true
  customers.value = await getCustomers('SOFTWARE_PHOTOBOOTH')
  loading.value = false
}
onMounted(loadCustomers)

const canSubmit = computed(
  () => Boolean(firstName.value && lastName.value && email.value && password.value && slug.value) && !submitting.value,
)

function publicLink(customerSlug: string) {
  const origin = import.meta.client ? window.location.origin : (config.public.apiUrl as string)
  return `${origin}/software-photobooth/${customerSlug}`
}

async function handleSubmit() {
  if (!canSubmit.value) return
  submitting.value = true
  formError.value = ''
  try {
    await registerCustomer(firstName.value, lastName.value, email.value, password.value, slug.value, 'SOFTWARE_PHOTOBOOTH')
    firstName.value = ''
    lastName.value = ''
    email.value = ''
    password.value = ''
    slug.value = ''
    await loadCustomers()
  } catch (err: any) {
    formError.value = err?.response?.data?.message || t('dashboard.softwareCustomers.saving')
  } finally {
    submitting.value = false
  }
}

async function copyLink(customer: CustomerItem) {
  try {
    await navigator.clipboard.writeText(publicLink(customer.slug))
    copiedId.value = customer.id
    setTimeout(() => {
      if (copiedId.value === customer.id) copiedId.value = ''
    }, 2000)
  } catch {
    // Clipboard API unavailable — the link is still visible to copy manually.
  }
}
</script>

<template>
  <div class="space-y-6">
    <p class="text-sm text-gray-500">{{ t('dashboard.softwareCustomers.subtitle') }}</p>

    <div class="rounded-2xl bg-white p-6 shadow-md">
      <h2 class="text-base font-semibold text-gray-800">{{ t('dashboard.softwareCustomers.newCustomer') }}</h2>

      <form class="mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2" @submit.prevent="handleSubmit">
        <input
          v-model="firstName"
          :placeholder="t('dashboard.customers.firstNamePlaceholder')"
          class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
        />
        <input
          v-model="lastName"
          :placeholder="t('dashboard.customers.lastNamePlaceholder')"
          class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
        />
        <input
          v-model="email"
          type="email"
          :placeholder="t('dashboard.customers.emailPlaceholder')"
          class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
        />
        <input
          v-model="password"
          type="password"
          :placeholder="t('dashboard.customers.passwordPlaceholder')"
          class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
        />
        <div class="sm:col-span-2">
          <input
            v-model="slug"
            :placeholder="t('dashboard.customers.slugPlaceholder')"
            class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
          <p class="mt-1 text-xs text-gray-400">{{ t('dashboard.customers.slugHint') }}</p>
        </div>

        <p v-if="formError" class="text-xs text-red-600 sm:col-span-2">{{ formError }}</p>

        <button
          type="submit"
          :disabled="!canSubmit"
          class="w-fit rounded-lg px-4 py-2 text-sm font-semibold text-white transition sm:col-span-2"
          :class="canSubmit ? 'bg-[#920f0f] hover:bg-[#7a0c0c]' : 'cursor-not-allowed bg-gray-300'"
        >
          {{ submitting ? t('dashboard.softwareCustomers.saving') : t('dashboard.softwareCustomers.addCustomer') }}
        </button>
      </form>
    </div>

    <div class="space-y-3">
      <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.softwareCustomers.loadingCustomers') }}</p>
      <p v-else-if="customers.length === 0" class="text-sm text-gray-400">{{ t('dashboard.softwareCustomers.empty') }}</p>

      <div v-for="customer in customers" :key="customer.id" class="rounded-2xl bg-white p-5 shadow-md">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <p class="font-semibold text-gray-800">{{ customer.firstName }} {{ customer.lastName }}</p>
            <p class="text-xs text-gray-400">{{ customer.email }}</p>
          </div>
        </div>
        <div class="mt-3 flex flex-wrap items-center gap-2 rounded-lg bg-gray-50 px-3 py-2 text-xs text-gray-600">
          <span class="min-w-0 flex-1 truncate">{{ publicLink(customer.slug) }}</span>
          <button type="button" class="shrink-0 font-semibold text-[#920f0f] hover:underline" @click="copyLink(customer)">
            {{ copiedId === customer.id ? t('dashboard.customers.linkCopied') : t('dashboard.customers.copyLink') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
