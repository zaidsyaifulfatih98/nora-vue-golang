<script setup lang="ts">
import type { PackageItem } from '~/composables/api/packages'

// ===== Teks (ubah di sini) =====
const EYEBROW = 'Investasi Kenangan'
const TITLE = 'Pilih Paket Sesuai Kebutuhan Acaramu'
const DESCRIPTION = 'Semua paket bisa disesuaikan lagi — konsultasikan kebutuhanmu dan kami bantu racik paket yang pas.'
const MOST_POPULAR = 'PALING POPULER'
const START_FROM = 'Start From'
const CTA = 'Pilih Paket Sekarang'
const WA_MESSAGE = 'Halo Nora Photobooth, saya ingin bertanya-tanya seputar paket photobooth.'

// Dipakai kalau data paket dari backend kosong/gagal dimuat.
const FALLBACK_PACKAGES = [
  { id: 'silver', name: 'Silver', price: '2500000', duration: '3 jam sesi', isPopular: false, features: ['1 backdrop tema pilihan', 'Cetak foto unlimited (4R)', '1 orang crew photobooth', 'Free props standar', 'Soft file semua foto'] },
  { id: 'gold', name: 'Gold', price: '4200000', duration: '5 jam sesi', isPopular: true, features: ['Backdrop custom sesuai tema', 'Cetak foto unlimited (4R & strip)', '2 orang crew photobooth', 'Props premium + digital filter', 'Soft file + album digital', 'Free galeri online 30 hari'] },
  { id: 'platinum', name: 'Platinum', price: '6500000', duration: '8 jam / full event', isPopular: false, features: ['Backdrop & frame full custom', 'Cetak foto unlimited semua ukuran', '3 orang crew + 1 fotografer lepas', 'Props premium + GIF & boomerang', 'Album digital + cetak mini album', 'Free galeri online 90 hari'] },
] as unknown as PackageItem[]
// ================================

const { trackCtaClick } = useAnalyticsApi()
const config = useRuntimeConfig()
const waHref = computed(() => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(WA_MESSAGE)}`)

function formatRupiah(value: string) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(value))
}

const fetched = await useServerFetch<PackageItem[]>('/packages', 'landing-packages')
const packages = computed(() => fetched.value ?? FALLBACK_PACKAGES)
</script>

<template>
  <section id="paket" class="relative bg-[#FFFFFF] py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ EYEBROW }}</span>
        <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ TITLE }}</h2>
        <p class="mt-4 font-poppins text-[#57607A]">{{ DESCRIPTION }}</p>
      </div>

      <!-- Mobile/tablet: carousel; lg ke atas: grid 3 kolom -->
      <PromoCarousel :items="packages" grid-from="lg" class="mt-16" track-class="pt-6 pb-4 lg:gap-8 lg:pt-0 lg:pb-0">
        <template #default="{ item: pkg }">
          <div
            class="relative flex h-full flex-col rounded-3xl p-8 transition hover:-translate-y-1"
            :class="pkg.isPopular ? 'bg-[#f7f3eb] text-black shadow-2xl shadow-[#1E2537]/30 lg:scale-105' : 'bg-white text-[#1E2537] shadow-sm ring-1 ring-[#E4E2DC]'"
          >
            <span v-if="pkg.isPopular" class="absolute -top-4 left-1/2 -translate-x-1/2 rounded-full bg-[#920f0f] px-4 py-1 text-xs font-bold tracking-wide text-white shadow-lg">
              {{ MOST_POPULAR }}
            </span>

            <h3 class="flex justify-center font-poppins text-2xl font-bold">{{ pkg.name }}</h3>
            <span class="mt-6 inline-flex w-fit items-center rounded-md bg-[#920f0f] px-3 py-1 text-xs font-bold tracking-wide text-white shadow-sm">
              {{ START_FROM }}
            </span>
            <p class="mt-2 font-poppins text-3xl font-bold">{{ formatRupiah(pkg.price) }}</p>

            <p v-if="pkg.description" class="mt-2 font-poppins text-sm" :class="pkg.isPopular ? 'text-black' : 'text-[#57607A]'">
              {{ pkg.description }}
            </p>

            <ul class="mt-8 flex-1 space-y-3">
              <li v-for="feature in pkg.features" :key="feature" class="flex items-start gap-2.5 text-sm">
                <Icon name="heroicons:check-circle" class="mt-0.5 shrink-0 text-lg text-[#920f0f]" />
                <span :class="pkg.isPopular ? 'text-black' : 'text-[#39445B]'">{{ feature }}</span>
              </li>
            </ul>
          </div>
        </template>
      </PromoCarousel>

      <div class="mt-12 flex justify-center">
        <a
          :href="waHref"
          target="_blank"
          rel="noopener"
          class="inline-flex items-center gap-2 rounded-full bg-[#920f0f] px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5 hover:bg-[#920f0f]"
          @click="trackCtaClick('promo_packages_cta')"
        >
          {{ CTA }}
        </a>
      </div>
    </div>
  </section>
</template>
