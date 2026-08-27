<script setup lang="ts">
const NAV_LINKS = [
  { label: 'Beranda', href: '#hero' },
  { label: 'Keunggulan', href: '#keunggulan' },
  { label: 'Paket', href: '#paket' },
  { label: 'Galeri', href: '#galeri' },
  { label: 'Testimoni', href: '#testimoni' },
]

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
        <a
          v-for="link in NAV_LINKS"
          :key="link.href"
          :href="link.href"
          class="text-sm font-medium text-[#000000] transition hover:text-[#920f0f]"
        >
          {{ link.label }}
        </a>
      </div>

      <button aria-label="Buka menu" class="text-2xl text-[#1E2537] lg:hidden" @click="open = !open">
        <Icon :name="open ? 'heroicons:x-mark' : 'heroicons:bars-3'" />
      </button>
    </nav>

    <div v-if="open" class="flex flex-col gap-1 border-t border-[#E4E2DC] bg-[#FAF9F6] px-6 py-4 lg:hidden">
      <a
        v-for="link in NAV_LINKS"
        :key="link.href"
        :href="link.href"
        class="rounded-lg px-3 py-2.5 text-sm font-medium text-[#39445B] hover:bg-[#F0EFEA]"
        @click="open = false"
      >
        {{ link.label }}
      </a>
      <a
        href="#booking"
        class="mt-2 rounded-full bg-[#1E2537] px-6 py-2.5 text-center text-sm font-semibold text-[#FAF9F6]"
        @click="open = false"
      >
        Booking Sekarang
      </a>
    </div>
  </header>
</template>
