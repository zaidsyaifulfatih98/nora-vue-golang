export interface GalleryPhotoItem {
  id: string
  url: string
  caption: string | null
  order: number
  isActive: boolean
}

export function useGalleryApi() {
  const axios = useAxios()

  const getGalleryPhotos = (all = false) =>
    axios.get(`/gallery${all ? '?all=true' : ''}`).then((r) => r.data.data as GalleryPhotoItem[])

  const uploadGalleryPhoto = (file: File, caption?: string) => {
    const formData = new FormData()
    formData.append('image', file)
    if (caption) formData.append('caption', caption)
    return axios
      .post('/gallery', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as GalleryPhotoItem)
  }

  const updateGalleryPhoto = (id: string, payload: Partial<Pick<GalleryPhotoItem, 'caption' | 'order' | 'isActive'>>) =>
    axios.patch(`/gallery/${id}`, payload).then((r) => r.data.data as GalleryPhotoItem)

  const deleteGalleryPhoto = (id: string) => axios.delete(`/gallery/${id}`)

  return { getGalleryPhotos, uploadGalleryPhoto, updateGalleryPhoto, deleteGalleryPhoto }
}
