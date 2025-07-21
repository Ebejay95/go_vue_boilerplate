const state = () => ({
	socket: null,
	isConnected: false,
	isConnecting: false,
	reconnectAttempts: 0,
	maxReconnectAttempts: 5,
	reconnectInterval: 5000,

	// Event listeners registry
	eventListeners: new Map(),

	// Message queue for when disconnected
	messageQueue: []
})

const mutations = {
	SET_SOCKET(state, socket) {
		state.socket = socket
	},

	SET_CONNECTED(state, status) {
		state.isConnected = status
	},

	SET_CONNECTING(state, status) {
		state.isConnecting = status
	},

	INCREMENT_RECONNECT_ATTEMPTS(state) {
		state.reconnectAttempts++
	},

	RESET_RECONNECT_ATTEMPTS(state) {
		state.reconnectAttempts = 0
	},

	ADD_EVENT_LISTENER(state, { eventType, callback }) {
		if (!state.eventListeners.has(eventType)) {
			state.eventListeners.set(eventType, [])
		}
		state.eventListeners.get(eventType).push(callback)
	},

	REMOVE_EVENT_LISTENER(state, { eventType, callback }) {
		if (state.eventListeners.has(eventType)) {
			const listeners = state.eventListeners.get(eventType)
			const index = listeners.indexOf(callback)
			if (index > -1) {
				listeners.splice(index, 1)
			}
		}
	},

	CLEAR_EVENT_LISTENERS(state) {
		state.eventListeners.clear()
	},

	ADD_TO_MESSAGE_QUEUE(state, message) {
		state.messageQueue.push(message)
	},

	CLEAR_MESSAGE_QUEUE(state) {
		state.messageQueue = []
	}
}

const actions = {
	// Connect to WebSocket
	connect({ commit, dispatch, state }) {
		if (state.isConnected || state.isConnecting) {
			return Promise.resolve()
		}

		return new Promise((resolve, reject) => {
			commit('SET_CONNECTING', true)

			// Use environment variable or default
			const wsUrl = `${process.env.VUE_APP_WS_URL || 'ws://localhost:8082'}/notifications`

			const socket = new WebSocket(wsUrl)

			socket.onopen = (event) => {
				commit('SET_SOCKET', socket)
				commit('SET_CONNECTED', true)
				commit('SET_CONNECTING', false)
				commit('RESET_RECONNECT_ATTEMPTS')

				// Process queued messages
				dispatch('processMessageQueue')

				// Start heartbeat
				dispatch('startHeartbeat')

				resolve(socket)
			}

			socket.onmessage = (event) => {
				try {
					const message = JSON.parse(event.data)
					dispatch('handleMessage', message)
				} catch (error) {
					console.error('Socket message parsing error:', error)
				}
			}

			socket.onclose = (event) => {
				commit('SET_CONNECTED', false)
				commit('SET_CONNECTING', false)
				commit('SET_SOCKET', null)

				// Auto-reconnect if not a clean close
				if (!event.wasClean && state.reconnectAttempts < state.maxReconnectAttempts) {
					dispatch('reconnect')
				}
			}

			socket.onerror = (error) => {
				console.error('💥 Socket error:', error)
				commit('SET_CONNECTING', false)
				reject(error)
			}
		})
	},

	// Disconnect from WebSocket
	disconnect({ commit, state }) {
		if (state.socket) {
			state.socket.close(1000, 'Manual disconnect')
			commit('SET_SOCKET', null)
			commit('SET_CONNECTED', false)
			commit('CLEAR_EVENT_LISTENERS')
		}
	},

	// Reconnect logic
	reconnect({ commit, dispatch, state }) {
		if (state.reconnectAttempts >= state.maxReconnectAttempts) {
			console.error('❌ Max reconnection attempts reached')
			return
		}

		commit('INCREMENT_RECONNECT_ATTEMPTS')

		setTimeout(() => {
			dispatch('connect')
		}, state.reconnectInterval)
	},

	// Send message
	send({ state, commit }, message) {
		if (state.isConnected && state.socket) {
			const messageStr = typeof message === 'string' ? message : JSON.stringify(message)
			state.socket.send(messageStr)
			return true
		} else {
			console.warn('Socket not connected, queueing message')
			commit('ADD_TO_MESSAGE_QUEUE', message)
			return false
		}
	},

	// Process queued messages
	processMessageQueue({ state, commit }) {
		if (state.messageQueue.length > 0 && state.isConnected) {
			state.messageQueue.forEach(message => {
				const messageStr = typeof message === 'string' ? message : JSON.stringify(message)
				state.socket.send(messageStr)
			})
			commit('CLEAR_MESSAGE_QUEUE')
		}
	},

	// Handle incoming messages
	handleMessage({ state }, message) {

		const { event, data } = message

		// Special system events
		if (event === 'pong') {
			// Heartbeat response
			return
		}

		// Trigger registered event listeners
		if (state.eventListeners.has(event)) {
			const listeners = state.eventListeners.get(event)
			listeners.forEach(callback => {
				try {
					callback(data, message)
				} catch (error) {
					console.error(`Error in event listener for '${event}':`, error)
				}
			})
		} else {
			console.warn(`No listeners registered for event: ${event}`)
		}
	},

	// Register event listener
	on({ commit }, { event, callback }) {
		commit('ADD_EVENT_LISTENER', { eventType: event, callback })
	},

	// Remove event listener
	off({ commit }, { event, callback }) {
		commit('REMOVE_EVENT_LISTENER', { eventType: event, callback })
	},

	// Emit event to server
	emit({ dispatch }, { event, data }) {
		return dispatch('send', { event, data })
	},

	// Start heartbeat
	startHeartbeat({ dispatch, state }) {
		if (state.heartbeatInterval) {
			clearInterval(state.heartbeatInterval)
		}

		state.heartbeatInterval = setInterval(() => {
			if (state.isConnected) {
				dispatch('emit', { event: 'ping' })
			}
		}, 30000) // Ping every 30 seconds
	},

	// Stop heartbeat
	stopHeartbeat({ state }) {
		if (state.heartbeatInterval) {
			clearInterval(state.heartbeatInterval)
			state.heartbeatInterval = null
		}
	}
}

const getters = {
	isConnected: (state) => state.isConnected,
	isConnecting: (state) => state.isConnecting,

	connectionStatus: (state) => {
		if (state.isConnected) return 'connected'
		if (state.isConnecting) return 'connecting'
		return 'disconnected'
	},

	hasQueuedMessages: (state) => state.messageQueue.length > 0,
	queuedMessageCount: (state) => state.messageQueue.length
}

export default {
	namespaced: true,
	state,
	mutations,
	actions,
	getters
}
