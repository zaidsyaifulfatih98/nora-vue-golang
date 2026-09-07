import { describe, expect, it } from 'vitest'
import { useDashboardSidebar } from '~/composables/useDashboardSidebar'

describe('useDashboardSidebar', () => {
  it('starts closed', () => {
    const sidebarOpen = useDashboardSidebar()
    expect(sidebarOpen.value).toBe(false)
  })

  it('shares the same state across every call (single source of truth)', () => {
    const first = useDashboardSidebar()
    const second = useDashboardSidebar()

    first.value = true

    expect(second.value).toBe(true)
  })
})
