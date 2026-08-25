import { z } from 'zod'

export const reviewSchema = z.object({
  name: z.string().min(1, 'Nama wajib diisi'),
  eventLabel: z.string().min(1, 'Label acara wajib diisi'),
  quote: z.string().min(1, 'Testimoni wajib diisi'),
  rating: z.coerce.number().int().min(1).max(5).default(5),
  isPublished: z.boolean().default(true),
  order: z.coerce.number().int().default(0),
})

export type ReviewFormData = z.infer<typeof reviewSchema>
