<script setup lang="ts">
const CANVAS_W = 400
const CANVAS_H = 300
const OUTPUT_SCALE = 2
const MIN_ZOOM = 0.5
const MAX_ZOOM = 4

const props = defineProps<{ file: File }>()
const emit = defineEmits<{ cancel: []; confirm: [file: File] }>()

const image = ref<HTMLImageElement | null>(null)
const zoom = ref(1)
const offset = reactive({ x: 0, y: 0 })
const dragState = ref<{ startX: number; startY: number; originX: number; originY: number } | null>(null)

let objectUrl: string | null = null

function clampZoom(value: number) {
  return Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, value))
}

onMounted(() => {
  objectUrl = URL.createObjectURL(props.file)
  const img = new Image()
  img.onload = () => {
    image.value = img
  }
  img.src = objectUrl
})

onBeforeUnmount(() => {
  if (objectUrl) URL.revokeObjectURL(objectUrl)
})

const baseScale = computed(() => {
  if (!image.value) return 1
  return Math.min(CANVAS_W / image.value.width, CANVAS_H / image.value.height)
})

const displayScale = computed(() => baseScale.value * zoom.value)

const imageStyle = computed(() => {
  if (!image.value) return {}
  return {
    width: `${image.value.width}px`,
    height: `${image.value.height}px`,
    transform: `translate(-50%, -50%) translate(${offset.x}px, ${offset.y}px) scale(${displayScale.value})`,
    transformOrigin: 'center',
  }
})

function handlePointerDown(e: PointerEvent) {
  ;(e.target as HTMLElement).setPointerCapture(e.pointerId)
  dragState.value = { startX: e.clientX, startY: e.clientY, originX: offset.x, originY: offset.y }
}

function handlePointerMove(e: PointerEvent) {
  if (!dragState.value) return
  const { startX, startY, originX, originY } = dragState.value
  offset.x = originX + (e.clientX - startX)
  offset.y = originY + (e.clientY - startY)
}

function handlePointerUp() {
  dragState.value = null
}

function handleWheel(e: WheelEvent) {
  e.preventDefault()
  zoom.value = clampZoom(zoom.value - e.deltaY * 0.001)
}

function handleConfirm() {
  if (!image.value) return

  const canvas = document.createElement('canvas')
  canvas.width = CANVAS_W * OUTPUT_SCALE
  canvas.height = CANVAS_H * OUTPUT_SCALE
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  const ratio = OUTPUT_SCALE
  const drawWidth = image.value.width * displayScale.value * ratio
  const drawHeight = image.value.height * displayScale.value * ratio
  const centerX = canvas.width / 2 + offset.x * ratio
  const centerY = canvas.height / 2 + offset.y * ratio

  ctx.drawImage(image.value, centerX - drawWidth / 2, centerY - drawHeight / 2, drawWidth, drawHeight)

  canvas.toBlob((blob) => {
    if (!blob) return
    const croppedFile = new File([blob], props.file.name.replace(/\.\w+$/, '.png'), { type: 'image/png' })
    emit('confirm', croppedFile)
  }, 'image/png')
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
    <div class="w-full max-w-md rounded-2xl bg-white p-5 shadow-xl">
      <h3 class="text-sm font-semibold text-gray-800">Atur Gambar</h3>
      <p class="mt-1 text-xs text-gray-500">
        Geser gambar untuk mengatur posisi, gunakan slider atau scroll untuk zoom in/out. Area kosong akan diisi putih.
      </p>

      <div
        class="relative mt-4 touch-none overflow-hidden rounded-lg border border-gray-200 bg-white select-none"
        :style="{ width: `${CANVAS_W}px`, height: `${CANVAS_H}px`, maxWidth: '100%', margin: '0 auto' }"
        @pointerdown="handlePointerDown"
        @pointermove="handlePointerMove"
        @pointerup="handlePointerUp"
        @pointerleave="handlePointerUp"
        @wheel="handleWheel"
      >
        <img
          v-if="image"
          :src="image.src"
          alt="preview"
          draggable="false"
          class="pointer-events-none absolute top-1/2 left-1/2"
          :style="imageStyle"
        />
      </div>

      <div class="mt-4 flex items-center gap-3">
        <span class="text-xs text-gray-500">Zoom</span>
        <input
          type="range"
          :min="MIN_ZOOM"
          :max="MAX_ZOOM"
          step="0.01"
          :value="zoom"
          class="flex-1"
          @input="zoom = clampZoom(Number(($event.target as HTMLInputElement).value))"
        />
      </div>

      <div class="mt-5 flex justify-end gap-2">
        <button
          type="button"
          class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50"
          @click="emit('cancel')"
        >
          Batal
        </button>
        <button
          type="button"
          :disabled="!image"
          class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-gray-300"
          @click="handleConfirm"
        >
          Terapkan
        </button>
      </div>
    </div>
  </div>
</template>
