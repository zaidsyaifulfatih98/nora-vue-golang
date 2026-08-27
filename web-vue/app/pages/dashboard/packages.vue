<script setup lang="ts">
import type { PackageItem } from '~/composables/api/packages'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getPackages, createPackage, updatePackage, deletePackage } = usePackagesApi()

const packages = ref<PackageItem[]>([])
const loading = ref(true)
const editingId = ref<string | null>(null)
const showForm = ref(false)
const submitting = ref(false)

const emptyForm = {
  name: '',
  price: '',
  duration: '',
  description: '',
  features: '',
  isPopular: false,
  isActive: true,
}
const form = reactive({ ...emptyForm })

function formatRupiah(value: string | number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(
    Number(value),
  )
}

async function loadPackages() {
  loading.value = true
  packages.value = await getPackages(true)
  loading.value = false
}

onMounted(loadPackages)

function openCreateForm() {
  editingId.value = null
  Object.assign(form, emptyForm)
  showForm.value = true
}

function openEditForm(pkg: PackageItem) {
  editingId.value = pkg.id
  Object.assign(form, {
    name: pkg.name,
    price: pkg.price,
    duration: pkg.duration,
    description: pkg.description,
    features: pkg.features.join('\n'),
    isPopular: pkg.isPopular,
    isActive: pkg.isActive,
  })
  showForm.value = true
}

async function handleSubmit() {
  submitting.value = true
  try {
    const payload = {
      name: form.name,
      price: Number(form.price),
      duration: form.duration,
      description: form.description,
      features: form.features.split('\n').map((f) => f.trim()).filter(Boolean),
      isPopular: form.isPopular,
      isActive: form.isActive,
    }

    if (editingId.value) {
      await updatePackage(editingId.value, payload)
    } else {
      await createPackage(payload)
    }

    showForm.value = false
    await loadPackages()
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: string) {
  await deletePackage(id)
  await loadPackages()
}

async function handleToggleActive(pkg: PackageItem) {
  await updatePackage(pkg.id, { isActive: !pkg.isActive })
  await loadPackages()
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <p class="text-sm text-gray-500">Kelola harga, deskripsi, dan fitur paket yang tampil di landing page.</p>
      <button
        class="flex items-center gap-2 rounded-lg bg-[#920f0f] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#7a0c0c]"
        @click="openCreateForm"
      >
        <Icon name="fe:plus" />
        Tambah Paket
      </button>
    </div>

    <div v-if="showForm" class="rounded-2xl bg-white p-6 shadow-md">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-semibold text-gray-800">{{ editingId ? 'Edit Paket' : 'Tambah Paket Baru' }}</h2>
        <button class="text-gray-400 hover:text-gray-600" @click="showForm = false">
          <Icon name="fe:close" />
        </button>
      </div>

      <form class="mt-4 grid gap-4 sm:grid-cols-2" @submit.prevent="handleSubmit">
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">Nama Paket</label>
          <input
            v-model="form.name"
            required
            class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">Harga (Rp)</label>
          <input
            v-model="form.price"
            required
            type="number"
            min="0"
            class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">Durasi</label>
          <input
            v-model="form.duration"
            required
            placeholder="mis. 5 jam sesi"
            class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="flex items-center gap-6 pt-6">
          <label class="flex items-center gap-2 text-sm text-gray-600">
            <input v-model="form.isPopular" type="checkbox" />
            Tandai Populer
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-600">
            <input v-model="form.isActive" type="checkbox" />
            Aktif / Tampilkan
          </label>
        </div>

        <div class="flex flex-col gap-1.5 sm:col-span-2">
          <label class="text-sm font-medium text-gray-700">Deskripsi</label>
          <textarea
            v-model="form.description"
            required
            rows="2"
            class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="flex flex-col gap-1.5 sm:col-span-2">
          <label class="text-sm font-medium text-gray-700">Fitur (satu baris = satu fitur)</label>
          <textarea
            v-model="form.features"
            required
            rows="5"
            placeholder="Backdrop custom sesuai tema&#10;Cetak foto unlimited&#10;2 orang crew photobooth"
            class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <button
          type="submit"
          :disabled="submitting"
          class="rounded-lg bg-[#920f0f] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#7a0c0c] disabled:opacity-60 sm:col-span-2"
        >
          {{ submitting ? 'Menyimpan...' : 'Simpan Paket' }}
        </button>
      </form>
    </div>

    <div class="grid gap-6 lg:grid-cols-3">
      <p v-if="loading" class="text-sm text-gray-400">Memuat paket...</p>

      <div v-for="pkg in packages" v-else :key="pkg.id" class="flex flex-col rounded-2xl bg-white p-6 shadow-md">
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-lg font-bold text-gray-900">{{ pkg.name }}</h3>
            <p class="text-xs text-gray-400">{{ pkg.duration }}</p>
          </div>
          <span v-if="pkg.isPopular" class="rounded-full bg-amber-100 px-2.5 py-1 text-xs font-semibold text-amber-700">
            Populer
          </span>
        </div>

        <p class="mt-3 text-xl font-bold text-gray-900">{{ formatRupiah(pkg.price) }}</p>
        <p class="mt-2 text-sm text-gray-500">{{ pkg.description }}</p>

        <ul class="mt-3 flex-1 space-y-1 text-sm text-gray-600">
          <li v-for="f in pkg.features" :key="f">• {{ f }}</li>
        </ul>

        <div class="mt-4 flex items-center justify-between border-t border-gray-100 pt-4">
          <label class="flex items-center gap-2 text-xs text-gray-500">
            <input type="checkbox" :checked="pkg.isActive" @change="handleToggleActive(pkg)" />
            Aktif
          </label>
          <div class="flex items-center gap-3">
            <button class="text-gray-400 hover:text-[#920f0f]" aria-label="Edit" @click="openEditForm(pkg)">
              <Icon name="fe:pencil" />
            </button>
            <button class="text-gray-400 hover:text-red-500" aria-label="Hapus" @click="handleDelete(pkg.id)">
              <Icon name="lucide:trash-2" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
