<script setup lang="ts">
const authStore = useAuthStore()
const { logout } = useAuthApi()
const route = useRoute()
const router = useRouter()

const menuMain = [
  { label: 'Ringkasan', href: '/dashboard', icon: 'fe:home' },
  { label: 'Keuangan', href: '/dashboard/finance', icon: 'fe:credit-card' },
  { label: 'Paket & Harga', href: '/dashboard/packages', icon: 'fe:package' },
  { label: 'Fitur Unggulan', href: '/dashboard/features', icon: 'fe:award' },
  { label: 'Template Frame', href: '/dashboard/frame-templates', icon: 'fe:layout' },
  { label: 'Backdrop', href: '/dashboard/backdrops', icon: 'fe:grid' },
  { label: 'Galeri', href: '/dashboard/gallery', icon: 'fe:image' },
  { label: 'Review', href: '/dashboard/reviews', icon: 'fe:star' },
]

function isActive(href: string) {
  return href === '/dashboard' ? route.path === '/dashboard' : route.path.startsWith(href)
}

async function handleLogout() {
  try {
    await logout()
  } finally {
    authStore.clearAuth()
    router.push('/login')
  }
}
</script>

<template>
  <aside class="flex h-screen flex-col border-r border-gray-100 bg-white">
    <div class="flex items-center gap-3 px-6 py-5">
      <div
        class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#C9A86A] to-[#E7C9A9] text-white font-bold"
      >
        N
      </div>
      <span class="text-lg font-semibold text-gray-800">Nora Photobooth</span>
    </div>

    <nav class="flex-1 px-4">
      <p class="mb-2 px-2 text-xs font-semibold uppercase tracking-wide text-gray-400">Menu</p>

      <ul class="space-y-1">
        <li v-for="item in menuMain" :key="item.href">
          <NuxtLink
            :to="item.href"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition"
            :class="isActive(item.href) ? 'bg-[#F1E4D6] text-[#8F6F3E]' : 'text-gray-600 hover:bg-gray-100'"
          >
            <Icon :name="item.icon" class="text-lg" />
            {{ item.label }}
          </NuxtLink>
        </li>
      </ul>
    </nav>

    <div class="border-t border-gray-100 px-4 py-4">
      <div class="flex items-center gap-3">
        <div
          class="flex h-10 w-10 items-center justify-center rounded-full bg-[#F1E4D6] text-sm font-bold text-[#8F6F3E]"
        >
          {{ authStore.user.firstName?.[0] }}{{ authStore.user.lastName?.[0] }}
        </div>
        <div class="flex-1">
          <p class="text-sm font-semibold text-gray-800">
            {{ authStore.user.firstName }} {{ authStore.user.lastName }}
          </p>
          <p class="text-xs text-gray-500">{{ authStore.user.role }}</p>
        </div>
        <button aria-label="Logout" class="text-gray-400 hover:text-red-500" @click="handleLogout">
          <Icon name="fe:log-out" />
        </button>
      </div>
    </div>
  </aside>
</template>
