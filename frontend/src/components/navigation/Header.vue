<template>
	<header class="bg-white shadow-sm border-b">
	  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex justify-between items-center h-16">

		  <!-- Logo/Brand -->
		  <div class="flex items-center space-x-8">
			<base-button link to="/">
				MyApp
			</base-button>
			<navigation-component></navigation-component>
		  </div>

		  <!-- Right Side -->
		  <div class="flex items-center space-x-4">

			<!-- Connection Status -->
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

			<!-- Notification Bell mit Badge -->
			<notification-list :show-badge="true" :show-panel="true"></notification-list>

			<!-- Settings/Profile (optional) -->
			<button class="text-gray-500 hover:text-gray-700 transition-colors">
			  <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
					  d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
			  </svg>
			</button>
		  </div>
		</div>
	  </div>

	  <!-- Notification Toast Area wird durch notification-list verwaltet -->
	</header>
  </template>

  <script>
  import { mapState, mapGetters, mapActions } from 'vuex'

  export default {
	name: 'HeaderComponent',

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
	  ...mapActions('notifications', ['success', 'error', 'warning', 'info']),

	  getStatusMessage() {
		if (this.healthy) {
		  return 'Verbunden'
		}

		if (this.status === 'connecting') {
		  return 'Verbindung...'
		}

		if (this.error) {
		  return `Verbindungsfehler`
		}

		return 'Getrennt'
	  },

	  async handleReconnect() {
		try {
		  this.info('Verbindung wird wiederhergestellt...')
		  const result = await this.reconnect()

		  if (result.status === 'healthy') {
			this.success('✅ Verbindung wiederhergestellt!')
		  } else {
			this.error(`❌ Wiederverbindung fehlgeschlagen: ${result.error}`, {
			  actions: [
				{
				  label: 'Erneut versuchen',
				  handler: () => this.handleReconnect()
				}
			  ]
			})
		  }
		} catch (error) {
		  console.error('❌ Reconnection error:', error)
		  this.error(`❌ Verbindungsfehler: ${error.message}`, {
			persistent: true,
			actions: [
			  {
				label: 'Erneut versuchen',
				handler: () => this.handleReconnect()
			  }
			]
		  })
		}
	  },

	  // Watch für Connection-Änderungen
	  handleConnectionChange(newStatus, oldStatus) {
		if (oldStatus && newStatus !== oldStatus) {
		  if (newStatus === 'connected') {
			this.success('🔗 Verbindung hergestellt')
		  } else if (newStatus === 'error' && oldStatus === 'connected') {
			this.error('🔌 Verbindung verloren', {
			  persistent: true,
			  actions: [
				{
				  label: 'Neu verbinden',
				  handler: () => this.handleReconnect()
				}
			  ]
			})
		  }
		}
	  }
	},

	watch: {
	  // React auf Connection-Status-Änderungen
	  status: {
		handler: 'handleConnectionChange',
		immediate: false
	  },

	  // React auf Health-Änderungen
	  healthy(newHealthy, oldHealthy) {
		if (oldHealthy !== undefined && newHealthy !== oldHealthy) {
		  if (newHealthy) {
			this.success('✅ Verbindung stabil')
		  } else {
			this.warning('⚠️ Verbindung instabil')
		  }
		}
	  }
	},

	mounted() {
	  console.log('🎯 Header component mounted - monitoring connection status')

	  // Optional: Periodische Health-Checks
	  if (process.env.NODE_ENV === 'development') {
		this.healthCheckInterval = setInterval(async () => {
		  try {
			await this.checkHealth()
		  } catch (error) {
			// Silent fail für Background-Checks
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
