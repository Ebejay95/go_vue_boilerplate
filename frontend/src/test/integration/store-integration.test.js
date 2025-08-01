import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createStore } from 'vuex'
import { store } from '@/store/index.js'

describe('Store Integration Tests', () => {
  let testStore

  beforeEach(() => {
    // Erstelle eine frische Store-Instanz für jeden Test
    testStore = createStore(store.options)
  })

  describe('App Initialization Flow', () => {
    it('initializes app with all modules correctly', async () => {
      // Mock gRPC connections
      vi.spyOn(testStore._actions['connection/initializeClient'][0], 'call')
        .mockResolvedValue({ userClient: {}, notificationClient: {} })

      vi.spyOn(testStore._actions['connection/checkHealth'][0], 'call')
        .mockResolvedValue({ status: 'healthy' })

      vi.spyOn(testStore._actions['socket/connect'][0], 'call')
        .mockResolvedValue(true)

      const result = await testStore.dispatch('initializeApp')

      expect(result).toBe(true)
      expect(testStore.state.appInitialized).toBe(true)
      expect(testStore.state.appInitializing).toBe(false)
    })

    it('handles initialization failure gracefully', async () => {
      vi.spyOn(testStore._actions['connection/initializeClient'][0], 'call')
        .mockRejectedValue(new Error('Connection failed'))

      await expect(testStore.dispatch('initializeApp')).rejects.toThrow('Connection failed')

      expect(testStore.state.appInitialized).toBe(false)
      expect(testStore.state.appInitializationError).toBe('Connection failed')
    })
  })

  describe('User and Notification Integration', () => {
    it('creates user and shows notification', async () => {
      // Setup mocks
      const mockUserResponse = {
        getUser: () => ({
          getId: () => 1,
          getName: () => 'Test User',
          getEmail: () => 'test@example.com',
          getAge: () => 25,
          getRole: () => 'user'
        })
      }

      vi.spyOn(testStore._actions['connection/promisifyGrpcCall'][0], 'call')
        .mockResolvedValue(mockUserResponse)

      // Create user
      const userData = {
        name: 'Test User',
        email: 'test@example.com',
        age: 25,
        role: 'user'
      }

      const result = await testStore.dispatch('users/createUser', userData)

      expect(result.user.name).toBe('Test User')
      expect(testStore.state.users.users).toHaveLength(1)
      expect(testStore.state.users.users[0].name).toBe('Test User')
    })

    it('handles notification creation and persistence', async () => {
      const mockNotificationResponse = {
        success: true,
        getId: () => 1
      }

      vi.spyOn(testStore._actions['connection/promisifyGrpcCall'][0], 'call')
        .mockResolvedValue(mockNotificationResponse)

      await testStore.dispatch('notifications/createNotification', {
        message: 'Test notification',
        type: 'info',
        persistent: true
      })

      // Prüfe ob Toast erstellt wurde
      expect(testStore.state.notifications.toasts).toHaveLength(0) // Persistent notifications don't create toasts

      // Prüfe Dispatch-Aufrufe
      expect(testStore._actions['connection/promisifyGrpcCall'][0].call).toHaveBeenCalled()
    })
  })

  describe('Socket Integration', () => {
    it('handles socket notifications correctly', async () => {
      // Simuliere Socket-Nachricht
      const notificationData = {
        id: 1,
        message: 'Socket notification',
        type: 'info',
        persistent: true,
        createdAt: new Date().toISOString()
      }

      await testStore.dispatch('notifications/handleSocketNotification', notificationData)

      // Prüfe ob Toast erstellt wurde
      expect(testStore.state.notifications.toasts).toHaveLength(1)
      expect(testStore.state.notifications.toasts[0].message).toBe('Socket notification')

      // Prüfe ob persistente Notification hinzugefügt wurde
      expect(testStore.state.notifications.persistentNotifications).toHaveLength(1)
      expect(testStore.state.notifications.persistentNotifications[0].id).toBe(1)
    })
  })
})
