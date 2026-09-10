<script setup lang="ts">
import type { PhotoboothFrameItem } from '~/composables/api/photoboothFrames'

defineProps<{ frames: PhotoboothFrameItem[] }>()
const emit = defineEmits<{ select: [frame: PhotoboothFrameItem] }>()

const { t } = useI18n()
</script>

<template>
  <div>
    <h3 class="text-center font-dm-serif text-2xl font-bold text-[#000000]">{{ t('tryModal.frame.title') }}</h3>
    <p class="mt-1 text-center font-poppins text-sm text-[#57607A]">{{ t('tryModal.frame.subtitle') }}</p>

    <div v-if="frames.length > 0" class="mt-6 grid grid-cols-2 gap-4">
      <button
        v-for="frame in frames"
        :key="frame.id"
        class="group overflow-hidden rounded-2xl bg-white text-left shadow-sm ring-1 ring-[#E4E2DC] transition hover:-translate-y-1 hover:shadow-lg"
        @click="emit('select', frame)"
      >
        <div class="aspect-[3/5] w-full overflow-hidden bg-[#F0EFEA]">
          <img
            :src="frame.imageUrl"
            :alt="frame.name"
            class="h-full w-full object-contain transition duration-300 group-hover:scale-105"
          />
        </div>
        <p class="px-3 py-2 font-poppins text-sm font-semibold text-[#1E2537]">{{ frame.name }}</p>
      </button>
    </div>
    <p v-else class="mt-8 text-center font-poppins text-sm text-[#57607A]">
      {{ t('tryModal.frame.empty') }}
    </p>
  </div>
</template>
