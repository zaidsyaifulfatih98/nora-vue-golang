<script setup lang="ts">
import type { PhotoboothFrameItem } from '~/composables/api/photoboothFrames'

const framesData = await useServerFetch<PhotoboothFrameItem[]>('/photobooth-frames', 'software-photobooth-frames')
const frames = computed(() => framesData.value ?? [])
</script>

<template>
  <Kiosk :frames="frames" />
</template>
