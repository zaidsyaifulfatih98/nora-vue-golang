<script setup lang="ts">
import QRCode from 'qrcode'
import type { PhotoboothFrameItem } from '~/composables/api/photoboothFrames'
import type { PhotoboothResultItem } from '~/composables/api/photoboothResults'

const props = defineProps<{ frames: PhotoboothFrameItem[] }>()
const emit = defineEmits<{ close: [] }>()

const { savePhotoboothResult } = usePhotoboothResultsApi()
const { uploadVoiceMessage } = useVoiceMessagesApi()

const DEFAULT_PHOTO_COUNT = 3
const MAX_CANVAS_SIDE = 1600
const CAPTURE_DELAY_SECONDS = 3

type Step = 'frame' | 'camera' | 'result' | 'voice'

const step = ref<Step>('frame')
const selectedFrame = ref<PhotoboothFrameItem | null>(null)

const videoRef = ref<HTMLVideoElement | null>(null)
const stream = ref<MediaStream | null>(null)
const cameraError = ref('')
const countdown = ref(0)
const capturing = ref(false)
const photos = ref<string[]>([])

const compositing = ref(false)
const resultImage = ref('')

const saving = ref(false)
const saveError = ref('')
const savedResult = ref<PhotoboothResultItem | null>(null)
const qrCodeDataUrl = ref('')

type VoiceStep = 'idle' | 'form' | 'recording' | 'recorded' | 'sending' | 'sent' | 'skipped'
const voiceStep = ref<VoiceStep>('idle')
const voiceGuestName = ref('')
const voiceError = ref('')
const voiceStream = ref<MediaStream | null>(null)
const mediaRecorder = ref<MediaRecorder | null>(null)
const voiceChunks = ref<Blob[]>([])
const voiceBlob = ref<Blob | null>(null)
const voiceAudioUrl = ref('')
const voiceSeconds = ref(0)
let voiceTimer: ReturnType<typeof setInterval> | null = null

// How many photos to take is driven by however many slots the chosen frame
// was configured with in the dashboard, not a fixed number.
const photoCount = computed(() => selectedFrame.value?.slots?.length || DEFAULT_PHOTO_COUNT)

const isCameraReady = computed(() => Boolean(stream.value) && !cameraError.value)
const photosLeft = computed(() => photoCount.value - photos.value.length)

function selectFrame(frame: PhotoboothFrameItem) {
  selectedFrame.value = frame
  step.value = 'camera'
  startCamera()
}

async function startCamera() {
  cameraError.value = ''
  photos.value = []
  try {
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user', width: { ideal: 1280 }, height: { ideal: 1280 } },
      audio: false,
    })
    await nextTick()
    if (videoRef.value) {
      videoRef.value.srcObject = stream.value
      await videoRef.value.play()
    }
  } catch {
    cameraError.value = 'Tidak bisa mengakses kamera. Pastikan kamu mengizinkan akses kamera pada browser.'
  }
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
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
      photos.value.push(canvas.toDataURL('image/jpeg', 0.92))
    }
  }

  capturing.value = false

  if (photos.value.length >= photoCount.value) {
    stopCamera()
    await buildResult()
  }
}

async function startCaptureCountdown() {
  countdown.value = CAPTURE_DELAY_SECONDS
  await captureSequence()
}

function retakePhotos() {
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

    // The canvas is sized to the frame PNG's own aspect ratio so its
    // transparent cutout boxes line up exactly where they were designed,
    // instead of being cropped or stretched to fit an unrelated size.
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

    // Slots are the exact transparent cutout boxes marked in the dashboard
    // (fractions of the frame image's own dimensions). Fall back to equal
    // vertical bands only if a frame was never configured with slots.
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
      drawCover(
        ctx,
        img,
        slot.x * canvas.width,
        slot.y * canvas.height,
        slot.width * canvas.width,
        slot.height * canvas.height,
      )
    })

    ctx.drawImage(frameImg, 0, 0, canvas.width, canvas.height)
    resultImage.value = canvas.toDataURL('image/png')
    step.value = 'result'
  } catch {
    cameraError.value = 'Gagal memuat gambar frame. Silakan coba lagi.'
    step.value = 'camera'
  } finally {
    compositing.value = false
  }
}

function redoPhotos() {
  resultImage.value = ''
  photos.value = []
  savedResult.value = null
  qrCodeDataUrl.value = ''
  saveError.value = ''
  step.value = 'camera'
  startCamera()
}

async function saveResult() {
  if (!resultImage.value || saving.value) return
  saving.value = true
  saveError.value = ''
  try {
    const blob = await (await fetch(resultImage.value)).blob()
    const saved = await savePhotoboothResult(blob, `nora-digital-photobooth-${Date.now()}.png`)
    savedResult.value = saved
    qrCodeDataUrl.value = await QRCode.toDataURL(saved.downloadUrl, { width: 240, margin: 1 })
  } catch {
    saveError.value = 'Gagal menyimpan hasil. Silakan coba lagi.'
  } finally {
    saving.value = false
  }
}

function goToVoiceStep() {
  step.value = 'voice'
}

function openVoiceForm() {
  voiceError.value = ''
  voiceStep.value = 'form'
}

function skipVoiceMessage() {
  stopVoiceStream()
  voiceStep.value = 'skipped'
}

function stopVoiceStream() {
  voiceStream.value?.getTracks().forEach((track) => track.stop())
  voiceStream.value = null
  if (voiceTimer) {
    clearInterval(voiceTimer)
    voiceTimer = null
  }
}

async function startVoiceRecording() {
  voiceError.value = ''
  voiceChunks.value = []
  voiceSeconds.value = 0
  try {
    voiceStream.value = await navigator.mediaDevices.getUserMedia({ audio: true })
    const recorder = new MediaRecorder(voiceStream.value)
    mediaRecorder.value = recorder

    recorder.ondataavailable = (e) => {
      if (e.data.size > 0) voiceChunks.value.push(e.data)
    }
    recorder.onstop = () => {
      voiceBlob.value = new Blob(voiceChunks.value, { type: recorder.mimeType || 'audio/webm' })
      if (voiceAudioUrl.value) URL.revokeObjectURL(voiceAudioUrl.value)
      voiceAudioUrl.value = URL.createObjectURL(voiceBlob.value)
      stopVoiceStream()
      voiceStep.value = 'recorded'
    }

    recorder.start()
    voiceStep.value = 'recording'
    voiceTimer = setInterval(() => (voiceSeconds.value += 1), 1000)
  } catch {
    voiceError.value = 'Tidak bisa mengakses mikrofon. Pastikan kamu mengizinkan akses mikrofon pada browser.'
  }
}

function stopVoiceRecording() {
  mediaRecorder.value?.stop()
}

function retakeVoiceRecording() {
  voiceBlob.value = null
  if (voiceAudioUrl.value) URL.revokeObjectURL(voiceAudioUrl.value)
  voiceAudioUrl.value = ''
  voiceStep.value = 'form'
}

async function sendVoiceMessage() {
  if (!voiceBlob.value) return
  voiceStep.value = 'sending'
  voiceError.value = ''
  try {
    await uploadVoiceMessage(voiceBlob.value, voiceGuestName.value, savedResult.value?.viewUrl)
    voiceStep.value = 'sent'
  } catch {
    voiceError.value = 'Gagal mengirim pesan suara. Silakan coba lagi.'
    voiceStep.value = 'recorded'
  }
}

function resetVoiceMessage() {
  stopVoiceStream()
  if (voiceAudioUrl.value) URL.revokeObjectURL(voiceAudioUrl.value)
  voiceStep.value = 'idle'
  voiceGuestName.value = ''
  voiceError.value = ''
  voiceChunks.value = []
  voiceBlob.value = null
  voiceAudioUrl.value = ''
  voiceSeconds.value = 0
}

function tryAgain() {
  resultImage.value = ''
  photos.value = []
  savedResult.value = null
  qrCodeDataUrl.value = ''
  saveError.value = ''
  resetVoiceMessage()
  step.value = 'frame'
  selectedFrame.value = null
}

function handleClose() {
  stopCamera()
  stopVoiceStream()
  emit('close')
}

onBeforeUnmount(() => {
  stopCamera()
  stopVoiceStream()
})
</script>

<template>
  <div class="fixed inset-0 z-[60] flex items-center justify-center bg-black/70 sm:px-4 sm:py-8">
    <div
      class="relative flex h-full w-full flex-col overflow-hidden bg-[#FAF9F6] sm:h-auto sm:max-h-[90vh] sm:w-full sm:max-w-lg sm:rounded-3xl sm:shadow-2xl"
    >
      <button
        aria-label="Tutup"
        class="absolute top-4 right-4 z-10 flex h-9 w-9 items-center justify-center rounded-full bg-white/90 text-[#1E2537] shadow-sm transition hover:bg-white"
        @click="handleClose"
      >
        <Icon name="heroicons:x-mark" class="text-lg" />
      </button>

      <div class="flex-1 overflow-y-auto px-6 pt-6 pb-8">
        <div class="mb-6 flex items-center justify-center gap-2">
          <span
            v-for="(label, i) in ['Pilih Frame', 'Ambil Foto', 'Hasil', 'Pesan Suara']"
            :key="label"
            class="flex items-center gap-2 font-dm-sans text-xs font-semibold"
            :class="(['frame', 'camera', 'result', 'voice'][i] === step) ? 'text-[#920f0f]' : 'text-[#B8B2A6]'"
          >
            <span
              class="flex h-6 w-6 items-center justify-center rounded-full"
              :class="(['frame', 'camera', 'result', 'voice'][i] === step) ? 'bg-[#920f0f] text-white' : 'bg-[#E4E2DC] text-[#7A7568]'"
            >
              {{ i + 1 }}
            </span>
            <span class="hidden sm:inline">{{ label }}</span>
          </span>
        </div>

        <div v-if="step === 'frame'">
          <h3 class="text-center font-dm-serif text-2xl font-bold text-[#000000]">Pilih Frame Favoritmu</h3>
          <p class="mt-1 text-center font-poppins text-sm text-[#57607A]">Frame ini bisa diatur oleh tim kami lewat dashboard.</p>

          <div v-if="frames.length > 0" class="mt-6 grid grid-cols-2 gap-4">
            <button
              v-for="frame in frames"
              :key="frame.id"
              class="group overflow-hidden rounded-2xl bg-white text-left shadow-sm ring-1 ring-[#E4E2DC] transition hover:-translate-y-1 hover:shadow-lg"
              @click="selectFrame(frame)"
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
            Belum ada frame yang tersedia saat ini. Silakan hubungi kami untuk info lebih lanjut.
          </p>
        </div>

        <div v-else-if="step === 'camera'">
          <h3 class="text-center font-dm-serif text-2xl font-bold text-[#000000]">Ambil {{ photoCount }} Foto Terbaikmu</h3>
          <p class="mt-1 text-center font-poppins text-sm text-[#57607A]">
            Foto {{ photos.length }} dari {{ photoCount }}
          </p>

          <div class="relative mx-auto mt-4 aspect-square w-full max-w-sm overflow-hidden rounded-2xl bg-black">
            <video
              v-show="isCameraReady"
              ref="videoRef"
              class="h-full w-full -scale-x-100 object-cover"
              muted
              playsinline
            />
            <div v-if="!isCameraReady && !cameraError" class="flex h-full w-full items-center justify-center text-white/70">
              <Icon name="heroicons:video-camera" class="text-4xl animate-pulse" />
            </div>
            <div v-if="cameraError" class="flex h-full w-full flex-col items-center justify-center gap-3 px-6 text-center">
              <Icon name="heroicons:exclamation-triangle" class="text-3xl text-white" />
              <p class="font-poppins text-xs text-white">{{ cameraError }}</p>
            </div>
            <div v-if="countdown > 0" class="absolute inset-0 flex items-center justify-center bg-black/40">
              <span class="font-dm-serif text-7xl font-bold text-white">{{ countdown }}</span>
            </div>
          </div>

          <div v-if="photos.length > 0" class="mt-4 flex justify-center gap-2">
            <img v-for="(p, i) in photos" :key="i" :src="p" class="h-14 w-14 rounded-lg object-cover ring-2 ring-white shadow" />
          </div>

          <div class="mt-6 flex items-center justify-center gap-3">
            <button
              v-if="cameraError"
              class="rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
              @click="startCamera"
            >
              Coba Lagi
            </button>
            <template v-else>
              <button
                :disabled="!isCameraReady || capturing || photosLeft <= 0 || compositing"
                class="rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-50"
                @click="startCaptureCountdown"
              >
                {{ compositing ? 'Memproses...' : capturing ? 'Bersiap...' : `Ambil Foto (${photosLeft} lagi)` }}
              </button>
              <button
                v-if="photos.length > 0 && !capturing"
                class="rounded-full border border-[#920f0f] px-6 py-2.5 text-sm font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5"
                @click="retakePhotos"
              >
                Ulangi
              </button>
            </template>
          </div>
        </div>

        <div v-else-if="step === 'result'">
          <h3 class="text-center font-dm-serif text-2xl font-bold text-[#000000]">Hasil Digital Photobooth-mu</h3>
          <p class="mt-1 text-center font-poppins text-sm text-[#57607A]">Simpan dan bagikan momenmu.</p>

          <div class="mx-auto mt-6 max-w-[220px] overflow-hidden rounded-2xl shadow-lg ring-1 ring-[#E4E2DC]">
            <img :src="resultImage" alt="Hasil digital photobooth" class="w-full" />
          </div>

          <div class="mt-6 flex flex-col items-center gap-3">
            <div class="flex flex-wrap items-center justify-center gap-3">
              <button
                class="flex items-center gap-2 rounded-full border border-[#920f0f] px-8 py-3 text-sm font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5"
                @click="redoPhotos"
              >
                <Icon name="heroicons:arrow-path" />
                Ulangi
              </button>
              <button
                :disabled="saving || Boolean(savedResult)"
                class="flex items-center gap-2 rounded-full bg-[#920f0f] px-8 py-3 text-sm font-semibold text-white shadow-lg shadow-[#1E2537]/25 transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60"
                @click="saveResult"
              >
                <Icon name="heroicons:cloud-arrow-up" />
                {{ saving ? 'Menyimpan...' : savedResult ? 'Tersimpan' : 'Simpan' }}
              </button>
            </div>

            <p v-if="saveError" class="font-poppins text-xs text-red-600">{{ saveError }}</p>

            <div v-if="qrCodeDataUrl && savedResult" class="mt-2 flex flex-col items-center gap-3 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-[#E4E2DC]">
              <img :src="qrCodeDataUrl" alt="QR code download hasil" class="h-40 w-40" />
              <p class="max-w-[220px] text-center font-poppins text-xs text-[#57607A]">
                Sudah tersimpan. Scan QR ini atau klik button dibawah ini
              </p>
              <a
                :href="savedResult.downloadUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="flex items-center gap-2 rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
              >
                <Icon name="heroicons:arrow-down-tray" />
                Download
              </a>
            </div>

            <button
              class="mt-2 flex items-center gap-2 rounded-full bg-[#1E2537] px-8 py-3 text-sm font-semibold text-white shadow-lg transition hover:-translate-y-0.5"
              @click="goToVoiceStep"
            >
              Lanjut
              <Icon name="heroicons:arrow-right" />
            </button>
          </div>
        </div>

        <div v-else-if="step === 'voice'">
          <h3 class="text-center font-dm-serif text-2xl font-bold text-[#000000]">Kirim Pesan Suara</h3>
          <p class="mt-1 text-center font-poppins text-sm text-[#57607A]">
            Tinggalkan ucapan untuk pasangan, atau lewati jika tidak ingin.
          </p>

          <div class="mx-auto mt-6 w-full max-w-xs rounded-2xl bg-white p-4 shadow-sm ring-1 ring-[#E4E2DC]">
            <div v-if="voiceStep === 'idle'" class="flex flex-col items-center gap-3">
              <button
                class="flex w-full items-center justify-center gap-2 rounded-full bg-[#920f0f] px-6 py-3 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
                @click="openVoiceForm"
              >
                <Icon name="heroicons:microphone" />
                Kirim Pesan Suara
              </button>
              <button class="font-poppins text-xs font-semibold text-[#57607A] underline underline-offset-4" @click="skipVoiceMessage">
                Lewati 
              </button>
            </div>

            <div v-else-if="voiceStep === 'form'" class="flex flex-col gap-3">
              <p class="text-center font-poppins text-sm font-semibold text-[#1E2537]">Rekam Pesan Suara</p>
              <input
                v-model="voiceGuestName"
                placeholder="Nama kamu (opsional)"
                class="rounded-lg border border-[#E4E2DC] px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none"
              />
              <button
                class="flex items-center justify-center gap-2 rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
                @click="startVoiceRecording"
              >
                <Icon name="heroicons:microphone" />
                Mulai Rekam
              </button>
              <button class="font-poppins text-xs font-semibold text-[#57607A] underline underline-offset-4" @click="skipVoiceMessage">
                Lewati (Opsional)
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
                Berhenti
              </button>
            </div>

            <div v-else-if="voiceStep === 'recorded'" class="flex flex-col items-center gap-3">
              <audio :src="voiceAudioUrl" controls class="w-full" />
              <div class="flex gap-3">
                <button
                  class="rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
                  @click="sendVoiceMessage"
                >
                  Kirim
                </button>
                <button
                  class="rounded-full border border-[#920f0f] px-6 py-2.5 text-sm font-semibold text-[#920f0f] transition hover:bg-[#920f0f]/5"
                  @click="retakeVoiceRecording"
                >
                  Rekam Ulang
                </button>
              </div>
            </div>

            <p v-else-if="voiceStep === 'sending'" class="text-center font-poppins text-sm text-[#57607A]">Mengirim...</p>

            <div v-else-if="voiceStep === 'sent'" class="flex flex-col items-center gap-2">
              <Icon name="heroicons:check-circle" class="text-3xl text-[#920f0f]" />
              <p class="text-center font-poppins text-sm text-[#57607A]">Pesan suara terkirim, terima kasih!</p>
            </div>

            <div v-else-if="voiceStep === 'skipped'" class="flex flex-col items-center gap-2">
              <Icon name="heroicons:heart" class="text-2xl text-[#920f0f]" />
              <p class="text-center font-poppins text-sm text-[#57607A]">Terima kasih sudah mencoba Digital Photobooth kami!</p>
            </div>

            <p v-if="voiceError" class="mt-2 text-center font-poppins text-xs text-red-600">{{ voiceError }}</p>
          </div>

          <div class="mt-6 flex justify-center">
            <button class="font-poppins text-sm font-semibold text-[#57607A] underline underline-offset-4" @click="tryAgain">
              Coba Frame Lain
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
