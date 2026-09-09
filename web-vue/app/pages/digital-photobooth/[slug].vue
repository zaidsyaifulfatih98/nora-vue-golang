<script setup lang="ts">
import type { PhotoboothFramesBySlug } from '~/composables/api/photoboothFrames'

const { t } = useI18n()
const route = useRoute()
const slug = computed(() => String(route.params.slug))

const data = await useServerFetch<PhotoboothFramesBySlug>(`/photobooth-frames/by-slug/${slug.value}`, `digital-photobooth-slug-${slug.value}`)

const ownerName = computed(() => {
  const owner = data.value?.owner
  return owner ? `${owner.firstName} ${owner.lastName}`.trim() : ''
})
const frames = computed(() => data.value?.frames ?? [])
const notFound = computed(() => !data.value)

const showTryModal = ref(false)

useHead({
  title: computed(() => (notFound.value ? t('digitalPhotobooth.customerPage.notFoundTitle') : t('digitalPhotobooth.customerPage.titleFor', { name: ownerName.value }))),
})
</script>

<template>
  <main class="overflow-x-hidden bg-white">
    <Navbar />

    <section v-if="notFound" class="relative flex min-h-[70vh] flex-col items-center justify-center bg-[#FAF9F6] px-6 pt-32 pb-16 text-center">
      <Icon name="heroicons:face-frown" class="text-5xl text-[#920f0f]" />
      <h1 class="mt-4 font-dm-serif text-3xl font-bold text-[#000000]">{{ t('digitalPhotobooth.customerPage.notFoundTitle') }}</h1>
      <p class="mt-3 max-w-md font-poppins text-sm text-[#57607A]">{{ t('digitalPhotobooth.customerPage.notFoundDescription') }}</p>
    </section>

    <section v-else class="relative bg-[#FAF9F6] pt-32 pb-16 lg:pt-40 lg:pb-20">
      <div class="mx-auto max-w-3xl px-6 text-center lg:px-10">
        <span class="inline-flex items-center gap-2 rounded-full bg-white px-4 py-1.5 font-dm-sans text-xs font-semibold tracking-wide text-[#920f0f] shadow-sm ring-1 ring-[#E4E2DC]">
          <Icon name="heroicons:device-phone-mobile" class="text-base" />
          {{ t('digitalPhotobooth.badge') }}
        </span>
        <h1 class="mt-6 font-dm-serif text-4xl leading-tight font-bold text-[#000000] sm:text-5xl">
          {{ t('digitalPhotobooth.customerPage.titleFor', { name: ownerName }) }}
        </h1>
        <p class="mt-4 font-poppins text-base leading-relaxed text-[#57607A] sm:text-lg">
          {{ t('digitalPhotobooth.customerPage.descriptionFor', { name: ownerName }) }}
        </p>
        <div class="mt-8 flex flex-wrap items-center justify-center gap-4">
          <button
            class="flex items-center gap-2 rounded-full bg-[#920f0f] px-8 py-3.5 text-sm font-semibold text-[#FAF9F6] shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5"
            @click="showTryModal = true"
          >
            <Icon name="heroicons:camera" />
            {{ t('digitalPhotobooth.ctaTry') }}
          </button>
        </div>
      </div>
    </section>

    <Footer />

    <TryModal v-if="showTryModal" :frames="frames" :owner-slug="slug" @close="showTryModal = false" />
  </main>
</template>
