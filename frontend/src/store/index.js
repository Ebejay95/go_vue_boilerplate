// store/index.js
import { createStore } from 'vuex'
import users from './modules/users/index.js'
import connection from './modules/connection/index.js'

export const store = createStore({
  modules: {
    users,
    connection
  },

  // Global actions for initialization
  actions: {
    async initializeApp({ dispatch }) {
      console.log('🚀 Initializing application...')

      try {
        // Initialize gRPC connection
        await dispatch('connection/initializeClient')

        // Check health
        const healthResult = await dispatch('connection/checkHealth')

        if (healthResult.status === 'healthy') {
          console.log('✅ Application initialized successfully')
          return true
        } else {
          console.warn('⚠️ Application initialized with connection issues')
          return false
        }
      } catch (error) {
        console.error('❌ Application initialization failed:', error)
        throw error
      }
    }
  },

  // Global getters for convenience
  getters: {
    isAppReady(state, getters) {
      return getters['connection/isConnected']
    },

    appStatus(state, getters) {
      return {
        connection: getters['connection/connectionInfo'],
        usersLoaded: getters['users/hasUsers'],
        userCount: getters['users/userCount']
      }
    }
  }
})
