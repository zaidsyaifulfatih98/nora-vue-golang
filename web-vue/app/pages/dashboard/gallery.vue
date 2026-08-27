<script setup lang="ts">
import type { GalleryPhotoItem } from '~/composables/api/gallery'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getGalleryPhotos, uploadGalleryPhoto, updateGalleryPhoto, deleteGalleryPhoto } = useGalleryApi()

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
  if (!confirm('Hapus foto ini dari galeri?')) return
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
      <h2 class="text-base font-semibold text-gray-800">Unggah Foto Baru</h2>
      <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center">
        <input v-model="caption" placeholder="Keterangan foto (opsional)" class="flex-1 rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500" />
        <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700">
          <Icon name="fe:upload" />
          {{ uploading ? 'Mengunggah...' : 'Pilih & Unggah Foto' }}
          <input ref="fileInput" type="file" accept="image/*" :disabled="uploading" class="hidden" @change="handleFileChange" />
        </label>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      <p v-if="loading" class="text-sm text-gray-400">Memuat galeri...</p>
      <p v-else-if="photos.length === 0" class="text-sm text-gray-400">Belum ada foto di galeri.</p>

      <div v-for="photo in photos" :key="photo.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
        <img :src="photo.url" :alt="photo.caption ?? 'Galeri Nora Photobooth'" class="h-40 w-full object-cover" />
        <button
          aria-label="Hapus foto"
          class="absolute right-2 bottom-2 flex h-6 w-6 items-center justify-center rounded-full bg-white/90 text-xs text-gray-500 shadow-md transition hover:bg-red-600 hover:text-white"
          @click="handleDelete(photo.id)"
        >
          <Icon name="lucide:trash-2" />
        </button>
        <div class="p-3">
          <p class="truncate text-xs text-gray-500">{{ photo.caption || 'Tanpa keterangan' }}</p>
          <div class="mt-2 flex items-center justify-between">
            <label class="flex items-center gap-1.5 text-xs text-gray-500">
              <input type="checkbox" :checked="photo.isActive" @change="handleToggleActive(photo)" /> Aktif
            </label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
