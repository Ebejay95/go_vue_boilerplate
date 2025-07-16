// store/modules/connection/index.js
import { UserServiceClient } from '../../../proto/user_grpc_web_pb'
import { ListUsersRequest } from '../../../proto/user_pb'

export default {
	namespaced: true,

	state() {
		return {
			status: 'disconnected', // 'connecting', 'connected', 'disconnected', 'error'
			healthy: false,
			message: 'Not connected',
			grpcWebUrl: process.env.VUE_APP_GRPC_WEB_URL,
			client: null,
			lastChecked: null,
			error: null
		}
	},

	mutations: {
		SET_STATUS(state, status) {
			state.status = status
		},

		SET_HEALTHY(state, healthy) {
			state.healthy = healthy
		},

		SET_MESSAGE(state, message) {
			state.message = message
		},

		SET_CLIENT(state, client) {
			state.client = client
		},

		SET_LAST_CHECKED(state, timestamp) {
			state.lastChecked = timestamp
		},

		SET_ERROR(state, error) {
			state.error = error
		},

		CLEAR_ERROR(state) {
			state.error = null
		}
	},

	actions: {
		initializeClient({ commit, state }) {
			console.log(`🔧 Initializing gRPC-Web client for ${state.grpcWebUrl}`)

			try {
				const client = new UserServiceClient(state.grpcWebUrl, null, null)
				commit('SET_CLIENT', client)
				commit('SET_STATUS', 'initialized')
				console.log('✅ gRPC-Web client initialized')
				return client
			} catch (error) {
				console.error('❌ Failed to initialize gRPC-Web client:', error)
				commit('SET_STATUS', 'error')
				commit('SET_ERROR', `Client initialization failed: ${error.message}`)
				throw error
			}
		},

		async checkHealth({ commit, state, dispatch }) {
			console.log('🏥 Checking gRPC-Web connection health...')

			commit('SET_STATUS', 'connecting')
			commit('CLEAR_ERROR')

			try {
				// Initialize client if not exists
				if (!state.client) {
					await dispatch('initializeClient')
				}

				// Test connection with a simple ListUsers call
				const request = new ListUsersRequest()
				const response = await dispatch('promisifyGrpcCall', {
					method: state.client.listUsers,
					request: request
				})

				const userCount = response.getUsersList().length

				commit('SET_STATUS', 'connected')
				commit('SET_HEALTHY', true)
				commit('SET_MESSAGE', `✅ Connected via gRPC-Web with ${userCount} users`)
				commit('SET_LAST_CHECKED', new Date().toISOString())

				console.log(`✅ gRPC-Web health check successful - ${userCount} users found`)

				return {
					status: 'healthy',
					connection: 'grpc-web',
					users: userCount,
					url: state.grpcWebUrl
				}

			} catch (error) {
				console.error('❌ gRPC-Web health check failed:', error)

				commit('SET_STATUS', 'error')
				commit('SET_HEALTHY', false)
				commit('SET_MESSAGE', `❌ gRPC-Web server not reachable: ${error.message}`)
				commit('SET_ERROR', error.message)

				return {
					status: 'unhealthy',
					connection: 'grpc-web',
					error: error.message,
					url: state.grpcWebUrl
				}
			}
		},

		promisifyGrpcCall({ state }, { method, request, timeout = 10000 }) {
			return new Promise((resolve, reject) => {
				const call = method.call(state.client, request, {}, (err, response) => {
					if (err) {
						console.error('❌ gRPC-Web Error Details:', {
							code: err.code,
							message: err.message,
							details: err.details,
							metadata: err.metadata
						})

						// User-friendly error messages
						let userFriendlyMessage = err.message || 'gRPC-Web connection failed'

						if (err.code === 14 || err.message?.includes('UNAVAILABLE')) {
							userFriendlyMessage = `gRPC server not reachable at ${state.grpcWebUrl}. Is the backend running?`
						} else if (err.code === 12 || err.message?.includes('UNIMPLEMENTED')) {
							userFriendlyMessage = 'gRPC method not implemented on server'
						}

						reject(new Error(userFriendlyMessage))
					} else {
						resolve(response)
					}
				})

				// Timeout handling
				setTimeout(() => {
					if (call && typeof call.cancel === 'function') {
						call.cancel()
					}
					reject(new Error(`gRPC call timeout after ${timeout}ms to ${state.grpcWebUrl}`))
				}, timeout)
			})
		},

		async reconnect({ dispatch }) {
			console.log('🔄 Attempting to reconnect...')
			return await dispatch('checkHealth')
		}
	},

	getters: {
		isConnected(state) {
			return state.status === 'connected' && state.healthy
		},

		isConnecting(state) {
			return state.status === 'connecting'
		},

		connectionInfo(state) {
			return {
				status: state.status,
				healthy: state.healthy,
				message: state.message,
				connection: 'grpc-web',
				url: state.grpcWebUrl,
				lastChecked: state.lastChecked
			}
		},

		hasError(state) {
			return !!state.error
		},

		grpcClient(state) {
			return state.client
		}
	}
}