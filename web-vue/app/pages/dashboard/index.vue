<script setup lang="ts">
import type { FinanceSummary } from '~/composables/api/finance'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getFinanceSummary } = useFinanceApi()
const { getPackages } = usePackagesApi()
const { getGalleryPhotos } = useGalleryApi()
const { getReviews } = useReviewsApi()

function formatRupiah(value: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
}

const summary = ref<FinanceSummary | null>(null)
const counts = reactive({ packages: 0, photos: 0, reviews: 0 })
const loading = ref(true)

onMounted(async () => {
  try {
    const [financeSummary, packages, photos, reviews] = await Promise.all([
      getFinanceSummary(),
      getPackages(true),
      getGalleryPhotos(true),
      getReviews(true),
    ])
    summary.value = financeSummary
    counts.packages = packages.length
    counts.photos = photos.length
    counts.reviews = reviews.length
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <div class="rounded-2xl border-gray-100 bg-white p-5 shadow-md">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-green-100">
          <Icon name="fe:dollar" class="text-lg text-green-600" />
        </div>
        <p class="mt-2 text-sm text-gray-500">Total Pemasukan</p>
        <p class="text-2xl font-bold text-gray-900">{{ loading ? '...' : formatRupiah(summary?.income ?? 0) }}</p>
      </div>

      <div class="rounded-2xl border-gray-100 bg-white p-5 shadow-md">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-red-100">
          <Icon name="fe:trending-down" class="text-lg text-red-500" />
        </div>
        <p class="mt-2 text-sm text-gray-500">Total Pengeluaran</p>
        <p class="text-2xl font-bold text-gray-900">{{ loading ? '...' : formatRupiah(summary?.expense ?? 0) }}</p>
      </div>

      <div class="rounded-2xl border-gray-100 bg-white p-5 shadow-md">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#E9EAF0]">
          <Icon name="fe:dollar" class="text-lg text-[#1E2537]" />
        </div>
        <p class="mt-2 text-sm text-gray-500">Saldo</p>
        <p class="text-2xl font-bold text-gray-900">{{ loading ? '...' : formatRupiah(summary?.balance ?? 0) }}</p>
      </div>
    </div>

    <div class="mt-6 grid gap-6 sm:grid-cols-3">
      <NuxtLink to="/dashboard/packages" class="flex items-center gap-4 rounded-2xl border-gray-100 bg-white p-5 shadow-md transition hover:-translate-y-0.5">
        <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-[#F5E1E1]">
          <Icon name="fe:package" class="text-lg text-[#920f0f]" />
        </div>
        <div>
          <p class="text-sm text-gray-500">Paket Aktif</p>
          <p class="text-xl font-bold text-gray-900">{{ loading ? '...' : counts.packages }}</p>
        </div>
      </NuxtLink>

      <NuxtLink to="/dashboard/gallery" class="flex items-center gap-4 rounded-2xl border-gray-100 bg-white p-5 shadow-md transition hover:-translate-y-0.5">
        <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-purple-100">
          <Icon name="fe:image" class="text-lg text-purple-600" />
        </div>
        <div>
          <p class="text-sm text-gray-500">Foto Galeri</p>
          <p class="text-xl font-bold text-gray-900">{{ loading ? '...' : counts.photos }}</p>
        </div>
      </NuxtLink>

      <NuxtLink to="/dashboard/reviews" class="flex items-center gap-4 rounded-2xl border-gray-100 bg-white p-5 shadow-md transition hover:-translate-y-0.5">
        <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-amber-100">
          <Icon name="fe:star" class="text-lg text-amber-600" />
        </div>
        <div>
          <p class="text-sm text-gray-500">Review</p>
          <p class="text-xl font-bold text-gray-900">{{ loading ? '...' : counts.reviews }}</p>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>
