const DEVICE_STORAGE_KEY = 'photobooth_camera_device_id'
const CAPTURE_DELAY_SECONDS = 3

export interface UsePhotoboothCameraOptions {
  photoCount: Ref<number> | ComputedRef<number>
  filterCss: Ref<string> | ComputedRef<string>
  // Fired once photos.length reaches photoCount — callers composite the
  // result from the returned photos array.
  onAllCaptured?: (photos: string[]) => void
  // Seeds cameraError on creation (e.g. a compositing failure that sent the
  // guest back to this step with a message to show).
  initialError?: string
  // i18n key for the "can't access the camera" message — differs between
  // the digital photobooth camera step and the software photobooth kiosk.
  errorMessageKey?: string
}

// Live camera capture: device enumeration/selection (remembered per-browser
// so a kiosk only needs to be chosen once, e.g. a Sony A6400 exposed as a
// webcam via Sony's "Imaging Edge Webcam" app instead of whatever the
// browser treats as the default), stream lifecycle, and the
// countdown-then-snapshot capture sequence.
export function usePhotoboothCamera(options: UsePhotoboothCameraOptions) {
  const { t } = useI18n()
  const errorMessageKey = options.errorMessageKey ?? 'tryModal.camera.error'

  const videoRef = ref<HTMLVideoElement | null>(null)
  const stream = ref<MediaStream | null>(null)
  const cameraError = ref(options.initialError ?? '')

  const videoDevices = ref<MediaDeviceInfo[]>([])
  const selectedDeviceId = ref(localStorage.getItem(DEVICE_STORAGE_KEY) || '')

  const countdown = ref(0)
  const capturing = ref(false)
  const photos = ref<string[]>([])

  const isCameraReady = computed(() => Boolean(stream.value) && !cameraError.value)
  const photosLeft = computed(() => options.photoCount.value - photos.value.length)

  async function refreshDevices() {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices()
      videoDevices.value = devices.filter((d) => d.kind === 'videoinput')
    } catch {
      // Ignore — the picker just stays empty/hidden if this isn't supported.
    }
  }

  async function startCamera() {
    cameraError.value = ''
    photos.value = []
    try {
      stream.value = await navigator.mediaDevices.getUserMedia({
        video: selectedDeviceId.value
          ? { deviceId: { exact: selectedDeviceId.value }, width: { ideal: 1280 }, height: { ideal: 1280 } }
          : { facingMode: 'user', width: { ideal: 1280 }, height: { ideal: 1280 } },
        audio: false,
      })
      await nextTick()
      if (videoRef.value) {
        videoRef.value.srcObject = stream.value
        await videoRef.value.play()
      }

      // Labels are only populated once permission has been granted, so this
      // is the earliest point the picker can show real device names.
      await refreshDevices()
      if (!selectedDeviceId.value) {
        selectedDeviceId.value = stream.value.getVideoTracks()[0]?.getSettings().deviceId ?? ''
      }
    } catch {
      if (selectedDeviceId.value) {
        // Stored device (e.g. a webcam app that isn't running anymore) is
        // no longer available — fall back to the default camera instead of
        // getting stuck on an error the guest can't fix.
        selectedDeviceId.value = ''
        localStorage.removeItem(DEVICE_STORAGE_KEY)
        await startCamera()
        return
      }
      cameraError.value = t(errorMessageKey)
    }
  }

  function stopCamera() {
    stream.value?.getTracks().forEach((track) => track.stop())
    stream.value = null
  }

  async function switchCamera(deviceId: string) {
    if (!deviceId || deviceId === selectedDeviceId.value) return
    selectedDeviceId.value = deviceId
    localStorage.setItem(DEVICE_STORAGE_KEY, deviceId)
    stopCamera()
    await startCamera()
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
        ctx.filter = options.filterCss.value
        ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
        photos.value.push(canvas.toDataURL('image/jpeg', 0.92))
      }
    }

    capturing.value = false

    if (photos.value.length >= options.photoCount.value) {
      stopCamera()
      options.onAllCaptured?.(photos.value)
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

  onBeforeUnmount(stopCamera)

  return {
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
    stopCamera,
    switchCamera,
    startCaptureCountdown,
    retakePhotos,
  }
}
