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
  ownerId?: string
}

export interface PhotoboothFramesBySlug {
  owner: { firstName: string; lastName: string }
  frames: PhotoboothFrameItem[]
}

export function usePhotoboothFramesApi() {
  const axios = useAxios()

  const getPhotoboothFrames = (all = false) =>
    axios.get(`/photobooth-frames${all ? '?all=true' : ''}`).then((r) => r.data.data as PhotoboothFrameItem[])

  // A DIGITAL_PHOTOBOOTH customer's own frames, for their dashboard.
  const getMyPhotoboothFrames = () => axios.get('/photobooth-frames/mine').then((r) => r.data.data as PhotoboothFrameItem[])

  const getPhotoboothFramesBySlug = (slug: string) =>
    axios.get(`/photobooth-frames/by-slug/${slug}`).then((r) => r.data.data as PhotoboothFramesBySlug)

  const uploadPhotoboothFrame = (file: File, name: string, slots: FrameSlot[], ownerId?: string) => {
    const formData = new FormData()
    formData.append('image', file)
    formData.append('name', name)
    formData.append('slots', JSON.stringify(slots))
    if (ownerId) formData.append('ownerId', ownerId)
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

  return {
    getPhotoboothFrames,
    getMyPhotoboothFrames,
    getPhotoboothFramesBySlug,
    uploadPhotoboothFrame,
    updatePhotoboothFrame,
    deletePhotoboothFrame,
  }
}
