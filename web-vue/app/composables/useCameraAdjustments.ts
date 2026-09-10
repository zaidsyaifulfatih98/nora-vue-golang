const STORAGE_KEY = 'photobooth_camera_adjustments'

export type FilterPreset = 'none' | 'bw' | 'sepia' | 'vintage' | 'cool' | 'warm'

// Brightness/contrast are always applied via CSS filter() functions;
// presets only add color-grading on top (saturate/hue-rotate/sepia/
// grayscale) so the two controls never fight each other.
export const FILTER_PRESETS: Record<FilterPreset, string> = {
  none: '',
  bw: 'grayscale(1)',
  sepia: 'sepia(0.7)',
  vintage: 'sepia(0.3) saturate(1.4) hue-rotate(-10deg)',
  cool: 'hue-rotate(15deg) saturate(1.1)',
  warm: 'hue-rotate(-15deg) saturate(1.2)',
}
export const FILTER_PRESET_KEYS = Object.keys(FILTER_PRESETS) as FilterPreset[]

function loadStoredAdjustments() {
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || 'null')
    return {
      brightness: saved?.brightness ?? 100,
      contrast: saved?.contrast ?? 100,
      filterPreset: FILTER_PRESET_KEYS.includes(saved?.filterPreset) ? saved.filterPreset : 'none',
    }
  } catch {
    return { brightness: 100, contrast: 100, filterPreset: 'none' as FilterPreset }
  }
}

// Shared brightness/contrast/filter calibration for both the digital
// photobooth camera step and the software photobooth kiosk — same
// localStorage key, so a booth only needs to be calibrated once regardless
// of which flow guests use.
export function useCameraAdjustments() {
  const stored = loadStoredAdjustments()
  const brightness = ref(stored.brightness)
  const contrast = ref(stored.contrast)
  const filterPreset = ref<FilterPreset>(stored.filterPreset)

  // Combined into one CSS filter() string used both for the live <video>
  // style and (via canvas ctx.filter) baked into the captured photo — CSS
  // filters on a <video> element don't carry over to drawImage() on their
  // own.
  const cameraFilterCss = computed(() =>
    `brightness(${brightness.value}%) contrast(${contrast.value}%) ${FILTER_PRESETS[filterPreset.value]}`.trim(),
  )

  watch([brightness, contrast, filterPreset], () => {
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({ brightness: brightness.value, contrast: contrast.value, filterPreset: filterPreset.value }),
    )
  })

  function resetCameraAdjustments() {
    brightness.value = 100
    contrast.value = 100
    filterPreset.value = 'none'
  }

  return { brightness, contrast, filterPreset, cameraFilterCss, resetCameraAdjustments }
}
