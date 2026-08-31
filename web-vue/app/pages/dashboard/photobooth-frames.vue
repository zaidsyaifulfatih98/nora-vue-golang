<script setup lang="ts">
import type { FrameSlot, PhotoboothFrameItem } from '~/composables/api/photoboothFrames'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getPhotoboothFrames, uploadPhotoboothFrame, updatePhotoboothFrame, deletePhotoboothFrame } = usePhotoboothFramesApi()

function defaultSlots(): FrameSlot[] {
  const margin = 0.045
  const gap = 0.027
  const height = (1 - margin * 2 - gap * 2) / 3
  return [0, 1, 2].map((i) => ({ x: margin, y: margin + i * (height + gap), width: 1 - margin * 2, height }))
}

const frames = ref<PhotoboothFrameItem[]>([])
const loading = ref(true)
const submitting = ref(false)

const editingId = ref<string | null>(null)
const name = ref('')
const file = ref<File | null>(null)
const filePreviewUrl = ref('')
const editingImageUrl = ref('')
const slots = ref<FrameSlot[]>([])
const fileInput = ref<HTMLInputElement | null>(null)

const previewUrl = computed(() => filePreviewUrl.value || editingImageUrl.value)

async function loadFrames() {
  loading.value = true
  frames.value = await getPhotoboothFrames(true)
  loading.value = false
}
onMounted(loadFrames)

function resetForm() {
  editingId.value = null
  name.value = ''
  file.value = null
  if (filePreviewUrl.value) URL.revokeObjectURL(filePreviewUrl.value)
  filePreviewUrl.value = ''
  editingImageUrl.value = ''
  slots.value = []
  if (fileInput.value) fileInput.value.value = ''
}

function openEditForm(frame: PhotoboothFrameItem) {
  editingId.value = frame.id
  name.value = frame.name
  file.value = null
  if (filePreviewUrl.value) URL.revokeObjectURL(filePreviewUrl.value)
  filePreviewUrl.value = ''
  editingImageUrl.value = frame.imageUrl
  slots.value = frame.slots?.length ? frame.slots.map((s) => ({ ...s })) : defaultSlots()
  if (fileInput.value) fileInput.value.value = ''
}

function onFileChange(e: Event) {
  const selected = (e.target as HTMLInputElement).files?.[0] ?? null
  if (!selected) return
  file.value = selected
  if (filePreviewUrl.value) URL.revokeObjectURL(filePreviewUrl.value)
  filePreviewUrl.value = URL.createObjectURL(selected)
  slots.value = defaultSlots()
}

const canSubmit = computed(() => Boolean(name.value && (editingId.value || file.value) && !submitting.value))

async function handleSubmit() {
  if (!canSubmit.value) return
  submitting.value = true
  try {
    if (editingId.value) {
      await updatePhotoboothFrame(editingId.value, { name: name.value, slots: slots.value }, file.value ?? undefined)
    } else if (file.value) {
      await uploadPhotoboothFrame(file.value, name.value, slots.value)
    }
    resetForm()
    await loadFrames()
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: string) {
  await deletePhotoboothFrame(id)
  if (editingId.value === id) resetForm()
  await loadFrames()
}

async function handleToggleActive(frame: PhotoboothFrameItem) {
  await updatePhotoboothFrame(frame.id, { isActive: !frame.isActive })
  await loadFrames()
}
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-2xl bg-white p-6 shadow-md">
      <div class="flex items-start justify-between gap-4">
        <div>
          <h2 class="text-base font-semibold text-gray-800">
            {{ editingId ? 'Edit Frame Digital Photobooth' : 'Tambah Frame Digital Photobooth' }}
          </h2>
          <p class="mt-1 text-xs text-gray-500">
            Berbeda dari "Template Frame" di landing page. Frame ini dipakai pada fitur "Coba Digital Photobooth" —
            unggah PNG transparan, lalu atur jumlah, ukuran, dan posisi kotak foto tepat di atas lubang transparan
            pada gambar di bawah ini. Jumlah kotak menentukan berapa kali tamu akan mengambil foto.
          </p>
        </div>
        <button v-if="editingId" class="shrink-0 text-gray-400 hover:text-gray-600" @click="resetForm"><Icon name="fe:close" /></button>
      </div>

      <form class="mt-4 flex flex-col gap-3" @submit.prevent="handleSubmit">
        <input
          v-model="name"
          placeholder="Nama Frame (contoh: Frame Emas Elegan)"
          class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
        />

        <div class="flex w-fit items-center gap-2">
          <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 transition hover:bg-gray-50">
            <Icon name="fe:upload" />
            {{ file ? file.name : editingId ? 'Ganti PNG (opsional)' : 'Pilih PNG Transparan' }}
            <input ref="fileInput" type="file" accept="image/png" class="hidden" @change="onFileChange" />
          </label>
        </div>

        <div v-if="previewUrl" class="mt-1">
          <p class="mb-2 text-xs font-medium text-gray-600">
            Geser kotak untuk memindahkan, tarik sudut kanan-bawah untuk mengubah ukuran, sampai ketiganya pas di lubang transparan frame:
          </p>
          <FrameSlotEditor v-model="slots" :image-url="previewUrl" />
        </div>

        <button
          type="submit"
          :disabled="!canSubmit"
          class="w-fit rounded-lg px-4 py-2 text-sm font-semibold text-white transition"
          :class="canSubmit ? 'bg-[#920f0f] hover:bg-[#7a0c0c]' : 'cursor-not-allowed bg-gray-300'"
        >
          {{ submitting ? 'Menyimpan...' : editingId ? 'Simpan Perubahan' : 'Tambah Frame' }}
        </button>
      </form>
    </div>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      <p v-if="loading" class="text-sm text-gray-400">Memuat frame...</p>
      <p v-else-if="frames.length === 0" class="text-sm text-gray-400">Belum ada frame digital photobooth.</p>

      <div v-for="frame in frames" :key="frame.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
        <div
          class="aspect-[3/5] w-full overflow-hidden"
          style="background-image: linear-gradient(45deg, #e5e7eb 25%, transparent 25%), linear-gradient(-45deg, #e5e7eb 25%, transparent 25%), linear-gradient(45deg, transparent 75%, #e5e7eb 75%), linear-gradient(-45deg, transparent 75%, #e5e7eb 75%); background-size: 16px 16px; background-position: 0 0, 0 8px, 8px -8px, -8px 0px;"
        >
          <img :src="frame.imageUrl" :alt="frame.name" class="h-full w-full object-contain" />
        </div>
        <div class="p-3">
          <p class="truncate text-sm font-semibold text-gray-800">{{ frame.name }}</p>
          <div class="mt-2 flex items-center justify-between">
            <label class="flex items-center gap-1.5 text-xs text-gray-500">
              <input type="checkbox" :checked="frame.isActive" @change="handleToggleActive(frame)" /> Aktif
            </label>
            <div class="flex items-center gap-3">
              <button aria-label="Edit" class="text-gray-400 hover:text-[#920f0f]" @click="openEditForm(frame)"><Icon name="fe:pencil" /></button>
              <button aria-label="Hapus" class="text-gray-400 hover:text-red-500" @click="handleDelete(frame.id)"><Icon name="lucide:trash-2" /></button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
