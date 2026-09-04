<script setup lang="ts">
const route = useRoute()
const sidebarOpen = useDashboardSidebar()
const { t } = useI18n()

const TITLE_KEYS: Record<string, string> = {
  '/dashboard': 'dashboard.titles.index',
  '/dashboard/analytics': 'dashboard.titles.analytics',
  '/dashboard/finance': 'dashboard.titles.finance',
  '/dashboard/packages': 'dashboard.titles.packages',
  '/dashboard/features': 'dashboard.titles.features',
  '/dashboard/frame-templates': 'dashboard.titles.frameTemplates',
  '/dashboard/photobooth-frames': 'dashboard.titles.photoboothFrames',
  '/dashboard/photobooth-results': 'dashboard.titles.photoboothResults',
  '/dashboard/voice-messages': 'dashboard.titles.voiceMessages',
  '/dashboard/backdrops': 'dashboard.titles.backdrops',
  '/dashboard/gallery': 'dashboard.titles.gallery',
  '/dashboard/reviews': 'dashboard.titles.reviews',
}

const title = computed(() => t(TITLE_KEYS[route.path] ?? 'dashboard.titles.default'))
</script>

<template>
  <header class="flex h-16 items-center justify-between border-b border-gray-100 bg-white px-4 sm:px-6">
    <div class="flex min-w-0 items-center gap-3">
      <button :aria-label="t('nav.openMenu')" class="text-gray-500 hover:text-gray-700 lg:hidden" @click="sidebarOpen = true">
        <Icon name="heroicons:bars-3" class="text-2xl" />
      </button>
      <h1 class="truncate text-base font-semibold text-gray-800 sm:text-lg">{{ title }}</h1>
    </div>

    <div class="flex shrink-0 items-center gap-4">
      <LanguageSwitcher />
      <button class="relative rounded-lg p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700">
        <Icon name="fe:bell" class="text-xl" />
      </button>
    </div>
  </header>
</template>
