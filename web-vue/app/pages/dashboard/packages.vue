<script setup lang="ts">
import type { PackageItem } from '~/composables/api/packages'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getPackages, createPackage, updatePackage, deletePackage } = usePackagesApi()
const { t } = useI18n()

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
  nameEn: '',
  durationEn: '',
  descriptionEn: '',
  featuresEn: '',
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
    nameEn: pkg.nameEn ?? '',
    durationEn: pkg.durationEn ?? '',
    descriptionEn: pkg.descriptionEn ?? '',
    featuresEn: (pkg.featuresEn ?? []).join('\n'),
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
      nameEn: form.nameEn || undefined,
      durationEn: form.durationEn || undefined,
      descriptionEn: form.descriptionEn || undefined,
      featuresEn: form.featuresEn ? form.featuresEn.split('\n').map((f) => f.trim()).filter(Boolean) : undefined,
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
      <p class="text-sm text-gray-500">{{ t('dashboard.packages.subtitle') }}</p>
      <button
        class="flex items-center gap-2 rounded-lg bg-[#920f0f] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#7a0c0c]"
        @click="openCreateForm"
      >
        <Icon name="fe:plus" />
        {{ t('dashboard.packages.addPackage') }}
      </button>
    </div>

    <div v-if="showForm" class="rounded-2xl bg-white p-6 shadow-md">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-semibold text-gray-800">{{ editingId ? t('dashboard.packages.editPackage') : t('dashboard.packages.newPackage') }}</h2>
        <button class="text-gray-400 hover:text-gray-600" @click="showForm = false">
          <Icon name="fe:close" />
        </button>
      </div>

      <form class="mt-4 grid gap-4 sm:grid-cols-2" @submit.prevent="handleSubmit">
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.packages.nameLabel') }}</label>
          <input
            v-model="form.name"
            required
            class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.packages.priceLabel') }}</label>
          <input
            v-model="form.price"
            required
            type="number"
            min="0"
            class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.packages.durationLabel') }}</label>
          <input
            v-model="form.duration"
            required
            :placeholder="t('dashboard.packages.durationPlaceholder')"
            class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="flex items-center gap-6 pt-6">
          <label class="flex items-center gap-2 text-sm text-gray-600">
            <input v-model="form.isPopular" type="checkbox" />
            {{ t('dashboard.packages.markPopular') }}
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-600">
            <input v-model="form.isActive" type="checkbox" />
            {{ t('dashboard.packages.activeShow') }}
          </label>
        </div>

        <div class="flex flex-col gap-1.5 sm:col-span-2">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.packages.descriptionLabel') }}</label>
          <textarea
            v-model="form.description"
            required
            rows="2"
            class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="flex flex-col gap-1.5 sm:col-span-2">
          <label class="text-sm font-medium text-gray-700">{{ t('dashboard.packages.featuresLabel') }}</label>
          <textarea
            v-model="form.features"
            required
            rows="5"
            :placeholder="t('dashboard.packages.featuresPlaceholder')"
            class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
          />
        </div>

        <div class="rounded-lg border border-dashed border-gray-200 p-3 sm:col-span-2">
          <p class="text-xs font-semibold text-gray-500">{{ t('dashboard.translationLabel') }}</p>
          <p class="mt-0.5 text-xs text-gray-400">{{ t('dashboard.translationHint') }}</p>
          <div class="mt-2 grid gap-3 sm:grid-cols-2">
            <input v-model="form.nameEn" :placeholder="t('dashboard.packages.nameEnPlaceholder')" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
            <input v-model="form.durationEn" :placeholder="t('dashboard.packages.durationEnPlaceholder')" class="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
            <div class="flex flex-col gap-1.5 sm:col-span-2">
              <label class="text-xs font-medium text-gray-600">{{ t('dashboard.packages.descriptionEnLabel') }}</label>
              <textarea v-model="form.descriptionEn" rows="2" class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]" />
            </div>
            <div class="flex flex-col gap-1.5 sm:col-span-2">
              <label class="text-xs font-medium text-gray-600">{{ t('dashboard.packages.featuresEnLabel') }}</label>
              <textarea
                v-model="form.featuresEn"
                rows="5"
                :placeholder="t('dashboard.packages.featuresEnPlaceholder')"
                class="resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-[#920f0f] focus:outline-none focus:ring-1 focus:ring-[#920f0f]"
              />
            </div>
          </div>
        </div>

        <button
          type="submit"
          :disabled="submitting"
          class="rounded-lg bg-[#920f0f] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#7a0c0c] disabled:opacity-60 sm:col-span-2"
        >
          {{ submitting ? t('dashboard.packages.saving') : t('dashboard.packages.savePackage') }}
        </button>
      </form>
    </div>

    <div class="grid gap-6 lg:grid-cols-3">
      <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.packages.loadingPackages') }}</p>

      <div v-for="pkg in packages" v-else :key="pkg.id" class="flex flex-col rounded-2xl bg-white p-6 shadow-md">
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-lg font-bold text-gray-900">{{ pkg.name }}</h3>
            <p class="text-xs text-gray-400">{{ pkg.duration }}</p>
          </div>
          <span v-if="pkg.isPopular" class="rounded-full bg-amber-100 px-2.5 py-1 text-xs font-semibold text-amber-700">
            {{ t('dashboard.packages.popular') }}
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
            {{ t('dashboard.packages.active') }}
          </label>
          <div class="flex items-center gap-3">
            <button class="text-gray-400 hover:text-[#920f0f]" :aria-label="t('dashboard.packages.editAria')" @click="openEditForm(pkg)">
              <Icon name="fe:pencil" />
            </button>
            <button class="text-gray-400 hover:text-red-500" :aria-label="t('dashboard.packages.deleteAria')" @click="handleDelete(pkg.id)">
              <Icon name="lucide:trash-2" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
