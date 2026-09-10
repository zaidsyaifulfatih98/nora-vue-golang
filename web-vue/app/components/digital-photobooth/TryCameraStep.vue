<script setup lang="ts">
const props = defineProps<{ photoCount: number; frameLoadError?: string; compositing?: boolean }>()
const emit = defineEmits<{ captured: [photos: string[]] }>()

const { t } = useI18n()
const { brightness, contrast, filterPreset, cameraFilterCss, resetCameraAdjustments } = useCameraAdjustments()

const photoCount = computed(() => props.photoCount)
const {
  videoRef,
  cameraError,
  isCameraReady,
  videoDevices,
  selectedDeviceId,
  countdown,
  capturing,
  photos,
  photosLeft,
  startCamera,
  switchCamera,
  startCaptureCountdown,
  retakePhotos,
} = usePhotoboothCamera({
  photoCount,
  filterCss: cameraFilterCss,
  initialError: props.frameLoadError,
  onAllCaptured: (captured) => emit('captured', captured),
})

const showCameraSettings = ref(false)

// Don't auto-start when arriving with a compositing failure to report —
// startCamera() clears cameraError immediately, which would wipe the
// message before the guest ever sees it. They resume via the "Coba Lagi"
// button instead, matching the pre-refactor behavior.
onMounted(() => {
  if (!props.frameLoadError) startCamera()
})
</script>

<template>
  <div>
    <div
      class="fixed inset-0 z-20 overflow-hidden bg-black sm:relative sm:z-auto sm:mx-auto sm:mt-4 sm:aspect-square sm:w-full sm:max-w-sm sm:rounded-2xl"
    >
      <video
        v-show="isCameraReady"
        ref="videoRef"
        class="h-full w-full -scale-x-100 object-cover"
        :style="{ filter: cameraFilterCss }"
        muted
        playsinline
      />

      <button
        v-if="isCameraReady"
        type="button"
        :aria-label="t('tryModal.camera.cameraSettings')"
        class="absolute top-2 right-2 z-10 flex h-8 w-8 items-center justify-center rounded-full bg-black/50 text-white transition hover:bg-black/70"
        @click="showCameraSettings = !showCameraSettings"
      >
        <Icon name="heroicons:adjustments-horizontal" class="text-sm" />
      </button>

      <div
        v-if="showCameraSettings"
        class="absolute top-11 right-2 z-10 w-56 rounded-xl bg-white p-3 text-left shadow-xl ring-1 ring-[#E4E2DC]"
      >
        <div class="flex items-center justify-between">
          <h4 class="font-poppins text-xs font-semibold text-[#1E2537]">{{ t('tryModal.camera.cameraSettings') }}</h4>
          <button type="button" :aria-label="t('tryModal.closeAria')" class="text-gray-400 hover:text-gray-600" @click="showCameraSettings = false">
            <Icon name="heroicons:x-mark" class="text-sm" />
          </button>
        </div>

        <div v-if="videoDevices.length > 1" class="mt-2">
          <p class="font-poppins text-[11px] font-medium text-gray-600">{{ t('tryModal.camera.cameraDevice') }}</p>
          <div class="mt-1.5 flex flex-col gap-1">
            <button
              v-for="device in videoDevices"
              :key="device.deviceId"
              type="button"
              class="w-full rounded-lg px-1.5 py-1 text-left font-poppins text-[10px] font-semibold transition"
              :class="selectedDeviceId === device.deviceId ? 'bg-[#920f0f] text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
              @click="switchCamera(device.deviceId)"
            >
              {{ device.label || t('tryModal.camera.unnamedCamera') }}
            </button>
          </div>
        </div>

        <div class="mt-2 flex flex-col gap-1">
          <label class="flex items-center justify-between font-poppins text-[11px] font-medium text-gray-600">
            <span>{{ t('tryModal.camera.brightness') }}</span>
            <span>{{ brightness }}%</span>
          </label>
          <input v-model.number="brightness" type="range" min="50" max="150" class="w-full accent-[#920f0f]" />
        </div>

        <div class="mt-2 flex flex-col gap-1">
          <label class="flex items-center justify-between font-poppins text-[11px] font-medium text-gray-600">
            <span>{{ t('tryModal.camera.contrast') }}</span>
            <span>{{ contrast }}%</span>
          </label>
          <input v-model.number="contrast" type="range" min="50" max="150" class="w-full accent-[#920f0f]" />
        </div>

        <div class="mt-2">
          <p class="font-poppins text-[11px] font-medium text-gray-600">{{ t('tryModal.camera.filter') }}</p>
          <div class="mt-1.5 grid grid-cols-3 gap-1.5">
            <button
              v-for="preset in FILTER_PRESET_KEYS"
              :key="preset"
              type="button"
              class="rounded-lg px-1.5 py-1 font-poppins text-[10px] font-semibold transition"
              :class="filterPreset === preset ? 'bg-[#920f0f] text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
              @click="filterPreset = preset"
            >
              {{ t(`tryModal.camera.filters.${preset}`) }}
            </button>
          </div>
        </div>

        <button
          type="button"
          class="mt-2 w-full rounded-lg border border-gray-200 py-1 font-poppins text-[10px] font-medium text-gray-500 transition hover:bg-gray-50"
          @click="resetCameraAdjustments"
        >
          {{ t('tryModal.camera.resetFilters') }}
        </button>
      </div>

      <div v-if="!isCameraReady && !cameraError" class="flex h-full w-full items-center justify-center text-white/70">
        <Icon name="heroicons:video-camera" class="text-4xl animate-pulse" />
      </div>
      <div v-if="cameraError" class="flex h-full w-full flex-col items-center justify-center gap-3 px-6 text-center">
        <Icon name="heroicons:exclamation-triangle" class="text-3xl text-white" />
        <p class="font-poppins text-xs text-white">{{ cameraError }}</p>
        <button
          class="rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5"
          @click="startCamera"
        >
          {{ t('tryModal.camera.retry') }}
        </button>
      </div>
      <div v-if="countdown > 0" class="absolute inset-0 flex items-center justify-center bg-black/40">
        <span class="font-dm-serif text-7xl font-bold text-white">{{ countdown }}</span>
      </div>

      <div v-if="!cameraError" class="absolute inset-x-0 bottom-0 flex flex-col items-center gap-3 bg-gradient-to-t from-black/70 to-transparent px-4 pt-8 pb-5">
        <div v-if="photos.length > 0" class="flex justify-center gap-2">
          <img v-for="(p, i) in photos" :key="i" :src="p" class="h-14 w-14 rounded-lg object-cover ring-2 ring-white shadow" />
        </div>

        <div class="flex items-center justify-center gap-3">
          <button
            :disabled="!isCameraReady || capturing || photosLeft <= 0 || compositing"
            class="rounded-full bg-[#920f0f] px-6 py-2.5 text-sm font-semibold text-white shadow transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-50"
            @click="startCaptureCountdown"
          >
            {{ compositing ? t('tryModal.camera.processing') : capturing ? t('tryModal.camera.preparing') : t('tryModal.camera.takePhoto', { left: photosLeft }) }}
          </button>
          <button
            v-if="photos.length > 0 && !capturing"
            class="rounded-full border border-white/60 px-6 py-2.5 text-sm font-semibold text-white transition hover:bg-white/10"
            @click="retakePhotos"
          >
            {{ t('tryModal.camera.retake') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
