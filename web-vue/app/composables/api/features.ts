export interface FeatureItem {
  id: string
  title: string
  description: string
  imageUrl: string
  order: number
  isActive: boolean
}

export function useFeaturesApi() {
  const axios = useAxios()

  const getFeatures = (all = false) =>
    axios.get(`/features${all ? '?all=true' : ''}`).then((r) => r.data.data as FeatureItem[])

  const uploadFeature = (file: File, title: string, description: string) => {
    const formData = new FormData()
    formData.append('image', file)
    formData.append('title', title)
    formData.append('description', description)
    return axios
      .post('/features', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as FeatureItem)
  }

  const updateFeature = (
    id: string,
    payload: Partial<Pick<FeatureItem, 'title' | 'description' | 'order' | 'isActive'>>,
    file?: File,
  ) => {
    if (!file) {
      return axios.patch(`/features/${id}`, payload).then((r) => r.data.data as FeatureItem)
    }

    const formData = new FormData()
    Object.entries(payload).forEach(([key, value]) => {
      if (value !== undefined) formData.append(key, String(value))
    })
    formData.append('image', file)

    return axios
      .patch(`/features/${id}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as FeatureItem)
  }

  const deleteFeature = (id: string) => axios.delete(`/features/${id}`)

  return { getFeatures, uploadFeature, updateFeature, deleteFeature }
}
