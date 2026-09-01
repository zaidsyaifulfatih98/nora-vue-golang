<script setup lang="ts">
import type { PackageItem } from '~/composables/api/packages'

const { t, tm } = useI18n()
const { tf, tfList } = useLocalizedField()

const FALLBACK_PACKAGES = computed<PackageItem[]>(() => [
  {
    id: 'silver', name: 'Silver', price: '2500000', duration: t('landing.packages.silver.duration'), description: '', isPopular: false, isActive: true, order: 1,
    features: tm('landing.packages.silver.features') as unknown as string[],
    nameEn: null, durationEn: null, descriptionEn: null, featuresEn: null,
  },
  {
    id: 'gold', name: 'Gold', price: '4200000', duration: t('landing.packages.gold.duration'), description: '', isPopular: true, isActive: true, order: 2,
    features: tm('landing.packages.gold.features') as unknown as string[],
    nameEn: null, durationEn: null, descriptionEn: null, featuresEn: null,
  },
  {
    id: 'platinum', name: 'Platinum', price: '6500000', duration: t('landing.packages.platinum.duration'), description: '', isPopular: false, isActive: true, order: 3,
    features: tm('landing.packages.platinum.features') as unknown as string[],
    nameEn: null, durationEn: null, descriptionEn: null, featuresEn: null,
  },
])

function formatRupiah(value: string) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(value))
}

const config = useRuntimeConfig()
const waHref = computed(
  () => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(t('landing.packages.waMessage'))}`,
)

const fetched = await useServerFetch<PackageItem[]>('/packages', 'landing-packages')
const packages = computed(() => fetched.value ?? FALLBACK_PACKAGES.value)
</script>

<template>
  <section id="paket" class="relative bg-[#FFFFFF] py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ t('landing.packages.eyebrow') }}</span>
        <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ t('landing.packages.title') }}</h2>
        <p class="mt-4 font-poppins text-[#57607A]">
          {{ t('landing.packages.description') }}
        </p>
      </div>

      <div class="mt-16 grid gap-8 lg:grid-cols-3">
        <div
          v-for="pkg in packages"
          :key="pkg.id"
          class="relative flex flex-col rounded-3xl p-8 transition hover:-translate-y-1"
          :class="pkg.isPopular ? 'bg-[#f7f3eb] text-black shadow-2xl shadow-[#1E2537]/30 lg:scale-105' : 'bg-white text-[#1E2537] shadow-sm ring-1 ring-[#E4E2DC]'"
        >
          <span v-if="pkg.isPopular" class="absolute -top-4 left-1/2 -translate-x-1/2 rounded-full bg-[#920f0f] px-4 py-1 text-xs font-bold tracking-wide text-white shadow-lg">
            {{ t('landing.packages.mostPopular') }}
          </span>

          <h3 class="font-poppins text-2xl font-bold">{{ tf(pkg.name, pkg.nameEn) }}</h3>
          <p class="mt-1 font-poppins text-sm" :class="pkg.isPopular ? 'text-black' : 'text-[#6C7686]'">{{ tf(pkg.duration, pkg.durationEn) }}</p>

          <p class="mt-6 font-poppins text-3xl font-bold">{{ formatRupiah(pkg.price) }}</p>

          <p v-if="pkg.description" class="mt-2 font-poppins text-sm" :class="pkg.isPopular ? 'text-black' : 'text-[#57607A]'">
            {{ tf(pkg.description, pkg.descriptionEn) }}
          </p>

          <ul class="mt-8 flex-1 space-y-3">
            <li v-for="feature in tfList(pkg.features, pkg.featuresEn)" :key="feature" class="flex items-start gap-2.5 text-sm">
              <Icon name="heroicons:check-circle" class="mt-0.5 shrink-0 text-lg text-[#920f0f]" />
              <span :class="pkg.isPopular ? 'text-black' : 'text-[#39445B]'">{{ feature }}</span>
            </li>
          </ul>
        </div>
      </div>

      <div class="mt-12 flex justify-center">
        <a
          :href="waHref"
          target="_blank"
          rel="noopener"
          class="inline-flex items-center gap-2 rounded-full bg-[#920f0f] px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5 hover:bg-[#920f0f]"
        >
          {{ t('landing.packages.cta') }}
        </a>
      </div>
    </div>
  </section>
</template>
