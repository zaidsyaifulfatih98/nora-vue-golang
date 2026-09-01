<script setup lang="ts">
import type { ReviewItem } from '~/composables/api/reviews'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getReviews, createReview, updateReview, deleteReview } = useReviewsApi()
const { t } = useI18n()

const reviews = ref<ReviewItem[]>([])
const loading = ref(true)
const showForm = ref(false)
const editingId = ref<string | null>(null)
const submitting = ref(false)

const emptyForm = { name: '', eventLabel: '', quote: '', eventLabelEn: '', quoteEn: '', rating: 5, isPublished: true }
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
    eventLabelEn: review.eventLabelEn ?? '',
    quoteEn: review.quoteEn ?? '',
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
      <p class="text-sm text-gray-500">{{ t('dashboard.reviews.subtitle') }}</p>
      <button
        class="flex items-center gap-2 rounded-lg bg-[#920f0f] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#7a0c0c]"
        @click="openCreateForm"
      >
        <Icon name="fe:plus" /> {{ t('dashboard.reviews.addReview') }}
      </button>
    </div>

    <div v-if="showForm" class="rounded-2xl bg-white p-6 shadow-md">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-semibold text-gray-800">{{ editingId ? t('dashboard.reviews.editReview') : t('dashboard.reviews.newReview') }}</h2>
        <button class="text-gray-400 hover:text-gray-600" @click="showForm = false"><Icon name="fe:close" /></button>
      </div>

      <form class="mt-4 grid gap-4 sm:grid-cols-2" @submit.prevent="handleSubmit">
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.reviews.coupleNameLabel') }}</label>
          <input v-model="form.name" required :placeholder="t('dashboard.reviews.coupleNamePlaceholder')" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
        </div>
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.reviews.eventInfoLabel') }}</label>
          <input v-model="form.eventLabel" required :placeholder="t('dashboard.reviews.eventInfoPlaceholder')" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
        </div>
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.reviews.ratingLabel') }}</label>
          <select v-model.number="form.rating" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]">
            <option v-for="r in [5, 4, 3, 2, 1]" :key="r" :value="r">{{ r }} {{ t('dashboard.reviews.star') }}</option>
          </select>
        </div>
        <label class="flex items-center gap-2 pt-6 text-sm text-gray-600">
          <input v-model="form.isPublished" type="checkbox" /> {{ t('dashboard.reviews.showOnLanding') }}
        </label>
        <div class="flex flex-col gap-1.5 sm:col-span-2">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.reviews.quoteLabel') }}</label>
          <textarea v-model="form.quote" required rows="3" class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
        </div>

        <div class="rounded-lg border border-dashed border-gray-200 p-3 sm:col-span-2">
          <p class="text-xs font-semibold text-gray-500">{{ t('dashboard.translationLabel') }}</p>
          <p class="mt-0.5 text-xs text-gray-400">{{ t('dashboard.translationHint') }}</p>
          <div class="mt-2 flex flex-col gap-3">
            <input v-model="form.eventLabelEn" :placeholder="t('dashboard.reviews.eventInfoEnPlaceholder')" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium text-gray-600">{{ t('dashboard.reviews.quoteEnLabel') }}</label>
              <textarea v-model="form.quoteEn" rows="3" class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
            </div>
          </div>
        </div>

        <button type="submit" :disabled="submitting" class="rounded-lg bg-[#920f0f] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#7a0c0c] disabled:opacity-60 sm:col-span-2">
          {{ submitting ? t('dashboard.reviews.saving') : t('dashboard.reviews.saveReview') }}
        </button>
      </form>
    </div>

    <div class="grid gap-6 lg:grid-cols-3">
      <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.reviews.loadingReviews') }}</p>
      <p v-else-if="reviews.length === 0" class="text-sm text-gray-400">{{ t('dashboard.reviews.empty') }}</p>

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
            <input type="checkbox" :checked="review.isPublished" @change="handleTogglePublished(review)" /> {{ t('dashboard.reviews.show') }}
          </label>
          <div class="flex items-center gap-3">
            <button class="text-gray-400 hover:text-[#920f0f]" :aria-label="t('dashboard.reviews.editAria')" @click="openEditForm(review)"><Icon name="fe:pencil" /></button>
            <button class="text-gray-400 hover:text-red-500" :aria-label="t('dashboard.reviews.deleteAria')" @click="handleDelete(review.id)"><Icon name="lucide:trash-2" /></button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
