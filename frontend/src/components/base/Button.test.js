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
        plugins: [router],
        stubs: {
          'router-link': {
            template: '<a><slot></slot></a>',
            props: ['to']
          }
        }
      },
      ...options
    })
  }

  it('renders correctly as router-link', () => {
    const wrapper = createWrapper({
      link: true,
      to: '/test'
    }, {
      slots: { default: 'Navigate' }
    })

    // Check for the stubbed router-link (rendered as 'a' tag)
    expect(wrapper.find('a').exists()).toBe(true)
    expect(wrapper.text()).toBe('Navigate')
  })
})
