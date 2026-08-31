<script setup lang="ts">
import type { PhotoboothFrameItem } from '~/composables/api/photoboothFrames'

const DEFAULT_MESSAGE = 'Halo Nora Photobooth, saya ingin bertanya-tanya seputar Digital Photobooth.'
const config = useRuntimeConfig()
const waHref = computed(
  () => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(DEFAULT_MESSAGE)}`,
)

const photoboothFramesData = await useServerFetch<PhotoboothFrameItem[]>('/photobooth-frames', 'digital-photobooth-frames')
const photoboothFrames = computed(() => photoboothFramesData.value ?? [])

const showTryModal = ref(false)

interface Step {
  number: string
  title: string
  variant: 'qr' | 'frame' | 'camera' | 'voice' | 'saved' | 'share'
}

const steps: Step[] = [
  { number: '01', title: 'Scan QR pada kartu yang tersedia', variant: 'qr' },
  { number: '02', title: 'Masukkan nama dan pilih frame favorit', variant: 'frame' },
  { number: '03', title: 'Ambil foto terbaik kalian', variant: 'camera' },
  { number: '04', title: 'Rekam voice message dan upload foto', variant: 'voice' },
  { number: '05', title: 'Semua foto dan ucapan akan tersimpan', variant: 'saved' },
  { number: '06', title: 'Bagikan atau Download kisah', variant: 'share' },
]

useHead({
  title: 'Digital Photobooth — Nora Photobooth',
  meta: [
    {
      name: 'description',
      content: 'Kenali cara kerja Digital Photobooth Nora: scan QR, ambil foto, rekam ucapan, hingga bagikan kisah acaramu.',
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
          <Icon name="heroicons:device-phone-mobile" class="text-base" />
          Fitur Terbaru
        </span>
        <h1 class="mt-6 font-dm-serif text-4xl leading-tight font-bold text-[#000000] sm:text-5xl">
          Digital Photobooth
        </h1>
        <p class="mt-4 font-poppins text-base leading-relaxed text-[#57607A] sm:text-lg">
          Cukup scan QR code di kartu meja, tamu bisa langsung mengambil foto, merekam ucapan suara, dan menyimpan
          kenangan acaramu — semuanya lewat ponsel mereka sendiri.
        </p>
        <div class="mt-8 flex flex-wrap items-center justify-center gap-4">
          <button
            class="flex items-center gap-2 rounded-full bg-[#920f0f] px-8 py-3.5 text-sm font-semibold text-[#FAF9F6] shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5"
            @click="showTryModal = true"
          >
            <Icon name="heroicons:camera" />
            Coba Digital Photobooth
          </button>
          <a
            :href="waHref"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-full border border-[#920f0f] px-8 py-3.5 text-sm font-semibold text-[#920f0f] transition hover:-translate-y-0.5 hover:bg-[#920f0f]/5"
          >
            Tanya-Tanya via WhatsApp
          </a>
        </div>
      </div>
    </section>

    <section class="relative bg-white py-20 lg:py-28">
      <div class="mx-auto max-w-7xl px-6 lg:px-10">
        <div class="mx-auto max-w-2xl text-center">
          <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">Cara Kerja</span>
          <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">6 Langkah Mudah Digital Photobooth</h2>
        </div>

        <div class="relative mt-20 grid grid-cols-2 gap-x-4 gap-y-16 sm:grid-cols-3 lg:flex lg:items-start lg:justify-between lg:gap-6">
          <div class="pointer-events-none absolute top-[22px] right-6 left-6 hidden border-t-2 border-dashed border-[#D8CFC2] lg:block" />

          <div v-for="step in steps" :key="step.number" class="relative z-10 flex flex-col items-center text-center">
            <span
              class="mb-5 flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-[#F3C9D3] font-dm-sans text-sm font-bold text-[#920f0f] shadow-sm"
            >
              {{ step.number }}
            </span>

            <div class="relative h-[300px] w-[148px] overflow-hidden rounded-[26px] border-[6px] border-black bg-black shadow-xl">
              <span class="absolute top-0 left-1/2 z-20 h-4 w-16 -translate-x-1/2 rounded-b-lg bg-black" />

              <div class="relative flex h-full w-full flex-col overflow-hidden bg-gradient-to-b from-[#7a1420] to-[#4c0d15] px-2.5 pt-6 pb-2.5">
                <template v-if="step.variant === 'qr'">
                  <p class="font-dm-sans text-[8px] font-semibold tracking-[0.15em] text-[#f3d9c6] uppercase">Wedding Memories Of</p>
                  <p class="mt-1 font-aloja text-lg text-white">Fatimah &amp; Fatah</p>
                  <div class="mt-3 flex flex-1 items-center justify-center">
                    <div class="flex h-16 w-16 items-center justify-center rounded-lg bg-white">
                      <Icon name="heroicons:qr-code" class="text-3xl text-black" />
                    </div>
                  </div>
                  <span class="mb-1 rounded-full bg-white/90 py-1.5 text-center font-poppins text-[8px] font-semibold text-[#7a1420]">
                    Ketuk Untuk Membuka
                  </span>
                </template>

                <template v-else-if="step.variant === 'frame'">
                  <p class="font-dm-sans text-[8px] font-semibold tracking-[0.15em] text-[#f3d9c6] uppercase">Wedding Memories Of</p>
                  <p class="mt-1 font-aloja text-lg text-white">Fatimah &amp; Fatah</p>
                  <div class="mt-3 flex flex-1 flex-col items-center justify-center gap-1.5">
                    <div class="h-14 w-24 rounded-md border border-white/30 bg-black/60" />
                    <div class="h-14 w-24 rounded-md border border-white/30 bg-black/60" />
                  </div>
                  <span class="mb-1 rounded-full bg-white/90 py-1.5 text-center font-poppins text-[8px] font-semibold text-[#7a1420]">
                    Pilih Frame
                  </span>
                </template>

                <template v-else-if="step.variant === 'camera'">
                  <Icon name="heroicons:arrow-left" class="absolute top-6 left-2.5 z-10 text-xs text-white" />
                  <div class="absolute inset-0 flex items-center justify-center bg-[#2b2b2b]">
                    <Icon name="heroicons:user-circle" class="text-6xl text-white/70" />
                  </div>
                  <span class="relative z-10 mt-auto mb-1 rounded-full bg-white/90 py-1.5 text-center font-poppins text-[8px] font-semibold text-[#7a1420]">
                    Lanjut
                  </span>
                </template>

                <template v-else-if="step.variant === 'voice'">
                  <p class="font-dm-sans text-[8px] font-semibold tracking-[0.15em] text-[#f3d9c6] uppercase">Untuk Mempelai</p>
                  <p class="mt-1 font-poppins text-[7px] leading-tight text-white/80">
                    Sampaikan harapan terbaikmu untuk mereka
                  </p>
                  <div class="flex flex-1 flex-col items-center justify-center gap-2">
                    <div class="flex h-10 w-10 items-center justify-center rounded-full bg-white/15 ring-2 ring-white/40">
                      <Icon name="heroicons:microphone" class="text-lg text-white" />
                    </div>
                    <span class="font-dm-sans text-[9px] text-white/80">00:07</span>
                  </div>
                  <span class="mb-1 rounded-full bg-[#920f0f] py-1.5 text-center font-poppins text-[7px] font-semibold text-white ring-1 ring-white/30">
                    Rekam Ucapan &amp; Harapan
                  </span>
                </template>

                <template v-else-if="step.variant === 'saved'">
                  <p class="font-dm-sans text-[8px] font-semibold tracking-[0.15em] text-[#f3d9c6] uppercase">Wedding Memories Of</p>
                  <p class="mt-1 font-aloja text-base text-white">Fatimah &amp; Fatah</p>
                  <div class="mt-2 flex flex-1 items-center justify-center">
                    <div class="flex flex-col gap-1 rounded-md bg-white p-1.5 shadow">
                      <div class="h-8 w-16 rounded-sm bg-[#c9a27a]" />
                      <div class="h-8 w-16 rounded-sm bg-[#a8b8a0]" />
                    </div>
                  </div>
                  <span class="mb-1 rounded-full bg-white/90 py-1.5 text-center font-poppins text-[7px] font-semibold text-[#7a1420]">
                    Bagikan Momen &amp; Harapan
                  </span>
                </template>

                <template v-else>
                  <p class="font-dm-sans text-[8px] font-semibold tracking-[0.15em] text-[#f3d9c6] uppercase">Semua Kisah</p>
                  <div class="mt-2 grid flex-1 grid-cols-3 content-start gap-1">
                    <div v-for="n in 6" :key="n" class="aspect-square rounded-sm bg-white/20" />
                  </div>
                  <div class="mb-1 flex items-center justify-center gap-3 text-white">
                    <Icon name="heroicons:heart" class="text-sm" />
                    <Icon name="heroicons:share" class="text-sm" />
                    <Icon name="heroicons:arrow-down-tray" class="text-sm" />
                  </div>
                </template>
              </div>
            </div>

            <p class="mt-5 max-w-[160px] font-poppins text-sm leading-snug text-[#39445B]">
              {{ step.title }}
            </p>
          </div>
        </div>
      </div>
    </section>

    <section class="bg-[#FAF9F6] py-16">
      <div class="mx-auto flex max-w-3xl flex-col items-center gap-4 px-6 text-center lg:px-10">
        <h2 class="font-dm-serif text-2xl font-bold text-[#000000] sm:text-3xl">Ingin Menambahkan Digital Photobooth ke Acaramu?</h2>
        <p class="font-poppins text-[#57607A]">Konsultasikan kebutuhan acaramu bersama tim kami sekarang.</p>
        <div class="mt-2 flex flex-wrap items-center justify-center gap-4">
          <button
            class="flex items-center gap-2 rounded-full bg-[#920f0f] px-8 py-3.5 text-sm font-semibold text-[#FAF9F6] shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5"
            @click="showTryModal = true"
          >
            <Icon name="heroicons:camera" />
            Coba Digital Photobooth
          </button>
          <a
            :href="waHref"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-full border border-[#920f0f] px-8 py-3.5 text-sm font-semibold text-[#920f0f] transition hover:-translate-y-0.5 hover:bg-[#920f0f]/5"
          >
            Booking Sekarang
          </a>
        </div>
      </div>
    </section>

    <Footer />
    <WhatsappCta />

    <TryModal v-if="showTryModal" :frames="photoboothFrames" @close="showTryModal = false" />
  </main>
</template>
