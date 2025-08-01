import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createStore } from 'vuex'
import { createRouter, createWebHistory } from 'vue-router'
import Welcome from '@/pages/general/Welcome.vue'

// Mock der Form-Komponente
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

  const createWrapper = () => {
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

    router = createRouter({
      history: createWebHistory(),
      routes: [{ path: '/', component: Welcome }]
    })

    return mount(Welcome, {
      global: {
        plugins: [store, router]
      }
    })
  }

  beforeEach(() => {
    wrapper = createWrapper()
  })

  it('renders the welcome page correctly', () => {
    expect(wrapper.find('h2').text()).toContain('Neuen Benutzer erstellen')
    expect(wrapper.find('[data-testid="mock-form"]').exists()).toBe(true)
  })

  it('loads users on mount', async () => {
    await flushPromises()

    expect(store.dispatch).toHaveBeenCalledWith('notifications/loadPersistentNotifications')
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

    expect(store.dispatch).toHaveBeenCalledWith('users/createUser', userData)
    expect(store.dispatch).toHaveBeenCalledWith('notifications/success', expect.any(String))
  })

  it('shows empty state when no users', () => {
    expect(wrapper.find('.text-center').text()).toContain('Keine Benutzer gefunden')
  })

  it('displays users in table when available', async () => {
    // Update store state to have users
    store.state.users.users = [
      { id: 1, name: 'John Doe', email: 'john@example.com', age: 30, role: 'user' }
    ]

    await wrapper.vm.$nextTick()

    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.find('tbody tr').exists()).toBe(true)
  })
})
