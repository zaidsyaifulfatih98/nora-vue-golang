import { z } from 'zod'

export const loginSchema = z.object({
  email: z.string().min(1, 'Email wajib diisi').email('Format email tidak valid'),
  password: z
    .string()
    .min(5, 'Password minimal 5 karakter')
    .max(25, 'Password maksimal 25 karakter'),
})

export type LoginFormData = z.infer<typeof loginSchema>
