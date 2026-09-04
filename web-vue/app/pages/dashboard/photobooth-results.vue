<script setup lang="ts">
import type { PhotoboothResultItem } from '~/composables/api/photoboothResults'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getPhotoboothResults, deletePhotoboothResult } = usePhotoboothResultsApi()
const { t, locale } = useI18n()

const results = ref<PhotoboothResultItem[]>([])
const loading = ref(true)

// Frames come in different shapes: a "4x6" composite (ratio ~1.5) vs a
// narrow "2x6" strip (ratio ~3) meant to be paired up and cut from one 4x6
// sheet. Resolved up front (before first render) so results land in the
// right section immediately instead of jumping between sections as each
// thumbnail loads.
const aspectRatios = reactive<Record<string, number>>({})
const STRIP_RATIO_THRESHOLD = 2.2

function isStrip(id: string) {
  return (aspectRatios[id] ?? 0) > STRIP_RATIO_THRESHOLD
}

const normalResults = computed(() => results.value.filter((r) => !isStrip(r.id)))
const stripResults = computed(() => results.value.filter((r) => isStrip(r.id)))

const selectedForMerge = ref<string[]>([])

function toggleMergeSelect(id: string) {
  const i = selectedForMerge.value.indexOf(id)
  if (i !== -1) {
    selectedForMerge.value.splice(i, 1)
    return
  }
  if (selectedForMerge.value.length >= 2) selectedForMerge.value.shift()
  selectedForMerge.value.push(id)
}

function clearMergeSelection() {
  selectedForMerge.value = []
}

function loadImage(src: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => resolve(img)
    img.onerror = reject
    img.src = src
  })
}

async function loadResults() {
  loading.value = true
  results.value = await getPhotoboothResults()
  await Promise.all(
    results.value.map(async (r) => {
      try {
        const img = await loadImage(r.viewUrl)
        aspectRatios[r.id] = img.height / img.width
      } catch {
        // Leave unresolved; falls into the "4R" section by default.
      }
    }),
  )
  loading.value = false
}
onMounted(loadResults)

async function handleDelete(id: string) {
  await deletePhotoboothResult(id)
  await loadResults()
}

const merging = ref(false)
const mergeError = ref('')

async function buildMergedImage(): Promise<string | null> {
  const [idA, idB] = selectedForMerge.value
  const resultA = results.value.find((r) => r.id === idA)
  const resultB = results.value.find((r) => r.id === idB)
  if (!resultA || !resultB) return null

  const [imgA, imgB] = await Promise.all([loadImage(resultA.viewUrl), loadImage(resultB.viewUrl)])

  const height = Math.max(imgA.height, imgB.height)
  const canvas = document.createElement('canvas')
  canvas.width = imgA.width + imgB.width
  canvas.height = height
  const ctx = canvas.getContext('2d')
  if (!ctx) return null

  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, canvas.width, canvas.height)
  ctx.drawImage(imgA, 0, 0)
  ctx.drawImage(imgB, imgA.width, 0)

  return canvas.toDataURL('image/png')
}

async function handleMergeAndPrint() {
  if (selectedForMerge.value.length !== 2) return
  merging.value = true
  mergeError.value = ''
  try {
    const dataUrl = await buildMergedImage()
    if (!dataUrl) {
      mergeError.value = t('dashboard.photoboothResults.mergeError')
      return
    }
    handlePrint(dataUrl)
  } catch {
    mergeError.value = t('dashboard.photoboothResults.mergeError')
  } finally {
    merging.value = false
  }
}

function handlePrint(url: string) {
  const printWindow = window.open('', '_blank', 'width=800,height=1000')
  if (!printWindow) return

  printWindow.document.write(`
    <html>
      <head>
        <title>Print</title>
        <style>
          @page { margin: 0; }
          html, body { margin: 0; padding: 0; height: 100%; display: flex; align-items: center; justify-content: center; }
          img { max-width: 100%; max-height: 100vh; }
        </style>
      </head>
      <body>
        <img src="${url}" />
      </body>
    </html>
  `)
  printWindow.document.close()

  const image = printWindow.document.querySelector('img')
  const triggerPrint = () => {
    printWindow.focus()
    printWindow.print()
  }
  if (image?.complete) triggerPrint()
  else image?.addEventListener('load', triggerPrint)
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString(locale.value === 'en' ? 'en-US' : 'id-ID', { dateStyle: 'medium', timeStyle: 'short' })
}
</script>

<template>
  <div class="space-y-8 pb-20">
    <p class="text-sm text-gray-500">
      {{ t('dashboard.photoboothResults.subtitle') }}
    </p>

    <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.photoboothResults.loadingResults') }}</p>
    <p v-else-if="results.length === 0" class="text-sm text-gray-400">{{ t('dashboard.photoboothResults.empty') }}</p>

    <template v-else>
      <section class="space-y-3">
        <h2 class="text-sm font-semibold text-gray-700">{{ t('dashboard.photoboothResults.sectionRegular') }}</h2>
        <p v-if="normalResults.length === 0" class="text-sm text-gray-400">{{ t('dashboard.photoboothResults.sectionEmpty') }}</p>
        <div v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
          <div v-for="result in normalResults" :key="result.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
            <div class="aspect-[3/5] w-full overflow-hidden bg-gray-50">
              <img :src="result.viewUrl" :alt="t('dashboard.photoboothResults.resultAlt')" class="h-full w-full object-cover" />
            </div>
            <div class="p-3">
              <p class="text-xs text-gray-400">{{ formatDate(result.createdAt) }}</p>
              <div class="mt-2 flex items-center justify-between">
                <a
                  :href="result.downloadUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex items-center gap-1.5 text-xs font-semibold text-[#920f0f] hover:underline"
                >
                  <Icon name="lucide:download" />
                  {{ t('dashboard.photoboothResults.download') }}
                </a>
                <div class="flex items-center gap-3">
                  <button
                    :aria-label="t('dashboard.photoboothResults.printAria')"
                    class="text-gray-400 hover:text-[#920f0f]"
                    @click="handlePrint(result.viewUrl)"
                  >
                    <Icon name="lucide:printer" />
                  </button>
                  <button :aria-label="t('dashboard.photoboothResults.deleteAria')" class="text-gray-400 hover:text-red-500" @click="handleDelete(result.id)">
                    <Icon name="lucide:trash-2" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="space-y-3">
        <div>
          <h2 class="text-sm font-semibold text-gray-700">{{ t('dashboard.photoboothResults.sectionStrip') }}</h2>
          <p class="text-xs text-gray-400">{{ t('dashboard.photoboothResults.mergeHint') }}</p>
        </div>
        <p v-if="stripResults.length === 0" class="text-sm text-gray-400">{{ t('dashboard.photoboothResults.sectionEmpty') }}</p>
        <div v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
          <div v-for="result in stripResults" :key="result.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
            <div class="relative aspect-[3/5] w-full overflow-hidden bg-gray-50">
              <img :src="result.viewUrl" :alt="t('dashboard.photoboothResults.resultAlt')" class="h-full w-full object-cover" />
              <label class="absolute top-2 left-2 flex h-6 w-6 items-center justify-center rounded-full bg-white/90 shadow-md ring-1 ring-[#E4E2DC]">
                <input
                  type="checkbox"
                  class="h-4 w-4 accent-[#920f0f]"
                  :checked="selectedForMerge.includes(result.id)"
                  @change="toggleMergeSelect(result.id)"
                />
              </label>
            </div>
            <div class="p-3">
              <p class="text-xs text-gray-400">{{ formatDate(result.createdAt) }}</p>
              <div class="mt-2 flex items-center justify-between">
                <a
                  :href="result.downloadUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex items-center gap-1.5 text-xs font-semibold text-[#920f0f] hover:underline"
                >
                  <Icon name="lucide:download" />
                  {{ t('dashboard.photoboothResults.download') }}
                </a>
                <div class="flex items-center gap-3">
                  <button
                    :aria-label="t('dashboard.photoboothResults.printAria')"
                    class="text-gray-400 hover:text-[#920f0f]"
                    @click="handlePrint(result.viewUrl)"
                  >
                    <Icon name="lucide:printer" />
                  </button>
                  <button :aria-label="t('dashboard.photoboothResults.deleteAria')" class="text-gray-400 hover:text-red-500" @click="handleDelete(result.id)">
                    <Icon name="lucide:trash-2" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </template>

    <div
      v-if="selectedForMerge.length > 0"
      class="fixed inset-x-0 bottom-0 z-40 flex items-center justify-between gap-4 border-t border-gray-100 bg-white px-6 py-4 shadow-[0_-4px_12px_rgba(0,0,0,0.06)] sm:left-64"
    >
      <div>
        <p class="text-sm font-semibold text-gray-800">
          {{ t('dashboard.photoboothResults.mergeSelectedCount', { count: selectedForMerge.length }) }}
        </p>
        <p v-if="mergeError" class="mt-0.5 text-xs text-red-600">{{ mergeError }}</p>
      </div>
      <div class="flex items-center gap-3">
        <button class="text-sm font-medium text-gray-500 hover:text-gray-700" @click="clearMergeSelection">
          {{ t('dashboard.photoboothResults.mergeClear') }}
        </button>
        <button
          type="button"
          :disabled="selectedForMerge.length !== 2 || merging"
          class="flex items-center gap-2 rounded-lg bg-[#920f0f] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#7a0c0c] disabled:cursor-not-allowed disabled:opacity-50"
          @click="handleMergeAndPrint"
        >
          <Icon name="lucide:printer" />
          {{ merging ? t('dashboard.photoboothResults.mergingInProgress') : t('dashboard.photoboothResults.mergeAndPrint') }}
        </button>
      </div>
    </div>
  </div>
</template>
