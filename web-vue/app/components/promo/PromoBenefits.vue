<script setup lang="ts">
// ===== Konten (ubah di sini) =====
// Untuk mengganti gambar: timpa file di public/promo/benefit-1.png ... benefit-3.png
// (nama sama), atau ubah path `image` di bawah. Gunakan cutout PNG transparan agar hasilnya rapi.
// Tulis **teks** untuk membuat teks tebal di deskripsi.
const HEADING_LINE_1 = 'Yang Kamu Dapatkan'
const HEADING_LINE_2 = 'di Nora Photobooth' // baris ini diberi garis highlight

const BENEFITS = [
  {
    image: '/promo/benefit-1.jpeg',
    title: 'Staf Profesional',
    description: 'Staf kami siap menyambutmu dengan **hangat**, membantu merapikan gayamu, dan memastikan setiap tawa tulus terbingkai sempurna dalam **kenangan**.',
  },
  {
    image: '/promo/benefit-2.jpeg',
    title: 'Teknologi Terbaru & Premium',
    description: 'Nikmati pengalaman photobooth **berstandar tinggi**: kamera **Sony A6400**, **layar iPad** yang responsif, serta cetakan instan tahan lama dari **Printer DNP**.',
  },
  {
    image: '/promo/benefit-3.jpeg',
    title: 'Properti & Aksesoris yang Seru',
    description: 'Ikin fotomu makin hidup dan ekspresif! Beragam **props unik, lucu, dan elegan** yang bisa disesuaikan dengan tema acara.',
  },
]
// =================================

// Memecah "teks **tebal** teks" menjadi segmen agar bisa dirender tanpa v-html.
function segments(text: string) {
  return text.split('**').map((part, i) => ({ text: part, bold: i % 2 === 1 }))
}
</script>

<template>
  <section class="relative bg-[#faf6f0] py-16 md:py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <h2 class="text-center font-poppins text-3xl leading-snug font-bold text-[#1E2537] sm:text-4xl">
        {{ HEADING_LINE_1 }}
        <br />
        <span class="inline-block bg-[linear-gradient(transparent_62%,#f3c6c6_62%)] px-1">{{ HEADING_LINE_2 }}</span>
      </h2>

      <!-- Mobile: carousel; sm ke atas: grid -->
      <PromoCarousel :items="BENEFITS" grid-from="sm" class="mt-12 sm:mt-16" track-class="sm:gap-10">
        <template #default="{ item }">
          <div class="flex h-full flex-col items-center text-center">
            <div class="flex h-48 w-full items-end justify-center sm:h-56">
              <img :src="item.image" :alt="item.title" class="max-h-full max-w-full object-contain" />
            </div>
            <h3 class="mt-8 font-poppins text-2xl font-bold text-[#1E2537]">{{ item.title }}</h3>
            <p class="mt-4 max-w-sm font-poppins text-base leading-relaxed text-[#1E2537]">
              <template v-for="(seg, i) in segments(item.description)" :key="i">
                <strong v-if="seg.bold" class="font-bold">{{ seg.text }}</strong>
                <template v-else>{{ seg.text }}</template>
              </template>
            </p>
          </div>
        </template>
      </PromoCarousel>
    </div>
  </section>
</template>
