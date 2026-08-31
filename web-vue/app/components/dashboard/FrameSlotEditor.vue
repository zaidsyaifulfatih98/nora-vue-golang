<script setup lang="ts">
import type { FrameSlot } from '~/composables/api/photoboothFrames'

const MIN_SLOTS = 1
const MAX_SLOTS = 10

const props = defineProps<{ imageUrl: string; modelValue: FrameSlot[] }>()
const emit = defineEmits<{ 'update:modelValue': [FrameSlot[]] }>()

const containerRef = ref<HTMLDivElement | null>(null)
const imgAspect = ref<number | null>(null)

const canAddSlot = computed(() => props.modelValue.length < MAX_SLOTS)
const canRemoveSlot = computed(() => props.modelValue.length > MIN_SLOTS)

function clamp(v: number, min: number, max: number) {
  return Math.min(Math.max(v, min), max)
}

function addSlot() {
  if (!canAddSlot.value) return
  const offset = 0.04 * props.modelValue.length
  const size = 0.35
  emit('update:modelValue', [
    ...props.modelValue,
    { x: clamp(0.1 + offset, 0, 1 - size), y: clamp(0.1 + offset, 0, 1 - size), width: size, height: size },
  ])
}

function removeSlot(index: number) {
  if (!canRemoveSlot.value) return
  emit('update:modelValue', props.modelValue.filter((_, i) => i !== index))
}

function onImgLoad(e: Event) {
  const img = e.target as HTMLImageElement
  if (img.naturalWidth && img.naturalHeight) imgAspect.value = img.naturalWidth / img.naturalHeight
}

type DragMode = 'move' | 'resize'
interface DragState {
  index: number
  mode: DragMode
  startX: number
  startY: number
  orig: FrameSlot
}
let dragState: DragState | null = null

function onPointerDown(e: PointerEvent, index: number, mode: DragMode) {
  e.preventDefault()
  const slot = props.modelValue[index]
  if (!slot) return
  dragState = { index, mode, startX: e.clientX, startY: e.clientY, orig: { ...slot } }
  window.addEventListener('pointermove', onPointerMove)
  window.addEventListener('pointerup', onPointerUp)
}

function onPointerMove(e: PointerEvent) {
  if (!dragState || !containerRef.value) return
  const rect = containerRef.value.getBoundingClientRect()
  const dx = (e.clientX - dragState.startX) / rect.width
  const dy = (e.clientY - dragState.startY) / rect.height
  const { orig, index, mode } = dragState

  const next = props.modelValue.map((s, i) => (i === index ? { ...s } : s))
  const slot = next[index]
  if (!slot) return

  if (mode === 'move') {
    slot.x = clamp(orig.x + dx, 0, 1 - orig.width)
    slot.y = clamp(orig.y + dy, 0, 1 - orig.height)
  } else {
    slot.width = clamp(orig.width + dx, 0.05, 1 - orig.x)
    slot.height = clamp(orig.height + dy, 0.05, 1 - orig.y)
  }

  emit('update:modelValue', next)
}

function onPointerUp() {
  dragState = null
  window.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('pointerup', onPointerUp)
}

onBeforeUnmount(() => {
  window.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('pointerup', onPointerUp)
})
</script>

<template>
  <div
    ref="containerRef"
    class="relative mx-auto w-full max-w-[280px] touch-none overflow-hidden rounded-xl bg-[#111] select-none"
    :style="{ aspectRatio: imgAspect ? String(imgAspect) : '3 / 5' }"
  >
    <img :src="imageUrl" class="pointer-events-none absolute inset-0 h-full w-full object-cover" @load="onImgLoad" />

    <div
      v-for="(slot, i) in modelValue"
      :key="i"
      class="absolute cursor-move border-2 border-dashed border-[#f4b13a] bg-[#f4b13a]/25"
      :style="{
        left: `${slot.x * 100}%`,
        top: `${slot.y * 100}%`,
        width: `${slot.width * 100}%`,
        height: `${slot.height * 100}%`,
      }"
      @pointerdown="onPointerDown($event, i, 'move')"
    >
      <span class="absolute top-1 left-1 rounded bg-black/60 px-1.5 py-0.5 text-[10px] font-bold text-white">{{ i + 1 }}</span>
      <button
        v-if="canRemoveSlot"
        type="button"
        aria-label="Hapus kotak"
        class="absolute top-1 right-1 flex h-5 w-5 items-center justify-center rounded-full bg-black/60 text-white"
        @pointerdown.stop
        @click.stop="removeSlot(i)"
      >
        <Icon name="heroicons:x-mark" class="text-xs" />
      </button>
      <span
        class="absolute right-0 bottom-0 h-4 w-4 cursor-nwse-resize rounded-tl bg-[#f4b13a]"
        @pointerdown.stop="onPointerDown($event, i, 'resize')"
      />
    </div>
  </div>

  <div class="mt-2 flex items-center justify-between gap-2">
    <p class="text-xs text-gray-500">{{ modelValue.length }} kotak foto</p>
    <button
      type="button"
      :disabled="!canAddSlot"
      class="flex items-center gap-1 rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-600 transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40"
      @click="addSlot"
    >
      <Icon name="heroicons:plus" />
      Tambah Kotak
    </button>
  </div>
</template>
