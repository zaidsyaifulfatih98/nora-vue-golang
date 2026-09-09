<script setup lang="ts">
import type { PhotoboothFramesBySlug } from '~/composables/api/photoboothFrames'

const { t } = useI18n()
const route = useRoute()
const slug = computed(() => String(route.params.slug))

const data = await useServerFetch<PhotoboothFramesBySlug>(`/photobooth-frames/by-slug/${slug.value}`, `software-photobooth-slug-${slug.value}`)

const frames = computed(() => data.value?.frames ?? [])
const notFound = computed(() => !data.value)

useHead({
  title: computed(() => (notFound.value ? t('softwarePhotobooth.customerPage.notFoundTitle') : t('softwarePhotobooth.metaTitle'))),
})
</script>

<template>
  <div v-if="notFound" class="fixed inset-0 z-50 flex flex-col items-center justify-center gap-4 bg-[#FAF9F6] px-6 text-center">
    <Icon name="heroicons:face-frown" class="text-5xl text-[#920f0f]" />
    <h1 class="font-dm-serif text-3xl font-bold text-[#000000]">{{ t('softwarePhotobooth.customerPage.notFoundTitle') }}</h1>
    <p class="max-w-md font-poppins text-sm text-[#57607A]">{{ t('softwarePhotobooth.customerPage.notFoundDescription') }}</p>
  </div>

  <Kiosk v-else :frames="frames" :owner-slug="slug" />
</template>
