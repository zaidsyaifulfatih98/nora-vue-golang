export interface BackdropItem {
  id: string
  name: string
  imageUrl: string
  order: number
  isActive: boolean
}

export function useBackdropsApi() {
  const axios = useAxios()

  const getBackdrops = (all = false) =>
    axios.get(`/backdrops${all ? '?all=true' : ''}`).then((r) => r.data.data as BackdropItem[])

  const uploadBackdrop = (file: File, name: string) => {
    const formData = new FormData()
    formData.append('image', file)
    formData.append('name', name)
    return axios
      .post('/backdrops', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as BackdropItem)
  }

  const updateBackdrop = (
    id: string,
    payload: Partial<Pick<BackdropItem, 'name' | 'order' | 'isActive'>>,
    file?: File,
  ) => {
    if (!file) {
      return axios.patch(`/backdrops/${id}`, payload).then((r) => r.data.data as BackdropItem)
    }

    const formData = new FormData()
    Object.entries(payload).forEach(([key, value]) => {
      if (value !== undefined) formData.append(key, String(value))
    })
    formData.append('image', file)

    return axios
      .patch(`/backdrops/${id}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as BackdropItem)
  }

  const deleteBackdrop = (id: string) => axios.delete(`/backdrops/${id}`)

  return { getBackdrops, uploadBackdrop, updateBackdrop, deleteBackdrop }
}
