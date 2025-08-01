import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createStore } from 'vuex'
import userStore from '@/store/modules/user/index.js'
import connectionStore from '@/store/modules/connection/index.js'
import notificationStore from '@/store/modules/notification/index.js'
import socketStore from '@/store/modules/socket/index.js'

describe('Store Integration Tests', () => {
  let testStore

  beforeEach(() => {
    // Create store with properly mocked actions
    testStore = createStore({
      modules: {
        users: {
          ...userStore,
          actions: {
            ...userStore.actions,
            fetchUsers: vi.fn().mockResolvedValue({ users: [], count: 0 })
          }
        },
        connection: {
          ...connectionStore,
          actions: {
            ...connectionStore.actions,
            initializeClient: vi.fn().mockResolvedValue({ userClient: {}, notificationClient: {} }),
            checkHealth: vi.fn().mockResolvedValue({ status: 'healthy' })
          }
        },
        notifications: {
          ...notificationStore,
          actions: {
            ...notificationStore.actions,
            initialize: vi.fn().mockResolvedValue()
          }
        },
        socket: {
          ...socketStore,
          actions: {
            ...socketStore.actions,
            connect: vi.fn().mockResolvedValue(true)
          }
        }
      },

      state() {
        return {
          appInitialized: false,
          appInitializing: false,
          appInitializationError: null
        }
      },

      mutations: {
        SET_APP_INITIALIZED(state, value) {
          state.appInitialized = value
        },
        SET_APP_INITIALIZING(state, value) {
          state.appInitializing = value
        },
        SET_APP_INITIALIZATION_ERROR(state, error) {
          state.appInitializationError = error
        }
      },

      actions: {
        async initializeApp({ commit, dispatch }) {
          commit('SET_APP_INITIALIZING', true)
          try {
            await dispatch('connection/initializeClient')
            await dispatch('connection/checkHealth')
            await dispatch('socket/connect')
            await dispatch('notifications/initialize')
            commit('SET_APP_INITIALIZED', true)
            return true
          } catch (error) {
            commit('SET_APP_INITIALIZATION_ERROR', error.message)
            throw error
          } finally {
            commit('SET_APP_INITIALIZING', false)
          }
        }
      }
    })
  })

  describe('App Initialization Flow', () => {
    it('initializes app with all modules correctly', async () => {
      const result = await testStore.dispatch('initializeApp')

      expect(result).toBe(true)
      expect(testStore.state.appInitialized).toBe(true)
      expect(testStore.state.appInitializing).toBe(false)
    })

    it('handles initialization failure gracefully', async () => {
      // Mock a failing action
      testStore._actions['connection/initializeClient'][0] = vi.fn().mockRejectedValue(new Error('Connection failed'))

      await expect(testStore.dispatch('initializeApp')).rejects.toThrow('Connection failed')

      expect(testStore.state.appInitialized).toBe(false)
      expect(testStore.state.appInitializationError).toBe('Connection failed')
    })
  })
})
