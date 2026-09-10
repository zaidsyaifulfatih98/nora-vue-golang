export type VoiceRecorderStep = 'idle' | 'form' | 'recording' | 'recorded' | 'sending' | 'sent'

export interface UseVoiceRecorderOptions {
  // Read lazily at send time, same as the original inline implementation —
  // savedResult/ownerSlug may not be settled yet when the composable is
  // created.
  getPhotoUrl?: () => string | undefined
  getOwnerSlug?: () => string | undefined
}

// Guest voice-greeting recording + upload, for the digital photobooth "send
// a voice message" step.
export function useVoiceRecorder(options: UseVoiceRecorderOptions = {}) {
  const { t } = useI18n()
  const { uploadVoiceMessage } = useVoiceMessagesApi()

  const voiceStep = ref<VoiceRecorderStep>('idle')
  const voiceGuestName = ref('')
  const voiceError = ref('')
  const voiceStream = ref<MediaStream | null>(null)
  const mediaRecorder = ref<MediaRecorder | null>(null)
  const voiceChunks = ref<Blob[]>([])
  const voiceBlob = ref<Blob | null>(null)
  const voiceAudioUrl = ref('')
  const voiceSeconds = ref(0)
  let voiceTimer: ReturnType<typeof setInterval> | null = null

  function openVoiceForm() {
    voiceError.value = ''
    voiceStep.value = 'form'
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
      voiceError.value = t('tryModal.voice.error')
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
      await uploadVoiceMessage(voiceBlob.value, voiceGuestName.value, options.getPhotoUrl?.(), options.getOwnerSlug?.())
      voiceStep.value = 'sent'
    } catch {
      voiceError.value = t('tryModal.voice.sendError')
      voiceStep.value = 'recorded'
    }
  }

  onBeforeUnmount(() => {
    stopVoiceStream()
    if (voiceAudioUrl.value) URL.revokeObjectURL(voiceAudioUrl.value)
  })

  return {
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
  }
}
