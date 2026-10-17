<script setup lang="ts">
import type { GalleryPhotoItem } from '~/composables/api/gallery'

// ===== Teks (ubah di sini) =====
const EYEBROW = 'Our Gallery'
const TITLE = 'Sekilas Keseruan di Setiap Bidikan'
const DESCRIPTION = 'Cuplikan suasana photobooth dari berbagai resepsi yang pernah kami temani.'
const ALT_FALLBACK = 'Momen Nora Photobooth'
// ================================

const PLACEHOLDER_TILES = [
  { gradient: 'from-[#F3D9C0] to-[#E7B98F]' },
  { gradient: 'from-[#F3D9E0] to-[#E3AFC2]' },
  { gradient: 'from-[#920f0f] to-[#8F6F3E]' },
  { gradient: 'from-[#D7D4CC] to-[#B99B76]' },
  { gradient: 'from-[#E7C9A9] to-[#D1A374]' },
  { gradient: 'from-[#D9E5CE] to-[#7F9B7E]' },
  { gradient: 'from-[#F3D9E0] to-[#EAC3CE]' },
  { gradient: 'from-[#F0EFEA] to-[#DCB98A]' },
]

const photos = await useServerFetch<GalleryPhotoItem[]>('/gallery', 'landing-gallery')
const hasPhotos = computed(() => Boolean(photos.value && photos.value.length > 0))
const itemCount = computed(() => (hasPhotos.value ? photos.value!.length : PLACEHOLDER_TILES.length))

const track = ref<HTMLElement | null>(null)

function scrollByPage(dir: 1 | -1) {
  track.value?.scrollBy({ left: dir * track.value.clientWidth, behavior: 'smooth' })
}
</script>

<template>
  <!-- Mobile: background putih, hanya eyebrow, foto besar 1 per slide. Desktop: seperti homepage. -->
  <section id="galeri" class="relative bg-white py-10 md:bg-[#f7f3eb] md:py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ EYEBROW }}</span>
        <h2 class="mt-3 hidden font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl md:block">{{ TITLE }}</h2>
        <p class="mt-4 hidden font-poppins text-[#57607A] md:block">{{ DESCRIPTION }}</p>
      </div>

      <div class="relative mt-6 md:mt-16">
        <div ref="track" class="flex snap-x snap-mandatory gap-4 overflow-x-auto scroll-smooth pb-2 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
          <template v-if="hasPhotos">
            <div
              v-for="photo in photos"
              :key="photo.id"
              class="group relative aspect-[2/3] w-[85%] flex-none snap-start overflow-hidden rounded-2xl shadow-sm ring-1 ring-[#E4E2DC]/70 sm:w-[calc(33.333%-0.667rem)] lg:w-[calc(25%-0.75rem)]"
            >
              <img :src="photo.url" :alt="photo.caption ?? ALT_FALLBACK" class="h-full w-full object-cover transition duration-300 group-hover:scale-105" />
            </div>
          </template>
          <template v-else>
            <div
              v-for="(tile, i) in PLACEHOLDER_TILES"
              :key="i"
              class="group relative aspect-[2/3] w-[85%] flex-none snap-start overflow-hidden rounded-2xl bg-gradient-to-br shadow-sm ring-1 ring-[#E4E2DC]/70 sm:w-[calc(33.333%-0.667rem)] lg:w-[calc(25%-0.75rem)]"
              :class="tile.gradient"
            >
              <div class="absolute inset-0 flex items-center justify-center bg-black/10 opacity-0 transition group-hover:opacity-100">
                <Icon name="fa6-solid:heart" class="text-3xl text-white drop-shadow" />
              </div>
            </div>
          </template>
        </div>

        <!-- Panah: selalu ada di mobile; di desktop hanya kalau foto > 4 -->
        <button
          v-if="itemCount > 1"
          type="button"
          aria-label="Sebelumnya"
          class="absolute top-1/2 left-2 flex -translate-y-1/2 rounded-full bg-white p-3 text-[#920f0f] shadow-md ring-1 ring-[#E4E2DC] transition hover:bg-[#920f0f] hover:text-white sm:-left-4"
          :class="itemCount > 4 ? '' : 'sm:hidden'"
          @click="scrollByPage(-1)"
        >
          <Icon name="fa6-solid:chevron-left" />
        </button>
        <button
          v-if="itemCount > 1"
          type="button"
          aria-label="Selanjutnya"
          class="absolute top-1/2 right-2 flex -translate-y-1/2 rounded-full bg-white p-3 text-[#920f0f] shadow-md ring-1 ring-[#E4E2DC] transition hover:bg-[#920f0f] hover:text-white sm:-right-4"
          :class="itemCount > 4 ? '' : 'sm:hidden'"
          @click="scrollByPage(1)"
        >
          <Icon name="fa6-solid:chevron-right" />
        </button>
      </div>
    </div>
  </section>
</template>
