<script setup lang="ts">
const route = useRoute()
const sidebarOpen = useDashboardSidebar()

const TITLES: Record<string, string> = {
  '/dashboard': 'Ringkasan',
  '/dashboard/finance': 'Keuangan',
  '/dashboard/packages': 'Paket & Harga',
  '/dashboard/features': 'Fitur Unggulan',
  '/dashboard/frame-templates': 'Template Frame',
  '/dashboard/photobooth-frames': 'Frame Digital Photobooth',
  '/dashboard/photobooth-results': 'Hasil Digital Photobooth',
  '/dashboard/voice-messages': 'Pesan Suara',
  '/dashboard/backdrops': 'Backdrop',
  '/dashboard/gallery': 'Galeri',
  '/dashboard/reviews': 'Review',
}

const title = computed(() => TITLES[route.path] ?? 'Dashboard')
</script>

<template>
  <header class="flex h-16 items-center justify-between border-b border-gray-100 bg-white px-4 sm:px-6">
    <div class="flex min-w-0 items-center gap-3">
      <button aria-label="Buka menu" class="text-gray-500 hover:text-gray-700 lg:hidden" @click="sidebarOpen = true">
        <Icon name="heroicons:bars-3" class="text-2xl" />
      </button>
      <h1 class="truncate text-base font-semibold text-gray-800 sm:text-lg">{{ title }}</h1>
    </div>

    <div class="flex shrink-0 items-center gap-4">
      <button class="relative rounded-lg p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700">
        <Icon name="fe:bell" class="text-xl" />
      </button>
    </div>
  </header>
</template>
