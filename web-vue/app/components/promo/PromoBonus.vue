<script setup lang="ts">
// ===== Teks (ubah di sini) =====
const BADGE = 'Bonus Spesial Hari Ini'
const BIG_NUMBER = '+1'
const BIG_UNIT = 'JAM'
const TITLE = 'Bonus Tambahan 1 Jam Cetak Unlimited'
const DESCRIPTION = 'Pesan hari ini dan dapatkan tambahan 1 jam sesi photobooth dengan cetak foto sepuasnya — tamu makin puas, kenangan makin banyak!'
const PERKS = ['Cetak foto unlimited selama 1 jam tambahan', 'Berlaku untuk semua paket', 'Tanpa biaya tambahan']
const COUNTDOWN_LABEL = 'Bonus berakhir dalam'
const CTA = 'Klaim Bonus Sekarang'
const NOTE = '*Berlaku untuk pemesanan yang dikonfirmasi hari ini. Promo diperbarui setiap hari.'
const WA_MESSAGE = 'Halo Nora Photobooth, saya ingin klaim bonus tambahan 1 jam cetak unlimited untuk pemesanan hari ini.'
// ================================

const { trackCtaClick } = useAnalyticsApi()
const config = useRuntimeConfig()
const waHref = computed(() => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(WA_MESSAGE)}`)

const { units } = useDailyCountdown()
</script>

<template>
  <section id="promo-bonus" class="relative bg-white py-16 md:py-24">
    <div class="mx-auto max-w-5xl px-6 lg:px-10">
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-[#faf6f0] via-white to-[#f7e1e1] p-6 shadow-xl ring-1 ring-[#920f0f]/15 sm:p-10 lg:p-14">
        <div class="pointer-events-none absolute -top-20 -right-20 h-64 w-64 rounded-full bg-[#920f0f]/10 blur-3xl" />

        <div class="relative grid items-center gap-10 lg:grid-cols-[auto_1fr] lg:gap-14">
          <!-- lencana besar -->
          <div class="mx-auto flex h-44 w-44 flex-col items-center justify-center rounded-full bg-[#920f0f] text-white shadow-2xl shadow-[#920f0f]/40 ring-8 ring-[#920f0f]/10 sm:h-52 sm:w-52">
            <span class="font-dm-serif text-6xl leading-none font-bold sm:text-7xl">{{ BIG_NUMBER }}</span>
            <span class="mt-1 font-poppins text-2xl font-bold tracking-[0.25em] sm:text-3xl">{{ BIG_UNIT }}</span>
            <span class="mt-1 font-poppins text-[11px] font-semibold tracking-widest uppercase text-white/80">Unlimited</span>
          </div>

          <div class="text-center lg:text-left">
            <span class="inline-flex items-center gap-2 rounded-full bg-[#920f0f]/10 px-4 py-1.5 font-dm-sans text-xs font-semibold tracking-wide text-[#920f0f]">
              <Icon name="heroicons:gift" class="text-base" />
              {{ BADGE }}
            </span>
            <h2 class="mt-4 font-dm-serif text-3xl leading-tight font-bold text-[#000000] sm:text-4xl">{{ TITLE }}</h2>
            <p class="mt-3 font-poppins text-sm leading-relaxed text-[#57607A] sm:text-base">{{ DESCRIPTION }}</p>

            <ul class="mt-5 space-y-2 text-left">
              <li v-for="perk in PERKS" :key="perk" class="flex items-start gap-2.5 font-poppins text-sm text-[#39445B]">
                <Icon name="heroicons:check-circle" class="mt-0.5 shrink-0 text-lg text-[#920f0f]" />
                {{ perk }}
              </li>
            </ul>

            <!-- countdown -->
            <p class="mt-7 font-poppins text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ COUNTDOWN_LABEL }}</p>
            <div class="mt-3 flex items-start justify-center gap-2.5 lg:justify-start" role="timer" aria-live="off">
              <template v-for="(u, i) in units" :key="u.label">
                <div class="flex w-[4.5rem] flex-col items-center rounded-xl bg-[#920f0f] px-2 py-3 text-white shadow-md sm:w-20">
                  <span class="font-poppins text-2xl leading-none font-bold tabular-nums sm:text-3xl">{{ u.value }}</span>
                  <span class="mt-1.5 font-poppins text-[10px] font-semibold tracking-wide uppercase text-white/80 sm:text-xs">{{ u.label }}</span>
                </div>
                <span v-if="i < units.length - 1" class="mt-2.5 font-poppins text-2xl font-bold text-[#920f0f]/60 sm:text-3xl">:</span>
              </template>
            </div>

            <a
              :href="waHref"
              target="_blank"
              rel="noopener noreferrer"
              class="mt-8 inline-flex items-center gap-2 rounded-full bg-[#920f0f] px-8 py-3.5 text-sm font-bold text-white shadow-lg shadow-[#920f0f]/30 transition hover:-translate-y-0.5"
              @click="trackCtaClick('promo_bonus')"
            >
              <Icon name="fa6-brands:whatsapp" class="text-lg" />
              {{ CTA }}
            </a>
            <p class="mt-4 font-poppins text-xs text-[#57607A]">{{ NOTE }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
