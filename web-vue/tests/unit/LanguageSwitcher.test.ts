import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import LanguageSwitcher from '~/components/LanguageSwitcher.vue'

// The real <Icon> component comes from the @nuxt/icon module, which isn't
// registered outside a running Nuxt app — stub it so mount() doesn't warn
// about an unknown element.
const globalStubs = { Icon: true }

describe('LanguageSwitcher', () => {
  it('renders one button per available locale', () => {
    const wrapper = mount(LanguageSwitcher, { global: { stubs: globalStubs } })
    expect(wrapper.findAll('button')).toHaveLength(2)
  })

  it('highlights the currently active locale', () => {
    const wrapper = mount(LanguageSwitcher, { global: { stubs: globalStubs } })
    const buttons = wrapper.findAll('button')

    // Default locale is 'id' (first button), per tests/setup.ts.
    expect(buttons[0]!.classes()).toContain('ring-2')
    expect(buttons[1]!.classes()).toContain('opacity-50')
  })

  it('switches locale when a flag is clicked', async () => {
    const wrapper = mount(LanguageSwitcher, { global: { stubs: globalStubs } })
    const { locale } = useI18n()

    await wrapper.findAll('button')[1]!.trigger('click')

    expect(locale.value).toBe('en')
  })
})
