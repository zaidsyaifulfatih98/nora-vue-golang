<script setup lang="ts">
interface ReviewData {
  id: string
  name: string
  eventLabel: string
  quote: string
  eventLabelEn?: string | null
  quoteEn?: string | null
  rating: number
}

const { t } = useI18n()

const FALLBACK_REVIEWS = computed<ReviewData[]>(() => [
  {
    id: 'sarah-bima',
    name: 'Sarah & Bima',
    eventLabel: t('landing.testimonials.sarahBima.eventLabel'),
    quote: t('landing.testimonials.sarahBima.quote'),
    rating: 5,
  },
  {
    id: 'dinda-raka',
    name: 'Dinda & Raka',
    eventLabel: t('landing.testimonials.dindaRaka.eventLabel'),
    quote: t('landing.testimonials.dindaRaka.quote'),
    rating: 5,
  },
  {
    id: 'putri-aldi',
    name: 'Putri & Aldi',
    eventLabel: t('landing.testimonials.putriAldi.eventLabel'),
    quote: t('landing.testimonials.putriAldi.quote'),
    rating: 5,
  },
])

const fetched = await useServerFetch<ReviewData[]>('/reviews', 'landing-reviews')
const reviews = computed(() => fetched.value ?? FALLBACK_REVIEWS.value)
</script>

<template>
  <section id="testimoni" class="relative bg-white py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ t('landing.testimonials.eyebrow') }}</span>
        <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ t('landing.testimonials.title') }}</h2>
      </div>

      <TestimonialsMarquee :reviews="reviews" />
    </div>
  </section>
</template>
