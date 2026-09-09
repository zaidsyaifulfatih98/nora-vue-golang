<script setup lang="ts">
import type { FrameSlot, PhotoboothFrameItem } from '~/composables/api/photoboothFrames'
import type { CustomerItem } from '~/composables/api/auth'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getPhotoboothFrames, getMyPhotoboothFrames, uploadPhotoboothFrame, updatePhotoboothFrame, deletePhotoboothFrame } =
  usePhotoboothFramesApi()
const { getCustomers } = useAuthApi()
const { t } = useI18n()
const authStore = useAuthStore()

// A customer (DIGITAL_PHOTOBOOTH or SOFTWARE_PHOTOBOOTH) manages their own
// frames (scoped server-side to their account) — only the "assign to
// customer" picker below is admin-only, since a customer's uploads are
// always their own.
const isCustomer = computed(() => ['DIGITAL_PHOTOBOOTH', 'SOFTWARE_PHOTOBOOTH'].includes(authStore.user.role))

const customers = ref<CustomerItem[]>([])
const ownerId = ref('')
if (!isCustomer.value) {
  Promise.all([getCustomers('DIGITAL_PHOTOBOOTH'), getCustomers('SOFTWARE_PHOTOBOOTH')]).then(([digital, software]) => {
    customers.value = [...digital, ...software]
  })
}

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
  frames.value = isCustomer.value ? await getMyPhotoboothFrames() : await getPhotoboothFrames(true)
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
  ownerId.value = ''
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
      await uploadPhotoboothFrame(file.value, name.value, slots.value, ownerId.value || undefined)
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
            {{ editingId ? t('dashboard.photoboothFrames.editFrame') : t('dashboard.photoboothFrames.newFrame') }}
          </h2>
          <p class="mt-1 text-xs text-gray-500">
            {{ isCustomer ? t('dashboard.photoboothFrames.myFramesHint') : t('dashboard.photoboothFrames.helpText') }}
          </p>
        </div>
        <button v-if="editingId" class="shrink-0 text-gray-400 hover:text-gray-600" @click="resetForm"><Icon name="fe:close" /></button>
      </div>

      <form class="mt-4 flex flex-col gap-3" @submit.prevent="handleSubmit">
        <input
          v-model="name"
          :placeholder="t('dashboard.photoboothFrames.namePlaceholder')"
          class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
        />

        <select
          v-if="!isCustomer && !editingId"
          v-model="ownerId"
          class="w-fit rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-600 focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
        >
          <option value="">{{ t('dashboard.photoboothFrames.assignCustomerNone') }}</option>
          <option v-for="customer in customers" :key="customer.id" :value="customer.id">
            {{ customer.firstName }} {{ customer.lastName }}
          </option>
        </select>

        <div class="flex w-fit items-center gap-2">
          <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 transition hover:bg-gray-50">
            <Icon name="fe:upload" />
            {{ file ? file.name : editingId ? t('dashboard.photoboothFrames.changePngOptional') : t('dashboard.photoboothFrames.choosePng') }}
            <input ref="fileInput" type="file" accept="image/png" class="hidden" @change="onFileChange" />
          </label>
        </div>

        <div v-if="previewUrl" class="mt-1">
          <p class="mb-2 text-xs font-medium text-gray-600">
            {{ t('dashboard.photoboothFrames.slotHint') }}
          </p>
          <FrameSlotEditor v-model="slots" :image-url="previewUrl" />
        </div>

        <button
          type="submit"
          :disabled="!canSubmit"
          class="w-fit rounded-lg px-4 py-2 text-sm font-semibold text-white transition"
          :class="canSubmit ? 'bg-[#920f0f] hover:bg-[#7a0c0c]' : 'cursor-not-allowed bg-gray-300'"
        >
          {{ submitting ? t('dashboard.photoboothFrames.saving') : editingId ? t('dashboard.photoboothFrames.saveChanges') : t('dashboard.photoboothFrames.addFrame') }}
        </button>
      </form>
    </div>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.photoboothFrames.loadingFrames') }}</p>
      <p v-else-if="frames.length === 0" class="text-sm text-gray-400">{{ t('dashboard.photoboothFrames.empty') }}</p>

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
              <input type="checkbox" :checked="frame.isActive" @change="handleToggleActive(frame)" /> {{ t('dashboard.photoboothFrames.active') }}
            </label>
            <div class="flex items-center gap-3">
              <button :aria-label="t('dashboard.photoboothFrames.editAria')" class="text-gray-400 hover:text-[#920f0f]" @click="openEditForm(frame)"><Icon name="fe:pencil" /></button>
              <button :aria-label="t('dashboard.photoboothFrames.deleteAria')" class="text-gray-400 hover:text-red-500" @click="handleDelete(frame.id)"><Icon name="lucide:trash-2" /></button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
