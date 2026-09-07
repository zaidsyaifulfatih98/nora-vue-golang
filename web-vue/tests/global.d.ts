// Ambient declarations for the Nuxt auto-imports stubbed in tests/setup.ts,
// so test files can reference them as bare globals (matching how the real
// app calls them) without a TypeScript "not defined" error.
import type { Ref, ComputedRef } from 'vue'

declare global {
  function useState<T>(key: string, init: () => T): Ref<T>
  function useI18n(): {
    locale: Ref<string>
    locales: ComputedRef<{ code: string; name: string; flag: string }[]>
    setLocale: (code: string) => void
    t: (key: string) => string
    te: (key: string) => boolean
    tm: (key: string) => unknown[]
    rt: (value: unknown) => string
  }
}

export {}
