<script setup lang="ts">
// ===== Teks (ubah di sini) =====
const ARIA_LABEL = 'Chat via WhatsApp'
const WA_MESSAGE = 'Halo Nora Photobooth, saya ingin bertanya-tanya seputar paket photobooth.'
// ================================

const { trackCtaClick } = useAnalyticsApi()
const config = useRuntimeConfig()
const href = computed(() => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(WA_MESSAGE)}`)
</script>

<template>
  <a
    :href="href"
    target="_blank"
    rel="noopener noreferrer"
    :aria-label="ARIA_LABEL"
    class="fixed right-6 bottom-6 z-50 flex h-14 w-14 items-center justify-center rounded-full bg-[#25D366] text-white shadow-lg transition-transform hover:scale-110 hover:shadow-xl"
    @click="trackCtaClick('promo_whatsapp_float')"
  >
    <Icon name="fa6-brands:whatsapp" class="text-3xl" />
  </a>
</template>
