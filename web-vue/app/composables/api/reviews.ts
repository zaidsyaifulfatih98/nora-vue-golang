export interface ReviewItem {
  id: string
  name: string
  eventLabel: string
  quote: string
  eventLabelEn: string | null
  quoteEn: string | null
  rating: number
  isPublished: boolean
  order: number
}

export interface ReviewPayload {
  name: string
  eventLabel: string
  quote: string
  eventLabelEn?: string
  quoteEn?: string
  rating?: number
  isPublished?: boolean
  order?: number
}

export function useReviewsApi() {
  const axios = useAxios()

  const getReviews = (all = false) =>
    axios.get(`/reviews${all ? '?all=true' : ''}`).then((r) => r.data.data as ReviewItem[])

  const createReview = (payload: ReviewPayload) =>
    axios.post('/reviews', payload).then((r) => r.data.data as ReviewItem)

  const updateReview = (id: string, payload: Partial<ReviewPayload>) =>
    axios.patch(`/reviews/${id}`, payload).then((r) => r.data.data as ReviewItem)

  const deleteReview = (id: string) => axios.delete(`/reviews/${id}`)

  return { getReviews, createReview, updateReview, deleteReview }
}
