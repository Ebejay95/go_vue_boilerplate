import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createStore } from 'vuex'
import { createRouter, createWebHistory } from 'vue-router'
import Welcome from '@/pages/general/Welcome.vue'

// Mock the Form component
vi.mock('@/components/base/Form.vue', () => ({
  default: {
    name: 'Form',
    template: '<div data-testid="mock-form"><slot></slot></div>',
    props: ['protoMessage', 'fieldConfig', 'customValidators'],
    emits: ['submit', 'error']
  }
}))

describe('Welcome Page', () => {
  let wrapper
  let store
  let router
  let mockDispatch

  const createWrapper = () => {
    mockDispatch = vi.fn()

    store = createStore({
      modules: {
        users: {
          namespaced: true,
          state: {
            users: [],
            loading: false,
            creating: false,
            error: null
          },
          getters: {
            hasUsers: () => false,
            userCount: () => 0
          },
          actions: {
            fetchUsers: vi.fn().mockResolvedValue({ users: [], count: 0 }),
            createUser: vi.fn().mockResolvedValue({ user: { id: 1, name: 'Test' } }),
            clearError: vi.fn()
          }
        },
        connection: {
          namespaced: true,
          state: { grpcWebUrl: 'http://localhost:8081' },
          getters: {
            isConnected: () => true,
            connectionInfo: () => ({ status: 'connected' })
          },
          actions: {
            checkHealth: vi.fn(),
            reconnect: vi.fn()
          }
        },
        notifications: {
          namespaced: true,
          getters: {
            persistentNotifications: () => [],
            unreadCount: () => 0,
            hasUnreadPersistent: () => false
          },
          actions: {
            info: vi.fn(),
            success: vi.fn(),
            error: vi.fn(),
            warning: vi.fn(),
            loadPersistentNotifications: vi.fn()
          }
        },
        socket: {
          namespaced: true,
          getters: {
            connectionStatus: () => 'connected'
          },
          actions: {
            on: vi.fn(),
            emit: vi.fn()
          }
        }
      }
    })

    // Mock the dispatch method
    store.dispatch = mockDispatch

    router = createRouter({
      history: createWebHistory(),
      routes: [{ path: '/', component: Welcome }]
    })

    return mount(Welcome, {
      global: {
        plugins: [store, router],
        stubs: {
          'base-section': { template: '<div><slot></slot></div>' },
          'base-card': { template: '<div><slot></slot></div>' },
          'base-button': { template: '<button><slot></slot></button>' }
        }
      }
    })
  }

  beforeEach(() => {
    wrapper = createWrapper()
  })

  it('loads users on mount', async () => {
    await flushPromises()

    expect(mockDispatch).toHaveBeenCalledWith('notifications/loadPersistentNotifications')
  })

  it('handles user creation', async () => {
    const formComponent = wrapper.findComponent({ name: 'Form' })
    const userData = {
      name: 'John Doe',
      email: 'john@example.com',
      age: 30,
      role: 'user'
    }

    await formComponent.vm.$emit('submit', {}, userData)
    await flushPromises()

    expect(mockDispatch).toHaveBeenCalledWith('users/createUser', userData)
    expect(mockDispatch).toHaveBeenCalledWith('notifications/success', expect.any(Object))
  })

  it('displays users in table when available', async () => {
    // Update store state to have users
    store.state.users.users = [
      { id: 1, name: 'John Doe', email: 'john@example.com', age: 30, role: 'user' }
    ]

    // Update getter to return true
    store.getters['users/hasUsers'] = true

    await wrapper.vm.$nextTick()

    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.find('tbody tr').exists()).toBe(true)
  })
})
