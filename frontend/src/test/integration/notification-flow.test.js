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

  describe('Complete Notification Lifecycle', () => {
    it('handles full notification lifecycle', async () => {
      // 1. Erstelle persistente Notification
      await store.dispatch('notifications/createNotification', {
        message: 'Test notification',
        type: 'info',
        persistent: true
      })

      // 2. Simuliere Socket-Event für die erstellte Notification
      await store.dispatch('notifications/handleSocketNotification', {
        id: 1,
        message: 'Test notification',
        type: 'info',
        persistent: true,
        read: false,
        createdAt: new Date().toISOString()
      })

      // 3. Prüfe State
      expect(store.state.notifications.toasts).toHaveLength(1)
      expect(store.state.notifications.persistentNotifications).toHaveLength(1)
      expect(store.getters['notifications/unreadCount']).toBe(1)

      // 4. Markiere als gelesen
      await store.dispatch('notifications/markAsRead', 1)

      // 5. Prüfe Updated State
      expect(store.getters['notifications/unreadCount']).toBe(0)

      // 6. Lösche Notification
      await store.dispatch('notifications/removePersistentNotification', 1)

      // 7. Prüfe Final State
      expect(store.state.notifications.persistentNotifications).toHaveLength(0)
    })

    it('handles bulk operations', async () => {
      // Füge mehrere Notifications hinzu
      const notifications = [
        { id: 1, message: 'Notification 1', read: false },
        { id: 2, message: 'Notification 2', read: false },
        { id: 3, message: 'Notification 3', read: true }
      ]

      notifications.forEach(notification => {
        store.commit('notifications/ADD_PERSISTENT_NOTIFICATION', {
          ...notification,
          type: 'info',
          persistent: true,
          createdAt: new Date().toISOString()
        })
      })

      expect(store.getters['notifications/unreadCount']).toBe(2)

      // Markiere alle als gelesen
      await store.dispatch('notifications/markAllAsRead')
      expect(store.getters['notifications/unreadCount']).toBe(0)

      // Lösche alle
      await store.dispatch('notifications/clearAllPersistent')
      expect(store.state.notifications.persistentNotifications).toHaveLength(0)
    })
  })

  describe('Socket Event Handling', () => {
    it('processes different socket events correctly', async () => {
      const events = [
        {
          event: 'notification',
          data: { id: 1, message: 'New notification', type: 'info', persistent: true }
        },
        {
          event: 'notification_updated',
          data: { id: 1, message: 'Updated notification', read: true }
        },
        {
          event: 'notification_deleted',
          data: { id: 1 }
        }
      ]

      // Registriere Event-Listener
      await store.dispatch('notifications/setupSocketListeners')

      // Simuliere Events
      for (const { event, data } of events) {
        if (event === 'notification') {
          await store.dispatch('notifications/handleSocketNotification', data)
          expect(store.state.notifications.persistentNotifications).toHaveLength(1)
        } else if (event === 'notification_updated') {
          await store.dispatch('notifications/handleNotificationUpdate', data)
        } else if (event === 'notification_deleted') {
          await store.dispatch('notifications/handleNotificationDeleted', data)
          expect(store.state.notifications.persistentNotifications).toHaveLength(0)
        }
      }
    })

    it('handles socket reconnection', async () => {
      // Simuliere Verbindungsabbruch
      store.commit('socket/SET_CONNECTED', false)

      // Simuliere Wiederverbindung
      store.commit('socket/SET_CONNECTED', true)

      // Socket-Reconnect sollte Notifications neu laden
      const loadSpy = vi.spyOn(store._actions['notifications/loadPersistentNotifications'][0], 'call')

      await store.dispatch('notifications/onSocketConnected')

      expect(loadSpy).toHaveBeenCalled()
    })
  })

  describe('Error Handling', () => {
    it('handles API errors gracefully', async () => {
      // Mock API-Fehler
      store.dispatch = vi.fn().mockImplementation((action) => {
        if (action === 'connection/promisifyGrpcCall') {
          return Promise.reject(new Error('API Error'))
        }
        return notificationStore.actions[action.split('/')[1]](
          {
            commit: store.commit,
            dispatch: store.dispatch,
            rootGetters: store.getters,
            state: store.state.notifications
          }
        )
      })

      // Versuche persistente Notification zu erstellen
      await store.dispatch('notifications/createNotification', {
        message: 'Test notification',
        type: 'info',
        persistent: true
      })

      // Sollte trotzdem lokale Notification erstellen
      expect(store.state.notifications.persistentNotifications).toHaveLength(1)
      expect(store.state.notifications.persistentNotifications[0].persistent).toBe(false) // Fallback
    })

    it('handles socket errors', async () => {
      // Simuliere Socket-Fehler
      const errorHandler = vi.fn()

      // Mock console.error für Test
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

      try {
        await store.dispatch('notifications/handleSocketNotification', {
          // Ungültige Daten
          invalid: 'data'
        })
      } catch (error) {
        errorHandler(error)
      }

      // Fehler sollte abgefangen werden
      expect(consoleSpy).not.toHaveBeenCalled() // Keine unbehandelten Fehler

      consoleSpy.mockRestore()
    })
  })
})
