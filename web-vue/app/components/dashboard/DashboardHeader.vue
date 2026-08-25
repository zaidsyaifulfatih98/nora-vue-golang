<script setup lang="ts">
const route = useRoute()

const TITLES: Record<string, string> = {
  '/dashboard': 'Ringkasan',
  '/dashboard/finance': 'Keuangan',
  '/dashboard/packages': 'Paket & Harga',
  '/dashboard/features': 'Fitur Unggulan',
  '/dashboard/frame-templates': 'Template Frame',
  '/dashboard/backdrops': 'Backdrop',
  '/dashboard/gallery': 'Galeri',
  '/dashboard/reviews': 'Review',
}

const title = computed(() => TITLES[route.path] ?? 'Dashboard')
</script>

<template>
  <header class="flex h-16 items-center justify-between border-b border-gray-100 bg-white px-6">
    <h1 class="text-lg font-semibold text-gray-800">{{ title }}</h1>

    <div class="flex items-center gap-4">
      <button class="relative rounded-lg p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700">
        <Icon name="fe:bell" class="text-xl" />
      </button>
    </div>
  </header>
</template>
