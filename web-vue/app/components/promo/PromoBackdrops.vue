<script setup lang="ts">
import type { BackdropItem } from '~/composables/api/backdrops'

// ===== Teks (ubah di sini) =====
const TITLE = 'Pilihan Backdrop'
// ================================

const data = await useServerFetch<BackdropItem[]>('/backdrops', 'landing-backdrops')
const backdrops = computed(() => data.value ?? [])

const selectedId = ref<string | null>(null)
const selected = computed(() => backdrops.value.find((b) => b.id === selectedId.value) ?? backdrops.value[0] ?? null)
</script>

<template>
  <section v-if="backdrops.length > 0" id="backdrop" class="relative bg-[#f7f3eb] py-16 md:py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ TITLE }}</h2>
      </div>

      <!-- Tripod stand + backdrop -->
      <div class="relative mx-auto mt-10 aspect-[4/3] w-full max-w-3xl px-0 sm:mt-16 sm:aspect-[16/11]">
        <!-- lantai: bayangan lembut besar + bayangan pekat di bawah tiap kaki -->
        <div class="absolute right-[2%] -bottom-3 left-[2%] h-8 rounded-[50%] bg-black/20 blur-xl" />
        <div class="absolute bottom-0 left-[-2%] h-3 w-[18%] rounded-[50%] bg-black/50 blur-md" />
        <div class="absolute right-[-2%] bottom-0 h-3 w-[18%] rounded-[50%] bg-black/50 blur-md" />

        <!-- tiang kiri & kanan (silinder logam dengan highlight + bayangan jatuh) -->
        <div
          v-for="side in ['left-[2%]', 'right-[2%]']"
          :key="side"
          class="absolute top-0 bottom-3 w-2 rounded-full bg-[linear-gradient(90deg,#050505_0%,#5a5a5a_30%,#151515_55%,#000_100%)] drop-shadow-[5px_8px_6px_rgba(0,0,0,0.35)] sm:w-3"
          :class="side"
        >
          <!-- ring sambungan teleskopik -->
          <span class="absolute top-[40%] -left-[3px] -mx-px h-3 sm:h-4 w-[calc(100%+8px)] rounded-[3px] bg-[linear-gradient(90deg,#000,#444_40%,#0a0a0a)] shadow-md" />
          <span class="absolute top-[41%] left-[3px] h-px w-[calc(100%-6px)] bg-white/20" />
        </div>

        <!-- batang atas -->
        <div class="absolute top-0 right-[1.2%] left-[1.2%] h-2 rounded-full bg-[linear-gradient(180deg,#5a5a5a_0%,#111_45%,#000_100%)] drop-shadow-[0_6px_5px_rgba(0,0,0,0.35)] sm:h-3" />
        <!-- klem sudut atas -->
        <div class="absolute top-[-3px] left-[0.6%] h-3.5 w-3.5 rounded-[4px] bg-[linear-gradient(135deg,#555,#000)] shadow-md sm:h-5 sm:w-5" />
        <div class="absolute top-[-3px] right-[0.6%] h-3.5 w-3.5 rounded-[4px] bg-[linear-gradient(225deg,#555,#000)] shadow-md sm:h-5 sm:w-5" />

        <!-- kaki dudukan -->
        <div class="absolute bottom-0 left-[-2%] h-3 w-[18%] -skew-x-12 rounded-md bg-[linear-gradient(180deg,#4a4a4a_0%,#0d0d0d_40%,#000_100%)] shadow-[0_4px_6px_rgba(0,0,0,0.5)] sm:h-4" />
        <div class="absolute right-[-2%] bottom-0 h-3 w-[18%] skew-x-12 rounded-md bg-[linear-gradient(180deg,#4a4a4a_0%,#0d0d0d_40%,#000_100%)] shadow-[0_4px_6px_rgba(0,0,0,0.5)] sm:h-4" />

        <!-- kain backdrop -->
        <div class="absolute top-[4%] right-[4%] bottom-[6%] left-[4%] sm:right-[3.5%] sm:bottom-[5%] sm:left-[3.5%] overflow-hidden bg-white">
          <img
            v-for="bd in backdrops"
            :key="bd.id"
            :src="bd.imageUrl"
            :alt="bd.name"
            class="absolute inset-0 h-full w-full object-cover transition-opacity duration-500"
            :class="bd.id === selected?.id ? 'opacity-100' : 'opacity-0'"
          />
          <!-- mata ayam di sudut -->
          <span class="absolute top-1.5 left-1.5 h-2 w-2 rounded-full bg-[#b5b5b5] shadow-inner ring-1 ring-black/40" />
          <span class="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-[#b5b5b5] shadow-inner ring-1 ring-black/40" />
          <span class="absolute bottom-1.5 left-1.5 h-2 w-2 rounded-full bg-[#b5b5b5] shadow-inner ring-1 ring-black/40" />
          <span class="absolute right-1.5 bottom-1.5 h-2 w-2 rounded-full bg-[#b5b5b5] shadow-inner ring-1 ring-black/40" />
        </div>
      </div>

      <!-- Pilihan backdrop (swatch) -->
      <div class="mt-8 flex flex-col items-center gap-3 sm:mt-10">
        <div class="flex flex-wrap items-center justify-center gap-3.5" role="radiogroup" :aria-label="TITLE">
          <button
            v-for="bd in backdrops"
            :key="bd.id"
            type="button"
            role="radio"
            :aria-checked="bd.id === selected?.id"
            :aria-label="bd.name"
            :title="bd.name"
            class="h-10 w-10 overflow-hidden rounded-full ring-1 ring-[#E4E2DC] transition hover:scale-110"
            :class="bd.id === selected?.id ? 'ring-2 ring-[#1E2537] ring-offset-2 ring-offset-[#f7f3eb]' : ''"
            @click="selectedId = bd.id"
          >
            <img :src="bd.imageUrl" :alt="bd.name" class="h-full w-full object-cover" />
          </button>
        </div>
        <p class="font-poppins text-sm font-semibold text-[#1E2537]">{{ selected?.name }}</p>
      </div>
    </div>
  </section>
</template>
