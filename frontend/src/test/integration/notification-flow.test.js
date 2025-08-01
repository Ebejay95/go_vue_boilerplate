import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createStore } from 'vuex'
import notificationStore from '@/store/modules/notification/index.js'
import socketStore from '@/store/modules/socket/index.js'

describe('Notification Flow Integration', () => {
  let store
  let mockSocket

  beforeEach(() => {
    mockSocket = {
      send: vi.fn(),
      close: vi.fn(),
      readyState: 1 // OPEN
    }

    store = createStore({
      modules: {
        notifications: notificationStore,
        socket: {
          ...socketStore,
          state: {
            ...socketStore.state(),
            socket: mockSocket,
            isConnected: true
          }
        },
        connection: {
          namespaced: true,
          getters: {
            isConnected: () => true,
            notificationClient: () => ({
              createNotification: vi.fn(),
              listNotifications: vi.fn(),
              deleteNotification: vi.fn()
            })
          },
          actions: {
            promisifyGrpcCall: vi.fn().mockResolvedValue({ success: true })
          }
        }
      }
    })
  })

  describe('Socket Event Handling', () => {
    it('handles socket reconnection', async () => {
      // Mock the loadPersistentNotifications action
      const loadSpy = vi.fn()
      store._actions['notifications/loadPersistentNotifications'] = [loadSpy]

      // Simulate connection loss
      store.commit('socket/SET_CONNECTED', false)

      // Simulate reconnection
      store.commit('socket/SET_CONNECTED', true)

      // Trigger the socket connected handler
      await store.dispatch('notifications/onSocketConnected')

      expect(loadSpy).toHaveBeenCalled()
    })
  })

  describe('Error Handling', () => {
    it('handles API errors gracefully', async () => {
      // Mock API error properly
      store.dispatch = vi.fn().mockImplementation((action, payload) => {
        if (action === 'connection/promisifyGrpcCall') {
          return Promise.reject(new Error('API Error'))
        }
        return notificationStore.actions[action.split('/')[1]](
          {
            commit: store.commit,
            dispatch: store.dispatch,
            rootGetters: store.getters,
            state: store.state.notifications
          },
          payload // Pass the payload correctly
        )
      })

      // Test with proper payload structure
      await store.dispatch('notifications/createNotification', {
        message: 'Test notification',
        type: 'info',
        persistent: true
      })

      // Should create local notification as fallback
      expect(store.state.notifications.persistentNotifications).toHaveLength(1)
      expect(store.state.notifications.persistentNotifications[0].persistent).toBe(false) // Fallback
    })
  })
})
