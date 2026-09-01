<script setup lang="ts">
import type { GalleryPhotoItem } from '~/composables/api/gallery'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getGalleryPhotos, uploadGalleryPhoto, updateGalleryPhoto, deleteGalleryPhoto } = useGalleryApi()
const { t } = useI18n()

const photos = ref<GalleryPhotoItem[]>([])
const loading = ref(true)
const uploading = ref(false)
const caption = ref('')
const fileInput = ref<HTMLInputElement | null>(null)

async function loadPhotos() {
  loading.value = true
  photos.value = await getGalleryPhotos(true)
  loading.value = false
}
onMounted(loadPhotos)

async function handleFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return

  uploading.value = true
  try {
    await uploadGalleryPhoto(file, caption.value || undefined)
    caption.value = ''
    if (fileInput.value) fileInput.value.value = ''
    await loadPhotos()
  } finally {
    uploading.value = false
  }
}

async function handleDelete(id: string) {
  if (!confirm(t('dashboard.gallery.confirmDelete'))) return
  await deleteGalleryPhoto(id)
  await loadPhotos()
}

async function handleToggleActive(photo: GalleryPhotoItem) {
  await updateGalleryPhoto(photo.id, { isActive: !photo.isActive })
  await loadPhotos()
}
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-2xl bg-white p-6 shadow-md">
      <h2 class="text-base font-semibold text-gray-800">{{ t('dashboard.gallery.uploadNewPhoto') }}</h2>
      <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center">
        <input v-model="caption" :placeholder="t('dashboard.gallery.captionPlaceholder')" class="flex-1 rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
        <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-[#920f0f] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#7a0c0c]">
          <Icon name="fe:upload" />
          {{ uploading ? t('dashboard.gallery.uploading') : t('dashboard.gallery.chooseUpload') }}
          <input ref="fileInput" type="file" accept="image/*" :disabled="uploading" class="hidden" @change="handleFileChange" />
        </label>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.gallery.loadingGallery') }}</p>
      <p v-else-if="photos.length === 0" class="text-sm text-gray-400">{{ t('dashboard.gallery.empty') }}</p>

      <div v-for="photo in photos" :key="photo.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
        <img :src="photo.url" :alt="photo.caption ?? t('dashboard.gallery.altFallback')" class="h-40 w-full object-cover" />
        <button
          :aria-label="t('dashboard.gallery.deletePhotoAria')"
          class="absolute right-2 bottom-2 flex h-6 w-6 items-center justify-center rounded-full bg-white/90 text-xs text-gray-500 shadow-md transition hover:bg-red-600 hover:text-white"
          @click="handleDelete(photo.id)"
        >
          <Icon name="lucide:trash-2" />
        </button>
        <div class="p-3">
          <p class="truncate text-xs text-gray-500">{{ photo.caption || t('dashboard.gallery.noCaption') }}</p>
          <div class="mt-2 flex items-center justify-between">
            <label class="flex items-center gap-1.5 text-xs text-gray-500">
              <input type="checkbox" :checked="photo.isActive" @change="handleToggleActive(photo)" /> {{ t('dashboard.gallery.active') }}
            </label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
