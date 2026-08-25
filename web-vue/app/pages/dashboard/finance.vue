<script setup lang="ts">
import type { FinanceEntryItem, FinanceSummary } from '~/composables/api/finance'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getFinanceEntries, getFinanceSummary, createFinanceEntry, updateFinanceEntry, deleteFinanceEntry } =
  useFinanceApi()

function formatRupiah(value: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
}

const initialForm = {
  type: 'INCOME' as 'INCOME' | 'EXPENSE',
  category: '',
  amount: '',
  date: new Date().toISOString().slice(0, 10),
  description: '',
}

const entries = ref<FinanceEntryItem[]>([])
const summary = ref<FinanceSummary | null>(null)
const loading = ref(true)
const submitting = ref(false)
const form = reactive({ ...initialForm })
const editingId = ref<string | null>(null)
const search = ref('')
const page = ref(1)
const pageSize = 10

async function loadData() {
  loading.value = true
  const [entriesData, summaryData] = await Promise.all([getFinanceEntries(), getFinanceSummary()])
  entries.value = entriesData
  summary.value = summaryData
  loading.value = false
}
onMounted(loadData)

async function handleSubmit() {
  submitting.value = true
  try {
    const payload = {
      type: form.type,
      category: form.category,
      amount: Number(form.amount),
      date: new Date(form.date).toISOString(),
      description: form.description || undefined,
    }
    if (editingId.value) {
      await updateFinanceEntry(editingId.value, payload)
      editingId.value = null
    } else {
      await createFinanceEntry(payload)
    }
    Object.assign(form, initialForm)
    await loadData()
  } finally {
    submitting.value = false
  }
}

function handleEdit(entry: FinanceEntryItem) {
  editingId.value = entry.id
  Object.assign(form, {
    type: entry.type,
    category: entry.category,
    amount: String(entry.amount),
    date: entry.date.slice(0, 10),
    description: entry.description ?? '',
  })
}

function handleCancelEdit() {
  editingId.value = null
  Object.assign(form, initialForm)
}

async function handleDelete(id: string) {
  await deleteFinanceEntry(id)
  if (editingId.value === id) handleCancelEdit()
  await loadData()
}

const filteredEntries = computed(() => {
  const keyword = search.value.trim().toLowerCase()
  if (!keyword) return entries.value
  return entries.value.filter((entry) => {
    const createdByName = `${entry.createdBy?.firstName ?? ''} ${entry.createdBy?.lastName ?? ''}`
    return (
      entry.category.toLowerCase().includes(keyword) ||
      (entry.description ?? '').toLowerCase().includes(keyword) ||
      createdByName.toLowerCase().includes(keyword)
    )
  })
})

const totalPages = computed(() => Math.max(1, Math.ceil(filteredEntries.value.length / pageSize)))
const currentPage = computed(() => Math.min(page.value, totalPages.value))
const paginatedEntries = computed(() =>
  filteredEntries.value.slice((currentPage.value - 1) * pageSize, currentPage.value * pageSize),
)

watch(search, () => {
  page.value = 1
})

function handleExportCsv() {
  const header = ['Tanggal', 'Kategori', 'Tipe', 'Nominal', 'Catatan', 'Dicatat Oleh']
  const escapeCsv = (value: string) => `"${value.replace(/"/g, '""')}"`
  const rows = filteredEntries.value.map((entry) => [
    new Date(entry.date).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }),
    entry.category,
    entry.type === 'INCOME' ? 'Pemasukan' : 'Pengeluaran',
    Number(entry.amount).toString(),
    entry.description ?? '',
    `${entry.createdBy?.firstName ?? ''} ${entry.createdBy?.lastName ?? ''}`.trim(),
  ])

  const csvContent = [header, ...rows].map((row) => row.map((cell) => escapeCsv(String(cell))).join(',')).join('\n')

  const blob = new Blob([`﻿${csvContent}`], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `riwayat-transaksi-${new Date().toISOString().slice(0, 10)}.csv`
  link.click()
  URL.revokeObjectURL(url)
}
</script>

<template>
  <div class="space-y-6">
    <div class="grid gap-6 sm:grid-cols-3">
      <div class="rounded-2xl bg-white p-5 shadow-md">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-green-100">
          <Icon name="fe:dollar" class="text-lg text-green-600" />
        </div>
        <p class="mt-2 text-sm text-gray-500">Total Pemasukan</p>
        <p class="text-xl font-bold text-gray-900">{{ formatRupiah(summary?.income ?? 0) }}</p>
      </div>
      <div class="rounded-2xl bg-white p-5 shadow-md">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-red-100">
          <Icon name="fe:trending-down" class="text-lg text-red-500" />
        </div>
        <p class="mt-2 text-sm text-gray-500">Total Pengeluaran</p>
        <p class="text-xl font-bold text-gray-900">{{ formatRupiah(summary?.expense ?? 0) }}</p>
      </div>
      <div class="rounded-2xl bg-white p-5 shadow-md">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#F1E4D6]">
          <Icon name="fe:dollar" class="text-lg text-[#8F6F3E]" />
        </div>
        <p class="mt-2 text-sm text-gray-500">Saldo</p>
        <p class="text-xl font-bold text-gray-900">{{ formatRupiah(summary?.balance ?? 0) }}</p>
      </div>
    </div>

    <div class="rounded-2xl bg-white p-6 shadow-md">
      <h2 class="text-base font-semibold text-gray-800">{{ editingId ? 'Edit Transaksi' : 'Catat Transaksi Baru' }}</h2>
      <form class="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-5" @submit.prevent="handleSubmit">
        <select v-model="form.type" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500">
          <option value="INCOME">Pemasukan</option>
          <option value="EXPENSE">Pengeluaran</option>
        </select>

        <input v-model="form.category" required placeholder="Kategori (mis. Booking Paket Gold)" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:col-span-2" />
        <input v-model="form.amount" required type="number" min="0" placeholder="Nominal (Rp)" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500" />
        <input v-model="form.date" required type="date" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500" />
        <input v-model="form.description" placeholder="Catatan (opsional)" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:col-span-2 lg:col-span-4" />

        <div class="flex items-center gap-2">
          <button type="submit" :disabled="submitting" class="flex items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700 disabled:opacity-60">
            <Icon :name="editingId ? 'fe:pencil' : 'fe:plus'" />
            {{ submitting ? 'Menyimpan...' : editingId ? 'Update' : 'Simpan' }}
          </button>
          <button v-if="editingId" type="button" class="flex items-center justify-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-semibold text-gray-600 transition hover:bg-gray-50" @click="handleCancelEdit">
            <Icon name="fe:close" /> Batal
          </button>
        </div>
      </form>
    </div>

    <div class="rounded-2xl bg-white shadow-md">
      <div class="flex flex-col gap-3 border-b border-gray-100 px-6 py-4 sm:flex-row sm:items-center sm:justify-between">
        <h2 class="text-base font-semibold text-gray-800">Riwayat Transaksi</h2>
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <div class="relative w-full sm:w-64">
            <Icon name="fe:search" class="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-gray-400" />
            <input v-model="search" type="text" placeholder="Cari kategori, catatan, atau pencatat..." class="w-full rounded-lg border border-gray-200 py-2 pr-3 pl-9 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500" />
          </div>
          <button type="button" :disabled="filteredEntries.length === 0" class="flex items-center justify-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-semibold text-gray-600 transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50" @click="handleExportCsv">
            <Icon name="fe:download" /> Export CSV
          </button>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-gray-100 text-gray-500">
              <th class="px-6 py-3 font-medium">Tanggal</th>
              <th class="px-6 py-3 font-medium">Kategori</th>
              <th class="px-6 py-3 font-medium">Tipe</th>
              <th class="px-6 py-3 font-medium">Nominal</th>
              <th class="px-6 py-3 font-medium">Dicatat oleh</th>
              <th class="px-6 py-3 font-medium"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="6" class="px-6 py-8 text-center text-gray-400">Memuat data...</td>
            </tr>
            <tr v-else-if="filteredEntries.length === 0">
              <td colspan="6" class="px-6 py-8 text-center text-gray-400">
                {{ search ? 'Tidak ada transaksi yang cocok dengan pencarian.' : 'Belum ada transaksi tercatat.' }}
              </td>
            </tr>
            <tr v-for="entry in paginatedEntries" v-else :key="entry.id" class="border-b border-gray-50">
              <td class="px-6 py-3 text-gray-700">
                {{ new Date(entry.date).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }) }}
              </td>
              <td class="px-6 py-3 text-gray-700">
                {{ entry.category }}
                <p v-if="entry.description" class="text-xs text-gray-400">{{ entry.description }}</p>
              </td>
              <td class="px-6 py-3">
                <span class="rounded-full px-2.5 py-1 text-xs font-semibold" :class="entry.type === 'INCOME' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-600'">
                  {{ entry.type === 'INCOME' ? 'Pemasukan' : 'Pengeluaran' }}
                </span>
              </td>
              <td class="px-6 py-3 font-semibold" :class="entry.type === 'INCOME' ? 'text-green-600' : 'text-red-500'">
                {{ entry.type === 'INCOME' ? '+' : '-' }}{{ formatRupiah(Number(entry.amount)) }}
              </td>
              <td class="px-6 py-3 text-gray-500">{{ entry.createdBy?.firstName }} {{ entry.createdBy?.lastName }}</td>
              <td class="px-6 py-3 text-right">
                <div class="flex items-center justify-end gap-3">
                  <button aria-label="Edit" class="text-gray-400 hover:text-blue-600" @click="handleEdit(entry)"><Icon name="fe:pencil" /></button>
                  <button aria-label="Hapus" class="text-gray-400 hover:text-red-500" @click="handleDelete(entry.id)"><Icon name="fe:trash-2" /></button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="!loading && filteredEntries.length > 0" class="flex flex-col items-center justify-between gap-3 border-t border-gray-100 px-6 py-4 sm:flex-row">
        <p class="text-xs text-gray-500">
          Menampilkan {{ (currentPage - 1) * pageSize + 1 }}-{{ Math.min(currentPage * pageSize, filteredEntries.length) }}
          dari {{ filteredEntries.length }} transaksi
        </p>
        <div class="flex items-center gap-2">
          <button
            type="button"
            :disabled="currentPage === 1"
            aria-label="Halaman sebelumnya"
            class="flex h-8 w-8 items-center justify-center rounded-lg border border-gray-200 text-gray-500 transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40"
            @click="page = Math.max(1, currentPage - 1)"
          >
            <Icon name="fe:arrow-left" />
          </button>
          <span class="text-xs font-medium text-gray-600">Halaman {{ currentPage }} dari {{ totalPages }}</span>
          <button
            type="button"
            :disabled="currentPage === totalPages"
            aria-label="Halaman berikutnya"
            class="flex h-8 w-8 items-center justify-center rounded-lg border border-gray-200 text-gray-500 transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40"
            @click="page = Math.min(totalPages, currentPage + 1)"
          >
            <Icon name="fe:arrow-right" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
