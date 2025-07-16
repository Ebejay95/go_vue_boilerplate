<!-- src/pages/general/Welcome.vue - Complete Corrected Version -->
<template>
	<base-section>
	  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">

		<!-- Dynamic Form Card -->
		<base-card>
		  <h2 class="text-lg font-medium text-gray-900 mb-4">Neuen Benutzer erstellen</h2>

		  <Form
			:protoMessage="CreateUserRequest"
			:fieldConfig="userFieldConfig"
			:customValidators="userValidators"
			submitText="Benutzer erstellen"
			submitLoadingText="Erstelle Benutzer..."
			@submit="handleCreateUser"
			@error="handleFormError"
		  />
		</base-card>

		<!-- Users List Card -->
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

		<!-- Persistent Notifications Card -->
		<base-card>
			<h3 class="text-lg font-medium text-gray-900 mb-4">Persistent Notifications</h3>

			<div v-if="persistentNotifications.length === 0" class="text-center py-8">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
				</svg>
				<p class="mt-2 text-sm text-gray-500">Keine persistenten Benachrichtigungen</p>
			</div>

			<ul v-else class="space-y-2">
				<li v-for="pNote in persistentNotifications" :key="pNote.id"
					class="p-3 bg-gray-50 rounded-lg border border-gray-200">
					<div class="flex justify-between items-start">
						<div class="flex-1">
							<p class="text-sm font-medium text-gray-900">{{ pNote.message }}</p>
							<p class="text-xs text-gray-500 mt-1">
								{{ pNote.type }} - {{ formatTime(pNote.createdAt) }}
							</p>
						</div>
						<div class="flex items-center space-x-2">
							<span v-if="!pNote.read" class="w-2 h-2 bg-blue-500 rounded-full"></span>
							<span :class="[
								'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
								pNote.type === 'success' ? 'bg-green-100 text-green-800' :
								pNote.type === 'error' ? 'bg-red-100 text-red-800' :
								pNote.type === 'warning' ? 'bg-yellow-100 text-yellow-800' :
								'bg-blue-100 text-blue-800'
							]">
								{{ pNote.type }}
							</span>
						</div>
					</div>
				</li>
			</ul>
		</base-card>
	  </div>

	  <!-- Debug and Test Buttons -->
	  <div class="mt-8 bg-gray-50 rounded-lg p-4">
		<h3 class="text-lg font-medium text-gray-900 mb-4">Debug & Test Controls</h3>
		<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
		  <button @click="testPersistentNotification" class="btn primary text-sm">
			Test Persistent
		  </button>
		  <button @click="testToastNotification" class="btn secondary text-sm">
			Test Toast
		  </button>
		  <button @click="refreshNotifications" class="btn flat text-sm">
			🔄 Refresh
		  </button>
		  <button @click="testNotificationConnection" class="btn flat text-sm">
			🧪 Test gRPC
		  </button>
		  <button @click="debugNotifications" class="btn flat text-sm">
			🔍 Debug
		  </button>
		  <button @click="testBackendServices" class="btn flat text-sm">
			🛠️ Test Backend
		  </button>
		</div>
	  </div>
	</base-section>
</template>

<script>
import { mapState, mapGetters, mapActions } from 'vuex'
import Form from '../../components/base/Form.vue'
import { CreateUserRequest } from '../../proto/user_pb'

export default {
	name: 'WelcomePage',

	components: {
		Form
	},

	data() {
		return {
			lastUpdated: null,
			CreateUserRequest, // Proto message class

			userFieldConfig: {
				name: {
					label: 'Name',
					placeholder: 'Vollständiger Name',
					required: true,
					help: 'Vor- und Nachname eingeben'
				},

				email: {
					label: 'Email',
					inputType: 'email',
					placeholder: 'benutzer@example.com',
					required: true,
					help: 'Gültige E-Mail-Adresse eingeben'
				},

				age: {
					label: 'Alter',
					inputType: 'range',
					min: 18,
					max: 100,
					step: 1,
					required: true,
					help: 'Schieberegler zum Auswählen des Alters verwenden'
				},

				role: {
					label: 'Rolle',
					inputType: 'select',
					options: [
						{ value: 'user', label: 'Benutzer' },
						{ value: 'admin', label: 'Administrator' }
					],
					required: true,
					placeholder: 'Rolle auswählen',
					help: 'Entsprechende Benutzerrolle auswählen'
				}
			},

			userValidators: {
				name: (value) => {
					if (value.length < 2) {
						return 'Name muss mindestens 2 Zeichen lang sein'
					}
					if (!/^[a-zA-ZäöüÄÖÜß\s]+$/.test(value)) {
						return 'Name darf nur Buchstaben und Leerzeichen enthalten'
					}
					return null
				},

				email: (value) => {
					const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
					if (!emailRegex.test(value)) {
						return 'Bitte gültige E-Mail-Adresse eingeben'
					}
					return null
				},

				age: (value) => {
					if (value < 18) {
						return 'Mindestalter ist 18 Jahre'
					}
					if (value > 100) {
						return 'Alter darf nicht über 100 Jahre sein'
					}
					return null
				},

				role: (value) => {
					const validRoles = ['user', 'admin', 'moderator']
					if (!validRoles.includes(value)) {
						return 'Ungültige Rolle ausgewählt'
					}
					return null
				}
			}
		}
	},

	computed: {
		...mapState('users', ['users', 'loading', 'creating', 'error']),
		...mapState('connection', ['grpcWebUrl']),
		...mapGetters('users', ['hasUsers', 'userCount']),
		...mapGetters('connection', ['isConnected', 'connectionInfo']),
		...mapGetters('notifications', [
			'persistentNotifications',
			'unreadCount',
			'hasUnreadPersistent'
		]),

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
		...mapActions('notifications', [
			'dispatchSuccess',
			'dispatchError',
			'dispatchWarning',
			'dispatchInfo',
			'loadPersistentNotifications'
		]),

		async debugNotifications() {
			console.log('🔍 Debug Notifications:')
			console.log('Store State:', {
				persistent: this.persistentNotifications,
				persistentLength: this.persistentNotifications?.length,
				unreadCount: this.unreadCount,
				hasUnread: this.hasUnreadPersistent,
				isConnected: this.$store.getters['connection/isConnected'],
				hasNotificationClient: !!this.$store.getters['connection/notificationClient']
			})

			// Try to load fresh data from gRPC
			try {
				console.log('🔄 Loading fresh notifications from gRPC...')
				const result = await this.loadPersistentNotifications()
				console.log('Fresh notifications result:', result)
			} catch (error) {
				console.error('❌ Failed to load fresh notifications:', error)
			}
		},

		async testPersistentNotification() {
			console.log('🧪 Testing persistent notification...')

			try {
				// Check connection first
				if (!this.$store.getters['connection/isConnected']) {
					console.error('❌ Not connected to gRPC')
					this.dispatchError('Not connected to backend')
					return
				}

				console.log('✅ Connected, creating notification...')

				// Create notification with correct action name
				const result = await this.dispatchSuccess('Test Persistent Notification from gRPC!', {
					persistent: true,
					actions: [
						{
							label: 'Test Action',
							handler: () => console.log('Action clicked!')
						}
					]
				})

				console.log('📝 Notification creation result:', result)

				// Force refresh after creation
				setTimeout(async () => {
					console.log('🔄 Force refreshing notifications...')
					await this.refreshNotifications()
					this.debugNotifications()
				}, 1000)

			} catch (error) {
				console.error('❌ Error in testPersistentNotification:', error)
				this.dispatchError(`Failed to create notification: ${error.message}`)
			}
		},

		testToastNotification() {
			console.log('🧪 Testing toast notification...')
			this.dispatchSuccess('Test Toast Notification (local only)')
		},

		async refreshNotifications() {
			try {
				console.log('🔄 Manually refreshing notifications...')
				const result = await this.loadPersistentNotifications()
				console.log('Refresh result:', result)

				if (result && result.count !== undefined) {
					this.dispatchSuccess(`📋 ${result.count} notifications loaded`)
				} else {
					this.dispatchInfo('📋 Notifications refreshed')
				}
			} catch (error) {
				console.error('❌ Failed to refresh notifications:', error)
				this.dispatchError(`❌ Failed to refresh notifications: ${error.message}`)
			}
		},

		async testNotificationConnection() {
			try {
				console.log('🧪 Testing notification gRPC connection...')

				const client = this.$store.getters['connection/notificationClient']
				if (!client) {
					throw new Error('Notification client not initialized')
				}

				console.log('✅ Notification client exists, testing ListNotifications...')

				// Import the request class
				const { ListNotificationsRequest } = await import('../../proto/notification_pb')
				const request = new ListNotificationsRequest()

				const response = await this.$store.dispatch('connection/promisifyGrpcCall', {
					method: client.listNotifications,
					request: request
				})

				const notifications = response.getNotificationsList()
				console.log(`✅ Successfully loaded ${notifications.length} notifications from gRPC`)

				this.dispatchSuccess(`✅ Notification gRPC working! Found ${notifications.length} notifications`)

				// Update store with loaded notifications
				if (notifications.length > 0) {
					await this.loadPersistentNotifications()
				}

			} catch (error) {
				console.error('❌ Notification gRPC test failed:', error)
				this.dispatchError(`❌ Notification gRPC failed: ${error.message}`)
			}
		},

		async testBackendServices() {
			console.log('🧪 Testing all backend services...')

			const tests = [
				{
					name: 'User Service - ListUsers',
					test: async () => {
						const { ListUsersRequest } = await import('../../proto/user_pb')
						const client = this.$store.getters['connection/grpcClient']
						const request = new ListUsersRequest()
						return await this.$store.dispatch('connection/promisifyGrpcCall', {
							method: client.listUsers,
							request: request
						})
					}
				},
				{
					name: 'Notification Service - ListNotifications',
					test: async () => {
						const { ListNotificationsRequest } = await import('../../proto/notification_pb')
						const client = this.$store.getters['connection/notificationClient']
						const request = new ListNotificationsRequest()
						return await this.$store.dispatch('connection/promisifyGrpcCall', {
							method: client.listNotifications,
							request: request
						})
					}
				},
				{
					name: 'Notification Service - CreateNotification',
					test: async () => {
						const { CreateNotificationRequest } = await import('../../proto/notification_pb')
						const client = this.$store.getters['connection/notificationClient']
						const request = new CreateNotificationRequest()
						request.setMessage('Backend test notification')
						request.setType('info')
						return await this.$store.dispatch('connection/promisifyGrpcCall', {
							method: client.createNotification,
							request: request
						})
					}
				}
			]

			for (const test of tests) {
				try {
					console.log(`🧪 Testing: ${test.name}`)
					const result = await test.test()
					console.log(`✅ ${test.name} - SUCCESS:`, result)
					this.dispatchSuccess(`✅ ${test.name} - SUCCESS`)
				} catch (error) {
					console.error(`❌ ${test.name} - FAILED:`, error.message)
					this.dispatchError(`❌ ${test.name} - FAILED: ${error.message}`)
				}
			}
		},

		async loadUsers() {
			try {
				this.dispatchInfo('Lade Benutzer...')
				const result = await this.fetchUsers()
				this.lastUpdated = new Date().toLocaleTimeString()

				if (!this.hasUsers) {
					this.dispatchInfo('📋 Keine Benutzer in der Datenbank gefunden')
				} else {
					this.dispatchSuccess(`📋 ${result.count} Benutzer erfolgreich geladen`, { persistent: true })
				}
			} catch (error) {
				console.error('❌ Error loading users:', error)
				this.dispatchError(`❌ Fehler beim Laden: ${error.message}`, {
					actions: [
						{
							label: 'Erneut versuchen',
							handler: () => this.loadUsers()
						}
					]
				})
			}
		},

		async handleCreateUser(protoMessage, formData) {
			if (!this.isConnected) {
				this.dispatchWarning('⚠️ Backend-Verbindung nicht verfügbar')
				return
			}

			try {
				console.log('Creating user with proto message:', protoMessage)
				console.log('Form data:', formData)

				const response = await this.createUser(formData)

				// Update timestamp
				this.lastUpdated = new Date().toLocaleTimeString()

				this.dispatchSuccess(`🎉 Benutzer "${formData.name}" erfolgreich erstellt!`, {
					actions: [
						{
							label: 'Liste aktualisieren',
							handler: () => this.loadUsers()
						}
					]
				})

				// Auto-refresh user list
				setTimeout(() => this.loadUsers(), 500)

			} catch (error) {
				console.error('❌ Error creating user:', error)
				this.dispatchError(`❌ Fehler beim Erstellen: ${error.message}`, {
					persistent: true,
					actions: [
						{
							label: 'Erneut versuchen',
							handler: () => this.handleCreateUser(protoMessage, formData)
						}
					]
				})
			}
		},

		handleFormError(error) {
			console.error('❌ Form error:', error)
			this.dispatchError(`❌ Formular-Fehler: ${error.message}`)
		},

		formatTime(timestamp) {
			const now = new Date()
			const time = new Date(timestamp)
			const diffMs = now - time
			const diffMins = Math.floor(diffMs / 60000)

			if (diffMins < 1) return 'Jetzt'
			if (diffMins < 60) return `vor ${diffMins}m`

			const diffHours = Math.floor(diffMins / 60)
			if (diffHours < 24) return `vor ${diffHours}h`

			const diffDays = Math.floor(diffHours / 24)
			return `vor ${diffDays}d`
		}
	},

	async mounted() {
		console.log('🎯 Welcome component mounted')

		// Wait a bit for the app to be fully initialized
		setTimeout(async () => {
			if (this.$store.getters['connection/isConnected']) {
				console.log('🔄 Loading initial notifications...')
				try {
					await this.loadPersistentNotifications()
				} catch (error) {
					console.warn('⚠️ Could not load initial notifications:', error.message)
				}
			}
		}, 1000)
	},

	watch: {
		'$store.getters["connection/isConnected"]': {
			handler(newConnected, oldConnected) {
				if (newConnected && !oldConnected) {
					// Connection restored - reload notifications
					console.log('🔄 Connection restored, reloading notifications...')
					this.loadPersistentNotifications().catch(error => {
						console.warn('⚠️ Could not reload notifications:', error.message)
					})
				}
			},
			immediate: false
		}
	}
}
</script>
