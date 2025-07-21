import {
	GetNotificationRequest,
	CreateNotificationRequest,
	ListNotificationsRequest
} from '../../../proto/notification_pb'

const mutations = {
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

	// Persistent notifications mutations
	SET_PERSISTENT_NOTIFICATIONS(state, notifications) {
		state.persistentNotifications = notifications || []
	},

	ADD_PERSISTENT_NOTIFICATION(state, notification) {
		console.log('=== MUTATION: ADD_PERSISTENT_NOTIFICATION ===')
		console.log('Adding notification:', notification)

		const exists = state.persistentNotifications.find(n => n.id === notification.id)
		if (!exists) {
			state.persistentNotifications.unshift(notification)
			console.log('✅ Notification added, new count:', state.persistentNotifications.length)
		} else {
			console.log('⚠️ Notification already exists, skipping')
		}
	},

	REMOVE_PERSISTENT(state, id) {
		state.persistentNotifications = state.persistentNotifications.filter(n => n.id !== id)
	},

	UPDATE_PERSISTENT_NOTIFICATION(state, updatedNotification) {
		const index = state.persistentNotifications.findIndex(n => n.id === updatedNotification.id)
		if (index !== -1) {
			state.persistentNotifications.splice(index, 1, updatedNotification)
		}
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

	// Loading states
	SET_LOADING(state, loading) {
		state.loading = loading
	},

	SET_ERROR(state, error) {
		state.error = error
	}
}

const actions = {
	// Initialize notification system
	async initialize({ dispatch }) {
		try {
			// Setup socket event listeners for notifications
			await dispatch('setupSocketListeners')

			// Socket connection is handled by main app initialization
			// Load existing persistent notifications via gRPC
			await dispatch('loadPersistentNotifications')

		} catch (error) {
			console.error('❌ Failed to initialize notification system:', error)
			throw error
		}
	},

	// Setup socket event listeners for notifications
	setupSocketListeners({ dispatch }) {
		// Listen for notification events from backend
		dispatch('socket/on', {
			event: 'notification',
			callback: (data) => dispatch('handleSocketNotification', data)
		}, { root: true })

		// Listen for notification updates
		dispatch('socket/on', {
			event: 'notification_updated',
			callback: (data) => dispatch('handleNotificationUpdate', data)
		}, { root: true })

		// Listen for notification deletions
		dispatch('socket/on', {
			event: 'notification_deleted',
			callback: (data) => dispatch('handleNotificationDeleted', data)
		}, { root: true })

		// Listen for connection status
		dispatch('socket/on', {
			event: 'connected',
			callback: () => dispatch('onSocketConnected')
		}, { root: true })
	},

	// Handle socket notifications from backend
	handleSocketNotification({ commit, dispatch }, notificationData) {

		// Always show as toast for immediate visibility
		dispatch('showToast', {
			message: notificationData.message,
			type: notificationData.type || 'info',
			duration: notificationData.type === 'error' ? 6000 : 4000
		})

		// If persistent, add to persistent store
		if (notificationData.persistent) {
			commit('ADD_PERSISTENT_NOTIFICATION', {
				...notificationData,
				read: false,
				createdAt: notificationData.createdAt || new Date().toISOString()
			})
		}
	},

	// Handle notification updates
	handleNotificationUpdate({ commit }, notificationData) {
		commit('UPDATE_PERSISTENT_NOTIFICATION', notificationData)
	},

	// Handle notification deletions
	handleNotificationDeleted({ commit }, { id }) {
		commit('REMOVE_PERSISTENT', id)
	},

	// Handle socket reconnection
	onSocketConnected({ dispatch }) {
		dispatch('loadPersistentNotifications')
	},

	// Load persistent notifications via gRPC
	async loadPersistentNotifications({ commit, dispatch, rootGetters }) {
		commit('SET_LOADING', true)
		commit('SET_ERROR', null)

		try {
			if (!rootGetters['connection/isConnected']) {
				console.warn('No gRPC connection available for notifications')
				return
			}

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
				userId: notification.getUserId() || null
			}))

			commit('SET_PERSISTENT_NOTIFICATIONS', notifications)

			return { notifications, count: notifications.length }

		} catch (error) {
			console.error('❌ Error loading persistent notifications:', error)
			commit('SET_ERROR', error.message)
		} finally {
			commit('SET_LOADING', false)
		}
	},

	// Create notification via gRPC
	async createNotification({ dispatch, rootGetters }, { message, type = 'info', persistent = false, userId = null }) {
		try {
			if (!rootGetters['connection/isConnected']) {
				console.warn('No gRPC connection - showing local toast only')
				return dispatch('showToast', { message, type })
			}

			const notificationClient = rootGetters['connection/notificationClient']
			if (!notificationClient) {
				throw new Error('Notification client not initialized')
			}

			const request = new CreateNotificationRequest()
			request.setMessage(message)
			request.setType(type)
			request.setPersistent(persistent)
			if (userId) {
				request.setUserId(userId)
			}

			const response = await dispatch('connection/promisifyGrpcCall', {
				method: notificationClient.createNotification,
				request: request
			}, { root: true })

			console.log('✅ Notification created via gRPC')
			return response

		} catch (error) {
			console.error('❌ Error creating notification:', error)
			// Fallback: show local toast
			return dispatch('showToast', { message, type, duration: 4000 })
		}
	},

	// Show toast notification
	showToast({ commit }, { message, type = 'info', duration = 2000, actions = [] }) {
		const toast = {
			message,
			type,
			duration,
			actions,
			createdAt: new Date().toISOString()
		}

		commit('ADD_TOAST', toast)

		if (duration > 0) {
			setTimeout(() => {
				commit('REMOVE_TOAST', toast.id)
			}, duration)
		}

		return toast.id
	},

	// Toast management
	dismissToast({ commit }, id) {
		commit('REMOVE_TOAST', id)
	},

	clearAllToasts({ commit }) {
		commit('CLEAR_ALL_TOASTS')
	},

	// Mark notification as read via gRPC
	async markAsRead({ commit, dispatch, rootGetters }, notificationId) {
		try {
			const notificationClient = rootGetters['connection/notificationClient']
			if (!notificationClient) {
				throw new Error('Notification client not available')
			}

			await dispatch('connection/promisifyGrpcCall', {
				method: notificationClient.markNotificationAsRead,
				request: { id: notificationId, userId: 1 } // TODO: get from auth
			}, { root: true })

			// Update local state immediately (optimistic update)
			commit('MARK_AS_READ', notificationId)

		} catch (error) {
			console.error('❌ Failed to mark notification as read:', error)
		}
	},

	// Remove persistent notification via gRPC
	async removePersistentNotification({ commit, dispatch, rootGetters }, id) {
		try {
			const notificationClient = rootGetters['connection/notificationClient']
			if (!notificationClient) {
				throw new Error('Notification client not available')
			}

			await dispatch('connection/promisifyGrpcCall', {
				method: notificationClient.deleteNotification,
				request: { id }
			}, { root: true })

		} catch (error) {
			console.error('❌ Failed to delete notification:', error)
		}
	},

	// Convenience methods for different notification types

	// Success notifications
	async success({ dispatch }, message, options = {}) {
		const toast = await dispatch('showToast', {
			message,
			type: 'success',
			duration: options.duration || 4000,
			actions: options.actions || []
		})

		if (options.persistent) {
			try {
				await dispatch('createNotification', {
					message,
					type: 'success',
					persistent: true,
					userId: options.userId
				})
			} catch (error) {
				console.error('❌ Failed to create persistent notification:', error)
			}
		}

		return toast
	},

	// Error notifications
	async error({ dispatch }, message, options = {}) {
		if (options.persistent) {
			return dispatch('createNotification', {
				message,
				type: 'error',
				persistent: true,
				...options
			})
		} else {
			return dispatch('showToast', {
				message,
				type: 'error',
				duration: options.duration || 5000,
				...options
			})
		}
	},

	// Warning notifications
	async warning({ dispatch }, message, options = {}) {
		if (options.persistent) {
			return dispatch('createNotification', {
				message,
				type: 'warning',
				persistent: true,
				...options
			})
		} else {
			return dispatch('showToast', {
				message,
				type: 'warning',
				duration: options.duration || 4000,
				...options
			})
		}
	},

	// Info notifications
	async info({ dispatch }, message, options = {}) {
		if (options.persistent) {
			return dispatch('createNotification', {
				message,
				type: 'info',
				persistent: true,
				...options
			})
		} else {
			return dispatch('showToast', {
				message,
				type: 'info',
				duration: options.duration || 3000,
				...options
			})
		}
	},

	// Send custom event via socket
	async sendSocketEvent({ dispatch }, { event, data }) {
		return dispatch('socket/emit', { event, data }, { root: true })
	}
}

const getters = {
	// Toast getters
	toasts(state) {
		return state.toasts
	},

	hasToasts(state) {
		return state.toasts.length > 0
	},

	// Persistent notification getters
	persistentNotifications(state) {
		return state.persistentNotifications.sort((a, b) => {
			// Unread first
			if (a.read !== b.read) {
				return a.read - b.read
			}
			// Then by creation time (newest first)
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

	// System status getters
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
