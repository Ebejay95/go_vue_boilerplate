<template>
	<header class="bg-white shadow-sm border-b border-gray-200 px-4 py-3">
	  <div class="flex items-center justify-between">
		<!-- Left: Logo/Brand -->
		<div class="flex items-center space-x-4">
		  <router-link
			to="/"
			class="text-xl font-semibold text-gray-900 hover:text-blue-600 transition-colors"
		  >
			Header
		  </router-link>
		  <navigation-component></navigation-component>
		</div>

		<!-- Right: Connection Status -->
		<div class="flex items-center space-x-3">
		  <!-- Connection Indicator Dot -->
		  <div class="flex items-center space-x-2">
			<div
			  :class="[
				'w-3 h-3 rounded-full transition-colors duration-300',
				connectionStatus.healthy ? 'bg-green-500' : 'bg-red-500'
			  ]"
			  :title="connectionStatus.message"
			></div>
			<span class="text-sm text-gray-600 hidden sm:inline">
			  {{ connectionStatus.message }}
			</span>
		  </div>

		  <!-- Reconnect Button (nur bei Problemen) -->
		  <button
			v-if="!connectionStatus.healthy"
			@click="handleReconnect"
			:disabled="isReconnecting"
			class="inline-flex items-center px-3 py-1 border border-gray-300 shadow-sm text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
		  >
			<svg
			  v-if="isReconnecting"
			  class="animate-spin -ml-1 mr-1 h-3 w-3"
			  xmlns="http://www.w3.org/2000/svg"
			  fill="none"
			  viewBox="0 0 24 24"
			>
			  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
			  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			{{ isReconnecting ? 'Reconnecting...' : 'Reconnect' }}
		  </button>

		  <!-- Debug Info (nur im Development) -->
		  <div
			v-if="isDevelopment && connectionStatus.url"
			class="text-xs text-gray-400 hidden md:block"
			:title="`gRPC-Web URL: ${connectionStatus.url}`"
		  >
			{{ connectionStatus.url.split('//')[1] }}
		  </div>
		</div>
	  </div>

	  <!-- Status Message (erscheint unter dem Header) -->
	  <div
		v-if="message"
		:class="[
		  'mt-3 p-3 rounded-md transition-all duration-300',
		  message.type === 'success'
			? 'bg-green-50 text-green-800 border border-green-200'
			: 'bg-red-50 text-red-800 border border-red-200'
		]"
	  >
		<div class="flex items-center justify-between">
		  <span class="text-sm">{{ message.text }}</span>
		  <button
			@click="clearMessage"
			class="text-gray-400 hover:text-gray-600 ml-2"
		  >
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
			</svg>
		  </button>
		</div>
	  </div>
	</header>
  </template>

  <script>
  import { mapState, mapGetters, mapActions } from 'vuex'

  export default {
	name: 'HeaderComponent',

	data() {
	  return {
		message: null,
		isReconnecting: false,
		isDevelopment: process.env.NODE_ENV === 'development'
	  }
	},

	computed: {
	  ...mapState('connection', ['status', 'healthy', 'message', 'grpcWebUrl', 'lastChecked', 'error']),
	  ...mapGetters('connection', ['isConnected', 'connectionInfo']),

	  connectionStatus() {
		return {
		  healthy: this.healthy,
		  message: this.getStatusMessage(),
		  url: this.grpcWebUrl,
		  lastChecked: this.lastChecked
		}
	  }
	},

	methods: {
	  ...mapActions('connection', ['checkHealth', 'reconnect']),

	  getStatusMessage() {
		if (this.healthy) {
		  return 'Connected'
		}

		if (this.status === 'connecting') {
		  return 'Connecting...'
		}

		if (this.error) {
		  return `Connection Error`
		}

		return 'Disconnected'
	  },

	  async handleReconnect() {
		this.isReconnecting = true
		this.clearMessage()

		try {
		  const result = await this.reconnect()

		  if (result.status === 'healthy') {
			this.showMessage('✅ Connection restored!', 'success')
		  } else {
			this.showMessage(`❌ Reconnection failed: ${result.error}`, 'error')
		  }
		} catch (error) {
		  console.error('❌ Reconnection error:', error)
		  this.showMessage(`❌ Reconnection failed: ${error.message}`, 'error')
		} finally {
		  this.isReconnecting = false
		}
	  },

	  showMessage(text, type = 'success') {
		this.message = { text, type }

		// Auto-hide success messages after 3 seconds
		if (type === 'success') {
		  setTimeout(() => {
			if (this.message && this.message.type === 'success') {
			  this.message = null
			}
		  }, 3000)
		}
	  },

	  clearMessage() {
		this.message = null
	  },

	  // Watch for connection changes
	  handleConnectionChange(newStatus, oldStatus) {
		if (oldStatus && newStatus !== oldStatus) {
		  if (newStatus === 'connected') {
			this.showMessage('✅ Connection established', 'success')
		  } else if (newStatus === 'error' && oldStatus === 'connected') {
			this.showMessage('⚠️ Connection lost', 'error')
		  }
		}
	  }
	},

	watch: {
	  // React to connection status changes
	  status: {
		handler: 'handleConnectionChange',
		immediate: false
	  },

	  // React to health changes
	  healthy(newHealthy, oldHealthy) {
		if (oldHealthy !== undefined && newHealthy !== oldHealthy) {
		  if (newHealthy) {
			this.showMessage('✅ Connection healthy', 'success')
		  } else {
			this.showMessage('⚠️ Connection unhealthy', 'error')
		  }
		}
	  }
	},

	mounted() {
	  console.log('🎯 Header component mounted - monitoring connection status')

	  // Optional: Periodic health check every 30 seconds
	  if (this.isDevelopment) {
		this.healthCheckInterval = setInterval(async () => {
		  try {
			await this.checkHealth()
		  } catch (error) {
			// Silent fail for background checks
			console.warn('⚠️ Background health check failed:', error.message)
		  }
		}, 30000)
	  }
	},

	beforeUnmount() {
	  if (this.healthCheckInterval) {
		clearInterval(this.healthCheckInterval)
	  }
	}
  }
  </script>

