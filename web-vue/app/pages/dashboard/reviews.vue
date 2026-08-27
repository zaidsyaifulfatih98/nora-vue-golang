<script setup lang="ts">
import type { ReviewItem } from '~/composables/api/reviews'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getReviews, createReview, updateReview, deleteReview } = useReviewsApi()

const reviews = ref<ReviewItem[]>([])
const loading = ref(true)
const showForm = ref(false)
const editingId = ref<string | null>(null)
const submitting = ref(false)

const emptyForm = { name: '', eventLabel: '', quote: '', rating: 5, isPublished: true }
const form = reactive({ ...emptyForm })

async function loadReviews() {
  loading.value = true
  reviews.value = await getReviews(true)
  loading.value = false
}
onMounted(loadReviews)

function openCreateForm() {
  editingId.value = null
  Object.assign(form, emptyForm)
  showForm.value = true
}

function openEditForm(review: ReviewItem) {
  editingId.value = review.id
  Object.assign(form, {
    name: review.name,
    eventLabel: review.eventLabel,
    quote: review.quote,
    rating: review.rating,
    isPublished: review.isPublished,
  })
  showForm.value = true
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (editingId.value) await updateReview(editingId.value, form)
    else await createReview(form)
    showForm.value = false
    await loadReviews()
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: string) {
  await deleteReview(id)
  await loadReviews()
}

async function handleTogglePublished(review: ReviewItem) {
  await updateReview(review.id, { isPublished: !review.isPublished })
  await loadReviews()
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <p class="text-sm text-gray-500">Kelola testimoni pasangan yang tampil di landing page.</p>
      <button
        class="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700"
        @click="openCreateForm"
      >
        <Icon name="fe:plus" /> Tambah Review
      </button>
    </div>

    <div v-if="showForm" class="rounded-2xl bg-white p-6 shadow-md">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-semibold text-gray-800">{{ editingId ? 'Edit Review' : 'Tambah Review Baru' }}</h2>
        <button class="text-gray-400 hover:text-gray-600" @click="showForm = false"><Icon name="fe:close" /></button>
      </div>

      <form class="mt-4 grid gap-4 sm:grid-cols-2" @submit.prevent="handleSubmit">
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">Nama Pasangan</label>
          <input v-model="form.name" required placeholder="mis. Sarah & Bima" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500" />
        </div>
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">Info Acara</label>
          <input v-model="form.eventLabel" required placeholder="mis. Resepsi, Mei 2026" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500" />
        </div>
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">Rating</label>
          <select v-model.number="form.rating" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500">
            <option v-for="r in [5, 4, 3, 2, 1]" :key="r" :value="r">{{ r }} bintang</option>
          </select>
        </div>
        <label class="flex items-center gap-2 pt-6 text-sm text-gray-600">
          <input v-model="form.isPublished" type="checkbox" /> Tampilkan di landing page
        </label>
        <div class="flex flex-col gap-1.5 sm:col-span-2">
          <label class="text-sm font-medium text-gray-700">Testimoni</label>
          <textarea v-model="form.quote" required rows="3" class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500" />
        </div>
        <button type="submit" :disabled="submitting" class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700 disabled:opacity-60 sm:col-span-2">
          {{ submitting ? 'Menyimpan...' : 'Simpan Review' }}
        </button>
      </form>
    </div>

    <div class="grid gap-6 lg:grid-cols-3">
      <p v-if="loading" class="text-sm text-gray-400">Memuat review...</p>
      <p v-else-if="reviews.length === 0" class="text-sm text-gray-400">Belum ada review.</p>

      <div v-for="review in reviews" :key="review.id" class="flex flex-col rounded-2xl bg-white p-6 shadow-md">
        <div class="flex items-center gap-1 text-amber-500">
          <Icon v-for="i in review.rating" :key="i" name="fe:star" class="text-sm" />
        </div>
        <p class="mt-3 flex-1 text-sm text-gray-600 italic">&ldquo;{{ review.quote }}&rdquo;</p>
        <div class="mt-4 border-t border-gray-100 pt-3">
          <p class="text-sm font-bold text-gray-900">{{ review.name }}</p>
          <p class="text-xs text-gray-400">{{ review.eventLabel }}</p>
        </div>
        <div class="mt-4 flex items-center justify-between border-t border-gray-100 pt-3">
          <label class="flex items-center gap-2 text-xs text-gray-500">
            <input type="checkbox" :checked="review.isPublished" @change="handleTogglePublished(review)" /> Tampilkan
          </label>
          <div class="flex items-center gap-3">
            <button class="text-gray-400 hover:text-blue-600" aria-label="Edit" @click="openEditForm(review)"><Icon name="fe:pencil" /></button>
            <button class="text-gray-400 hover:text-red-500" aria-label="Hapus" @click="handleDelete(review.id)"><Icon name="lucide:trash-2" /></button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
