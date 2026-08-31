<script setup lang="ts">
const authStore = useAuthStore()
const { logout } = useAuthApi()
const route = useRoute()
const router = useRouter()
const sidebarOpen = useDashboardSidebar()

const menuMain = [
  { label: 'Ringkasan', href: '/dashboard', icon: 'lucide:layout-dashboard' },
  { label: 'Keuangan', href: '/dashboard/finance', icon: 'lucide:wallet' },
  { label: 'Paket & Harga', href: '/dashboard/packages', icon: 'lucide:package' },
  { label: 'Fitur Unggulan', href: '/dashboard/features', icon: 'lucide:award' },
  { label: 'Template Frame', href: '/dashboard/frame-templates', icon: 'lucide:layout-template' },
  { label: 'Frame Digital Photobooth', href: '/dashboard/photobooth-frames', icon: 'lucide:scan-face' },
  { label: 'Hasil Digital Photobooth', href: '/dashboard/photobooth-results', icon: 'lucide:images' },
  { label: 'Pesan Suara', href: '/dashboard/voice-messages', icon: 'lucide:mic' },
  { label: 'Backdrop', href: '/dashboard/backdrops', icon: 'lucide:layout-grid' },
  { label: 'Galeri', href: '/dashboard/gallery', icon: 'lucide:image' },
  { label: 'Review', href: '/dashboard/reviews', icon: 'lucide:star' },
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
  <div
    v-if="sidebarOpen"
    class="fixed inset-0 z-30 bg-black/50 lg:hidden"
    aria-hidden="true"
    @click="sidebarOpen = false"
  />

  <aside
    class="fixed inset-y-0 left-0 z-40 flex h-screen w-72 shrink-0 -translate-x-full flex-col border-r border-gray-100 bg-white transition-transform duration-300 lg:static lg:z-auto lg:w-64 lg:translate-x-0"
    :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <div class="flex items-center justify-between gap-3 px-6 py-5">
      <div class="flex items-center gap-3">
        <span class="relative flex h-10 w-10 items-center justify-center overflow-hidden rounded-full shadow-md">
          <img src="/nora_logo.jpg" alt="Nora Photobooth" class="h-full w-full scale-75 object-contain" />
        </span>
        <span class="font-aloja text-xl tracking-wide text-[#1E2537]">Nora Photobooth</span>
      </div>
      <button aria-label="Tutup menu" class="text-gray-400 hover:text-gray-600 lg:hidden" @click="sidebarOpen = false">
        <Icon name="heroicons:x-mark" class="text-xl" />
      </button>
    </div>

    <nav class="flex-1 overflow-y-auto px-4">
      <p class="mb-2 px-2 text-xs font-semibold uppercase tracking-wide text-gray-400">Menu</p>

      <ul class="space-y-1">
        <li v-for="item in menuMain" :key="item.href">
          <NuxtLink
            :to="item.href"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition"
            :class="isActive(item.href) ? 'bg-[#f7f3eb] text-[#920f0f]' : 'text-gray-600 hover:bg-gray-100'"
            @click="sidebarOpen = false"
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
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-[#f7f3eb] text-sm font-bold text-[#920f0f]"
        >
          {{ authStore.user.firstName?.[0] }}{{ authStore.user.lastName?.[0] }}
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-semibold text-gray-800">
            {{ authStore.user.firstName }} {{ authStore.user.lastName }}
          </p>
          <p class="text-xs text-gray-500">{{ authStore.user.role }}</p>
        </div>
        <button aria-label="Logout" class="shrink-0 text-gray-400 hover:text-red-500" @click="handleLogout">
          <Icon name="lucide:log-out" />
        </button>
      </div>
    </div>
  </aside>
</template>
