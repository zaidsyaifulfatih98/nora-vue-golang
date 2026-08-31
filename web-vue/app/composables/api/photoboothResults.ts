export interface PhotoboothResultItem {
  id: string
  createdAt: string
  viewUrl: string
  downloadUrl: string
}

export function usePhotoboothResultsApi() {
  const axios = useAxios()

  const savePhotoboothResult = (file: Blob, name: string) => {
    const formData = new FormData()
    formData.append('image', file, name)
    return axios
      .post('/photobooth-results', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      .then((r) => r.data.data as PhotoboothResultItem)
  }

  const getPhotoboothResults = () => axios.get('/photobooth-results').then((r) => r.data.data as PhotoboothResultItem[])

  const deletePhotoboothResult = (id: string) => axios.delete(`/photobooth-results/${id}`)

  return { savePhotoboothResult, getPhotoboothResults, deletePhotoboothResult }
}
