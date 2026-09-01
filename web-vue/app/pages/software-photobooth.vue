<script setup lang="ts">
const { t } = useI18n()

useHead({
  title: computed(() => t('softwarePhotobooth.metaTitle')),
  meta: [
    {
      name: 'description',
      content: computed(() => t('softwarePhotobooth.metaDescription')),
    },
  ],
})
</script>

<template>
  <main class="overflow-x-hidden bg-white">
    <Navbar />

    <section class="relative flex min-h-screen items-center justify-center bg-[#FAF9F6] px-6 pt-32 pb-20 lg:px-10">
      <div class="mx-auto max-w-xl text-center">
        <span class="inline-flex items-center gap-2 rounded-full bg-white px-4 py-1.5 font-dm-sans text-xs font-semibold tracking-wide text-[#920f0f] shadow-sm ring-1 ring-[#E4E2DC]">
          <Icon name="heroicons:wrench-screwdriver" class="text-base" />
          {{ t('softwarePhotobooth.badge') }}
        </span>

        <h1 class="mt-6 font-dm-serif text-4xl leading-tight font-bold text-[#000000] sm:text-5xl">
          {{ t('softwarePhotobooth.title') }}
        </h1>

        <p class="mt-4 font-poppins text-base leading-relaxed text-[#57607A] sm:text-lg">
          {{ t('softwarePhotobooth.description') }}
        </p>

        <div class="mt-10 flex justify-center">
          <NuxtLink
            to="/"
            class="flex items-center gap-2 rounded-full bg-[#920f0f] px-8 py-3.5 text-sm font-semibold text-[#FAF9F6] shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5"
          >
            <Icon name="heroicons:arrow-left" />
            {{ t('softwarePhotobooth.backCta') }}
          </NuxtLink>
        </div>
      </div>
    </section>

    <Footer />
    <WhatsappCta />
  </main>
</template>
