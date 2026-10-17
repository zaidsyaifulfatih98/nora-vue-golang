<script setup lang="ts">
import type { FeatureItem } from '~/composables/api/features'

// ===== Teks (ubah di sini) =====
const EYEBROW = 'Kenapa Nora Photobooth'
const TITLE = 'Buat Momenmu Lebih Berkesan'
const DESCRIPTION = 'Kami memadukan teknologi, kreativitas, dan sentuhan personal agar setiap tamu punya cerita untuk dibawa pulang.'
const EMPTY = 'Fitur unggulan akan segera hadir.'
// ================================

const features = await useServerFetch<FeatureItem[]>('/features', 'landing-features')
const hasFeatures = computed(() => Boolean(features.value && features.value.length > 0))
</script>

<template>
  <section id="keunggulan" class="relative bg-[#f7f3eb] py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ EYEBROW }}</span>
        <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ TITLE }}</h2>
        <p class="mt-4 font-poppins text-[#57607A]">{{ DESCRIPTION }}</p>
      </div>

      <!-- Mobile: carousel; sm ke atas: grid -->
      <PromoCarousel v-if="hasFeatures" :items="features ?? []" grid-from="sm" class="mt-16" track-class="sm:gap-6">
        <template #default="{ item: feature }">
          <div class="group flex h-full rounded-2xl transition duration-300 hover:-translate-y-1">
            <div class="flex w-full flex-col overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-[#E4E2DC]/70 transition duration-300 group-hover:shadow-xl">
              <div class="aspect-[4/3] w-full overflow-hidden">
                <img :src="feature.imageUrl" :alt="feature.title" class="h-full w-full object-cover transition duration-300 group-hover:scale-105" />
              </div>
              <div class="flex flex-1 flex-col p-7">
                <h3 class="font-poppins text-lg font-bold text-[#1E2537]">{{ feature.title }}</h3>
                <p class="mt-2 font-poppins text-sm leading-relaxed text-[#57607A]">{{ feature.description }}</p>
              </div>
            </div>
          </div>
        </template>
      </PromoCarousel>
      <p v-else class="mt-16 text-center font-poppins text-sm text-[#57607A]">{{ EMPTY }}</p>
    </div>
  </section>
</template>
