export interface VoiceMessageItem {
  id: string
  createdAt: string
  guestName: string
  audioUrl: string
  photoUrl: string
}

export function useVoiceMessagesApi() {
  const axios = useAxios()

  const getVoiceMessages = () => axios.get('/voice-messages').then((r) => r.data.data as VoiceMessageItem[])

  const uploadVoiceMessage = (file: Blob, guestName: string, photoUrl?: string) => {
    const extension = file.type.includes('mp4') ? 'mp4' : file.type.includes('ogg') ? 'ogg' : 'webm'
    const formData = new FormData()
    formData.append('audio', file, `voice-message.${extension}`)
    if (guestName) formData.append('guestName', guestName)
    if (photoUrl) formData.append('photoUrl', photoUrl)
    return axios
      .post('/voice-messages', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as VoiceMessageItem)
  }

  const deleteVoiceMessage = (id: string) => axios.delete(`/voice-messages/${id}`)

  return { getVoiceMessages, uploadVoiceMessage, deleteVoiceMessage }
}
