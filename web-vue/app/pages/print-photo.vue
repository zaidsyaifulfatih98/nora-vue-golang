<script setup lang="ts">
const { t } = useI18n()
const config = useRuntimeConfig()

const waHref = computed(
  () => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(t('printPhoto.waMessage'))}`,
)

interface Feature {
  icon: string
  title: string
  description: string
}

const FEATURES = computed<Feature[]>(() => [
  { icon: 'heroicons:bolt', title: t('printPhoto.features.quickDry.title'), description: t('printPhoto.features.quickDry.description') },
  { icon: 'heroicons:shield-check', title: t('printPhoto.features.scratchResistant.title'), description: t('printPhoto.features.scratchResistant.description') },
  { icon: 'heroicons:beaker', title: t('printPhoto.features.waterResistant.title'), description: t('printPhoto.features.waterResistant.description') },
  { icon: 'heroicons:squares-plus', title: t('printPhoto.features.unlimited.title'), description: t('printPhoto.features.unlimited.description') },
  { icon: 'heroicons:truck', title: t('printPhoto.features.fastShipping.title'), description: t('printPhoto.features.fastShipping.description') },
])

interface SizeOption {
  label: string
  dimensions: string
  description: string
}

const SIZES = computed<SizeOption[]>(() => [
  { label: t('printPhoto.sizes.r4.label'), dimensions: '10.2 × 15.2 cm', description: t('printPhoto.sizes.r4.description') },
  { label: t('printPhoto.sizes.r2.label'), dimensions: '6.4 × 9 cm', description: t('printPhoto.sizes.r2.description') },
])

useHead({
  title: computed(() => t('printPhoto.metaTitle')),
  meta: [
    {
      name: 'description',
      content: computed(() => t('printPhoto.metaDescription')),
    },
  ],
})
</script>

<template>
  <main class="overflow-x-hidden bg-white">
    <Navbar />

    <section class="relative bg-[#FAF9F6] pt-32 pb-16 lg:pt-40 lg:pb-20">
      <div class="mx-auto max-w-3xl px-6 text-center lg:px-10">
        <span class="inline-flex items-center gap-2 rounded-full bg-white px-4 py-1.5 font-dm-sans text-xs font-semibold tracking-wide text-[#920f0f] shadow-sm ring-1 ring-[#E4E2DC]">
          <Icon name="heroicons:printer" class="text-base" />
          {{ t('printPhoto.badge') }}
        </span>
        <h1 class="mt-6 font-dm-serif text-4xl leading-tight font-bold text-[#000000] sm:text-5xl">
          {{ t('printPhoto.title') }}
        </h1>
        <p class="mt-4 font-poppins text-base leading-relaxed text-[#57607A] sm:text-lg">
          {{ t('printPhoto.description') }}
        </p>
        <div class="mt-8 flex flex-wrap items-center justify-center gap-4">
          <a
            :href="waHref"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-2 rounded-full bg-[#920f0f] px-8 py-3.5 text-sm font-semibold text-[#FAF9F6] shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5"
          >
            <Icon name="fa6-brands:whatsapp" />
            {{ t('printPhoto.ctaOrder') }}
          </a>
        </div>
      </div>
    </section>

    <section class="relative bg-white py-20 lg:py-28">
      <div class="mx-auto max-w-7xl px-6 lg:px-10">
        <div class="mx-auto max-w-2xl text-center">
          <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ t('printPhoto.printerEyebrow') }}</span>
          <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ t('printPhoto.printerTitle') }}</h2>
          <p class="mt-4 font-poppins text-[#57607A]">{{ t('printPhoto.printerDescription') }}</p>
        </div>

        <div class="mt-16 grid gap-6 sm:grid-cols-2 lg:grid-cols-5">
          <div v-for="feature in FEATURES" :key="feature.title" class="rounded-2xl bg-[#f7f3eb] p-6 shadow-sm ring-1 ring-[#E4E2DC]/70">
            <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-white shadow-sm">
              <Icon :name="feature.icon" class="text-xl text-[#920f0f]" />
            </div>
            <h3 class="mt-4 font-poppins text-base font-bold text-[#1E2537]">{{ feature.title }}</h3>
            <p class="mt-2 font-poppins text-sm leading-relaxed text-[#57607A]">{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="relative bg-[#f7f3eb] py-20 lg:py-28">
      <div class="mx-auto max-w-7xl px-6 lg:px-10">
        <div class="mx-auto max-w-2xl text-center">
          <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ t('printPhoto.sizesEyebrow') }}</span>
          <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ t('printPhoto.sizesTitle') }}</h2>
        </div>

        <div class="mt-16 grid gap-8 sm:grid-cols-2">
          <div
            v-for="size in SIZES"
            :key="size.label"
            class="flex flex-col items-center rounded-3xl bg-white p-8 text-center shadow-sm ring-1 ring-[#E4E2DC]"
          >
            <div class="flex h-16 w-11 items-center justify-center rounded-md bg-gradient-to-br from-amber-200 via-rose-200 to-orange-300 shadow-inner" />
            <h3 class="mt-6 font-poppins text-2xl font-bold text-[#1E2537]">{{ size.label }}</h3>
            <p class="mt-1 font-dm-sans text-sm font-semibold text-[#920f0f]">{{ size.dimensions }}</p>
            <p class="mt-3 font-poppins text-sm text-[#57607A]">{{ size.description }}</p>
          </div>
        </div>

        <div class="mt-12 flex justify-center">
          <a
            :href="waHref"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-2 rounded-full bg-[#920f0f] px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5"
          >
            <Icon name="fa6-brands:whatsapp" class="text-lg" />
            {{ t('printPhoto.ctaOrder') }}
          </a>
        </div>
      </div>
    </section>

    <Footer />
    <WhatsappCta />
  </main>
</template>
