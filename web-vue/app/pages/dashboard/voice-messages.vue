<script setup lang="ts">
import type { VoiceMessageItem } from '~/composables/api/voiceMessages'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getVoiceMessages, deleteVoiceMessage } = useVoiceMessagesApi()
const { t, locale } = useI18n()

const messages = ref<VoiceMessageItem[]>([])
const loading = ref(true)

async function loadMessages() {
  loading.value = true
  messages.value = await getVoiceMessages()
  loading.value = false
}
onMounted(loadMessages)

async function handleDelete(id: string) {
  await deleteVoiceMessage(id)
  await loadMessages()
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString(locale.value === 'en' ? 'en-US' : 'id-ID', { dateStyle: 'medium', timeStyle: 'short' })
}
</script>

<template>
  <div class="space-y-6">
    <p class="text-sm text-gray-500">
      {{ t('dashboard.voiceMessages.subtitle') }}
    </p>

    <div class="space-y-4">
      <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.voiceMessages.loadingMessages') }}</p>
      <p v-else-if="messages.length === 0" class="text-sm text-gray-400">{{ t('dashboard.voiceMessages.empty') }}</p>

      <div v-for="message in messages" :key="message.id" class="flex flex-col gap-4 rounded-2xl bg-white p-5 shadow-md sm:flex-row sm:items-center">
        <img
          v-if="message.photoUrl"
          :src="message.photoUrl"
          :alt="t('dashboard.voiceMessages.resultAlt')"
          class="h-20 w-16 shrink-0 rounded-lg object-cover ring-1 ring-gray-100"
        />
        <div v-else class="flex h-20 w-16 shrink-0 items-center justify-center rounded-lg bg-gray-50 ring-1 ring-gray-100">
          <Icon name="lucide:image-off" class="text-lg text-gray-300" />
        </div>

        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <p class="font-semibold text-gray-800">{{ message.guestName || t('dashboard.voiceMessages.guestFallback') }}</p>
            <span class="text-xs text-gray-400">{{ formatDate(message.createdAt) }}</span>
          </div>
          <audio :src="message.audioUrl" controls class="mt-2 h-9 w-full max-w-md" />
        </div>

        <button
          :aria-label="t('dashboard.voiceMessages.deleteAria')"
          class="shrink-0 self-start text-gray-400 hover:text-red-500 sm:self-center"
          @click="handleDelete(message.id)"
        >
          <Icon name="lucide:trash-2" />
        </button>
      </div>
    </div>
  </div>
</template>
