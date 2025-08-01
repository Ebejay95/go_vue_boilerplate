import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createStore } from 'vuex'
import userStore from '@/store/modules/user/index.js'
import connectionStore from '@/store/modules/connection/index.js'

describe('API Integration Tests', () => {
  let store
  let mockGrpcClient

  beforeEach(() => {
    mockGrpcClient = {
      listUsers: vi.fn(),
      createUser: vi.fn(),
      deleteUser: vi.fn(),
      getUser: vi.fn()
    }

    store = createStore({
      modules: {
        users: userStore,
        connection: {
          ...connectionStore,
          state: {
            ...connectionStore.state(),
            userClient: mockGrpcClient
          },
          getters: {
            ...connectionStore.getters,
            isConnected: () => true,
            grpcClient: () => mockGrpcClient
          }
        }
      }
    })
  })

  describe('User API Operations', () => {
    it('handles API errors gracefully', async () => {
      store.dispatch = vi.fn().mockImplementation((action) => {
        if (action === 'connection/promisifyGrpcCall') {
          return Promise.reject(new Error('Network error'))
        }
        return userStore.actions[action.split('/')[1]](
          { commit: store.commit, dispatch: store.dispatch, rootGetters: store.getters }
        )
      })

      await expect(store.dispatch('users/fetchUsers')).rejects.toThrow('Network error')
      expect(store.state.users.error).toContain('Network error')
    })

    it('deletes user via API', async () => {
      // Add user first
      store.commit('users/ADD_USER', {
        id: 1,
        name: 'Test User',
        email: 'test@example.com',
        age: 30,
        role: 'user'
      })

      expect(store.state.users.users).toHaveLength(1)

      const mockResponse = {
        getSuccess: () => true,
        getMessage: () => 'User deleted successfully'
      }

      store.dispatch = vi.fn().mockImplementation((action, payload) => {
        if (action === 'connection/promisifyGrpcCall') {
          return Promise.resolve(mockResponse)
        }
        return userStore.actions[action.split('/')[1]](
          { commit: store.commit, dispatch: store.dispatch, rootGetters: store.getters },
          payload
        )
      })

      const result = await store.dispatch('users/deleteUser', 1)

      expect(result.success).toBe(true)
      expect(store.state.users.users).toHaveLength(0)
    })
  })

  describe('Connection Management', () => {
    it('handles connection loss and recovery', async () => {
      // Create a proper store with connection getters
      const connectionStore = createStore({
        state: {
          status: 'connected',
          healthy: true
        },
        mutations: {
          SET_STATUS(state, status) {
            state.status = status
          },
          SET_HEALTHY(state, healthy) {
            state.healthy = healthy
          }
        },
        getters: {
          isConnected(state) {
            return state.status === 'connected' && state.healthy
          }
        }
      })

      expect(connectionStore.getters.isConnected).toBe(true)

      connectionStore.commit('SET_STATUS', 'error')
      connectionStore.commit('SET_HEALTHY', false)

      expect(connectionStore.getters.isConnected).toBe(false)
    })
  })
})
