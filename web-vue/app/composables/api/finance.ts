export type FinanceEntryType = 'INCOME' | 'EXPENSE'

export interface FinanceEntryItem {
  id: string
  type: FinanceEntryType
  category: string
  amount: number
  description: string | null
  date: string
  createdById: string
  createdBy?: { firstName: string; lastName: string }
}

export interface FinanceEntryPayload {
  type: FinanceEntryType
  category: string
  amount: number
  description?: string
  date: string
}

export interface FinanceSummary {
  income: number
  expense: number
  balance: number
}

export function useFinanceApi() {
  const axios = useAxios()

  const getFinanceEntries = (params?: { from?: string; to?: string }) =>
    axios.get('/finance', { params }).then((r) => r.data.data as FinanceEntryItem[])

  const getFinanceSummary = (params?: { from?: string; to?: string }) =>
    axios.get('/finance/summary', { params }).then((r) => r.data.data as FinanceSummary)

  const createFinanceEntry = (payload: FinanceEntryPayload) =>
    axios.post('/finance', payload).then((r) => r.data.data as FinanceEntryItem)

  const updateFinanceEntry = (id: string, payload: Partial<FinanceEntryPayload>) =>
    axios.patch(`/finance/${id}`, payload).then((r) => r.data.data as FinanceEntryItem)

  const deleteFinanceEntry = (id: string) => axios.delete(`/finance/${id}`)

  return { getFinanceEntries, getFinanceSummary, createFinanceEntry, updateFinanceEntry, deleteFinanceEntry }
}
