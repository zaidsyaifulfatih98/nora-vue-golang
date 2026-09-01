/**
 * Picks the English variant of a DB-sourced field when the active locale is
 * English and a translation was actually filled in, falling back to the
 * Indonesian (base) value otherwise.
 */
export function useLocalizedField() {
  const { locale } = useI18n()

  function tf(base: string, en?: string | null): string {
    return locale.value === 'en' && en ? en : base
  }

  function tfList(base: string[], en?: string[] | null): string[] {
    return locale.value === 'en' && en && en.length > 0 ? en : base
  }

  return { tf, tfList }
}
