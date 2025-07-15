<template>
	<base-section>
	  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">

		<base-card>
		  <h2 class="text-lg font-medium text-gray-900 mb-4">Neuen Benutzer erstellen</h2>

		  <form @submit.prevent="submitCreateUser" class="space-y-4">
			<div>
			  <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
			  <input
				type="text"
				id="name"
				v-model="newUser.name"
				required
				placeholder="Name eingeben"
				class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
			  />
			</div>

			<div>
			  <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
			  <input
				type="email"
				id="email"
				v-model="newUser.email"
				required
				placeholder="email@example.com"
				class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
			  />
			</div>

			<div>
			  <label for="age" class="block text-sm font-medium text-gray-700">Alter</label>
			  <input
				type="number"
				id="age"
				v-model.number="newUser.age"
				required
				min="1"
				max="150"
				placeholder="25"
				class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
			  />
			</div>

			<div>
			  <label for="role" class="block text-sm font-medium text-gray-700">Rolle</label>
			  <select
				id="role"
				v-model="newUser.role"
				class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
			  >
				<option value="user">Benutzer</option>
				<option value="admin">Administrator</option>
			  </select>
			</div>

			<button
			  type="submit"
			  :disabled="isCreating || !isConnected"
			  class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
			>
			  <svg v-if="isCreating" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			  </svg>
			  {{ isCreating ? 'Erstelle Benutzer...' : 'Benutzer erstellen' }}
			</button>
		  </form>
		</base-card>

		<!-- Users List -->
		<base-card>
		  <div class="flex justify-between items-center mb-4">
			<h2 class="text-lg font-medium text-gray-900">Alle Benutzer</h2>
			<button
			  @click="loadUsers"
			  :disabled="isLoading || !isConnected"
			  class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
			>
			  <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			  </svg>
			  {{ isLoading ? 'Lade...' : 'Aktualisieren' }}
			</button>
		  </div>

		  <div class="space-y-4">
			<!-- Loading State -->
			<div v-if="isLoading && !hasUsers" class="flex items-center justify-center py-8">
			  <div class="flex items-center">
				<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				<span class="text-gray-500">Lade Benutzer...</span>
			  </div>
			</div>

			<!-- Empty State -->
			<div v-else-if="!hasUsers && !isLoading" class="text-center py-8">
			  <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
				<path d="M34 40h10v-4a6 6 0 00-10.712-3.714M34 40H14m20 0v-4a9.971 9.971 0 00-.712-3.714M14 40H4v-4a6 6 0 0110.713-3.714M14 40v-4c0-1.313.253-2.566.713-3.714m0 0A10.003 10.003 0 0124 26c4.21 0 7.813 2.602 9.288 6.286M30 14a6 6 0 11-12 0 6 6 0 0112 0z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
			  </svg>
			  <h3 class="mt-2 text-sm font-medium text-gray-900">Keine Benutzer gefunden</h3>
			  <p class="mt-1 text-sm text-gray-500">Erstelle den ersten Benutzer um anzufangen.</p>
			</div>

			<!-- Users Table -->
			<div v-else class="overflow-hidden">
			  <table class="min-w-full divide-y divide-gray-200">
				<thead class="bg-gray-50">
				  <tr>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Benutzer</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Alter</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Rolle</th>
				  </tr>
				</thead>
				<tbody class="bg-white divide-y divide-gray-200">
				  <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50 transition-colors">
					<td class="px-6 py-4 whitespace-nowrap">
					  <div class="flex items-center">
						<div class="flex-shrink-0 h-10 w-10">
						  <div class="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center">
							<span class="text-white font-medium text-sm">{{ user.name.charAt(0).toUpperCase() }}</span>
						  </div>
						</div>
						<div class="ml-4">
						  <div class="text-sm font-medium text-gray-900">{{ user.name }}</div>
						  <div class="text-sm text-gray-500">ID: {{ user.id }}</div>
						</div>
					  </div>
					</td>
					<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.email }}</td>
					<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.age }}</td>
					<td class="px-6 py-4 whitespace-nowrap">
					  <span :class="[
						'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
						user.role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-green-100 text-green-800'
					  ]">
						{{ user.role === 'admin' ? 'Administrator' : 'Benutzer' }}
					  </span>
					</td>
				  </tr>
				</tbody>
			  </table>
			</div>

			<!-- Last Updated -->
			<div v-if="lastUpdated" class="text-center text-xs text-gray-500 mt-4">
			  Zuletzt aktualisiert: {{ lastUpdated }}
			</div>
		  </div>
		</base-card>
	  </div>

	  <!-- Debug Info -->
	  <div class="mt-8 bg-gray-100 p-4 rounded-lg">
		<h3 class="text-sm font-medium text-gray-900 mb-2">Debug-Informationen</h3>
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 text-xs text-gray-600">
		  <div>
			<strong>gRPC-Web URL:</strong> {{ connectionInfo.url }}
		  </div>
		  <div>
			<strong>Verbindungsstatus:</strong> {{ connectionInfo.status }}
		  </div>
		  <div>
			<strong>Verbindungshealth:</strong> {{ connectionInfo.healthy ? 'Gesund' : 'Ungesund' }}
		  </div>
		  <div>
			<strong>Benutzeranzahl:</strong> {{ userCount }}
		  </div>
		  <div>
			<strong>Zuletzt geprüft:</strong> {{ connectionInfo.lastChecked ? new Date(connectionInfo.lastChecked).toLocaleTimeString() : 'Nie' }}
		  </div>
		  <div>
			<strong>Fehler:</strong> {{ error || 'Keine' }}
		  </div>
		</div>
	  </div>
	</base-section>
  </template>

  <script>
  import { mapState, mapGetters, mapActions } from 'vuex'

  export default {
	name: 'WelcomePage',

	data() {
	  return {
		lastUpdated: null,
		newUser: {
		  name: '',
		  email: '',
		  age: '',
		  role: 'user'
		}
	  }
	},

	computed: {
	  ...mapState('users', ['users', 'loading', 'creating', 'error']),
	  ...mapState('connection', ['grpcWebUrl']),
	  ...mapGetters('users', ['hasUsers', 'userCount']),
	  ...mapGetters('connection', ['isConnected', 'connectionInfo']),

	  isLoading() {
		return this.loading
	  },

	  isCreating() {
		return this.creating
	  }
	},

	methods: {
	  ...mapActions('users', ['fetchUsers', 'createUser', 'clearError']),
	  ...mapActions('connection', ['checkHealth', 'reconnect']),
	  ...mapActions(['initializeApp']),
	  ...mapActions('notifications', ['success', 'error', 'warning', 'info']),

	  async loadUsers() {
		if (!this.isConnected) {
		  this.warning('⚠️ Backend-Verbindung nicht verfügbar', {
			actions: [
			  {
				label: 'Verbindung prüfen',
				handler: () => this.checkConnection()
			  }
			]
		  })
		  return
		}

		try {
		  this.info('Lade Benutzer...', { duration: 2000 })
		  const result = await this.fetchUsers()
		  this.lastUpdated = new Date().toLocaleTimeString()

		  if (!this.hasUsers) {
			this.info('📋 Keine Benutzer in der Datenbank gefunden')
		  } else {
			this.success(`📋 ${result.count} Benutzer erfolgreich geladen`)
		  }
		} catch (error) {
		  console.error('❌ Error loading users:', error)
		  this.error(`❌ Fehler beim Laden: ${error.message}`, {
			actions: [
			  {
				label: 'Erneut versuchen',
				handler: () => this.loadUsers()
			  }
			]
		  })
		}
	  },

	  async submitCreateUser() {
		if (!this.isConnected) {
		  this.warning('⚠️ Backend-Verbindung nicht verfügbar')
		  return
		}

		try {
		  const response = await this.createUser(this.newUser)

		  // Reset form
		  const userName = this.newUser.name
		  this.newUser = {
			name: '',
			email: '',
			age: '',
			role: 'user'
		  }

		  // Update timestamp
		  this.lastUpdated = new Date().toLocaleTimeString()

		  this.success(`🎉 Benutzer "${userName}" erfolgreich erstellt!`, {
			actions: [
			  {
				label: 'Liste aktualisieren',
				handler: () => this.loadUsers()
			  },
			  {
				label: 'Neuen erstellen',
				handler: () => this.$refs.nameInput?.focus()
			  }
			]
		  })
		} catch (error) {
		  console.error('❌ Error creating user:', error)
		  this.error(`❌ Fehler beim Erstellen: ${error.message}`, {
			persistent: true,
			actions: [
			  {
				label: 'Erneut versuchen',
				handler: () => this.submitCreateUser()
			  }
			]
		  })
		}
	  },

	  async checkConnection() {
		try {
		  this.info('Prüfe Verbindung...')
		  await this.checkHealth()
		  if (this.isConnected) {
			this.success('✅ Verbindung erfolgreich!')
			// Automatisch Benutzer laden nach erfolgreicher Verbindung
			setTimeout(() => this.loadUsers(), 500)
		  }
		} catch (error) {
		  this.error(`❌ Verbindungsfehler: ${error.message}`, {
			actions: [
			  {
				label: 'Erneut versuchen',
				handler: () => this.checkConnection()
			  }
			]
		  })
		}
	  }
	},

	async mounted() {
	  console.log('🚀 Welcome page mounted')

	  try {
		// Initialize the entire app (connection + health check)
		this.info('Initialisiere Anwendung...', { duration: 3000 })
		const initialized = await this.initializeApp()

		if (initialized && this.isConnected) {
		  // Load initial data
		  await this.loadUsers()
		  this.success('✅ Anwendung bereit!', { duration: 3000 })
		} else {
		  console.warn('⚠️ Application initialized with connection issues')
		  this.warning('⚠️ Verbindungsprobleme erkannt. Einige Funktionen sind möglicherweise nicht verfügbar.', {
			persistent: true,
			actions: [
			  {
				label: 'Verbindung prüfen',
				handler: () => this.checkConnection()
			  }
			]
		  })
		}
	  } catch (error) {
		console.error('❌ Failed to initialize application:', error)
		this.error(`❌ Initialisierung fehlgeschlagen: ${error.message}`, {
		  persistent: true,
		  actions: [
			{
			  label: 'Erneut versuchen',
			  handler: () => this.$router.go(0) // Reload page
			}
		  ]
		})
	  }
	}
  }
  </script>
