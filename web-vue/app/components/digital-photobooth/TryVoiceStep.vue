<script setup lang="ts">
const props = defineProps<{ photoUrl?: string; ownerSlug?: string }>()
const emit = defineEmits<{ done: [] }>()

const { t } = useI18n()
const {
  voiceStep,
  voiceGuestName,
  voiceError,
  voiceAudioUrl,
  voiceSeconds,
  openVoiceForm,
  startVoiceRecording,
  stopVoiceRecording,
  retakeVoiceRecording,
  sendVoiceMessage,
} = useVoiceRecorder({
  getPhotoUrl: () => props.photoUrl,
  getOwnerSlug: () => props.ownerSlug,
})

function skipVoiceMessage() {
  emit('done')
}
</script>

<template>
  <div>
    <h3 class="text-center font-dm-serif text-2xl font-bold text-[#000000]">{{ t('tryModal.voice.title') }}</h3>
    <p class="mt-1 text-center font-poppins text-sm text-[#57607A]">
      {{ t('tryModal.voice.subtitle') }}
    </p>

    <div class="mx-auto mt-6 w-full max-w-xs rounded-2xl bg-white p-4 shadow-sm ring-1 ring-[#E4E2DC]">
      <div v-if="voiceStep === 'idle'" class="flex flex-col items-center gap-3">
        <button
          class="flex w-full items-center justify-center gap-2 rounded-full bg-[#920f0f] px-6 py-3 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
          @click="openVoiceForm"
        >
          <Icon name="heroicons:microphone" />
          {{ t('tryModal.voice.sendCta') }}
        </button>
        <button class="font-poppins text-xs font-semibold text-[#57607A] underline underline-offset-4" @click="skipVoiceMessage">
          {{ t('tryModal.voice.skip') }}
        </button>
      </div>

      <div v-else-if="voiceStep === 'form'" class="flex flex-col gap-3">
        <p class="text-center font-poppins text-sm font-semibold text-[#1E2537]">{{ t('tryModal.voice.formTitle') }}</p>
        <input
          v-model="voiceGuestName"
          :placeholder="t('tryModal.voice.namePlaceholder')"
          class="rounded-lg border border-[#E4E2DC] px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none"
        />
        <button
          class="flex items-center justify-center gap-2 rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
          @click="startVoiceRecording"
        >
          <Icon name="heroicons:microphone" />
          {{ t('tryModal.voice.startRecording') }}
        </button>
        <button class="font-poppins text-xs font-semibold text-[#57607A] underline underline-offset-4" @click="skipVoiceMessage">
          {{ t('tryModal.voice.skipOptional') }}
        </button>
      </div>

      <div v-else-if="voiceStep === 'recording'" class="flex flex-col items-center gap-3">
        <div class="flex h-14 w-14 items-center justify-center rounded-full bg-[#920f0f]/10 ring-2 ring-[#920f0f]">
          <Icon name="heroicons:microphone" class="animate-pulse text-2xl text-[#920f0f]" />
        </div>
        <p class="font-dm-sans text-sm text-[#57607A]">{{ voiceSeconds }}s</p>
        <button
          class="rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
          @click="stopVoiceRecording"
        >
          {{ t('tryModal.voice.stopRecording') }}
        </button>
      </div>

      <div v-else-if="voiceStep === 'recorded'" class="flex flex-col items-center gap-3">
        <audio :src="voiceAudioUrl" controls class="w-full" />
        <div class="flex gap-3">
          <button
            class="rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
            @click="sendVoiceMessage"
          >
            {{ t('tryModal.voice.send') }}
          </button>
          <button
            class="rounded-full border border-[#920f0f] px-6 py-2.5 text-sm font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5"
            @click="retakeVoiceRecording"
          >
            {{ t('tryModal.voice.retake') }}
          </button>
        </div>
      </div>

      <p v-else-if="voiceStep === 'sending'" class="text-center font-poppins text-sm text-[#57607A]">{{ t('tryModal.voice.sending') }}</p>

      <div v-else-if="voiceStep === 'sent'" class="flex flex-col items-center gap-2">
        <Icon name="heroicons:check-circle" class="text-3xl text-[#920f0f]" />
        <p class="text-center font-poppins text-sm text-[#57607A]">{{ t('tryModal.voice.sent') }}</p>
      </div>

      <p v-if="voiceError" class="mt-2 text-center font-poppins text-xs text-red-600">{{ voiceError }}</p>
    </div>

    <div v-if="voiceStep === 'sent'" class="mt-6 flex justify-center">
      <button
        class="flex items-center gap-2 rounded-full bg-[#1E2537] px-8 py-3 text-sm font-semibold text-white shadow-lg transition hover:-translate-y-0.5"
        @click="emit('done')"
      >
        {{ t('tryModal.result.continue') }}
        <Icon name="heroicons:arrow-right" />
      </button>
    </div>
  </div>
</template>
