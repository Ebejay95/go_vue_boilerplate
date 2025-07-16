// src/store/modules/notifications/index.js
// Updated with gRPC integration for persistent notifications

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

	// NEW: Set persistent notifications from gRPC
	SET_PERSISTENT_NOTIFICATIONS(state, notifications) {
		state.persistentNotifications = notifications || []
	},

	ADD_PERSISTENT(state, notification) {
		// Check if notification already exists (avoid duplicates)
		const exists = state.persistentNotifications.find(n => n.id === notification.id)
		if (!exists) {
			state.persistentNotifications.unshift(notification)
		}
	},

	REMOVE_PERSISTENT(state, id) {
		state.persistentNotifications = state.persistentNotifications.filter(n => n.id !== id)
	},

	CLEAR_ALL_PERSISTENT(state) {
		state.persistentNotifications = []
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

	// NEW: Loading states for gRPC operations
	SET_LOADING(state, loading) {
		state.loading = loading
	},

	SET_ERROR(state, error) {
		state.error = error
	}
}

const actions = {
	// NEW: Load persistent notifications from gRPC
	async loadPersistentNotifications({ commit, dispatch, rootGetters }) {
		commit('SET_LOADING', true)
		commit('SET_ERROR', null)

		try {
			// Check if connection is available
			if (!rootGetters['connection/isConnected']) {
				console.warn('⚠️ No gRPC connection available for notifications')
				return
			}

			console.log('📋 Loading persistent notifications via gRPC...')

			// Get the notification client from connection module
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
				createdAt: new Date().toISOString(),
				read: false // You might want to add this to your proto
			}))

			commit('SET_PERSISTENT_NOTIFICATIONS', notifications)
			console.log(`✅ Loaded ${notifications.length} persistent notifications`)

			return { notifications, count: notifications.length }

		} catch (error) {
			console.error('❌ Error loading persistent notifications:', error)
			commit('SET_ERROR', error.message)
			throw error
		} finally {
			commit('SET_LOADING', false)
		}
	},

	// NEW: Create persistent notification via gRPC
	async createPersistentNotification({ commit, dispatch, rootGetters }, { message, type = 'info' }) {
		try {
			// Check if connection is available
			if (!rootGetters['connection/isConnected']) {
				console.warn('⚠️ No gRPC connection - storing notification locally only')
				return dispatch('addPersistentLocal', { message, type })
			}

			console.log('➕ Creating persistent notification via gRPC...')

			// Get the notification client from connection module
			const notificationClient = rootGetters['connection/notificationClient']
			if (!notificationClient) {
				throw new Error('Notification client not initialized')
			}

			const request = new CreateNotificationRequest()
			request.setMessage(message)
			request.setType(type)

			const response = await dispatch('connection/promisifyGrpcCall', {
				method: notificationClient.createNotification,
				request: request
			}, { root: true })

			const notification = response.getNotification()
			const createdNotification = {
				id: notification.getId(),
				message: notification.getMessage(),
				type: notification.getType(),
				createdAt: new Date().toISOString(),
				read: false
			}

			// Add to local state immediately
			commit('ADD_PERSISTENT', createdNotification)
			console.log(`✅ Persistent notification created: ${createdNotification.id}`)

			// Reload all persistent notifications to ensure sync
			setTimeout(() => {
				dispatch('loadPersistentNotifications').catch(error => {
					console.warn('⚠️ Could not reload notifications after creation:', error.message)
				})
			}, 100)

			return createdNotification

		} catch (error) {
			console.error('❌ Error creating persistent notification:', error)
			// Fallback to local storage
			return dispatch('addPersistentLocal', { message, type })
		}
	},

	// Base Toast Action (unchanged)
	showToast({ commit }, { message, type = 'info', duration = 2000, actions = [] }) {
		const toast = {
			id: Date.now() + Math.random(),
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

	// Local persistent action (fallback)
	addPersistentLocal({ commit }, { message, type = 'info', actions = [], priority = 'normal' }) {
		const notification = {
			id: Date.now() + Math.random(),
			message,
			type,
			actions,
			priority,
			createdAt: new Date().toISOString(),
			read: false
		}

		commit('ADD_PERSISTENT', notification)
		return notification.id
	},

	// Updated: Use gRPC for persistent notifications
	async addPersistent({ dispatch }, { message, type = 'info', actions = [], priority = 'normal' }) {
		try {
			return await dispatch('createPersistentNotification', { message, type })
		} catch (error) {
			console.warn('⚠️ Falling back to local persistent notification')
			return dispatch('addPersistentLocal', { message, type, actions, priority })
		}
	},

	// Toast Management (unchanged)
	dismissToast({ commit }, id) {
		commit('REMOVE_TOAST', id)
	},

	clearAllToasts({ commit }) {
		commit('CLEAR_ALL_TOASTS')
	},

	// Persistent Management
	removePersistentNotification({ commit }, id) {
		commit('REMOVE_PERSISTENT', id)
		// TODO: Also remove from gRPC backend if needed
	},

	clearAllPersistent({ commit }) {
		commit('CLEAR_ALL_PERSISTENT')
		// TODO: Also clear from gRPC backend if needed
	},

	markAsRead({ commit }, id) {
		commit('MARK_AS_READ', id)
		// TODO: Also update in gRPC backend if needed
	},

	markAllAsRead({ commit }) {
		commit('MARK_ALL_AS_READ')
		// TODO: Also update in gRPC backend if needed
	},

	// Updated: Unified Success Action
	dispatchSuccess({ dispatch }, message, options = {}) {
		console.log('🎉 Dispatching success:', message, options)

		if (options.persistent) {
			return dispatch('addPersistent', {
				message,
				type: 'success',
				...options
			})
		} else {
			return dispatch('showToast', {
				message,
				type: 'success',
				duration: 2000,
				...options
			})
		}
	},

	// Updated: Unified Error Action
	dispatchError({ dispatch }, message, options = {}) {
		console.log('❌ Dispatching error:', message, options)

		if (options.persistent) {
			return dispatch('addPersistent', {
				message,
				type: 'error',
				priority: 'high',
				...options
			})
		} else {
			return dispatch('showToast', {
				message,
				type: 'error',
				duration: 4000,
				...options
			})
		}
	},

	// Updated: Unified Warning Action
	dispatchWarning({ dispatch }, message, options = {}) {
		console.log('⚠️ Dispatching warning:', message, options)

		if (options.persistent) {
			return dispatch('addPersistent', {
				message,
				type: 'warning',
				priority: 'medium',
				...options
			})
		} else {
			return dispatch('showToast', {
				message,
				type: 'warning',
				duration: 3000,
				...options
			})
		}
	},

	// Updated: Unified Info Action
	dispatchInfo({ dispatch }, message, options = {}) {
		console.log('ℹ️ Dispatching info:', message, options)

		if (options.persistent) {
			return dispatch('addPersistent', {
				message,
				type: 'info',
				...options
			})
		} else {
			return dispatch('showToast', {
				message,
				type: 'info',
				duration: 2000,
				...options
			})
		}
	},

	// Updated: Critical notifications (always both)
	dispatchCritical({ dispatch }, message, options = {}) {
		console.log('🚨 Dispatching critical:', message, options)

		// Show toast immediately
		dispatch('showToast', {
			message,
			type: 'error',
			duration: 5000,
			...options
		})

		// Also add to persistent
		return dispatch('addPersistent', {
			message,
			type: 'error',
			priority: 'critical',
			...options
		})
	}
}

const getters = {
	toasts(state) {
		return state.toasts
	},

	hasToasts(state) {
		return state.toasts.length > 0
	},

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

	persistentByType: (state) => (type) => {
		return state.persistentNotifications.filter(n => n.type === type)
	},

	persistentByPriority: (state) => (priority) => {
		return state.persistentNotifications.filter(n => n.priority === priority)
	},

	criticalNotifications(state) {
		return state.persistentNotifications.filter(n => n.priority === 'critical' && !n.read)
	},

	// NEW: Loading and error states
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
