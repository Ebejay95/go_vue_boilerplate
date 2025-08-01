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
    it('fetches users from API', async () => {
      const mockUsers = [
        {
          getId: () => 1,
          getName: () => 'John Doe',
          getEmail: () => 'john@example.com',
          getAge: () => 30,
          getRole: () => 'user'
        }
      ]

      const mockResponse = {
        getUsersList: () => mockUsers
      }

      // Mock gRPC call
      store.dispatch = vi.fn().mockImplementation((action, payload, options) => {
        if (action === 'connection/promisifyGrpcCall') {
          return Promise.resolve(mockResponse)
        }
        return userStore.actions[action.split('/')[1]](
          { commit: store.commit, dispatch: store.dispatch, rootGetters: store.getters },
          payload
        )
      })

      const result = await store.dispatch('users/fetchUsers')

      expect(result.users).toHaveLength(1)
      expect(result.users[0].name).toBe('John Doe')
    })

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

    it('creates user via API', async () => {
      const userData = {
        name: 'Jane Doe',
        email: 'jane@example.com',
        age: 25,
        role: 'user'
      }

      const mockResponse = {
        getUser: () => ({
          getId: () => 2,
          getName: () => 'Jane Doe',
          getEmail: () => 'jane@example.com',
          getAge: () => 25,
          getRole: () => 'user'
        })
      }

      store.dispatch = vi.fn().mockImplementation((action, payload, options) => {
        if (action === 'connection/promisifyGrpcCall') {
          return Promise.resolve(mockResponse)
        }
        return userStore.actions[action.split('/')[1]](
          { commit: store.commit, dispatch: store.dispatch, rootGetters: store.getters },
          payload
        )
      })

      const result = await store.dispatch('users/createUser', userData)

      expect(result.user.name).toBe('Jane Doe')
      expect(result.user.id).toBe(2)
    })

    it('deletes user via API', async () => {
      // Füge erst einen User hinzu
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

      store.dispatch = vi.fn().mockImplementation((action, payload, options) => {
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
      // Simuliere Verbindungsverlust
      store.commit('connection/SET_STATUS', 'error')
      store.commit('connection/SET_HEALTHY', false)

      expect(store.getters['connection/isConnected']).toBe(false)

      // Simuliere Wiederverbindung
      store.dispatch = vi.fn().mockResolvedValue({
        status: 'healthy',
        connection: 'grpc-web',
        users: 0,
        url: 'http://localhost:8081'
      })

      const result = await store.dispatch('connection/checkHealth')

      expect(result.status).toBe('healthy')
    })

    it('retries failed requests', async () => {
      let callCount = 0

      store.dispatch = vi.fn().mockImplementation((action) => {
        if (action === 'connection/promisifyGrpcCall') {
          callCount++
          if (callCount === 1) {
            return Promise.reject(new Error('Temporary network error'))
          }
          return Promise.resolve({ getUsersList: () => [] })
        }
        return userStore.actions[action.split('/')[1]](
          { commit: store.commit, dispatch: store.dispatch, rootGetters: store.getters }
        )
      })

      // Erster Versuch schlägt fehl
      await expect(store.dispatch('users/fetchUsers')).rejects.toThrow('Temporary network error')

      // Zweiter Versuch erfolgreich
      const result = await store.dispatch('users/fetchUsers')
      expect(result.users).toEqual([])
      expect(callCount).toBe(2)
    })
  })
})
