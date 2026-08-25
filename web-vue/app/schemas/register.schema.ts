import { z } from 'zod'

export const registerSchema = z.object({
  firstName: z.string().min(5, 'Nama depan minimal 5 karakter').max(25, 'Nama depan maksimal 25 karakter'),
  lastName: z.string().min(5, 'Nama belakang minimal 5 karakter').max(25, 'Nama belakang maksimal 25 karakter'),
  email: z.string().min(1, 'Email wajib diisi').email('Format email tidak valid'),
  password: z.string().min(6, 'Password minimal 6 karakter'),
  role: z.enum(['SUPER_ADMIN', 'ADMIN']),
})

export type RegisterFormData = z.infer<typeof registerSchema>
