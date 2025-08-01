import { describe, it, expect, vi, beforeEach } from 'vitest'
import notificationStore from '@/store/modules/notification/index.js'

describe('Notification Store', () => {
  let store
  let mockRootGetters

  beforeEach(() => {
    mockRootGetters = {
      'connection/isConnected': true,
      'connection/notificationClient': {
        createNotification: vi.fn(),
        listNotifications: vi.fn(),
        deleteNotification: vi.fn()
      }
    }

    store = {
      state: notificationStore.state(),
      commit: vi.fn((mutation, payload) => {
        // Simuliere Mutations für Tests
        if (mutation === 'ADD_TOAST') {
          store.state.toasts.push(payload)
        }
      }),
      dispatch: vi.fn(),
      rootGetters: mockRootGetters
    }
  })

  describe('Actions', () => {
    it('showToast adds toast to state', async () => {
      const toastData = {
        message: 'Test message',
        type: 'success',
        duration: 3000
      }

      await notificationStore.actions.showToast(store, toastData)

      expect(store.commit).toHaveBeenCalledWith('ADD_TOAST', expect.objectContaining({
        message: 'Test message',
        type: 'success',
        duration: 3000
      }))
    })

    it('createNotification handles persistent notifications', async () => {
      const notificationData = {
        message: 'Persistent notification',
        type: 'info',
        persistent: true
      }

      store.dispatch.mockResolvedValue({ success: true })

      await notificationStore.actions.createNotification(store, notificationData)

      expect(store.dispatch).toHaveBeenCalledWith('connection/promisifyGrpcCall',
        expect.any(Object),
        { root: true }
      )
    })
  })

  describe('Getters', () => {
    it('unreadCount returns correct count', () => {
      store.state.persistentNotifications = [
        { id: 1, read: false },
        { id: 2, read: true },
        { id: 3, read: false }
      ]

      const count = notificationStore.getters.unreadCount(store.state)
      expect(count).toBe(2)
    })

    it('hasUnreadPersistent returns true when unread notifications exist', () => {
      store.state.persistentNotifications = [
        { id: 1, read: false }
      ]

      const hasUnread = notificationStore.getters.hasUnreadPersistent(store.state)
      expect(hasUnread).toBe(true)
    })
  })
})
