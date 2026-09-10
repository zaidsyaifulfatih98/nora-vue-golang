<script setup lang="ts">
defineProps<{ resultImage: string }>()
const emit = defineEmits<{ retry: []; continue: [] }>()

const { t } = useI18n()
</script>

<template>
  <div>
    <h3 class="text-center font-dm-serif text-2xl font-bold text-[#000000]">{{ t('tryModal.result.title') }}</h3>
    <p class="mt-1 text-center font-poppins text-sm text-[#57607A]">{{ t('tryModal.result.subtitle') }}</p>

    <div class="mx-auto mt-6 max-w-[220px] overflow-hidden rounded-2xl shadow-lg ring-1 ring-[#E4E2DC]">
      <img :src="resultImage" :alt="t('tryModal.result.alt')" class="w-full" />
    </div>

    <div class="mt-6 flex flex-col items-center gap-3">
      <div class="flex flex-wrap items-center justify-center gap-3">
        <button
          class="flex items-center gap-2 rounded-full border border-[#920f0f] px-8 py-3 text-sm font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5"
          @click="emit('retry')"
        >
          <Icon name="heroicons:arrow-path" />
          {{ t('tryModal.result.retry') }}
        </button>
        <button
          class="flex items-center gap-2 rounded-full bg-[#1E2537] px-8 py-3 text-sm font-semibold text-white shadow-lg transition hover:-translate-y-0.5"
          @click="emit('continue')"
        >
          {{ t('tryModal.result.continue') }}
          <Icon name="heroicons:arrow-right" />
        </button>
      </div>
    </div>
  </div>
</template>
