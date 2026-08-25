<script setup lang="ts">
const DEFAULT_MESSAGE = 'Halo Nora Photobooth, saya ingin bertanya-tanya seputar paket photobooth.'
const config = useRuntimeConfig()
const href = computed(
  () => `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(DEFAULT_MESSAGE)}`,
)
</script>

<template>
  <a
    :href="href"
    target="_blank"
    rel="noopener noreferrer"
    aria-label="Chat via WhatsApp"
    class="fixed bottom-6 right-6 z-50 flex h-14 w-14 items-center justify-center rounded-full bg-[#25D366] text-white shadow-lg transition-transform hover:scale-110 hover:shadow-xl"
  >
    <Icon name="fa6-brands:whatsapp" class="text-3xl" />
  </a>
</template>
