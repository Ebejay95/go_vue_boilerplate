import { createStore } from 'vuex'
import users from './modules/users/index.js'
import connection from './modules/connection/index.js'
import notifications from './modules/notifications/index.js'

export const store = createStore({
	modules: {
		users,
		connection,
		notifications
	},

	state() {
		return {
			appInitialized: false,
			appInitializing: false,
			appInitializationError: null
		}
	},

	mutations: {
		SET_APP_INITIALIZED(state, value) {
			state.appInitialized = value
		},

		SET_APP_INITIALIZING(state, value) {
			state.appInitializing = value
		},

		SET_APP_INITIALIZATION_ERROR(state, error) {
			state.appInitializationError = error
		}
	},

	actions: {
		async initializeApp({ commit, dispatch, state }) {
			// Verhindere mehrfache Initialisierung
			if (state.appInitialized || state.appInitializing) {
				console.log('🔄 App already initialized or initializing')
				return state.appInitialized
			}

			console.log('🚀 Starting application initialization...')
			commit('SET_APP_INITIALIZING', true)
			commit('SET_APP_INITIALIZATION_ERROR', null)

			try {
				// 1. Initialize gRPC connection
				console.log('📡 Initializing gRPC connection...')
				await dispatch('connection/initializeClient')

				// 2. Check health
				console.log('🏥 Checking connection health...')
				const healthResult = await dispatch('connection/checkHealth')

				if (healthResult.status === 'healthy') {
					console.log('✅ App initialization successful')
					commit('SET_APP_INITIALIZED', true)

					// Optional: Load initial data
					try {
						console.log('📋 Loading initial user data...')
						await dispatch('users/fetchUsers')
					} catch (userError) {
						console.warn('⚠️ Could not load initial users:', userError.message)
						// Don't fail initialization just because users couldn't be loaded
					}

					return true
				} else {
					console.warn('⚠️ App initialized with connection issues')
					commit('SET_APP_INITIALIZATION_ERROR', 'Connection health check failed')
					return false
				}
			} catch (error) {
				console.error('❌ App initialization failed:', error)
				commit('SET_APP_INITIALIZATION_ERROR', error.message)
				throw error
			} finally {
				commit('SET_APP_INITIALIZING', false)
			}
		},

		async retryInitialization({ dispatch }) {
			console.log('🔄 Retrying app initialization...')
			return await dispatch('initializeApp')
		}
	},

	getters: {
		isAppInitialized(state) {
			return state.appInitialized
		},

		isAppInitializing(state) {
			return state.appInitializing
		},

		appInitializationError(state) {
			return state.appInitializationError
		},

		appStatus(state) {
			if (state.appInitializing) return 'initializing'
			if (state.appInitialized) return 'ready'
			if (state.appInitializationError) return 'error'
			return 'not_initialized'
		}
	}
})