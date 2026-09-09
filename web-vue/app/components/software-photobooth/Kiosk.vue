<script setup lang="ts">
import QRCode from 'qrcode'
import type { PhotoboothFrameItem } from '~/composables/api/photoboothFrames'
import type { PhotoboothResultItem } from '~/composables/api/photoboothResults'

const props = defineProps<{ frames: PhotoboothFrameItem[]; ownerSlug?: string }>()

const { t } = useI18n()
const { savePhotoboothResult } = usePhotoboothResultsApi()

const MAX_CANVAS_SIDE = 2000
const TIMER_OPTIONS = [3, 5, 10]
const CAMERA_STORAGE_KEY = 'photobooth_camera_device_id'
const CAMERA_ADJUST_STORAGE_KEY = 'photobooth_camera_adjustments'

type FilterPreset = 'none' | 'bw' | 'sepia' | 'vintage' | 'cool' | 'warm'
// Brightness/contrast are always applied via CSS filter() functions;
// presets only add color-grading on top (saturate/hue-rotate/sepia/
// grayscale) so the two controls never fight each other.
const FILTER_PRESETS: Record<FilterPreset, string> = {
  none: '',
  bw: 'grayscale(1)',
  sepia: 'sepia(0.7)',
  vintage: 'sepia(0.3) saturate(1.4) hue-rotate(-10deg)',
  cool: 'hue-rotate(15deg) saturate(1.1)',
  warm: 'hue-rotate(-15deg) saturate(1.2)',
}
const FILTER_PRESET_KEYS = Object.keys(FILTER_PRESETS) as FilterPreset[]

const brightness = ref(100)
const contrast = ref(100)
const filterPreset = ref<FilterPreset>('none')
const showCameraSettings = ref(false)

// Combined into one CSS filter() string used both for the live <video> style
// and (via canvas ctx.filter) baked into the captured photo — CSS filters on
// a <video> element don't carry over to drawImage() on their own.
const cameraFilterCss = computed(() =>
  `brightness(${brightness.value}%) contrast(${contrast.value}%) ${FILTER_PRESETS[filterPreset.value]}`.trim(),
)

function resetCameraAdjustments() {
  brightness.value = 100
  contrast.value = 100
  filterPreset.value = 'none'
}

useHead({
  title: computed(() => t('softwarePhotobooth.metaTitle')),
  meta: [{ name: 'description', content: computed(() => t('softwarePhotobooth.metaDescription')) }],
})

const frames = computed(() => props.frames)

type Step = 'welcome' | 'frame' | 'session' | 'result'
const step = ref<Step>('welcome')

const countdownSeconds = ref(TIMER_OPTIONS[0])

const isFullscreen = ref(false)
async function toggleFullscreen() {
  try {
    if (!document.fullscreenElement) {
      await document.documentElement.requestFullscreen()
    } else {
      await document.exitFullscreen()
    }
  } catch {
    // Fullscreen can be denied by the browser/OS — no need to surface an
    // error, the button simply has no effect that time.
  }
}
onMounted(() => {
  selectedDeviceId.value = localStorage.getItem(CAMERA_STORAGE_KEY) || ''
  document.addEventListener('fullscreenchange', () => {
    isFullscreen.value = Boolean(document.fullscreenElement)
  })

  try {
    const saved = JSON.parse(localStorage.getItem(CAMERA_ADJUST_STORAGE_KEY) || 'null')
    if (saved) {
      brightness.value = saved.brightness ?? 100
      contrast.value = saved.contrast ?? 100
      filterPreset.value = FILTER_PRESET_KEYS.includes(saved.filterPreset) ? saved.filterPreset : 'none'
    }
  } catch {
    // Corrupt/missing stored value — just keep the defaults.
  }
})

// Kiosk operators calibrate this once for their venue's lighting, so it's
// worth remembering across sessions rather than resetting for every guest.
watch([brightness, contrast, filterPreset], () => {
  localStorage.setItem(
    CAMERA_ADJUST_STORAGE_KEY,
    JSON.stringify({ brightness: brightness.value, contrast: contrast.value, filterPreset: filterPreset.value }),
  )
})

function startSession() {
  step.value = 'frame'
}

function goBack() {
  if (step.value === 'frame') {
    step.value = 'welcome'
  } else if (step.value === 'session') {
    stopCamera()
    photos.value = []
    selectedFrame.value = null
    step.value = 'frame'
  } else {
    navigateTo('/')
  }
}

const selectedFrame = ref<PhotoboothFrameItem | null>(null)

const videoRef = ref<HTMLVideoElement | null>(null)
const stream = ref<MediaStream | null>(null)
const cameraError = ref('')
const videoDevices = ref<MediaDeviceInfo[]>([])
// Read on the client only — localStorage doesn't exist during SSR, and this
// whole page is desktop-kiosk-only anyway (no need for the value up front).
const selectedDeviceId = ref('')

const countdown = ref(0)
const capturing = ref(false)
const photos = ref<string[]>([])
const compositing = ref(false)
const resultImage = ref('')

const saving = ref(false)
const saveError = ref('')
const savedResult = ref<PhotoboothResultItem | null>(null)
const qrCodeDataUrl = ref('')
const showQrOverlay = ref(false)
const shareCopied = ref(false)

const DEFAULT_PHOTO_COUNT = 3
const photoCount = computed(() => selectedFrame.value?.slots?.length || DEFAULT_PHOTO_COUNT)
const isCameraReady = computed(() => Boolean(stream.value) && !cameraError.value)
const photosLeft = computed(() => photoCount.value - photos.value.length)
const allPhotosTaken = computed(() => photos.value.length >= photoCount.value && photos.value.length > 0)

function selectFrame(frame: PhotoboothFrameItem) {
  selectedFrame.value = frame
  step.value = 'session'
  startCamera()
}

async function refreshDevices() {
  try {
    const devices = await navigator.mediaDevices.enumerateDevices()
    videoDevices.value = devices.filter((d) => d.kind === 'videoinput')
  } catch {
    // Picker just stays empty/hidden if unsupported.
  }
}

async function startCamera() {
  cameraError.value = ''
  photos.value = []
  try {
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: selectedDeviceId.value
        ? { deviceId: { exact: selectedDeviceId.value }, width: { ideal: 1920 }, height: { ideal: 1080 } }
        : { facingMode: 'user', width: { ideal: 1920 }, height: { ideal: 1080 } },
      audio: false,
    })
    await nextTick()
    if (videoRef.value) {
      videoRef.value.srcObject = stream.value
      await videoRef.value.play()
    }
    await refreshDevices()
    if (!selectedDeviceId.value) {
      selectedDeviceId.value = stream.value.getVideoTracks()[0]?.getSettings().deviceId ?? ''
    }
  } catch {
    if (selectedDeviceId.value) {
      selectedDeviceId.value = ''
      localStorage.removeItem(CAMERA_STORAGE_KEY)
      await startCamera()
      return
    }
    cameraError.value = t('softwarePhotobooth.session.cameraError')
  }
}

async function switchCamera(deviceId: string) {
  if (!deviceId || deviceId === selectedDeviceId.value) return
  selectedDeviceId.value = deviceId
  localStorage.setItem(CAMERA_STORAGE_KEY, deviceId)
  stopCamera()
  await startCamera()
}

function stopCamera() {
  stream.value?.getTracks().forEach((track) => track.stop())
  stream.value = null
}

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

async function captureSequence() {
  if (capturing.value || !isCameraReady.value || photosLeft.value <= 0) return
  capturing.value = true

  for (let i = countdown.value; i > 0; i--) {
    countdown.value = i
    await sleep(1000)
  }
  countdown.value = 0

  const video = videoRef.value
  if (video) {
    const canvas = document.createElement('canvas')
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    const ctx = canvas.getContext('2d')
    if (ctx) {
      ctx.translate(canvas.width, 0)
      ctx.scale(-1, 1)
      ctx.filter = cameraFilterCss.value
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
      photos.value.push(canvas.toDataURL('image/jpeg', 0.92))
    }
  }

  capturing.value = false

  if (photos.value.length >= photoCount.value) {
    // Stop the camera once every shot is in, but let the operator confirm
    // with "Lanjutkan" (or retake) instead of compositing immediately —
    // gives a last look before committing to the frame.
    stopCamera()
  }
}

async function startCaptureCountdown() {
  countdown.value = countdownSeconds.value
  await captureSequence()
}

function retakeAllPhotos() {
  photos.value = []
  if (!stream.value) startCamera()
}

function loadImage(src: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => resolve(img)
    img.onerror = reject
    img.src = src
  })
}

function drawCover(ctx: CanvasRenderingContext2D, img: HTMLImageElement, x: number, y: number, w: number, h: number) {
  const imgRatio = img.width / img.height
  const boxRatio = w / h
  let sx = 0
  let sy = 0
  let sw = img.width
  let sh = img.height

  if (imgRatio > boxRatio) {
    sw = img.height * boxRatio
    sx = (img.width - sw) / 2
  } else {
    sh = img.width / boxRatio
    sy = (img.height - sh) / 2
  }

  ctx.drawImage(img, sx, sy, sw, sh, x, y, w, h)
}

async function buildResult() {
  if (!selectedFrame.value) return
  compositing.value = true
  savedResult.value = null
  qrCodeDataUrl.value = ''
  saveError.value = ''

  try {
    const frameImg = await loadImage(selectedFrame.value.imageUrl)

    const scale = Math.min(1, MAX_CANVAS_SIDE / Math.max(frameImg.width, frameImg.height))
    const canvas = document.createElement('canvas')
    canvas.width = Math.round(frameImg.width * scale)
    canvas.height = Math.round(frameImg.height * scale)
    const ctx = canvas.getContext('2d')
    if (!ctx) {
      compositing.value = false
      return
    }

    ctx.fillStyle = '#1a1a1a'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    const count = photoCount.value
    const slots = selectedFrame.value.slots?.length
      ? selectedFrame.value.slots
      : Array.from({ length: count }, (_, i) => ({
          x: 0.045,
          y: 0.045 + i * (1 / count),
          width: 0.91,
          height: 1 / count - 0.03,
        }))

    const photoImages = await Promise.all(photos.value.map((src) => loadImage(src)))
    photoImages.forEach((img, index) => {
      const slot = slots[index]
      if (!slot) return
      drawCover(ctx, img, slot.x * canvas.width, slot.y * canvas.height, slot.width * canvas.width, slot.height * canvas.height)
    })

    ctx.drawImage(frameImg, 0, 0, canvas.width, canvas.height)
    resultImage.value = canvas.toDataURL('image/png')
    step.value = 'result'
    saveResult()
  } catch {
    cameraError.value = t('softwarePhotobooth.session.frameLoadError')
    step.value = 'session'
  } finally {
    compositing.value = false
  }
}

async function saveResult() {
  if (!resultImage.value || saving.value) return
  saving.value = true
  saveError.value = ''
  try {
    const blob = await (await fetch(resultImage.value)).blob()
    const saved = await savePhotoboothResult(blob, `nora-software-photobooth-${Date.now()}.png`, props.ownerSlug)
    savedResult.value = saved
    qrCodeDataUrl.value = await QRCode.toDataURL(saved.downloadUrl, { width: 320, margin: 1 })
  } catch {
    saveError.value = t('softwarePhotobooth.result.saveError')
  } finally {
    saving.value = false
  }
}

async function shareResult() {
  if (!savedResult.value) return
  const nav = navigator as Navigator & { share?: (data: ShareData) => Promise<void> }
  if (nav.share) {
    try {
      await nav.share({ title: 'Nora Photobooth', url: savedResult.value.downloadUrl })
    } catch {
      // Guest cancelled the share sheet — nothing to do.
    }
    return
  }
  try {
    await navigator.clipboard.writeText(savedResult.value.downloadUrl)
    shareCopied.value = true
    setTimeout(() => (shareCopied.value = false), 2500)
  } catch {
    // Clipboard access denied — silently ignore, QR/print still work.
  }
}

function handlePrint(url: string) {
  const printWindow = window.open('', '_blank', 'width=800,height=1000')
  if (!printWindow) return

  printWindow.document.write(`
    <html>
      <head>
        <title>Print</title>
        <style>
          @page { margin: 0; }
          html, body { margin: 0; padding: 0; height: 100%; display: flex; align-items: center; justify-content: center; }
          img { max-width: 100%; max-height: 100vh; }
        </style>
      </head>
      <body>
        <img src="${url}" />
      </body>
    </html>
  `)
  printWindow.document.close()

  const image = printWindow.document.querySelector('img')
  const triggerPrint = () => {
    printWindow.focus()
    printWindow.print()
  }
  if (image?.complete) triggerPrint()
  else image?.addEventListener('load', triggerPrint)
}

function newSession() {
  resultImage.value = ''
  photos.value = []
  savedResult.value = null
  qrCodeDataUrl.value = ''
  saveError.value = ''
  showQrOverlay.value = false
  selectedFrame.value = null
  step.value = 'frame'
}

onBeforeUnmount(() => {
  stopCamera()
})
</script>

<template>
  <div class="fixed inset-0 z-50 flex flex-col overflow-y-auto bg-[#FAF9F6]">
    <button
      v-if="step !== 'result'"
      type="button"
      class="fixed top-4 left-4 z-30 flex items-center gap-2 rounded-full bg-white px-4 py-2 text-sm font-semibold text-[#39445B] shadow-sm ring-1 ring-[#E4E2DC] transition hover:bg-[#f7f3eb]"
      @click="goBack"
    >
      <Icon name="heroicons:arrow-left" />
      {{ t('softwarePhotobooth.backBtn') }}
    </button>

    <button
      v-if="step === 'welcome' || step === 'session'"
      type="button"
      :aria-label="t('softwarePhotobooth.session.cameraSettings')"
      class="fixed top-4 right-4 z-30 flex h-10 w-10 items-center justify-center rounded-full bg-white text-[#39445B] shadow-sm ring-1 ring-[#E4E2DC] transition hover:bg-[#f7f3eb]"
      @click="showCameraSettings = !showCameraSettings"
    >
      <Icon name="heroicons:adjustments-horizontal" />
    </button>

    <div
      v-if="showCameraSettings"
      class="fixed top-16 right-4 z-30 w-72 rounded-2xl bg-white p-4 text-left shadow-xl ring-1 ring-[#E4E2DC]"
    >
      <div class="flex items-center justify-between">
        <h4 class="font-poppins text-sm font-semibold text-[#1E2537]">{{ t('softwarePhotobooth.session.cameraSettings') }}</h4>
        <button type="button" :aria-label="t('tryModal.closeAria')" class="text-gray-400 hover:text-gray-600" @click="showCameraSettings = false">
          <Icon name="heroicons:x-mark" />
        </button>
      </div>

      <div class="mt-3">
        <p class="font-poppins text-xs font-medium text-gray-600">{{ t('softwarePhotobooth.welcome.timerLabel') }}</p>
        <div class="mt-2 grid grid-cols-3 gap-2">
          <button
            v-for="opt in TIMER_OPTIONS"
            :key="opt"
            type="button"
            class="rounded-lg px-2 py-1.5 font-poppins text-xs font-semibold transition"
            :class="countdownSeconds === opt ? 'bg-[#920f0f] text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
            @click="countdownSeconds = opt"
          >
            {{ t('softwarePhotobooth.welcome.timerOptionSeconds', { seconds: opt }) }}
          </button>
        </div>
      </div>

      <div v-if="videoDevices.length > 1" class="mt-3">
        <p class="font-poppins text-xs font-medium text-gray-600">{{ t('softwarePhotobooth.session.cameraDevice') }}</p>
        <div class="mt-2 flex flex-col gap-1.5">
          <button
            v-for="device in videoDevices"
            :key="device.deviceId"
            type="button"
            class="w-full rounded-lg px-2 py-1.5 text-left font-poppins text-xs font-semibold transition"
            :class="selectedDeviceId === device.deviceId ? 'bg-[#920f0f] text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
            @click="switchCamera(device.deviceId)"
          >
            {{ device.label || t('softwarePhotobooth.session.unnamedCamera') }}
          </button>
        </div>
      </div>

      <div class="mt-3 flex flex-col gap-1">
        <label class="flex items-center justify-between font-poppins text-xs font-medium text-gray-600">
          <span>{{ t('softwarePhotobooth.session.brightness') }}</span>
          <span>{{ brightness }}%</span>
        </label>
        <input v-model.number="brightness" type="range" min="50" max="150" class="w-full accent-[#920f0f]" />
      </div>

      <div class="mt-3 flex flex-col gap-1">
        <label class="flex items-center justify-between font-poppins text-xs font-medium text-gray-600">
          <span>{{ t('softwarePhotobooth.session.contrast') }}</span>
          <span>{{ contrast }}%</span>
        </label>
        <input v-model.number="contrast" type="range" min="50" max="150" class="w-full accent-[#920f0f]" />
      </div>

      <div class="mt-3">
        <p class="font-poppins text-xs font-medium text-gray-600">{{ t('softwarePhotobooth.session.filter') }}</p>
        <div class="mt-2 grid grid-cols-3 gap-2">
          <button
            v-for="preset in FILTER_PRESET_KEYS"
            :key="preset"
            type="button"
            class="rounded-lg px-2 py-1.5 font-poppins text-xs font-semibold transition"
            :class="filterPreset === preset ? 'bg-[#920f0f] text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
            @click="filterPreset = preset"
          >
            {{ t(`softwarePhotobooth.session.filters.${preset}`) }}
          </button>
        </div>
      </div>

      <button
        type="button"
        class="mt-3 w-full rounded-lg border border-gray-200 py-1.5 font-poppins text-xs font-medium text-gray-500 transition hover:bg-gray-50"
        @click="resetCameraAdjustments"
      >
        {{ t('softwarePhotobooth.session.resetFilters') }}
      </button>
    </div>

    <div v-if="step === 'welcome'" class="flex flex-1 flex-col items-center justify-center px-6 py-10 text-center">
      <span class="relative flex h-20 w-20 items-center justify-center overflow-hidden rounded-full shadow-lg">
        <img src="/nora_logo.jpg" alt="Nora Photobooth" class="h-full w-full scale-75 object-contain" />
      </span>
      <h1 class="mt-6 font-dm-serif text-4xl font-bold text-[#000000] sm:text-5xl">{{ t('softwarePhotobooth.welcome.title') }}</h1>
      <p class="mt-3 max-w-md font-poppins text-base text-[#57607A]">{{ t('softwarePhotobooth.welcome.subtitle') }}</p>

      <div class="mt-10 flex flex-wrap items-center justify-center gap-4">
        <button
          class="flex items-center gap-2 rounded-full bg-[#920f0f] px-10 py-4 text-base font-semibold text-white shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5"
          @click="startSession"
        >
          <Icon name="heroicons:play" />
          {{ t('softwarePhotobooth.welcome.startBtn') }}
        </button>
        <button
          class="flex items-center gap-2 rounded-full border border-[#920f0f] px-8 py-4 text-base font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5"
          @click="toggleFullscreen"
        >
          <Icon :name="isFullscreen ? 'heroicons:arrows-pointing-in' : 'heroicons:arrows-pointing-out'" />
          {{ isFullscreen ? t('softwarePhotobooth.welcome.exitFullscreenBtn') : t('softwarePhotobooth.welcome.fullscreenBtn') }}
        </button>
      </div>
    </div>

    <div v-else-if="step === 'frame'" class="flex flex-1 flex-col items-center px-6 py-10">
      <h2 class="font-dm-serif text-3xl font-bold text-[#000000]">{{ t('softwarePhotobooth.frame.title') }}</h2>
      <p class="mt-1 font-poppins text-sm text-[#57607A]">{{ t('softwarePhotobooth.frame.subtitle') }}</p>

      <div v-if="frames.length > 0" class="mt-8 grid w-full max-w-4xl grid-cols-2 gap-6 sm:grid-cols-3 lg:grid-cols-4">
        <button
          v-for="frame in frames"
          :key="frame.id"
          class="group overflow-hidden rounded-2xl bg-white text-left shadow-sm ring-1 ring-[#E4E2DC] transition hover:-translate-y-1 hover:shadow-lg"
          @click="selectFrame(frame)"
        >
          <div class="aspect-[3/5] w-full overflow-hidden bg-[#F0EFEA]">
            <img :src="frame.imageUrl" :alt="frame.name" class="h-full w-full object-contain transition duration-300 group-hover:scale-105" />
          </div>
          <p class="px-3 py-2 text-center font-poppins text-sm font-semibold text-[#1E2537]">{{ frame.name }}</p>
        </button>
      </div>
      <p v-else class="mt-10 font-poppins text-sm text-[#57607A]">{{ t('softwarePhotobooth.frame.empty') }}</p>
    </div>

    <div v-else-if="step === 'session'" class="absolute inset-0 z-0 bg-black">
      <video
        v-show="isCameraReady"
        ref="videoRef"
        class="h-full w-full -scale-x-100 object-cover"
        :style="{ filter: cameraFilterCss }"
        muted
        playsinline
      />

      <div v-if="!isCameraReady && !cameraError" class="flex h-full w-full items-center justify-center text-white/70">
        <Icon name="heroicons:video-camera" class="animate-pulse text-5xl" />
      </div>
      <div v-if="cameraError" class="flex h-full w-full flex-col items-center justify-center gap-3 px-6 text-center">
        <Icon name="heroicons:exclamation-triangle" class="text-4xl text-white" />
        <p class="font-poppins text-sm text-white">{{ cameraError }}</p>
        <button class="rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition" @click="startCamera">
          {{ t('softwarePhotobooth.session.retryCamera') }}
        </button>
      </div>
      <div v-if="countdown > 0" class="absolute inset-0 flex items-center justify-center bg-black/40">
        <span class="font-dm-serif text-9xl font-bold text-white">{{ countdown }}</span>
      </div>

      <div class="absolute inset-x-0 bottom-0 flex flex-col items-center gap-4 bg-gradient-to-t from-black/70 to-transparent px-4 pt-10 pb-6">
        <div class="flex justify-center gap-3">
          <div
            v-for="i in photoCount"
            :key="i"
            class="relative aspect-video w-20 shrink-0 overflow-hidden rounded-lg bg-white/90 ring-1 ring-white/40 sm:w-28"
          >
            <img v-if="photos[i - 1]" :src="photos[i - 1]" class="h-full w-full object-cover" />
            <div v-else class="flex h-full w-full items-center justify-center bg-white/20 text-white/70">
              <Icon name="heroicons:camera" class="text-lg" />
            </div>
            <span class="absolute top-1 left-1 rounded-full bg-black/60 px-1.5 py-0.5 text-[10px] font-bold text-white">{{ i }}</span>
          </div>
        </div>

        <p v-if="compositing" class="font-poppins text-sm text-white">{{ t('softwarePhotobooth.session.processing') }}</p>

        <div v-else-if="!cameraError && allPhotosTaken" class="flex items-center justify-center gap-3">
          <button
            class="rounded-full border border-white/60 px-6 py-3.5 text-sm font-semibold text-white transition hover:bg-white/10"
            @click="retakeAllPhotos"
          >
            {{ t('softwarePhotobooth.session.retakeBtn') }}
          </button>
          <button
            class="rounded-full bg-[#920f0f] px-8 py-3.5 text-base font-semibold text-white shadow-lg transition hover:-translate-y-0.5"
            @click="buildResult"
          >
            {{ t('softwarePhotobooth.session.continueBtn') }}
          </button>
        </div>

        <div v-else-if="!cameraError && photos.length > 0 && !capturing" class="flex items-center justify-center">
          <button
            class="rounded-full border border-white/60 px-6 py-3.5 text-sm font-semibold text-white transition hover:bg-white/10"
            @click="retakeAllPhotos"
          >
            {{ t('softwarePhotobooth.session.retakeBtn') }}
          </button>
        </div>
      </div>

      <button
        v-if="!cameraError && !allPhotosTaken"
        :aria-label="t('softwarePhotobooth.session.takePhotoBtn', { left: photosLeft })"
        :disabled="!isCameraReady || capturing || photosLeft <= 0"
        class="absolute top-1/2 right-6 z-10 h-20 w-20 shrink-0 -translate-y-1/2 rounded-full bg-white shadow-lg ring-4 ring-white/30 transition hover:scale-105 disabled:cursor-not-allowed disabled:opacity-40"
        @click="startCaptureCountdown"
      />
    </div>

    <div v-else-if="step === 'result'" class="flex flex-1 flex-col items-center px-6 py-10">
      <h2 class="font-dm-serif text-3xl font-bold text-[#000000]">{{ t('softwarePhotobooth.result.title') }}</h2>
      <p class="mt-1 font-poppins text-sm text-[#57607A]">{{ t('softwarePhotobooth.result.subtitle') }}</p>

      <div class="mx-auto mt-6 max-h-[55vh] max-w-2xl overflow-hidden rounded-2xl shadow-xl ring-1 ring-[#E4E2DC]">
        <img :src="resultImage" :alt="t('softwarePhotobooth.result.alt')" class="max-h-[55vh] w-full object-contain" />
      </div>

      <p v-if="saving" class="mt-4 font-poppins text-sm text-[#57607A]">{{ t('softwarePhotobooth.result.saving') }}</p>
      <p v-if="saveError" class="mt-4 font-poppins text-xs text-red-600">{{ saveError }}</p>

      <div class="mt-6 flex flex-wrap items-center justify-center gap-3">
        <button
          class="flex items-center gap-2 rounded-full bg-[#1E2537] px-6 py-3 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
          @click="shareResult"
        >
          <Icon name="heroicons:share" />
          {{ shareCopied ? t('softwarePhotobooth.result.shareCopied') : t('softwarePhotobooth.result.shareBtn') }}
        </button>
        <button
          :disabled="!savedResult"
          class="flex items-center gap-2 rounded-full border border-[#920f0f] px-6 py-3 text-sm font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5 disabled:cursor-not-allowed disabled:opacity-50"
          @click="showQrOverlay = true"
        >
          <Icon name="heroicons:qr-code" />
          {{ t('softwarePhotobooth.result.qrBtn') }}
        </button>
        <button
          class="flex items-center gap-2 rounded-full border border-[#920f0f] px-6 py-3 text-sm font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5"
          @click="handlePrint(resultImage)"
        >
          <Icon name="lucide:printer" />
          {{ t('softwarePhotobooth.result.printBtn') }}
        </button>
      </div>

      <button class="mt-8 font-poppins text-sm font-semibold text-[#57607A] underline underline-offset-4" @click="newSession">
        {{ t('softwarePhotobooth.result.newSessionBtn') }}
      </button>

      <div
        v-if="showQrOverlay"
        class="fixed inset-0 z-[70] flex items-center justify-center bg-black/70 px-4"
        @click.self="showQrOverlay = false"
      >
        <div class="relative flex flex-col items-center gap-4 rounded-3xl bg-white p-8 shadow-2xl">
          <button
            :aria-label="t('tryModal.closeAria')"
            class="absolute top-4 right-4 text-[#1E2537] hover:text-[#920f0f]"
            @click="showQrOverlay = false"
          >
            <Icon name="heroicons:x-mark" class="text-xl" />
          </button>
          <h3 class="font-dm-serif text-xl font-bold text-[#000000]">{{ t('softwarePhotobooth.result.qrOverlayTitle') }}</h3>
          <img v-if="qrCodeDataUrl" :src="qrCodeDataUrl" alt="QR code" class="h-64 w-64" />
        </div>
      </div>
    </div>
  </div>
</template>
