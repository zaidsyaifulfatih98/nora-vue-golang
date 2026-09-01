<script setup lang="ts">
import type { FeatureItem } from '~/composables/api/features'

const { t } = useI18n()
const { tf } = useLocalizedField()
const features = await useServerFetch<FeatureItem[]>('/features', 'landing-features')
const hasFeatures = computed(() => Boolean(features.value && features.value.length > 0))
</script>

<template>
  <section id="keunggulan" class="relative bg-[#f7f3eb] py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ t('landing.features.eyebrow') }}</span>
        <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ t('landing.features.title') }}</h2>
        <p class="mt-4 font-poppins text-[#57607A]">
          {{ t('landing.features.description') }}
        </p>
      </div>

      <div v-if="hasFeatures" class="mt-16 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="feature in features" :key="feature.id" class="group rounded-2xl transition duration-300 hover:-translate-y-1">
          <div class="overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-[#E4E2DC]/70 transition duration-300 group-hover:shadow-xl">
            <div class="aspect-[4/3] w-full overflow-hidden">
              <img :src="feature.imageUrl" :alt="tf(feature.title, feature.titleEn)" class="h-full w-full object-cover transition duration-300 group-hover:scale-105" />
            </div>
            <div class="p-7">
              <h3 class="font-poppins text-lg font-bold text-[#1E2537]">{{ tf(feature.title, feature.titleEn) }}</h3>
              <p class="mt-2 font-poppins text-sm leading-relaxed text-[#57607A]">{{ tf(feature.description, feature.descriptionEn) }}</p>
            </div>
          </div>
        </div>
      </div>
      <p v-else class="mt-16 text-center font-poppins text-sm text-[#57607A]">{{ t('landing.features.empty') }}</p>
    </div>
  </section>
</template>
