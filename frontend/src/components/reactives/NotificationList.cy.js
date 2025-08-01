import NotificationList from './NotificationList.vue'
import { createStore } from 'vuex'

describe('NotificationList Component', () => {
  const createMockStore = (initialState = {}) => {
    return createStore({
      modules: {
        notifications: {
          namespaced: true,
          state: {
            toasts: [],
            persistentNotifications: [],
            loading: false,
            error: null,
            ...initialState
          },
          getters: {
            unreadCount: (state) => state.persistentNotifications.filter(n => !n.read).length,
            hasUnreadPersistent: (state) => state.persistentNotifications.some(n => !n.read),
            isLoading: (state) => state.loading,
            error: (state) => state.error
          },
          actions: {
            dismissToast: cy.stub(),
            removePersistentNotification: cy.stub(),
            markAsRead: cy.stub(),
            markAllAsRead: cy.stub(),
            clearAllPersistent: cy.stub(),
            loadPersistentNotifications: cy.stub()
          }
        }
      }
    })
  }

  it('displays notification badge with correct count', () => {
    const store = createMockStore({
      persistentNotifications: [
        { id: 1, message: 'Test 1', read: false },
        { id: 2, message: 'Test 2', read: false },
        { id: 3, message: 'Test 3', read: true }
      ]
    })

    cy.mount(NotificationList, {
      global: {
        plugins: [store]
      }
    })

    cy.get('.bg-red-500').should('contain', '2')
  })

  it('opens and closes notification panel', () => {
    const store = createMockStore()

    cy.mount(NotificationList, {
      global: {
        plugins: [store]
      }
    })

    // Panel should be closed initially
    cy.get('[data-testid="notification-panel"]').should('not.exist')

    // Click bell to open
    cy.get('svg').first().click()
    cy.get('h3').should('contain', 'Benachrichtigungen')

    // Click close button
    cy.get('button').contains('×').click()
    cy.get('[data-testid="notification-panel"]').should('not.exist')
  })

  it('displays toast notifications', () => {
    const store = createMockStore({
      toasts: [
        {
          id: 1,
          message: 'Success message',
          type: 'success',
          duration: 3000
        }
      ]
    })

    cy.mount(NotificationList, {
      global: {
        plugins: [store]
      }
    })

    cy.get('.toast').should('be.visible')
    cy.get('.toast').should('contain', 'Success message')
    cy.get('.border-green-400').should('exist')
  })

  it('shows empty state when no notifications', () => {
    const store = createMockStore()

    cy.mount(NotificationList, {
      global: {
        plugins: [store]
      }
    })

    cy.get('svg').first().click()
    cy.get('h4').should('contain', 'Keine Benachrichtigungen')
  })

  it('marks notification as read when clicked', () => {
    const markAsRead = cy.stub()
    const store = createMockStore({
      persistentNotifications: [
        { id: 1, message: 'Unread notification', read: false, createdAt: new Date().toISOString() }
      ]
    })

    store._actions['notifications/markAsRead'] = [markAsRead]

    cy.mount(NotificationList, {
      global: {
        plugins: [store]
      }
    })

    cy.get('svg').first().click()
    cy.get('.bg-blue-50').click()

    cy.then(() => {
      expect(markAsRead).to.have.been.calledWith(1)
    })
  })

  it('removes notification when delete button clicked', () => {
    const removePersistent = cy.stub()
    const store = createMockStore({
      persistentNotifications: [
        { id: 1, message: 'Test notification', read: false, createdAt: new Date().toISOString() }
      ]
    })

    store._actions['notifications/removePersistentNotification'] = [removePersistent]

    cy.mount(NotificationList, {
      global: {
        plugins: [store]
      }
    })

    cy.get('svg').first().click()
    cy.get('[data-testid="delete-notification"]').click()

    cy.then(() => {
      expect(removePersistent).to.have.been.calledWith(1)
    })
  })
})
