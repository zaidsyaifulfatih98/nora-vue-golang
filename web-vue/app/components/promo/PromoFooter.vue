<script setup lang="ts">
// ===== Teks & link (ubah di sini) =====
const TAGLINE = 'Menemani hari bahagiamu dengan sesi foto yang elegan, seru, dan penuh kenangan.'
const EXPLORE_TITLE = 'Jelajahi'
const FOLLOW_TITLE = 'Ikuti Kami'
const EXPLORE_LINKS = [
  { label: 'Beranda', href: '#hero' },
  { label: 'Keunggulan', href: '#keunggulan' },
  { label: 'Paket', href: '#paket' },
  { label: 'Galeri', href: '#galeri' },
  { label: 'Testimoni', href: '#testimoni' },
]
const SOCIALS = [
  { label: 'WhatsApp', icon: 'fa6-brands:whatsapp', href: '#' },
  { label: 'Instagram', icon: 'fa6-brands:instagram', href: '#' },
  { label: 'TikTok', icon: 'fa6-brands:tiktok', href: '#' },
]
// ======================================

const year = new Date().getFullYear()
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
            <span class="text-2xl tracking-wide text-[#FAF9F6]">
              <span class="font-marykate text-3xl">NORA</span> <span class="font-quicksand text-xl font-bold text-[#FAF9F6]">PHOTOBOOTH</span>
            </span>
          </div>
          <p class="mt-4 font-poppins text-sm leading-relaxed">{{ TAGLINE }}</p>
        </div>

        <div>
          <h4 class="font-poppins text-base font-bold text-[#FAF9F6]">{{ EXPLORE_TITLE }}</h4>
          <ul class="mt-4 space-y-2 text-sm">
            <li v-for="link in EXPLORE_LINKS" :key="link.href"><a :href="link.href" class="hover:text-white">{{ link.label }}</a></li>
          </ul>
        </div>

        <div>
          <h4 class="font-poppins text-base font-bold text-[#FAF9F6]">{{ FOLLOW_TITLE }}</h4>
          <div class="mt-4 flex gap-3">
            <a
              v-for="s in SOCIALS"
              :key="s.label"
              :href="s.href"
              :aria-label="s.label"
              class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 transition hover:bg-[#920f0f] hover:text-white"
            >
              <Icon :name="s.icon" />
            </a>
          </div>
        </div>
      </div>

      <div class="mt-12 border-t border-white/10 pt-6 text-center text-xs">© {{ year }} Nora Photobooth. Semua hak cipta dilindungi.</div>
    </div>
  </footer>
</template>
