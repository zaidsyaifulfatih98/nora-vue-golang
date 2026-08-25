<script setup lang="ts">
interface ReviewData {
  id: string
  name: string
  eventLabel: string
  quote: string
  rating: number
}

const FALLBACK_REVIEWS: ReviewData[] = [
  {
    id: 'sarah-bima',
    name: 'Sarah & Bima',
    eventLabel: 'Resepsi, Mei 2026',
    quote: 'Tamu-tamu kami betah antre di booth-nya karena seru banget! Hasil cetakannya juga cantik, senada banget sama dekorasi kami.',
    rating: 5,
  },
  {
    id: 'dinda-raka',
    name: 'Dinda & Raka',
    eventLabel: 'Akad & Resepsi, Maret 2026',
    quote: 'Tim Nora responsif banget dari awal konsultasi sampai hari-H. Propsnya lucu-lucu, template fotonya juga bisa custom sesuai request.',
    rating: 5,
  },
  {
    id: 'putri-aldi',
    name: 'Putri & Aldi',
    eventLabel: 'Garden Party, Januari 2026',
    quote: 'Worth it banget buat kenang-kenangan tamu. Cetak fotonya cepat, hasilnya bagus, dan galeri online-nya memudahkan kami bagi-bagi foto.',
    rating: 5,
  },
]

const fetched = await useServerFetch<ReviewData[]>('/reviews', 'landing-reviews')
const reviews = computed(() => fetched.value ?? FALLBACK_REVIEWS)
</script>

<template>
  <section id="testimoni" class="relative bg-white py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">Testimoni</span>
        <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">Cerita dari Pasangan yang Kami Temani</h2>
      </div>

      <TestimonialsMarquee :reviews="reviews" />
    </div>
  </section>
</template>
