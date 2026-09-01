<script setup lang="ts">
const { locale, locales, setLocale } = useI18n()

const FLAG_ICONS: Record<string, string> = {
  id: 'circle-flags:id',
  en: 'circle-flags:gb',
}

const availableLocales = computed(() => locales.value as { code: string; name: string }[])
</script>

<template>
  <div class="flex items-center gap-1.5">
    <button
      v-for="loc in availableLocales"
      :key="loc.code"
      type="button"
      :aria-label="loc.name"
      :title="loc.name"
      class="flex h-7 w-7 items-center justify-center rounded-full transition"
      :class="locale === loc.code ? 'ring-2 ring-[#920f0f] ring-offset-1' : 'opacity-50 grayscale hover:opacity-100 hover:grayscale-0'"
      @click="setLocale(loc.code as 'id' | 'en')"
    >
      <Icon :name="FLAG_ICONS[loc.code] ?? 'heroicons:flag'" class="h-full w-full" />
    </button>
  </div>
</template>
