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
  variant: 'frame' | 'camera' | 'result' | 'voice'
}

// Mirrors the actual 4-step flow inside the "Coba Digital Photobooth" modal
// (TryModal.vue: frame -> camera -> result -> voice), so the marketing
// preview matches what guests really see instead of an older, different flow.
const steps: Step[] = [
  { number: '01', title: 'Pilih frame favoritmu', variant: 'frame' },
  { number: '02', title: 'Ambil foto sesuai frame yang dipilih', variant: 'camera' },
  { number: '03', title: 'Simpan hasil dan dapatkan QR untuk download', variant: 'result' },
  { number: '04', title: 'Kirim pesan suara untuk pasangan (opsional)', variant: 'voice' },
]

useHead({
  title: 'Digital Photobooth — Nora Photobooth',
  meta: [
    {
      name: 'description',
      content: 'Kenali cara kerja Digital Photobooth Nora: pilih frame, ambil foto, simpan hasilnya, hingga kirim pesan suara untuk pasangan.',
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
          Tamu cukup memilih frame, mengambil beberapa foto, menyimpan hasilnya, dan bisa mengirim pesan suara untuk
          pasangan — semuanya langsung dari ponsel mereka sendiri.
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
          <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">4 Langkah Mudah Digital Photobooth</h2>
        </div>

        <div class="relative mt-20 grid grid-cols-2 gap-x-4 gap-y-16 lg:grid-cols-4 lg:items-start lg:gap-4">
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
                <template v-if="step.variant === 'frame'">
                  <p class="text-center font-dm-serif text-[10px] font-bold text-white">Pilih Frame Favoritmu</p>
                  <div class="mt-3 grid flex-1 grid-cols-2 content-start gap-2">
                    <div class="flex aspect-[3/5] flex-col gap-1 rounded-md bg-white p-1 shadow ring-1 ring-black/5">
                      <div class="flex-1 rounded-sm bg-gradient-to-br from-amber-200 via-rose-200 to-orange-300" />
                      <div class="h-2 rounded-sm bg-[#f3c9d3]" />
                    </div>
                    <div class="flex aspect-[3/5] flex-col gap-1 rounded-md bg-white p-1 shadow ring-1 ring-black/5">
                      <div class="flex-1 rounded-sm bg-gradient-to-br from-sky-200 via-teal-100 to-emerald-200" />
                      <div class="h-2 rounded-sm bg-[#c9a27a]" />
                    </div>
                  </div>
                </template>

                <template v-else-if="step.variant === 'camera'">
                  <p class="text-center font-dm-serif text-[10px] font-bold text-white">Ambil Foto Terbaikmu</p>
                  <div class="relative mt-2 flex-1 overflow-hidden rounded-lg bg-gradient-to-br from-amber-200 via-rose-200 to-orange-300">
                    <svg viewBox="0 0 64 64" class="absolute inset-0 h-full w-full">
                      <circle cx="32" cy="24" r="12" fill="#8a5a44" fill-opacity="0.75" />
                      <path d="M6,66 C6,44 18,36 32,36 C46,36 58,44 58,66 Z" fill="#8a5a44" fill-opacity="0.75" />
                    </svg>
                  </div>
                  <span class="mt-2 mb-1 rounded-full bg-[#920f0f] py-1.5 text-center font-poppins text-[7px] font-semibold text-white ring-1 ring-white/30">
                    Ambil Foto
                  </span>
                </template>

                <template v-else-if="step.variant === 'result'">
                  <p class="text-center font-dm-serif text-[10px] font-bold text-white">Hasil Digital Photobooth-mu</p>
                  <div class="relative mx-auto mt-2 aspect-[3/5] w-[70%] flex-1 overflow-hidden rounded-md bg-gradient-to-br from-amber-200 via-rose-200 to-orange-300 shadow">
                    <svg viewBox="0 0 64 64" class="absolute inset-0 h-full w-full">
                      <circle cx="32" cy="22" r="11" fill="#8a5a44" fill-opacity="0.75" />
                      <path d="M8,62 C8,42 19,34 32,34 C45,34 56,42 56,62 Z" fill="#8a5a44" fill-opacity="0.75" />
                    </svg>
                    <div class="absolute right-1 bottom-1 flex h-8 w-8 items-center justify-center rounded bg-white shadow">
                      <Icon name="heroicons:qr-code" class="text-lg text-black" />
                    </div>
                    <div class="absolute top-1 right-1 flex h-4 w-4 items-center justify-center rounded-full bg-white/90 shadow">
                      <Icon name="heroicons:arrow-down-tray" class="text-[9px] text-[#7a1420]" />
                    </div>
                  </div>
                  <span class="mt-2 mb-1 rounded-full bg-[#920f0f] py-1.5 text-center font-poppins text-[7px] font-semibold text-white ring-1 ring-white/30">
                    Simpan
                  </span>
                </template>

                <template v-else-if="step.variant === 'voice'">
                  <p class="text-center font-dm-serif text-[10px] font-bold text-white">Kirim Pesan Suara</p>
                  <p class="mt-1 text-center font-poppins text-[7px] font-semibold tracking-wide text-[#f3d9c6] uppercase">(Opsional)</p>
                  <div class="flex flex-1 flex-col items-center justify-center gap-2">
                    <div class="flex h-10 w-10 items-center justify-center rounded-full bg-white/15 ring-2 ring-white/40">
                      <Icon name="heroicons:microphone" class="text-lg text-white" />
                    </div>
                  </div>
                  <span class="mb-1 rounded-full bg-[#920f0f] py-1.5 text-center font-poppins text-[7px] font-semibold text-white ring-1 ring-white/30">
                    Kirim Pesan Suara
                  </span>
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

    <Footer />
    <WhatsappCta />

    <TryModal v-if="showTryModal" :frames="photoboothFrames" @close="showTryModal = false" />
  </main>
</template>
