<script setup lang="ts">
const year = new Date().getFullYear()
</script>

<template>
  <footer class="bg-[#12161F] pt-16 pb-8 text-[#D7D4CC]">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="grid gap-10 sm:grid-cols-3 lg:grid-cols-3">
        <div>
          <img src="/nora_logo.png" alt="Nora Photobooth" class="h-9 w-auto" />
          <p class="mt-4 font-poppins text-sm leading-relaxed">
            Menemani hari bahagiamu dengan sesi foto yang elegan, seru, dan penuh kenangan.
          </p>
        </div>

        <div>
          <h4 class="font-poppins text-base font-bold text-[#FAF9F6]">Jelajahi</h4>
          <ul class="mt-4 space-y-2 text-sm">
            <li><a href="#hero" class="hover:text-white">Beranda</a></li>
            <li><a href="#keunggulan" class="hover:text-white">Keunggulan</a></li>
            <li><a href="#paket" class="hover:text-white">Paket</a></li>
            <li><a href="#galeri" class="hover:text-white">Galeri</a></li>
            <li><a href="#testimoni" class="hover:text-white">Testimoni</a></li>
          </ul>
        </div>

        <div>
          <h4 class="font-poppins text-base font-bold text-[#FAF9F6]">Ikuti Kami</h4>
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

      <div class="mt-12 border-t border-white/10 pt-6 text-center text-xs">© {{ year }} Nora Photobooth. Semua hak cipta dilindungi.</div>
    </div>
  </footer>
</template>
