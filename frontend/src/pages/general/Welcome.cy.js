import Welcome from './Welcome.vue'
import { createStore } from 'vuex'

describe('Welcome Page', () => {
  const createMockStore = () => {
    return createStore({
      modules: {
        users: {
          namespaced: true,
          state: {
            users: [
              { id: 1, name: 'John Doe', email: 'john@example.com', age: 30, role: 'user' }
            ],
            loading: false,
            creating: false,
            error: null
          },
          getters: {
            hasUsers: (state) => state.users.length > 0,
            userCount: (state) => state.users.length
          },
          actions: {
            fetchUsers: cy.stub().resolves({ users: [], count: 0 }),
            createUser: cy.stub().resolves({ user: { id: 2, name: 'Jane Doe' } }),
            clearError: cy.stub()
          }
        },
        connection: {
          namespaced: true,
          state: {
            grpcWebUrl: 'http://localhost:8081'
          },
          getters: {
            isConnected: () => true,
            connectionInfo: () => ({ status: 'connected' })
          },
          actions: {
            checkHealth: cy.stub(),
            reconnect: cy.stub()
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
            info: cy.stub(),
            success: cy.stub(),
            error: cy.stub(),
            warning: cy.stub(),
            loadPersistentNotifications: cy.stub()
          }
        },
        socket: {
          namespaced: true,
          getters: {
            connectionStatus: () => 'connected'
          },
          actions: {
            on: cy.stub(),
            emit: cy.stub()
          }
        }
      }
    })
  }

  it('renders welcome page correctly', () => {
    const store = createMockStore()

    cy.mount(Welcome, {
      global: {
        plugins: [store]
      }
    })

    cy.get('h2').should('contain', 'Neuen Benutzer erstellen')
    cy.get('h2').should('contain', 'Alle Benutzer')
  })

  it('displays users in table', () => {
    const store = createMockStore()

    cy.mount(Welcome, {
      global: {
        plugins: [store]
      }
    })

    cy.get('table').should('exist')
    cy.get('tbody tr').should('contain', 'John Doe')
    cy.get('tbody tr').should('contain', 'john@example.com')
  })

  it('shows loading state when loading users', () => {
    const store = createMockStore()
    store.state.users.loading = true
    store.state.users.users = []

    cy.mount(Welcome, {
      global: {
        plugins: [store]
      }
    })

    cy.get('.animate-spin').should('exist')
    cy.get('.text-gray-500').should('contain', 'Lade Benutzer')
  })

  it('shows empty state when no users', () => {
    const store = createMockStore()
    store.state.users.users = []

    cy.mount(Welcome, {
      global: {
        plugins: [store]
      }
    })

    cy.get('.text-center').should('contain', 'Keine Benutzer gefunden')
  })

  it('handles refresh button click', () => {
    const fetchUsers = cy.stub().resolves({ users: [], count: 0 })
    const store = createMockStore()
    store._actions['users/fetchUsers'] = [fetchUsers]

    cy.mount(Welcome, {
      global: {
        plugins: [store]
      }
    })

    cy.get('button').contains('Aktualisieren').click()

    cy.then(() => {
      expect(fetchUsers).to.have.been.called
    })
  })
})
