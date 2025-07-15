const mutations = {
	ADD_NOTIFICATION(state, notification) {
	  state.notifications.push({
		id: Date.now() + Math.random(), // Unique ID
		...notification,
		createdAt: new Date().toISOString()
	  })
	},

	REMOVE_NOTIFICATION(state, id) {
	  state.notifications = state.notifications.filter(n => n.id !== id)
	},

	CLEAR_ALL_NOTIFICATIONS(state) {
	  state.notifications = []
	},

	MARK_AS_READ(state, id) {
	  const notification = state.notifications.find(n => n.id === id)
	  if (notification) {
		notification.read = true
	  }
	},

	MARK_ALL_AS_READ(state) {
	  state.notifications.forEach(n => n.read = true)
	}
  }

  const actions = {
	// Hauptmethode zum Hinzufügen von Notifications
	addNotification({ commit }, { message, type = 'info', duration = 5000, persistent = false, actions = [] }) {
	  const notification = {
		message,
		type, // 'success', 'error', 'warning', 'info'
		duration,
		persistent, // Wenn true, wird nicht automatisch entfernt
		actions, // Array von Action-Objekten { label, handler }
		read: false
	  }

	  commit('ADD_NOTIFICATION', notification)

	  // Auto-remove nach duration (außer persistent)
	  if (!persistent && duration > 0) {
		setTimeout(() => {
		  commit('REMOVE_NOTIFICATION', notification.id)
		}, duration)
	  }

	  return notification.id
	},

	// Convenience-Methoden für verschiedene Typen
	success({ dispatch }, message, options = {}) {
	  return dispatch('addNotification', {
		message,
		type: 'success',
		duration: 4000,
		...options
	  })
	},

	error({ dispatch }, message, options = {}) {
	  return dispatch('addNotification', {
		message,
		type: 'error',
		duration: 8000, // Errors bleiben länger
		...options
	  })
	},

	warning({ dispatch }, message, options = {}) {
	  return dispatch('addNotification', {
		message,
		type: 'warning',
		duration: 6000,
		...options
	  })
	},

	info({ dispatch }, message, options = {}) {
	  return dispatch('addNotification', {
		message,
		type: 'info',
		duration: 5000,
		...options
	  })
	},

	// Notification entfernen
	removeNotification({ commit }, id) {
	  commit('REMOVE_NOTIFICATION', id)
	},

	// Alle Notifications löschen
	clearAllNotifications({ commit }) {
	  commit('CLEAR_ALL_NOTIFICATIONS')
	},

	// Als gelesen markieren
	markAsRead({ commit }, id) {
	  commit('MARK_AS_READ', id)
	},

	markAllAsRead({ commit }) {
	  commit('MARK_ALL_AS_READ')
	}
  }

  const getters = {
	allNotifications(state) {
	  return state.notifications
	},

	unreadNotifications(state) {
	  return state.notifications.filter(n => !n.read)
	},

	unreadCount(state) {
	  return state.notifications.filter(n => !n.read).length
	},

	notificationsByType: (state) => (type) => {
	  return state.notifications.filter(n => n.type === type)
	},

	latestNotifications(state) {
	  return state.notifications
		.slice()
		.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
		.slice(0, 5)
	},

	hasNotifications(state) {
	  return state.notifications.length > 0
	},

	hasUnreadNotifications(state) {
	  return state.notifications.some(n => !n.read)
	}
  }

  export default {
	namespaced: true,

	state() {
	  return {
		notifications: []
	  }
	},

	mutations,
	actions,
	getters
}
