// src/store/modules/notifications/index.js
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

	ADD_PERSISTENT(state, notification) {
		state.persistentNotifications.unshift(notification)
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
	}
}

const actions = {
	// Base Toast Action
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

	// Base Persistent Action
	addPersistent({ commit }, { message, type = 'info', actions = [], priority = 'normal' }) {
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

	// Toast Management
	dismissToast({ commit }, id) {
		commit('REMOVE_TOAST', id)
	},

	clearAllToasts({ commit }) {
		commit('CLEAR_ALL_TOASTS')
	},

	// Persistent Management
	removePersistentNotification({ commit }, id) {
		commit('REMOVE_PERSISTENT', id)
	},

	clearAllPersistent({ commit }) {
		commit('CLEAR_ALL_PERSISTENT')
	},

	markAsRead({ commit }, id) {
		commit('MARK_AS_READ', id)
	},

	markAllAsRead({ commit }) {
		commit('MARK_ALL_AS_READ')
	},

	// Unified Success Action
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

	// Unified Error Action
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

	// Unified Warning Action
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

	// Unified Info Action
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

	// Critical notifications (always both)
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
	}
}

export default {
	namespaced: true,

	state() {
		return {
			toasts: [],
			persistentNotifications: []
		}
	},

	mutations,
	actions,
	getters
}