export interface PhotoboothResultItem {
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

  return { savePhotoboothResult }
}
