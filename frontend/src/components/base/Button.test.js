import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import Button from '@/components/base/Button.vue'

describe('Button Component', () => {
  const createWrapper = (props = {}, options = {}) => {
    const router = createRouter({
      history: createWebHistory(),
      routes: [{ path: '/', component: { template: '<div>Home</div>' } }]
    })

    return mount(Button, {
      props,
      global: {
        plugins: [router]
      },
      ...options
    })
  }

  it('renders correctly as button', () => {
    const wrapper = createWrapper({}, {
      slots: { default: 'Click Me' }
    })

    expect(wrapper.find('button').exists()).toBe(true)
    expect(wrapper.text()).toBe('Click Me')
    expect(wrapper.classes()).toContain('btn')
  })

  it('renders correctly as router-link', () => {
    const wrapper = createWrapper({
      link: true,
      to: '/test'
    }, {
      slots: { default: 'Navigate' }
    })

    expect(wrapper.find('router-link').exists()).toBe(true)
    expect(wrapper.find('router-link').attributes('to')).toBe('/test')
  })

  it('applies correct CSS classes based on mode', async () => {
    const wrapper = createWrapper({ mode: 'primary' })
    expect(wrapper.classes()).toContain('primary')

    await wrapper.setProps({ mode: 'secondary' })
    expect(wrapper.classes()).toContain('secondary')
  })

  it('emits click event when clicked', async () => {
    const wrapper = createWrapper()
    await wrapper.find('button').trigger('click')

    expect(wrapper.emitted('click')).toBeTruthy()
  })
})
