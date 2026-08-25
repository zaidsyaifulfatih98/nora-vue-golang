import { z } from 'zod'

export const financeSchema = z.object({
  type: z.enum(['INCOME', 'EXPENSE']),
  category: z.string().min(1, 'Kategori wajib diisi'),
  amount: z.coerce.number().min(0, 'Nominal tidak boleh negatif'),
  description: z.string().optional(),
  date: z.string().min(1, 'Tanggal wajib diisi'),
})

export type FinanceFormData = z.infer<typeof financeSchema>
