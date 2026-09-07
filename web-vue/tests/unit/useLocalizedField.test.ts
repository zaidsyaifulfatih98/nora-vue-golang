import { describe, expect, it } from 'vitest'
import { useLocalizedField } from '~/composables/useLocalizedField'

describe('useLocalizedField', () => {
  describe('tf (single string)', () => {
    it('falls back to the Indonesian value when locale is "id"', () => {
      const { tf } = useLocalizedField()
      // The stub in tests/setup.ts defaults locale to 'id'.
      expect(tf('Nama Paket', 'Package Name')).toBe('Nama Paket')
    })

    it('returns the English value once the shared locale switches to "en"', () => {
      const { locale } = useI18n()
      const { tf } = useLocalizedField()

      locale.value = 'en'

      expect(tf('Nama Paket', 'Package Name')).toBe('Package Name')
    })

    it('falls back to the base value when no English translation was filled in', () => {
      const { tf } = useLocalizedField()
      expect(tf('Nama Paket', null)).toBe('Nama Paket')
      expect(tf('Nama Paket', undefined)).toBe('Nama Paket')
      // An empty string counts as "not translated yet", same as null.
      expect(tf('Nama Paket', '')).toBe('Nama Paket')
    })
  })

  describe('tfList (string array, e.g. package features)', () => {
    it('falls back to the base list when the English list is empty', () => {
      const { tfList } = useLocalizedField()
      const base = ['Backdrop custom', 'Cetak unlimited']
      expect(tfList(base, [])).toEqual(base)
      expect(tfList(base, null)).toEqual(base)
    })
  })
})
