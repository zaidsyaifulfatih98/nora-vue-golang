<script setup lang="ts">
const { t } = useI18n()

const NAV_LINKS = computed(() => [
  { label: t('nav.home'), href: '/#hero' },
  { label: t('nav.digitalPhotobooth'), href: '/digital-photobooth' },
  { label: t('nav.softwarePhotobooth'), href: '/software-photobooth' },
  { label: t('nav.printPhoto'), href: '/print-photo' },
])

const open = ref(false)
const scrolled = ref(false)

function onScroll() {
  scrolled.value = window.scrollY > 12
}

onMounted(() => window.addEventListener('scroll', onScroll))
onBeforeUnmount(() => window.removeEventListener('scroll', onScroll))
</script>

<template>
  <header
    class="fixed top-0 z-50 w-full transition-all duration-300"
    :class="scrolled ? 'bg-[#FFFFFF]/90 shadow-sm backdrop-blur-md' : 'bg-transparent'"
  >
    <nav class="mx-auto flex max-w-7xl items-center justify-between px-6 py-4 lg:px-10">
      <a href="#hero" class="flex items-center gap-2">
        <span class="relative flex h-9 w-9 items-center justify-center overflow-hidden rounded-full shadow-md">
          <img src="/nora_logo.jpg" alt="Nora Photobooth" class="h-full w-full scale-75 object-contain" />
        </span>
        <span class="font-aloja text-2xl tracking-wide text-[#1E2537]">
          Nora <span class="text-[#000000]">Photobooth</span>
        </span>
      </a>

      <div class="hidden items-center gap-8 lg:flex">
        <NuxtLink
          v-for="link in NAV_LINKS"
          :key="link.href"
          :to="link.href"
          class="text-sm font-medium text-[#000000] transition hover:text-[#920f0f]"
        >
          {{ link.label }}
        </NuxtLink>
        <LanguageSwitcher />
        <NuxtLink
          to="/login"
          class="flex items-center gap-2 rounded-full bg-[#1E2537] px-5 py-2 text-sm font-semibold text-[#FAF9F6] transition hover:bg-[#920f0f]"
        >
          <Icon name="lucide:log-in" />
          {{ t('nav.login') }}
        </NuxtLink>
      </div>

      <div class="flex items-center gap-3 lg:hidden">
        <LanguageSwitcher />
        <button :aria-label="t('nav.openMenu')" class="text-2xl text-[#1E2537]" @click="open = !open">
          <Icon :name="open ? 'heroicons:x-mark' : 'heroicons:bars-3'" />
        </button>
      </div>
    </nav>

    <div v-if="open" class="flex flex-col gap-1 border-t border-[#E4E2DC] bg-[#FAF9F6] px-6 py-4 lg:hidden">
      <NuxtLink
        v-for="link in NAV_LINKS"
        :key="link.href"
        :to="link.href"
        class="rounded-lg px-3 py-2.5 text-sm font-medium text-[#39445B] hover:bg-[#F0EFEA]"
        @click="open = false"
      >
        {{ link.label }}
      </NuxtLink>
      <NuxtLink
        to="/login"
        class="mt-2 flex items-center justify-center gap-2 rounded-full bg-[#1E2537] px-6 py-2.5 text-sm font-semibold text-[#FAF9F6] transition hover:bg-[#920f0f]"
        @click="open = false"
      >
        <Icon name="lucide:log-in" />
        {{ t('nav.login') }}
      </NuxtLink>
    </div>
  </header>
</template>
