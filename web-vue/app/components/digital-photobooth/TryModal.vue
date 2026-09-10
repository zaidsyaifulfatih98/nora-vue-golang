<script setup lang="ts">
import QRCode from 'qrcode'
import type { PhotoboothFrameItem } from '~/composables/api/photoboothFrames'
import type { PhotoboothResultItem } from '~/composables/api/photoboothResults'

const props = defineProps<{ frames: PhotoboothFrameItem[]; ownerSlug?: string }>()
const emit = defineEmits<{ close: [] }>()

const { savePhotoboothResult } = usePhotoboothResultsApi()
const { compositeFrame } = usePhotoboothCompositor()
const { printImage } = usePrintImage()
const { t } = useI18n()

const DEFAULT_PHOTO_COUNT = 3
const MAX_CANVAS_SIDE = 1600

type Step = 'frame' | 'camera' | 'result' | 'voice' | 'finish'

const step = ref<Step>('frame')
const selectedFrame = ref<PhotoboothFrameItem | null>(null)
const cameraFrameLoadError = ref('')

// How many photos to take is driven by however many slots the chosen frame
// was configured with in the dashboard, not a fixed number.
const photoCount = computed(() => selectedFrame.value?.slots?.length || DEFAULT_PHOTO_COUNT)

const compositing = ref(false)
const resultImage = ref('')

const saving = ref(false)
const saveError = ref('')
const savedResult = ref<PhotoboothResultItem | null>(null)
const qrCodeDataUrl = ref('')

function selectFrame(frame: PhotoboothFrameItem) {
  selectedFrame.value = frame
  cameraFrameLoadError.value = ''
  step.value = 'camera'
}

async function buildResult(photos: string[]) {
  if (!selectedFrame.value) return
  compositing.value = true
  savedResult.value = null
  qrCodeDataUrl.value = ''
  saveError.value = ''

  try {
    resultImage.value = await compositeFrame(
      selectedFrame.value.imageUrl,
      photos,
      selectedFrame.value.slots ?? [],
      photoCount.value,
      MAX_CANVAS_SIDE,
    )
    step.value = 'result'
    // Save/QR/download now live on the "finish" step after the voice
    // message, but the upload still needs to happen as soon as the photo is
    // ready — sendVoiceMessage() links the voice note to savedResult.viewUrl,
    // and the guest may reach that step before a click-triggered save would
    // have finished.
    saveResult()
  } catch {
    cameraFrameLoadError.value = t('tryModal.camera.frameLoadError')
    step.value = 'camera'
  } finally {
    compositing.value = false
  }
}

function redoPhotos() {
  resultImage.value = ''
  cameraFrameLoadError.value = ''
  step.value = 'camera'
}

async function saveResult() {
  if (!resultImage.value || saving.value) return
  saving.value = true
  saveError.value = ''
  try {
    const blob = await (await fetch(resultImage.value)).blob()
    const saved = await savePhotoboothResult(blob, `nora-digital-photobooth-${Date.now()}.png`, props.ownerSlug)
    savedResult.value = saved
    qrCodeDataUrl.value = await QRCode.toDataURL(saved.downloadUrl, { width: 240, margin: 1 })
  } catch {
    saveError.value = t('tryModal.result.saveError')
  } finally {
    saving.value = false
  }
}

function goToVoiceStep() {
  step.value = 'voice'
}

function tryAgain() {
  resultImage.value = ''
  savedResult.value = null
  qrCodeDataUrl.value = ''
  saveError.value = ''
  step.value = 'frame'
  selectedFrame.value = null
}

function handleClose() {
  emit('close')
}
</script>

<template>
  <div class="fixed inset-0 z-[60] flex items-center justify-center bg-black/70 sm:px-4 sm:py-8">
    <div
      class="relative flex h-full w-full flex-col overflow-hidden bg-[#FAF9F6] sm:h-auto sm:max-h-[90vh] sm:w-full sm:max-w-lg sm:rounded-3xl sm:shadow-2xl"
    >
      <button
        :aria-label="t('tryModal.closeAria')"
        class="absolute top-4 left-4 z-30 flex h-9 w-9 items-center justify-center rounded-full bg-white/90 text-[#1E2537] shadow-sm transition hover:bg-white"
        @click="handleClose"
      >
        <Icon name="heroicons:x-mark" class="text-lg" />
      </button>

      <div class="flex-1 overflow-y-auto px-6 pt-6 pb-8">
        <div class="mb-6 flex items-center justify-center gap-2">
          <span
            v-for="(label, i) in [t('tryModal.stepLabels.frame'), t('tryModal.stepLabels.camera'), t('tryModal.stepLabels.result'), t('tryModal.stepLabels.voice'), t('tryModal.stepLabels.finish')]"
            :key="label"
            class="flex items-center gap-2 font-dm-sans text-xs font-semibold"
            :class="(['frame', 'camera', 'result', 'voice', 'finish'][i] === step) ? 'text-[#920f0f]' : 'text-[#B8B2A6]'"
          >
            <span
              class="flex h-6 w-6 items-center justify-center rounded-full"
              :class="(['frame', 'camera', 'result', 'voice', 'finish'][i] === step) ? 'bg-[#920f0f] text-white' : 'bg-[#E4E2DC] text-[#7A7568]'"
            >
              {{ i + 1 }}
            </span>
            <span class="hidden sm:inline">{{ label }}</span>
          </span>
        </div>

        <TryFrameStep v-if="step === 'frame'" :frames="frames" @select="selectFrame" />

        <TryCameraStep
          v-else-if="step === 'camera'"
          :photo-count="photoCount"
          :frame-load-error="cameraFrameLoadError"
          :compositing="compositing"
          @captured="buildResult"
        />

        <TryResultStep v-else-if="step === 'result'" :result-image="resultImage" @retry="redoPhotos" @continue="goToVoiceStep" />

        <TryVoiceStep
          v-else-if="step === 'voice'"
          :photo-url="savedResult?.viewUrl"
          :owner-slug="ownerSlug"
          @done="step = 'finish'"
        />

        <TryFinishStep
          v-else-if="step === 'finish'"
          :result-image="resultImage"
          :saved-result="savedResult"
          :saving="saving"
          :save-error="saveError"
          :qr-code-data-url="qrCodeDataUrl"
          @save="saveResult"
          @print="printImage(resultImage)"
          @try-again="tryAgain"
        />
      </div>
    </div>
  </div>
</template>
