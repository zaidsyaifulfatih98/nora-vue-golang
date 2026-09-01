export interface PackageItem {
  id: string
  name: string
  price: string
  duration: string
  description: string
  features: string[]
  nameEn: string | null
  durationEn: string | null
  descriptionEn: string | null
  featuresEn: string[] | null
  isPopular: boolean
  isActive: boolean
  order: number
}

export interface PackagePayload {
  name: string
  price: number
  duration: string
  description: string
  features: string[]
  nameEn?: string
  durationEn?: string
  descriptionEn?: string
  featuresEn?: string[]
  isPopular?: boolean
  isActive?: boolean
  order?: number
}

export function usePackagesApi() {
  const axios = useAxios()

  const getPackages = (all = false) =>
    axios.get(`/packages${all ? '?all=true' : ''}`).then((r) => r.data.data as PackageItem[])

  const createPackage = (payload: PackagePayload) =>
    axios.post('/packages', payload).then((r) => r.data.data as PackageItem)

  const updatePackage = (id: string, payload: Partial<PackagePayload>) =>
    axios.patch(`/packages/${id}`, payload).then((r) => r.data.data as PackageItem)

  const deletePackage = (id: string) => axios.delete(`/packages/${id}`)

  return { getPackages, createPackage, updatePackage, deletePackage }
}
