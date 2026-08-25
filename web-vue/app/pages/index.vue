<script setup lang="ts">
import type { FrameTemplateItem } from '~/composables/api/frameTemplates'
import type { BackdropItem } from '~/composables/api/backdrops'

const frameTemplatesData = await useServerFetch<FrameTemplateItem[]>('/frame-templates', 'landing-frame-templates')
const backdropsData = await useServerFetch<BackdropItem[]>('/backdrops', 'landing-backdrops')

const frameTemplates = computed(() => frameTemplatesData.value ?? [])
const backdrops = computed(() => backdropsData.value ?? [])
</script>

<template>
  <main class="overflow-x-hidden">
    <Navbar />
    <Hero />
    <Features />
    <FrameTemplates :templates="frameTemplates" />
    <Backdrops :backdrops="backdrops" />
    <Packages />
    <Gallery />
    <Testimonials />
    <Footer />
    <WhatsappCta />
  </main>
</template>
