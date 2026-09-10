<script setup lang="ts">
import type { PhotoboothResultItem } from '~/composables/api/photoboothResults'

defineProps<{
  resultImage: string
  savedResult: PhotoboothResultItem | null
  saving: boolean
  saveError: string
  qrCodeDataUrl: string
}>()
const emit = defineEmits<{ save: []; print: []; tryAgain: [] }>()

const { t } = useI18n()
</script>

<template>
  <div>
    <h3 class="text-center font-dm-serif text-2xl font-bold text-[#000000]">{{ t('tryModal.finish.title') }}</h3>
    <p class="mt-1 text-center font-poppins text-sm text-[#57607A]">{{ t('tryModal.finish.subtitle') }}</p>

    <div class="mx-auto mt-6 max-w-[220px] overflow-hidden rounded-2xl shadow-lg ring-1 ring-[#E4E2DC]">
      <img :src="resultImage" :alt="t('tryModal.result.alt')" class="w-full" />
    </div>

    <div class="mt-6 flex flex-col items-center gap-3">
      <p v-if="saving" class="font-poppins text-sm text-[#57607A]">{{ t('tryModal.result.saving') }}</p>
      <button
        v-else-if="!savedResult"
        class="flex items-center gap-2 rounded-full bg-[#920f0f] px-8 py-3 text-sm font-semibold text-white shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5"
        @click="emit('save')"
      >
        <Icon name="heroicons:cloud-arrow-up" />
        {{ t('tryModal.result.save') }}
      </button>

      <p v-if="saveError" class="font-poppins text-xs text-red-600">{{ saveError }}</p>

      <div v-if="qrCodeDataUrl && savedResult" class="mt-2 flex flex-col items-center gap-3 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-[#E4E2DC]">
        <img :src="qrCodeDataUrl" :alt="t('tryModal.result.downloadAlt')" class="h-40 w-40" />
        <p class="max-w-[220px] text-center font-poppins text-xs text-[#57607A]">
          {{ t('tryModal.result.downloadHint') }}
        </p>
        <div class="flex flex-wrap items-center justify-center gap-3">
          <a
            :href="savedResult.downloadUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-2 rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
          >
            <Icon name="heroicons:arrow-down-tray" />
            {{ t('tryModal.result.download') }}
          </a>
          <button
            type="button"
            class="flex items-center gap-2 rounded-full border border-[#920f0f] px-6 py-2.5 text-sm font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5"
            @click="emit('print')"
          >
            <Icon name="lucide:printer" />
            {{ t('tryModal.result.print') }}
          </button>
        </div>
      </div>

      <button class="mt-2 font-poppins text-sm font-semibold text-[#57607A] underline underline-offset-4" @click="emit('tryAgain')">
        {{ t('tryModal.voice.tryAnotherFrame') }}
      </button>
    </div>
  </div>
</template>
