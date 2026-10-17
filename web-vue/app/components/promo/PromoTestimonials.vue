<script setup lang="ts">
interface ReviewData {
  id: string
  name: string
  eventLabel: string
  quote: string
  rating: number
}

// ===== Teks (ubah di sini) =====
const EYEBROW = 'Testimoni'
const TITLE = 'Cerita dari event yang kami temani'

// Dipakai kalau data review dari backend kosong/gagal dimuat.
const FALLBACK_REVIEWS: ReviewData[] = [
  { id: 'sarah-bima', name: 'Sarah & Bima', eventLabel: 'Resepsi, Mei 2026', quote: 'Tamu-tamu kami betah antre di booth-nya karena seru banget! Hasil cetakannya juga cantik, senada banget sama dekorasi kami.', rating: 5 },
  { id: 'dinda-raka', name: 'Dinda & Raka', eventLabel: 'Akad & Resepsi, Maret 2026', quote: 'Tim Nora responsif banget dari awal konsultasi sampai hari-H. Propsnya lucu-lucu, template fotonya juga bisa custom sesuai request.', rating: 5 },
  { id: 'putri-aldi', name: 'Putri & Aldi', eventLabel: 'Garden Party, Januari 2026', quote: 'Worth it banget buat kenang-kenangan tamu. Cetak fotonya cepat, hasilnya bagus, dan galeri online-nya memudahkan kami bagi-bagi foto.', rating: 5 },
]
// ================================

const GRADIENTS = ['from-[#F3D9C0] to-[#E7B98F]', 'from-[#F3D9E0] to-[#E3AFC2]', 'from-[#920f0f] to-[#8F6F3E]']

const fetched = await useServerFetch<ReviewData[]>('/reviews', 'landing-reviews')
const reviews = computed(() => fetched.value ?? FALLBACK_REVIEWS)
// Digandakan agar marquee berjalan mulus tanpa putus.
const track = computed(() => [...reviews.value, ...reviews.value])

function getInitials(name: string) {
  return name
    .split(/\s*&\s*|\s+/)
    .map((part) => part[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()
}
</script>

<template>
  <section id="testimoni" class="relative bg-white py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <span class="font-dm-sans text-xs font-semibold tracking-[0.2em] text-[#920f0f] uppercase">{{ EYEBROW }}</span>
        <h2 class="mt-3 font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ TITLE }}</h2>
      </div>

      <div class="group mt-16 overflow-hidden [mask-image:linear-gradient(to_right,transparent,black_5%,black_95%,transparent)]">
        <div class="flex w-max animate-[marquee_40s_linear_infinite] gap-6 group-hover:[animation-play-state:paused]">
          <div
            v-for="(review, i) in track"
            :key="`${review.id}-${i}`"
            class="flex w-[320px] shrink-0 flex-col rounded-2xl bg-white p-8 shadow-sm ring-1 ring-[#E4E2DC]/70 sm:w-[380px]"
          >
            <Icon name="fa6-solid:quote-left" class="text-2xl text-[#920f0f]" />
            <p class="mt-4 flex-1 font-poppins text-sm leading-relaxed text-[#39445B] italic">&ldquo;{{ review.quote }}&rdquo;</p>

            <div class="mt-6 flex items-center gap-1 text-[#920f0f]">
              <Icon v-for="r in review.rating" :key="r" name="fa6-solid:star" class="text-sm" />
            </div>

            <div class="mt-4 flex items-center gap-3 border-t border-[#F0EFEA] pt-4">
              <div class="flex h-11 w-11 items-center justify-center rounded-full bg-gradient-to-br text-sm font-bold text-white" :class="GRADIENTS[i % GRADIENTS.length]">
                {{ getInitials(review.name) }}
              </div>
              <div>
                <p class="font-poppins text-sm font-bold text-[#1E2537]">{{ review.name }}</p>
                <p class="font-poppins text-xs text-[#6C7686]">{{ review.eventLabel }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
