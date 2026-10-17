<script setup lang="ts">
// Halaman ads (ads.noraphotobooth.online). Seluruh teks ada di components/promo/Promo*.vue
// dan tidak memakai i18n, jadi cukup ubah konstanta di bagian atas tiap komponen.
useHead({
  title: 'Sewa Photobooth untuk Acara Spesialmu — Nora Photobooth',
})
</script>

<template>
  <main class="flex flex-col overflow-x-hidden">
    <PromoHero class="order-1" />
    <!-- Di mobile galeri tampil tepat di bawah hero; di desktop tetap setelah paket. -->
    <PromoGallery class="order-2 md:order-7" />
    <PromoFeatures class="order-3" />
    <PromoFrameTemplates class="order-4" />
    <PromoBackdrops class="order-5" />
    <PromoPackages class="order-6" />
    <PromoTestimonials class="order-8" />
    <PromoFooter class="order-9" />
    <PromoWhatsappCta />
  </main>
</template>
