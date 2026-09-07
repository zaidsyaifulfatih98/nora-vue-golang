// Stubs for the Nuxt auto-imported composables used by the code under test.
// These are normally injected by Nuxt's build pipeline (#imports), which
// doesn't exist under plain Vitest — so anything a test touches that relies
// on one of these needs it declared here.
import * as vue from 'vue'
import { ref, computed } from 'vue'
import { beforeEach, vi } from 'vitest'

// Nuxt auto-imports Vue's Composition API (ref, computed, watch, onMounted,
// reactive, nextTick, ...) so component <script setup> blocks use them
// without an import statement. Plain Vitest doesn't do that, so mirror it
// here for every export Vue provides.
for (const [name, value] of Object.entries(vue)) {
  vi.stubGlobal(name, value)
}

// useState: Nuxt's SSR-safe shared ref. A plain ref is enough for unit tests
// since there's no server/client hydration boundary to worry about here.
const stateStore = new Map<string, ReturnType<typeof ref>>()
vi.stubGlobal('useState', (key: string, init: () => unknown) => {
  if (!stateStore.has(key)) stateStore.set(key, ref(init()))
  return stateStore.get(key)!
})

// useI18n: real locale switching without loading actual translation JSON.
// `locale` is a module-level singleton (like the real plugin) so every
// composable/component under test that calls useI18n() shares one reactive
// value — tests can flip `locale.value` and see the effect everywhere.
const locale = ref('id')
vi.stubGlobal('useI18n', () => ({
  locale,
  locales: computed(() => [
    { code: 'id', name: 'Indonesia', flag: '🇮🇩' },
    { code: 'en', name: 'English', flag: '🇬🇧' },
  ]),
  setLocale: vi.fn((code: string) => {
    locale.value = code
  }),
  t: (key: string) => key,
  te: () => true,
  tm: () => [],
  rt: (value: unknown) => value,
}))

// Reset shared fakes between tests so one test's `locale.value = 'en'`
// doesn't leak into the next.
beforeEach(() => {
  locale.value = 'id'
  stateStore.clear()
})
