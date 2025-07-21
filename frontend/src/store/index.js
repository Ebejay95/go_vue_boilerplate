import { createStore } from 'vuex'
import users from './modules/user/index.js'
import connection from './modules/connection/index.js'
import notifications from './modules/notification/index.js'
import socket from './modules/socket/index.js' // Renamed from websocket to socket

export const store = createStore({
	modules: {
		users,
		connection,
		notifications,
		socket // General socket system
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
			if (state.appInitialized || state.appInitializing) {
				return state.appInitialized
			}

			commit('SET_APP_INITIALIZING', true)
			commit('SET_APP_INITIALIZATION_ERROR', null)

			try {
				// 1. Initialize gRPC connection
				await dispatch('connection/initializeClient')
				const healthResult = await dispatch('connection/checkHealth')

				if (healthResult.status === 'healthy') {
					// 2. Connect to socket system
					await dispatch('socket/connect')

					// 3. Initialize notification system (sets up socket listeners)
					await dispatch('notifications/initialize')

					commit('SET_APP_INITIALIZED', true)

					// Load initial data
					const loadingPromises = []

					loadingPromises.push(
						dispatch('users/fetchUsers').catch(error => {
							console.warn('Could not load initial users:', error.message)
						})
					)

					loadingPromises.push(
						dispatch('notifications/loadPersistentNotifications').catch(error => {
							console.warn('Could not load persistent notifications:', error.message)
						})
					)

					// Wait for all initial data to load
					await Promise.allSettled(loadingPromises)

					return true
				} else {
					commit('SET_APP_INITIALIZATION_ERROR', 'Connection health check failed')
					return false
				}
			} catch (error) {
				console.error('App initialization failed:', error)
				commit('SET_APP_INITIALIZATION_ERROR', error.message)
				throw error
			} finally {
				commit('SET_APP_INITIALIZING', false)
			}
		},

		async retryInitialization({ dispatch }) {
			return await dispatch('initializeApp')
		},

		// Cleanup when app is closed
		async cleanupApp({ dispatch }) {
			await dispatch('socket/disconnect')
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
