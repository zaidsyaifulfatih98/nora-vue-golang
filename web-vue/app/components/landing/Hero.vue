
<script setup lang="ts">
import type { GalleryPhotoItem } from '~/composables/api/gallery'

const { t } = useI18n()
const config = useRuntimeConfig()

const photos = await useServerFetch<GalleryPhotoItem[]>('/gallery', 'hero-gallery')
const heroPhotos = computed(() => (photos.value ? photos.value.slice(0, 5) : null))
const waHref = computed(
  () => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(t('landing.whatsappCta.message'))}`,
)
</script>

<template>
  <section id="hero" class="relative overflow-hidden bg-white pt-32 pb-24 lg:pt-40 lg:pb-32">
    <div class="mx-auto grid max-w-7xl items-center gap-16 px-6 lg:grid-cols-2 lg:px-10">
      <div class="relative z-10">
        <span class="inline-flex items-center gap-2 rounded-full bg-white/70 px-4 py-1.5 font-dm-sans text-xs font-semibold tracking-wide text-[#920f0f] shadow-sm ring-1 ring-[#E4E2DC]">
          <Icon name="heroicons:sparkles" class="text-base" />
          {{ t('landing.hero.badge') }}
        </span>

        <h1 class="mt-6 font-dm-serif text-4xl leading-tight font-bold text-[#000000] sm:text-5xl lg:text-6xl">
          {{ t('landing.hero.title') }}
        </h1>

        <p class="mt-6 max-w-lg font-poppins text-base leading-relaxed text-[#57607A] sm:text-lg">
          {{ t('landing.hero.descriptionP1') }}
          <br />
          <br />
          {{ t('landing.hero.descriptionP2') }}
        </p>

        <div class="mt-8 flex flex-wrap items-center gap-4">
          <a
            :href="waHref"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-full bg-[#920f0f] px-8 py-3.5 text-sm font-semibold text-[#FAF9F6] shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5 hover:bg-[#920f0f]"
          >
            {{ t('landing.hero.cta') }}
          </a>
        </div>
      </div>

      <div class="relative z-10 flex justify-center lg:justify-end">
        <div class="relative w-full max-w-[400px]">
          <HeroCarousel :photos="heroPhotos" />
        </div>
      </div>
    </div>
  </section>
</template>
