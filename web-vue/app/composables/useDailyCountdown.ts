/**
 * Hitung mundur sampai akhir hari ini (23:59:59 waktu perangkat pengunjung);
 * setelah lewat tengah malam otomatis mulai lagi untuk hari berikutnya.
 * Nilai baru diisi setelah mounted agar render server dan client tidak mismatch,
 * sebelum itu `units` berisi "--".
 */
export function useDailyCountdown() {
  const remaining = ref<number | null>(null)
  let timer: ReturnType<typeof setInterval> | undefined

  function tick() {
    const now = new Date()
    const end = new Date(now)
    end.setHours(23, 59, 59, 999)
    remaining.value = Math.max(0, end.getTime() - now.getTime())
  }

  onMounted(() => {
    tick()
    timer = setInterval(tick, 1000)
  })
  onBeforeUnmount(() => timer && clearInterval(timer))

  const pad = (n: number) => String(n).padStart(2, '0')
  const units = computed(() => {
    const total = remaining.value
    if (total === null) {
      return [
        { label: 'Jam', value: '--' },
        { label: 'Menit', value: '--' },
        { label: 'Detik', value: '--' },
      ]
    }
    const s = Math.floor(total / 1000)
    return [
      { label: 'Jam', value: pad(Math.floor(s / 3600)) },
      { label: 'Menit', value: pad(Math.floor((s % 3600) / 60)) },
      { label: 'Detik', value: pad(s % 60) },
    ]
  })

  return { units }
}
