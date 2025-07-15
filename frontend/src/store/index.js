// src/store/index.js - Store Integration
import { createStore } from 'vuex'
import users from './modules/users/index.js'
import connection from './modules/connection/index.js'
import notifications from './modules/notifications/index.js' // NEU

export const store = createStore({
  modules: {
    users,
    connection,
    notifications // NEU
  },

  actions: {
    async initializeApp({ dispatch }) {
      console.log('🚀 Initializing application...')

      try {
        // Show loading notification
        dispatch('notifications/info', 'Anwendung wird initialisiert...')

        // Initialize gRPC connection
        await dispatch('connection/initializeClient')

        // Check health
        const healthResult = await dispatch('connection/checkHealth')

        if (healthResult.status === 'healthy') {
          console.log('✅ Application initialized successfully')
          dispatch('notifications/success', '✅ Anwendung erfolgreich gestartet!')
          return true
        } else {
          console.warn('⚠️ Application initialized with connection issues')
          dispatch('notifications/warning', '⚠️ Verbindungsprobleme erkannt')
          return false
        }
      } catch (error) {
        console.error('❌ Application initialization failed:', error)
        dispatch('notifications/error', `❌ Initialisierung fehlgeschlagen: ${error.message}`)
        throw error
      }
    }
  }
})
