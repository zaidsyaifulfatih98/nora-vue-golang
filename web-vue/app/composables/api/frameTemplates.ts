export type FrameTemplateFit = 'COVER' | 'CONTAIN'

export interface FrameTemplateItem {
  id: string
  name: string
  description: string
  nameEn: string | null
  descriptionEn: string | null
  imageUrl: string
  fit: FrameTemplateFit
  order: number
  isActive: boolean
}

export function useFrameTemplatesApi() {
  const axios = useAxios()

  const getFrameTemplates = (all = false) =>
    axios.get(`/frame-templates${all ? '?all=true' : ''}`).then((r) => r.data.data as FrameTemplateItem[])

  const uploadFrameTemplate = (
    file: File,
    name: string,
    description: string,
    fit: FrameTemplateFit,
    nameEn?: string,
    descriptionEn?: string,
  ) => {
    const formData = new FormData()
    formData.append('image', file)
    formData.append('name', name)
    formData.append('description', description)
    formData.append('fit', fit)
    if (nameEn) formData.append('nameEn', nameEn)
    if (descriptionEn) formData.append('descriptionEn', descriptionEn)
    return axios
      .post('/frame-templates', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as FrameTemplateItem)
  }

  const updateFrameTemplate = (
    id: string,
    payload: Partial<Pick<FrameTemplateItem, 'name' | 'description' | 'nameEn' | 'descriptionEn' | 'fit' | 'order' | 'isActive'>>,
    file?: File,
  ) => {
    if (!file) {
      return axios.patch(`/frame-templates/${id}`, payload).then((r) => r.data.data as FrameTemplateItem)
    }

    const formData = new FormData()
    Object.entries(payload).forEach(([key, value]) => {
      if (value !== undefined && value !== null) formData.append(key, String(value))
    })
    formData.append('image', file)

    return axios
      .patch(`/frame-templates/${id}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as FrameTemplateItem)
  }

  const deleteFrameTemplate = (id: string) => axios.delete(`/frame-templates/${id}`)

  return { getFrameTemplates, uploadFrameTemplate, updateFrameTemplate, deleteFrameTemplate }
}
