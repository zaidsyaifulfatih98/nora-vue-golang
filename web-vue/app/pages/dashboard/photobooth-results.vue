<script setup lang="ts">
import type { PhotoboothResultItem } from '~/composables/api/photoboothResults'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getPhotoboothResults, deletePhotoboothResult } = usePhotoboothResultsApi()

const results = ref<PhotoboothResultItem[]>([])
const loading = ref(true)

async function loadResults() {
  loading.value = true
  results.value = await getPhotoboothResults()
  loading.value = false
}
onMounted(loadResults)

async function handleDelete(id: string) {
  await deletePhotoboothResult(id)
  await loadResults()
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })
}
</script>

<template>
  <div class="space-y-6">
    <p class="text-sm text-gray-500">
      Kumpulan hasil foto yang disimpan tamu lewat fitur "Coba Digital Photobooth" di landing page.
    </p>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
      <p v-if="loading" class="text-sm text-gray-400">Memuat hasil...</p>
      <p v-else-if="results.length === 0" class="text-sm text-gray-400">Belum ada hasil digital photobooth.</p>

      <div v-for="result in results" :key="result.id" class="group relative overflow-hidden rounded-2xl bg-white shadow-md">
        <div class="aspect-[3/5] w-full overflow-hidden bg-gray-50">
          <img :src="result.viewUrl" alt="Hasil digital photobooth" class="h-full w-full object-cover" />
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
              Download
            </a>
            <button aria-label="Hapus" class="text-gray-400 hover:text-red-500" @click="handleDelete(result.id)">
              <Icon name="lucide:trash-2" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
