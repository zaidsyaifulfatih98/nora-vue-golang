<script setup lang="ts">
import type { AnalyticsSummary } from '~/composables/api/analytics'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { getAnalyticsSummary } = useAnalyticsApi()
const { t, te } = useI18n()

const summary = ref<AnalyticsSummary | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await getAnalyticsSummary()
    // Go's json.Marshal encodes a nil slice as `null`, not `[]`, when a
    // breakdown query matches zero rows — normalize so the template can
    // always assume arrays.
    summary.value = {
      ...data,
      ctaBreakdown: data.ctaBreakdown ?? [],
      topPages: data.topPages ?? [],
    }
  } finally {
    loading.value = false
  }
})

function ctaLabel(label: string) {
  return te(`dashboard.analytics.ctaLabels.${label}`) ? t(`dashboard.analytics.ctaLabels.${label}`) : label
}

const maxCtaCount = computed(() => Math.max(1, ...(summary.value?.ctaBreakdown.map((r) => r.count) ?? [1])))
const maxPageCount = computed(() => Math.max(1, ...(summary.value?.topPages.map((r) => r.count) ?? [1])))
</script>

<template>
  <div class="space-y-6">
    <p class="text-sm text-gray-500">{{ t('dashboard.analytics.subtitle') }}</p>

    <p v-if="loading" class="text-sm text-gray-400">{{ t('dashboard.analytics.loading') }}</p>

    <template v-else>
      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <div class="rounded-2xl bg-white p-5 shadow-md">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-100">
            <Icon name="lucide:users" class="text-lg text-blue-600" />
          </div>
          <p class="mt-2 text-sm text-gray-500">{{ t('dashboard.analytics.totalPageViews') }}</p>
          <p class="text-2xl font-bold text-gray-900">{{ summary?.totalPageViews ?? 0 }}</p>
        </div>
        <div class="rounded-2xl bg-white p-5 shadow-md">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#F5E1E1]">
            <Icon name="lucide:mouse-pointer-click" class="text-lg text-[#920f0f]" />
          </div>
          <p class="mt-2 text-sm text-gray-500">{{ t('dashboard.analytics.totalCtaClicks') }}</p>
          <p class="text-2xl font-bold text-gray-900">{{ summary?.totalCtaClicks ?? 0 }}</p>
        </div>
        <div class="rounded-2xl bg-white p-5 shadow-md">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-green-100">
            <Icon name="lucide:calendar-days" class="text-lg text-green-600" />
          </div>
          <p class="mt-2 text-sm text-gray-500">{{ t('dashboard.analytics.todayPageViews') }}</p>
          <p class="text-2xl font-bold text-gray-900">{{ summary?.todayPageViews ?? 0 }}</p>
        </div>
        <div class="rounded-2xl bg-white p-5 shadow-md">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-amber-100">
            <Icon name="lucide:zap" class="text-lg text-amber-600" />
          </div>
          <p class="mt-2 text-sm text-gray-500">{{ t('dashboard.analytics.todayCtaClicks') }}</p>
          <p class="text-2xl font-bold text-gray-900">{{ summary?.todayCtaClicks ?? 0 }}</p>
        </div>
      </div>

      <div class="grid gap-6 lg:grid-cols-2">
        <div class="rounded-2xl bg-white p-6 shadow-md">
          <h2 class="text-base font-semibold text-gray-800">{{ t('dashboard.analytics.ctaBreakdownTitle') }}</h2>
          <p v-if="!summary?.ctaBreakdown.length" class="mt-4 text-sm text-gray-400">{{ t('dashboard.analytics.emptyCta') }}</p>
          <ul v-else class="mt-4 space-y-3">
            <li v-for="row in summary.ctaBreakdown" :key="row.label">
              <div class="flex items-center justify-between text-sm">
                <span class="text-gray-700">{{ ctaLabel(row.label) }}</span>
                <span class="font-semibold text-gray-900">{{ row.count }}</span>
              </div>
              <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-100">
                <div class="h-full rounded-full bg-[#920f0f]" :style="{ width: `${(row.count / maxCtaCount) * 100}%` }" />
              </div>
            </li>
          </ul>
        </div>

        <div class="rounded-2xl bg-white p-6 shadow-md">
          <h2 class="text-base font-semibold text-gray-800">{{ t('dashboard.analytics.topPagesTitle') }}</h2>
          <p v-if="!summary?.topPages.length" class="mt-4 text-sm text-gray-400">{{ t('dashboard.analytics.emptyPages') }}</p>
          <ul v-else class="mt-4 space-y-3">
            <li v-for="row in summary.topPages" :key="row.label">
              <div class="flex items-center justify-between text-sm">
                <span class="font-mono text-gray-700">{{ row.label || '/' }}</span>
                <span class="font-semibold text-gray-900">{{ row.count }}</span>
              </div>
              <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-100">
                <div class="h-full rounded-full bg-blue-500" :style="{ width: `${(row.count / maxPageCount) * 100}%` }" />
              </div>
            </li>
          </ul>
        </div>
      </div>
    </template>
  </div>
</template>
