<script setup lang="ts">
import type { FeatureItem } from '~/composables/api/features'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getFeatures, uploadFeature, updateFeature, deleteFeature } = useFeaturesApi()
const { t } = useI18n()

const features = ref<FeatureItem[]>([])
const loading = ref(true)
const submitting = ref(false)

const editingId = ref<string | null>(null)
const title = ref('')
const description = ref('')
const titleEn = ref('')
const descriptionEn = ref('')
const file = ref<File | null>(null)
const pendingFile = ref<File | null>(null)
const editingImageUrl = ref<string | null>(null)
const loadingExistingImage = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

async function loadFeatures() {
  loading.value = true
  features.value = await getFeatures(true)
  loading.value = false
}
onMounted(loadFeatures)

function resetForm() {
  editingId.value = null
  title.value = ''
  description.value = ''
  titleEn.value = ''
  descriptionEn.value = ''
  file.value = null
  pendingFile.value = null
  editingImageUrl.value = null
  if (fileInput.value) fileInput.value.value = ''
}

function openEditForm(feature: FeatureItem) {
  editingId.value = feature.id
  title.value = feature.title
  description.value = feature.description
  titleEn.value = feature.titleEn ?? ''
  descriptionEn.value = feature.descriptionEn ?? ''
  file.value = null
  pendingFile.value = null
  editingImageUrl.value = feature.imageUrl
  if (fileInput.value) fileInput.value.value = ''
}

async function handleAdjustExistingImage() {
  if (!editingImageUrl.value) return
  loadingExistingImage.value = true
  try {
    const res = await fetch(editingImageUrl.value)
    const blob = await res.blob()
    pendingFile.value = new File([blob], 'current-image.png', { type: blob.type || 'image/png' })
  } finally {
    loadingExistingImage.value = false
  }
}

function onFileChange(e: Event) {
  const selected = (e.target as HTMLInputElement).files?.[0] ?? null
  if (selected) pendingFile.value = selected
}

const canSubmit = computed(
  () => Boolean(title.value && description.value && (editingId.value || file.value) && !submitting.value),
)

async function handleSubmit() {
  if (!canSubmit.value) return
  submitting.value = true
  try {
    if (editingId.value) {
      await updateFeature(
        editingId.value,
        { title: title.value, description: description.value, titleEn: titleEn.value || null, descriptionEn: descriptionEn.value || null },
        file.value ?? undefined,
      )
    } else if (file.value) {
      await uploadFeature(file.value, title.value, description.value, titleEn.value || undefined, descriptionEn.value || undefined)
    }
    resetForm()
    await loadFeatures()
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: string) {
  await deleteFeature(id)
  if (editingId.value === id) resetForm()
  await loadFeatures()
}

async function handleToggleActive(feature: FeatureItem) {
  await updateFeature(feature.id, { isActive: !feature.isActive })
  await loadFeatures()
}
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-2xl bg-white p-6 shadow-md">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-semibold text-gray-800">{{ editingId ? t('dashboard.features.editFeature') : t('dashboard.features.newFeature') }}</h2>
        <button v-if="editingId" class="text-gray-400 hover:text-gray-600" @click="resetForm"><Icon name="fe:close" /></button>
      </div>

      <form class="mt-4 flex flex-col gap-3" @submit.prevent="handleSubmit">
        <input v-model="title" :placeholder="t('dashboard.features.titlePlaceholder')" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
        <textarea v-model="description" :placeholder="t('dashboard.features.descriptionPlaceholder')" rows="2" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />

        <div class="rounded-lg border border-dashed border-gray-200 p-3">
          <p class="text-xs font-semibold text-gray-500">{{ t('dashboard.translationLabel') }}</p>
          <p class="mt-0.5 text-xs text-gray-400">{{ t('dashboard.translationHint') }}</p>
          <div class="mt-2 flex flex-col gap-2">
            <input v-model="titleEn" :placeholder="t('dashboard.features.titleEnPlaceholder')" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
            <textarea v-model="descriptionEn" :placeholder="t('dashboard.features.descriptionEnPlaceholder')" rows="2" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
          </div>
        </div>

        <div class="flex w-fit items-center gap-2">
          <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 transition hover:bg-gray-50">
            <Icon name="fe:upload" />
            {{ file ? file.name : editingId ? t('dashboard.features.changeImageOptional') : t('dashboard.features.chooseImage') }}
            <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onFileChange" />
          </label>
          <button v-if="file" type="button" class="text-sm font-medium text-[#920f0f] hover:text-[#7a0c0c]" @click="pendingFile = file">
            {{ t('dashboard.features.adjustZoom') }}
          </button>
          <button
            v-if="!file && editingId && editingImageUrl"
            type="button"
            :disabled="loadingExistingImage"
            class="text-sm font-medium text-[#920f0f] hover:text-[#7a0c0c] disabled:text-gray-400"
            @click="handleAdjustExistingImage"
          >
            {{ loadingExistingImage ? t('dashboard.features.loadingImage') : t('dashboard.features.adjustCurrentImage') }}
          </button>
        </div>

        <button
          type="submit"
          :disabled="!canSubmit"
          class="w-fit rounded-lg px-4 py-2 text-sm font-semibold text-white transition"
          :class="canSubmit ? 'bg-[#920f0f] hover:bg-[#7a0c0c]' : 'cursor-not-allowed bg-gray-300'"
        >
          {{ submitting ? t('dashboard.features.saving') : editingId ? t('dashboard.features.saveChanges') : t('dashboard.features.addFeature') }}
        </button>
      </form>
    </div>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.features.loadingFeatures') }}</p>
      <p v-else-if="features.length === 0" class="text-sm text-gray-400">{{ t('dashboard.features.empty') }}</p>

      <div v-for="feature in features" :key="feature.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
        <img :src="feature.imageUrl" :alt="feature.title" class="h-40 w-full object-cover" />
        <div class="p-3">
          <p class="truncate text-sm font-semibold text-gray-800">{{ feature.title }}</p>
          <p class="mt-1 line-clamp-2 text-xs text-gray-500">{{ feature.description }}</p>
          <div class="mt-2 flex items-center justify-between">
            <label class="flex items-center gap-1.5 text-xs text-gray-500">
              <input type="checkbox" :checked="feature.isActive" @change="handleToggleActive(feature)" /> {{ t('dashboard.features.active') }}
            </label>
            <div class="flex items-center gap-3">
              <button :aria-label="t('dashboard.features.editAria')" class="text-gray-400 hover:text-[#920f0f]" @click="openEditForm(feature)"><Icon name="fe:pencil" /></button>
              <button :aria-label="t('dashboard.features.deleteAria')" class="text-gray-400 hover:text-red-500" @click="handleDelete(feature.id)"><Icon name="lucide:trash-2" /></button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ImageCropModal
      v-if="pendingFile"
      :file="pendingFile"
      @cancel="pendingFile = null"
      @confirm="(cropped) => { file = cropped; pendingFile = null }"
    />
  </div>
</template>
