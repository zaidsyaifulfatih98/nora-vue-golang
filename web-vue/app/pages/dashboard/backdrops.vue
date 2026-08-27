<script setup lang="ts">
import type { BackdropItem } from '~/composables/api/backdrops'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getBackdrops, uploadBackdrop, updateBackdrop, deleteBackdrop } = useBackdropsApi()

const backdrops = ref<BackdropItem[]>([])
const loading = ref(true)
const submitting = ref(false)

const editingId = ref<string | null>(null)
const name = ref('')
const file = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

async function loadBackdrops() {
  loading.value = true
  backdrops.value = await getBackdrops(true)
  loading.value = false
}
onMounted(loadBackdrops)

function resetForm() {
  editingId.value = null
  name.value = ''
  file.value = null
  if (fileInput.value) fileInput.value.value = ''
}

function openEditForm(backdrop: BackdropItem) {
  editingId.value = backdrop.id
  name.value = backdrop.name
  file.value = null
  if (fileInput.value) fileInput.value.value = ''
}

function onFileChange(e: Event) {
  const selected = (e.target as HTMLInputElement).files?.[0] ?? null
  if (selected) file.value = selected
}

const canSubmit = computed(() => Boolean(name.value && (editingId.value || file.value) && !submitting.value))

async function handleSubmit() {
  if (!canSubmit.value) return
  submitting.value = true
  try {
    if (editingId.value) {
      await updateBackdrop(editingId.value, { name: name.value }, file.value ?? undefined)
    } else if (file.value) {
      await uploadBackdrop(file.value, name.value)
    }
    resetForm()
    await loadBackdrops()
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: string) {
  await deleteBackdrop(id)
  if (editingId.value === id) resetForm()
  await loadBackdrops()
}

async function handleToggleActive(backdrop: BackdropItem) {
  await updateBackdrop(backdrop.id, { isActive: !backdrop.isActive })
  await loadBackdrops()
}
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-2xl bg-white p-6 shadow-md">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-semibold text-gray-800">{{ editingId ? 'Edit Backdrop' : 'Tambah Backdrop' }}</h2>
        <button v-if="editingId" class="text-gray-400 hover:text-gray-600" @click="resetForm"><Icon name="fe:close" /></button>
      </div>

      <form class="mt-4 flex flex-col gap-3" @submit.prevent="handleSubmit">
        <input v-model="name" placeholder="Nama Backdrop (contoh: Black)" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />

        <div class="flex w-fit items-center gap-2">
          <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 transition hover:bg-gray-50">
            <Icon name="fe:upload" />
            {{ file ? file.name : editingId ? 'Ganti Gambar (opsional)' : 'Pilih Gambar' }}
            <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onFileChange" />
          </label>
        </div>

        <button
          type="submit"
          :disabled="!canSubmit"
          class="w-fit rounded-lg px-4 py-2 text-sm font-semibold text-white transition"
          :class="canSubmit ? 'bg-[#920f0f] hover:bg-[#7a0c0c]' : 'cursor-not-allowed bg-gray-300'"
        >
          {{ submitting ? 'Menyimpan...' : editingId ? 'Simpan Perubahan' : 'Tambah Backdrop' }}
        </button>
      </form>
    </div>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      <p v-if="loading" class="text-sm text-gray-400">Memuat backdrop...</p>
      <p v-else-if="backdrops.length === 0" class="text-sm text-gray-400">Belum ada backdrop.</p>

      <div v-for="backdrop in backdrops" :key="backdrop.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
        <div class="aspect-[4/3] w-full overflow-hidden bg-gray-100">
          <img :src="backdrop.imageUrl" :alt="backdrop.name" class="h-full w-full object-cover" />
        </div>
        <div class="p-3">
          <p class="truncate text-sm font-semibold text-gray-800">{{ backdrop.name }}</p>
          <div class="mt-2 flex items-center justify-between">
            <label class="flex items-center gap-1.5 text-xs text-gray-500">
              <input type="checkbox" :checked="backdrop.isActive" @change="handleToggleActive(backdrop)" /> Aktif
            </label>
            <div class="flex items-center gap-3">
              <button aria-label="Edit" class="text-gray-400 hover:text-[#920f0f]" @click="openEditForm(backdrop)"><Icon name="fe:pencil" /></button>
              <button aria-label="Hapus" class="text-gray-400 hover:text-red-500" @click="handleDelete(backdrop.id)"><Icon name="lucide:trash-2" /></button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
