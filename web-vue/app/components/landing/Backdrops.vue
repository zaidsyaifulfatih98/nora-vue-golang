<script setup lang="ts">
import type { BackdropItem } from '~/composables/api/backdrops'

defineProps<{ backdrops: BackdropItem[] }>()
const { t } = useI18n()
</script>

<template>
  <section v-if="backdrops.length > 0" id="backdrop" class="relative bg-[#f7f3eb] py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">{{ t('landing.backdrops.title') }}</h2>
      </div>

      <div class="mt-16">
        <Carousel :items="backdrops">
          <template #default="{ item: bd }">
            <div class="group rounded-2xl transition duration-300 hover:-translate-y-1">
              <div class="overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-[#E4E2DC]/70 transition duration-300">
                <div class="aspect-[4/3] w-full overflow-hidden bg-white">
                  <img :src="bd.imageUrl" :alt="bd.name" class="h-full w-full object-cover transition duration-300 group-hover:scale-105" />
                </div>
                <div class="p-4 text-center">
                  <h3 class="font-poppins text-base font-bold text-[#1E2537]">{{ bd.name }}</h3>
                </div>
              </div>
            </div>
          </template>
        </Carousel>
      </div>
    </div>
  </section>
</template>
