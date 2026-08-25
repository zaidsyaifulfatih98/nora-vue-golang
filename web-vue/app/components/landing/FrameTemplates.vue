<script setup lang="ts">
import type { FrameTemplateItem } from '~/composables/api/frameTemplates'

const props = defineProps<{ templates: FrameTemplateItem[] }>()

const config = useRuntimeConfig()
const waHref = computed(
  () =>
    `https://wa.me/${config.public.whatsappNumber}?text=${encodeURIComponent(
      'Halo, saya ingin request desain frame photobooth sendiri.',
    )}`,
)
</script>

<template>
  <section v-if="templates.length > 0" id="template-frame" class="relative bg-[#FFFFFF] py-24">
    <div class="mx-auto max-w-7xl px-6 lg:px-10">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="font-dm-serif text-3xl font-bold text-[#000000] sm:text-4xl">Template Menarik dan Bisa di Kustom</h2>
        <p class="mt-4 font-poppins text-[#57607A]">
          Pilih dari tampilan elegan, siap untuk acara. Kami menyesuaikan teks dan warna agar sesuai dengan tema acara.
        </p>
      </div>

      <div class="mt-16">
        <Carousel :items="templates">
          <template #default="{ item: tpl }">
            <div class="group rounded-2xl transition duration-300 hover:-translate-y-1">
              <div class="overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-[#E4E2DC]/70 transition duration-300">
                <div class="aspect-[4/3] w-full overflow-hidden bg-white">
                  <img
                    :src="tpl.imageUrl"
                    :alt="tpl.name"
                    class="h-full w-full transition duration-300 group-hover:scale-105"
                    :class="tpl.fit === 'CONTAIN' ? 'object-contain' : 'object-cover'"
                  />
                </div>
                <div class="p-6">
                  <h3 class="font-poppins text-lg font-bold text-[#1E2537]">{{ tpl.name }}</h3>
                  <p class="mt-1.5 font-poppins text-sm text-[#57607A]">{{ tpl.description }}</p>
                </div>
              </div>
            </div>
          </template>
        </Carousel>
      </div>

      <div class="mt-12 flex justify-center">
        <a
          :href="waHref"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 rounded-full bg-[#920f0f] px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5 hover:bg-[#920f0f]"
        >
          <Icon name="fa6-brands:whatsapp" class="text-lg" />
          Request Desain Sendiri
        </a>
      </div>
    </div>
  </section>
</template>
