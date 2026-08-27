<script setup lang="ts">
import type { FrameTemplateFit, FrameTemplateItem } from '~/composables/api/frameTemplates'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getFrameTemplates, uploadFrameTemplate, updateFrameTemplate, deleteFrameTemplate } = useFrameTemplatesApi()

const FIT_OPTIONS: { value: FrameTemplateFit; label: string; hint: string }[] = [
  { value: 'COVER', label: 'Isi Penuh (Cover)', hint: 'Gambar memenuhi kartu, sisi yang berlebih akan terpotong.' },
  { value: 'CONTAIN', label: 'Tampilkan Utuh (Contain)', hint: 'Seluruh gambar terlihat, area kosong diisi warna latar.' },
]

const templates = ref<FrameTemplateItem[]>([])
const loading = ref(true)
const submitting = ref(false)

const editingId = ref<string | null>(null)
const name = ref('')
const description = ref('')
const fit = ref<FrameTemplateFit>('COVER')
const file = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

async function loadTemplates() {
  loading.value = true
  templates.value = await getFrameTemplates(true)
  loading.value = false
}
onMounted(loadTemplates)

function resetForm() {
  editingId.value = null
  name.value = ''
  description.value = ''
  fit.value = 'COVER'
  file.value = null
  if (fileInput.value) fileInput.value.value = ''
}

function openEditForm(template: FrameTemplateItem) {
  editingId.value = template.id
  name.value = template.name
  description.value = template.description
  fit.value = template.fit
  file.value = null
  if (fileInput.value) fileInput.value.value = ''
}

function onFileChange(e: Event) {
  const selected = (e.target as HTMLInputElement).files?.[0] ?? null
  if (selected) file.value = selected
}

const canSubmit = computed(
  () => Boolean(name.value && description.value && (editingId.value || file.value) && !submitting.value),
)

async function handleSubmit() {
  if (!canSubmit.value) return
  submitting.value = true
  try {
    if (editingId.value) {
      await updateFrameTemplate(
        editingId.value,
        { name: name.value, description: description.value, fit: fit.value },
        file.value ?? undefined,
      )
    } else if (file.value) {
      await uploadFrameTemplate(file.value, name.value, description.value, fit.value)
    }
    resetForm()
    await loadTemplates()
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: string) {
  await deleteFrameTemplate(id)
  if (editingId.value === id) resetForm()
  await loadTemplates()
}

async function handleToggleActive(template: FrameTemplateItem) {
  await updateFrameTemplate(template.id, { isActive: !template.isActive })
  await loadTemplates()
}

async function handleToggleFit(template: FrameTemplateItem) {
  await updateFrameTemplate(template.id, { fit: template.fit === 'COVER' ? 'CONTAIN' : 'COVER' })
  await loadTemplates()
}
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-2xl bg-white p-6 shadow-md">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-semibold text-gray-800">{{ editingId ? 'Edit Template Frame' : 'Tambah Template Frame' }}</h2>
        <button v-if="editingId" class="text-gray-400 hover:text-gray-600" @click="resetForm"><Icon name="fe:close" /></button>
      </div>

      <form class="mt-4 flex flex-col gap-3" @submit.prevent="handleSubmit">
        <input v-model="name" placeholder="Nama Template (contoh: Flora)" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
        <textarea v-model="description" placeholder="Deskripsi singkat template" rows="2" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />

        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">Fit Ukuran Gambar</label>
          <p class="text-xs text-gray-400">
            Gambar tidak dipotong otomatis saat upload &mdash; pilihan ini yang menentukan tampilannya di kartu (ukuran kartu selalu seragam).
          </p>
          <div class="flex flex-col gap-2 sm:flex-row">
            <label
              v-for="option in FIT_OPTIONS"
              :key="option.value"
              class="flex flex-1 cursor-pointer flex-col gap-0.5 rounded-lg border px-3 py-2 text-sm transition"
              :class="fit === option.value ? 'border-[#920f0f] bg-[#f7f3eb] text-[#7a0c0c]' : 'border-gray-200 text-gray-600 hover:bg-gray-50'"
            >
              <span class="flex items-center gap-2 font-medium">
                <input type="radio" name="fit" :value="option.value" :checked="fit === option.value" @change="fit = option.value" />
                {{ option.label }}
              </span>
              <span class="pl-5 text-xs text-gray-400">{{ option.hint }}</span>
            </label>
          </div>
        </div>

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
          {{ submitting ? 'Menyimpan...' : editingId ? 'Simpan Perubahan' : 'Tambah Template' }}
        </button>
      </form>
    </div>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      <p v-if="loading" class="text-sm text-gray-400">Memuat template...</p>
      <p v-else-if="templates.length === 0" class="text-sm text-gray-400">Belum ada template frame.</p>

      <div v-for="template in templates" :key="template.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
        <div class="aspect-[4/3] w-full overflow-hidden bg-gray-100">
          <img :src="template.imageUrl" :alt="template.name" class="h-full w-full" :class="template.fit === 'CONTAIN' ? 'object-contain' : 'object-cover'" />
        </div>
        <div class="p-3">
          <div class="flex items-center gap-2">
            <p class="truncate text-sm font-semibold text-gray-800">{{ template.name }}</p>
            <button
              type="button"
              class="shrink-0 rounded-full bg-gray-100 px-2 py-0.5 text-[10px] font-medium text-gray-500 transition hover:bg-gray-200"
              title="Klik untuk ganti fit"
              @click="handleToggleFit(template)"
            >
              {{ template.fit === 'CONTAIN' ? 'Contain' : 'Cover' }}
            </button>
          </div>
          <p class="mt-1 line-clamp-2 text-xs text-gray-500">{{ template.description }}</p>
          <div class="mt-2 flex items-center justify-between">
            <label class="flex items-center gap-1.5 text-xs text-gray-500">
              <input type="checkbox" :checked="template.isActive" @change="handleToggleActive(template)" /> Aktif
            </label>
            <div class="flex items-center gap-3">
              <button aria-label="Edit" class="text-gray-400 hover:text-[#920f0f]" @click="openEditForm(template)"><Icon name="fe:pencil" /></button>
              <button aria-label="Hapus" class="text-gray-400 hover:text-red-500" @click="handleDelete(template.id)"><Icon name="lucide:trash-2" /></button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
