<script setup lang="ts">
// ===== Teks (ubah di sini) =====
const BADGE = 'Promo Hari Ini'
const TITLE = 'Diskon Sampai 47%'
const SUBTITLE = 'HARI INI !!'
const DESCRIPTION = 'Amankan tanggal acaramu sekarang dan nikmati potongan harga untuk semua paket photobooth'
const COUNTDOWN_LABEL = 'Promo berakhir dalam'
const CTA = 'Klaim Diskon Sekarang'
const NOTE = '*Berlaku untuk pemesanan yang dikonfirmasi hari ini. Promo diperbarui setiap hari.'
const WA_MESSAGE = 'Halo Nora Photobooth, saya ingin klaim promo diskon Rp500.000 untuk pemesanan hari ini.'
// ================================

const { trackCtaClick } = useAnalyticsApi()
const config = useRuntimeConfig()
const waHref = computed(() => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(WA_MESSAGE)}`)

const { units } = useDailyCountdown()
</script>

<template>
  <section id="promo-diskon" class="relative bg-[#f7f3eb] py-16 md:py-24">
    <div class="mx-auto max-w-4xl px-6 lg:px-10">
      <div class="promo-card relative overflow-hidden rounded-3xl bg-[#920f0f] px-6 py-12 text-center text-white shadow-2xl shadow-[#920f0f]/30 sm:px-12 sm:py-16">
        <!-- dekorasi -->
        <div class="pointer-events-none absolute -top-16 -left-16 h-48 w-48 rounded-full bg-white/10 blur-2xl" />
        <div class="pointer-events-none absolute -right-16 -bottom-16 h-56 w-56 rounded-full bg-[#f3c6c6]/20 blur-2xl" />

        <div class="relative">
          <span class="inline-flex items-center gap-2 rounded-full bg-white/15 px-4 py-1.5 font-dm-sans text-xs font-semibold tracking-wide ring-1 ring-white/30">
            <Icon name="heroicons:fire" class="text-base" />
            {{ BADGE }}
          </span>

          <h2 class="mt-6 font-dm-serif text-5xl leading-tight font-bold sm:text-6xl">{{ TITLE }}</h2>
          <p class="mt-2 font-dm-serif text-5xl leading-tight font-bold text-white">{{ SUBTITLE }}</p>
          <p class="mx-auto mt-5 max-w-xl font-poppins text-sm leading-relaxed text-white/85 sm:text-base">{{ DESCRIPTION }}</p>

          <!-- countdown -->
          <p class="mt-8 font-poppins text-xs font-semibold tracking-[0.2em] uppercase text-white/80">{{ COUNTDOWN_LABEL }}</p>
          <div class="mt-3 flex items-start justify-center gap-3 sm:gap-4" role="timer" aria-live="off">
            <template v-for="(u, i) in units" :key="u.label">
              <div class="promo-unit flex w-20 flex-col items-center rounded-2xl bg-white px-2 py-3 text-[#920f0f] shadow-lg sm:w-24 sm:py-4">
                <span class="font-poppins text-3xl leading-none font-bold tabular-nums sm:text-4xl">{{ u.value }}</span>
                <span class="mt-1.5 font-poppins text-[11px] font-semibold tracking-wide uppercase text-[#57607A] sm:text-xs">{{ u.label }}</span>
              </div>
              <span v-if="i < units.length - 1" class="mt-3 font-poppins text-3xl font-bold text-white/70 sm:mt-4 sm:text-4xl">:</span>
            </template>
          </div>

          <a
            :href="waHref"
            target="_blank"
            rel="noopener noreferrer"
            class="mt-10 inline-flex items-center gap-2 rounded-full bg-white px-8 py-3.5 text-sm font-bold text-[#920f0f] shadow-lg transition hover:-translate-y-0.5 hover:bg-[#faf6f0]"
            @click="trackCtaClick('promo_discount')"
          >
            <Icon name="fa6-brands:whatsapp" class="text-lg" />
            {{ CTA }}
          </a>
          <p class="mt-4 font-poppins text-xs text-white/70">{{ NOTE }}</p>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.promo-card {
  animation: promo-zoom 4s ease-in-out infinite;
}
.promo-unit {
  animation: promo-bob 2.4s ease-in-out infinite;
}
.promo-unit:nth-of-type(2) {
  animation-delay: 0.3s;
}
.promo-unit:nth-of-type(3) {
  animation-delay: 0.6s;
}

@keyframes promo-zoom {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.03); }
}
@keyframes promo-bob {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

@media (prefers-reduced-motion: reduce) {
  .promo-card,
  .promo-unit {
    animation: none;
  }
}
</style>
