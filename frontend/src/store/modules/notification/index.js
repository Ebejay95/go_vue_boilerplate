// src/store/modules/notification/index.js
import {
	GetNotificationRequest,
	CreateNotificationRequest,
	ListNotificationsRequest,
	DeleteNotificationRequest,
	MarkNotificationAsReadRequest
} from '../../../proto/notification_pb'

const mutations = {
	// Toast mutations
	ADD_TOAST(state, toast) {
		const newToast = {
			id: Date.now() + Math.random(),
			...toast,
			createdAt: new Date().toISOString()
		}
		state.toasts.push(newToast)
	},

	REMOVE_TOAST(state, id) {
		state.toasts = state.toasts.filter(t => t.id !== id)
	},

	CLEAR_ALL_TOASTS(state) {
		state.toasts = []
	},

	// Persistent notifications (Bell) mutations
	SET_PERSISTENT_NOTIFICATIONS(state, notifications) {
		console.log('=== SET_PERSISTENT_NOTIFICATIONS ===')
		console.log('Setting notifications count:', notifications.length)
		state.persistentNotifications = notifications || []
	},

	ADD_PERSISTENT_NOTIFICATION(state, notification) {
		console.log('=== ADD_PERSISTENT_NOTIFICATION ===')
		console.log('Adding to bell:', notification)

		const exists = state.persistentNotifications.find(n => n.id === notification.id)
		if (!exists) {
			state.persistentNotifications.unshift(notification)
			console.log('✅ Added to bell, total:', state.persistentNotifications.length)
		}
	},

	REMOVE_PERSISTENT_NOTIFICATION(state, id) {
		state.persistentNotifications = state.persistentNotifications.filter(n => n.id !== id)
		console.log('🗑️ Removed from bell, remaining:', state.persistentNotifications.length)
	},

	MARK_AS_READ(state, id) {
		const notification = state.persistentNotifications.find(n => n.id === id)
		if (notification) {
			notification.read = true
		}
	},

	MARK_ALL_AS_READ(state) {
		state.persistentNotifications.forEach(n => n.read = true)
	},

	CLEAR_ALL_PERSISTENT(state) {
		state.persistentNotifications = []
	},

	// Loading states
	SET_LOADING(state, loading) {
		state.loading = loading
	},

	SET_ERROR(state, error) {
		state.error = error
	}
}

const actions = {
	// Core toast action
	showToast({ commit }, { message, type = 'info', duration = 3000, actions = [] }) {
		const toast = {
			id: Date.now() + Math.random(),
			message,
			type,
			duration,
			actions,
			createdAt: new Date().toISOString()
		}

		commit('ADD_TOAST', toast)
		return toast.id
	},

	// Remove toast action (called by component)
	dismissToast({ commit }, toastId) {
		commit('REMOVE_TOAST', toastId)
	},

	// Main notification creation method
	async createNotification({ commit, dispatch, rootGetters }, {
		message,
		type = 'info',
		persistent = false,
		duration = null,
		actions = [],
		userId = null
	}) {
		console.log('📨 Creating notification:', { message, type, persistent })

		// 1. Always show toast (for both persistent and non-persistent)
		if (!persistent) {
				const toastDuration = duration || (type === 'error' ? 4000 : type === 'warning' ? 3000 : 2000)
				dispatch('showToast', {
				message,
				type,
				duration: toastDuration,
				actions
			})
		}

		// 2. If persistent, save to backend (but DON'T add to bell manually)
		if (persistent) {
			try {
				if (!rootGetters['connection/isConnected']) {
					console.warn('No gRPC connection - adding to bell only')
					// Add to bell even without backend connection
					commit('ADD_PERSISTENT_NOTIFICATION', {
						id: Date.now() + Math.random(),
						message,
						type,
						read: false,
						persistent: false, // Couldn't save to backend
						createdAt: new Date().toISOString(),
						userId
					})
					return
				}

				const notificationClient = rootGetters['connection/notificationClient']
				if (!notificationClient) {
					throw new Error('Notification client not initialized')
				}

				const request = new CreateNotificationRequest()
				request.setMessage(message)
				request.setType(type)
				request.setPersistent(true)
				if (userId) {
					request.setUserId(userId)
				}

				const response = await dispatch('connection/promisifyGrpcCall', {
					method: notificationClient.createNotification,
					request: request
				}, { root: true })

				console.log('💾 Notification saved to backend')

				// ❌ ENTFERNT: Nicht mehr manuell zur Bell hinzufügen!
				// Die Socket-Listener werden das automatisch machen

				return response

			} catch (error) {
				console.error('❌ Error saving notification:', error)
				// Fallback: add to bell only
				commit('ADD_PERSISTENT_NOTIFICATION', {
					id: Date.now() + Math.random(),
					message,
					type,
					read: false,
					persistent: false,
					createdAt: new Date().toISOString(),
					userId
				})
			}
		}
		// If not persistent, only the toast was shown - no further action needed
	},

	// Convenience methods with your desired API
	async info({ dispatch }, { message, ...options }) {
		return dispatch('createNotification', {
			message,
			type: 'info',
			duration: 3000,
			...options
		})
	},

	async success({ dispatch }, { message, ...options }) {
		return dispatch('createNotification', {
			message,
			type: 'success',
			duration: 2000,
			...options
		})
	},

	async error({ dispatch }, { message, ...options }) {
		return dispatch('createNotification', {
			message,
			type: 'error',
			duration: 6000,
			...options
		})
	},

	async warning({ dispatch }, { message, ...options }) {
		return dispatch('createNotification', {
			message,
			type: 'warning',
			duration: 4000,
			...options
		})
	},

	// Initialization
	async initialize({ dispatch }) {
		console.log('🚀 Initializing notification system...')

		try {
			await dispatch('setupSocketListeners')
			await dispatch('loadPersistentNotifications')
			console.log('✅ Notification system initialized')
		} catch (error) {
			console.error('❌ Failed to initialize notification system:', error)
			throw error
		}
	},

	// Socket listeners
	setupSocketListeners({ dispatch }) {
		console.log('📋 Setting up notification socket listeners...')

		dispatch('socket/on', {
			event: 'notification',
			callback: (data) => dispatch('handleSocketNotification', data)
		}, { root: true })

		dispatch('socket/on', {
			event: 'notification_updated',
			callback: (data) => dispatch('handleNotificationUpdate', data)
		}, { root: true })

		dispatch('socket/on', {
			event: 'notification_deleted',
			callback: (data) => dispatch('handleNotificationDeleted', data)
		}, { root: true })

		dispatch('socket/on', {
			event: 'connected',
			callback: () => dispatch('onSocketConnected')
		}, { root: true })
	},

	// Socket event handlers
	handleSocketNotification({ commit, dispatch }, notificationData) {
		console.log('🔔 Socket notification received:', notificationData)

		// Always show as toast
		dispatch('showToast', {
			message: notificationData.message,
			type: notificationData.type || 'info',
			duration: notificationData.type === 'error' ? 6000 : 4000
		})

		// Add to bell if persistent
		if (notificationData.persistent) {
			commit('ADD_PERSISTENT_NOTIFICATION', {
				...notificationData,
				read: false,
				createdAt: notificationData.createdAt || new Date().toISOString()
			})
		}
	},

	handleNotificationUpdate({ commit }, notificationData) {
		console.log('🔄 Notification updated:', notificationData)
		// Handle updates if needed
	},

	handleNotificationDeleted({ commit }, { id }) {
		console.log('🗑️ Notification deleted:', id)
		commit('REMOVE_PERSISTENT_NOTIFICATION', id)
	},

	onSocketConnected({ dispatch }) {
		console.log('🔌 Socket reconnected, reloading notifications...')
		dispatch('loadPersistentNotifications')
	},

	// Load persistent notifications from backend
	async loadPersistentNotifications({ commit, dispatch, rootGetters }) {
		commit('SET_LOADING', true)
		commit('SET_ERROR', null)

		try {
			if (!rootGetters['connection/isConnected']) {
				console.warn('No gRPC connection available')
				return { notifications: [], count: 0 }
			}

			console.log('📋 Loading persistent notifications from backend...')

			const notificationClient = rootGetters['connection/notificationClient']
			if (!notificationClient) {
				throw new Error('Notification client not initialized')
			}

			const request = new ListNotificationsRequest()
			const response = await dispatch('connection/promisifyGrpcCall', {
				method: notificationClient.listNotifications,
				request: request
			}, { root: true })

			const notifications = response.getNotificationsList().map(notification => ({
				id: notification.getId(),
				message: notification.getMessage(),
				type: notification.getType(),
				read: notification.getRead(),
				createdAt: notification.getCreatedAt(),
				userId: notification.getUserId() || null,
				persistent: true
			}))

			commit('SET_PERSISTENT_NOTIFICATIONS', notifications)
			console.log(`✅ Loaded ${notifications.length} persistent notifications`)

			return { notifications, count: notifications.length }

		} catch (error) {
			console.error('❌ Error loading persistent notifications:', error)
			commit('SET_ERROR', error.message)
			return { notifications: [], count: 0 }
		} finally {
			commit('SET_LOADING', false)
		}
	},

	// Management actions
	async markAllAsRead({ commit, dispatch, state, rootGetters }) {
		try {
			const notificationClient = rootGetters['connection/notificationClient']
			if (notificationClient) {
				const unreadNotifications = state.persistentNotifications.filter(n => !n.read)
				for (const notification of unreadNotifications) {
					const request = new MarkNotificationAsReadRequest()
					request.setId(notification.id)
					request.setUserId(1)

					await dispatch('connection/promisifyGrpcCall', {
						method: notificationClient.markNotificationAsRead,
						request: request
					}, { root: true })
				}
			}

			commit('MARK_ALL_AS_READ')

		} catch (error) {
			console.error('❌ Failed to mark all as read:', error)
		}
	},

	async clearAllPersistent({ commit, dispatch, state, rootGetters }) {
		try {
			const notificationClient = rootGetters['connection/notificationClient']
			if (notificationClient) {
				for (const notification of state.persistentNotifications) {
					const request = new DeleteNotificationRequest()
					request.setId(notification.id)

					await dispatch('connection/promisifyGrpcCall', {
						method: notificationClient.deleteNotification,
						request: request
					}, { root: true })
				}
			}

			commit('SET_PERSISTENT_NOTIFICATIONS', [])

		} catch (error) {
			console.error('❌ Failed to clear all persistent notifications:', error)
		}
	},

	async removePersistentNotification({ commit, dispatch, rootGetters }, notificationId) {
		try {
			const notificationClient = rootGetters['connection/notificationClient']
			if (notificationClient) {
				const request = new DeleteNotificationRequest()
				request.setId(notificationId)

				await dispatch('connection/promisifyGrpcCall', {
					method: notificationClient.deleteNotification,
					request: request
				}, { root: true })
			}

			commit('REMOVE_PERSISTENT_NOTIFICATION', notificationId)

		} catch (error) {
			console.error('❌ Failed to remove persistent notification:', error)
		}
	},

	async markAsRead({ commit, dispatch, rootGetters }, notificationId) {
		try {
			const notificationClient = rootGetters['connection/notificationClient']
			if (notificationClient) {
				const request = new MarkNotificationAsReadRequest()
				request.setId(notificationId)
				request.setUserId(1)

				await dispatch('connection/promisifyGrpcCall', {
					method: notificationClient.markNotificationAsRead,
					request: request
				}, { root: true })
			}

			commit('MARK_AS_READ', notificationId)

		} catch (error) {
			console.error('❌ Failed to mark as read:', error)
		}
	}
}

const getters = {
	toasts(state) {
		return state.toasts
	},

	hasToasts(state) {
		return state.toasts.length > 0
	},

	// Bell notifications (persistent)
	persistentNotifications(state) {
		return state.persistentNotifications.sort((a, b) => {
			if (a.read !== b.read) return a.read - b.read
			return new Date(b.createdAt) - new Date(a.createdAt)
		})
	},

	unreadNotifications(state) {
		return state.persistentNotifications.filter(n => !n.read)
	},

	unreadCount(state) {
		return state.persistentNotifications.filter(n => !n.read).length
	},

	hasUnreadPersistent(state) {
		return state.persistentNotifications.some(n => !n.read)
	},

	hasPersistentNotifications(state) {
		return state.persistentNotifications.length > 0
	},

	isLoading(state) {
		return state.loading
	},

	error(state) {
		return state.error
	}
}

export default {
	namespaced: true,

	state() {
		return {
			toasts: [],
			persistentNotifications: [],
			loading: false,
			error: null
		}
	},

	mutations,
	actions,
	getters
}
