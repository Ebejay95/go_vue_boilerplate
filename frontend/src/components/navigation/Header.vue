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
	  ...mapActions('notifications', {
		notifySuccess: 'success',
		notifyError: 'error',
		notifyWarning: 'warning',
		notifyInfo: 'info'
	  }),

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
		  this.notifyInfo('Verbindung wird wiederhergestellt...')
		  const result = await this.reconnect()

		  if (result.status === 'healthy') {
			this.notifySuccess('Verbindung wiederhergestellt!')
		  } else {
			this.notifyError(`❌ Wiederverbindung fehlgeschlagen: ${result.error}`, {
			  actions: [
				{
				  label: 'Erneut versuchen',
				  handler: () => this.handleReconnect()
				}
			  ]
			})
		  }
		} catch (error) {
		  console.error('Reconnection error:', error)
		  this.notifyError(`❌ Verbindungsfehler: ${error.message}`, {
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
			this.notifySuccess('🔗 Verbindung hergestellt')
		  } else if (newStatus === 'error' && oldStatus === 'connected') {
			this.notifyError('🔌 Verbindung verloren', {
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
			this.notifySuccess('Verbindung stabil')
		  } else {
			this.notifyWarning('⚠️ Verbindung instabil')
		  }
		}
	  }
	}
  }
  </script>
