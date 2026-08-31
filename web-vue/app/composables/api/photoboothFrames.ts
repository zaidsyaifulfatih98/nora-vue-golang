export interface FrameSlot {
  x: number
  y: number
  width: number
  height: number
}

export interface PhotoboothFrameItem {
  id: string
  name: string
  imageUrl: string
  slots: FrameSlot[]
  order: number
  isActive: boolean
}

export function usePhotoboothFramesApi() {
  const axios = useAxios()

  const getPhotoboothFrames = (all = false) =>
    axios.get(`/photobooth-frames${all ? '?all=true' : ''}`).then((r) => r.data.data as PhotoboothFrameItem[])

  const uploadPhotoboothFrame = (file: File, name: string, slots: FrameSlot[]) => {
    const formData = new FormData()
    formData.append('image', file)
    formData.append('name', name)
    formData.append('slots', JSON.stringify(slots))
    return axios
      .post('/photobooth-frames', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as PhotoboothFrameItem)
  }

  const updatePhotoboothFrame = (
    id: string,
    payload: Partial<Pick<PhotoboothFrameItem, 'name' | 'order' | 'isActive'>> & { slots?: FrameSlot[] },
    file?: File,
  ) => {
    if (!file) {
      return axios.patch(`/photobooth-frames/${id}`, payload).then((r) => r.data.data as PhotoboothFrameItem)
    }

    const formData = new FormData()
    Object.entries(payload).forEach(([key, value]) => {
      if (value === undefined) return
      formData.append(key, key === 'slots' ? JSON.stringify(value) : String(value))
    })
    formData.append('image', file)

    return axios
      .patch(`/photobooth-frames/${id}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as PhotoboothFrameItem)
  }

  const deletePhotoboothFrame = (id: string) => axios.delete(`/photobooth-frames/${id}`)

  return { getPhotoboothFrames, uploadPhotoboothFrame, updatePhotoboothFrame, deletePhotoboothFrame }
}
