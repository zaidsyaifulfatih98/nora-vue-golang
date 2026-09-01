<script setup lang="ts">
const { t } = useI18n()
const year = new Date().getFullYear()

const EXPLORE_LINKS = computed(() => [
  { label: t('nav.home'), href: '#hero' },
  { label: t('nav.features'), href: '#keunggulan' },
  { label: t('nav.packages'), href: '#paket' },
  { label: t('nav.gallery'), href: '#galeri' },
  { label: t('nav.testimonials'), href: '#testimoni' },
])
</script>

<template>
  <footer class="bg-[#920f0f] pt-16 pb-8 text-[#D7D4CC]">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="grid gap-10 sm:grid-cols-3 lg:grid-cols-3">
        <div>
          <div class="flex items-center gap-2">
            <span class="relative flex h-9 w-9 items-center justify-center overflow-hidden rounded-full shadow-md">
              <img src="/nora_logo.jpg" alt="Nora Photobooth" class="h-full w-full scale-75 object-contain" />
            </span>
            <span class="font-aloja text-2xl tracking-wide text-[#FAF9F6]">
              Nora <span class="text-[#FAF9F6]">Photobooth</span>
            </span>
          </div>
          <p class="mt-4 font-poppins text-sm leading-relaxed">
            {{ t('landing.footer.tagline') }}
          </p>
        </div>

        <div>
          <h4 class="font-poppins text-base font-bold text-[#FAF9F6]">{{ t('landing.footer.explore') }}</h4>
          <ul class="mt-4 space-y-2 text-sm">
            <li v-for="link in EXPLORE_LINKS" :key="link.href"><a :href="link.href" class="hover:text-white">{{ link.label }}</a></li>
          </ul>
        </div>

        <div>
          <h4 class="font-poppins text-base font-bold text-[#FAF9F6]">{{ t('landing.footer.followUs') }}</h4>
          <div class="mt-4 flex gap-3">
            <a href="#" aria-label="WhatsApp" class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 transition hover:bg-[#920f0f] hover:text-white">
              <Icon name="fa6-brands:whatsapp" />
            </a>
            <a href="#" aria-label="Instagram" class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 transition hover:bg-[#920f0f] hover:text-white">
              <Icon name="fa6-brands:instagram" />
            </a>
            <a href="#" aria-label="TikTok" class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 transition hover:bg-[#920f0f] hover:text-white">
              <Icon name="fa6-brands:tiktok" />
            </a>
          </div>
        </div>
      </div>

      <div class="mt-12 border-t border-white/10 pt-6 text-center text-xs">{{ t('landing.footer.rights', { year }) }}</div>
    </div>
  </footer>
</template>
