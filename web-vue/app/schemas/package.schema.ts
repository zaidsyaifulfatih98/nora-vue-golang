import { z } from 'zod'

export const packageSchema = z.object({
  name: z.string().min(1, 'Nama paket wajib diisi').max(50, 'Nama paket maksimal 50 karakter'),
  price: z.coerce.number().min(0, 'Harga tidak boleh negatif'),
  duration: z.string().min(1, 'Durasi wajib diisi'),
  description: z.string().min(1, 'Deskripsi wajib diisi'),
  features: z.array(z.string()).default([]),
  isPopular: z.boolean().default(false),
  isActive: z.boolean().default(true),
  order: z.coerce.number().int().default(0),
})

export type PackageFormData = z.infer<typeof packageSchema>
