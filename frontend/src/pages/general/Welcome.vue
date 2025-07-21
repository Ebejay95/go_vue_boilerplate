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
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
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
					<td class="px-6 py-4 whitespace-nowrap">
						<base-button
						  mode="primary"
						  @click="showDeleteDialog(user)"
						  :disabled="isLoading"
						>
						  Löschen
						</base-button>
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

	  <!-- Socket Status and Test Controls -->
	  <div class="mt-8 bg-gray-50 rounded-lg p-4">
		<h3 class="text-lg font-medium text-gray-900 mb-4">Socket System & Notifications</h3>

		<!-- Socket Status -->
		<div class="mb-4 p-3 bg-white rounded-lg border">
		  <div class="flex items-center space-x-3">
			<div class="flex items-center space-x-2">
			  <span
				:class="[
				  'inline-block w-3 h-3 rounded-full',
				  socketStatus === 'connected' ? 'bg-green-500' :
				  socketStatus === 'connecting' ? 'bg-yellow-500 animate-pulse' : 'bg-red-500'
				]"
			  ></span>
			  <span class="text-sm font-medium text-gray-700">
				Socket: {{ socketStatusText }}
			  </span>
			</div>

			<div class="flex items-center space-x-2">
			  <span class="text-sm text-gray-500">
				Persistent Notifications: {{ unreadCount }}
			  </span>
			  <span v-if="hasUnreadPersistent" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
				{{ unreadCount }} ungelesen
			  </span>
			</div>
		  </div>
		</div>

		<!-- Test Buttons -->
		<div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-2">
		  <button
			@click="testSuccessNotification"
			class="px-3 py-2 bg-green-100 text-green-800 rounded text-sm hover:bg-green-200 transition-colors"
		  >
			✅ Success
		  </button>

		  <button
			@click="testErrorNotification"
			class="px-3 py-2 bg-red-100 text-red-800 rounded text-sm hover:bg-red-200 transition-colors"
		  >
			❌ Error
		  </button>

		  <button
			@click="testWarningNotification"
			class="px-3 py-2 bg-yellow-100 text-yellow-800 rounded text-sm hover:bg-yellow-200 transition-colors"
		  >
			⚠️ Warning
		  </button>

		  <button
			@click="testInfoNotification"
			class="px-3 py-2 bg-blue-100 text-blue-800 rounded text-sm hover:bg-blue-200 transition-colors"
		  >
			ℹ️ Info
		  </button>

		  <button
			@click="testPersistentNotification"
			class="px-3 py-2 bg-purple-100 text-purple-800 rounded text-sm hover:bg-purple-200 transition-colors"
		  >
			💾 Persistent
		  </button>

		  <button
			@click="testCustomSocketEvent"
			class="px-3 py-2 bg-indigo-100 text-indigo-800 rounded text-sm hover:bg-indigo-200 transition-colors"
		  >
			🔌 Socket Event
		  </button>

		  <button
			@click="simulateBackendNotification"
			class="px-3 py-2 bg-orange-100 text-orange-800 rounded text-sm hover:bg-orange-200 transition-colors"
		  >
			📨 Backend Sim
		  </button>

		  <button
			@click="clearAllToasts"
			class="px-3 py-2 bg-gray-100 text-gray-800 rounded text-sm hover:bg-gray-200 transition-colors"
		  >
			🧹 Clear
		  </button>
		</div>
	  </div>

	  <!-- Delete Confirmation Dialog -->
	  <base-dialog
		mode="dialog"
		:visible="deleteDialog.visible"
		:primary-action="deleteDialog.primaryAction"
		close-label="Abbrechen"
		:close-on-overlay="true"
		@close="closeDeleteDialog"
		@primary-action="confirmDeleteUser"
	  >
		<div class="text-center">
		  <h3 class="text-lg font-medium text-gray-900 mb-2">
			Benutzer löschen?
		  </h3>

		  <p class="text-sm text-gray-500 mb-6">
			Möchten Sie den Benutzer <strong>"{{ deleteDialog.user?.name }}"</strong> wirklich löschen?
			Diese Aktion kann nicht rückgängig gemacht werden.
		  </p>

		  <div v-if="deleteDialog.user" class="bg-gray-50 rounded-lg p-3 mb-4">
			<div class="flex items-center space-x-3">
			  <div class="h-8 w-8 rounded-full bg-blue-500 flex items-center justify-center">
				<span class="text-white font-medium text-xs">{{ deleteDialog.user.name.charAt(0).toUpperCase() }}</span>
			  </div>
			  <div class="text-left">
				<div class="text-sm font-medium text-gray-900">{{ deleteDialog.user.name }}</div>
				<div class="text-xs text-gray-500">{{ deleteDialog.user.email }}</div>
			  </div>
			</div>
		  </div>
		</div>
	  </base-dialog>
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
			CreateUserRequest,

			// Delete Dialog State
			deleteDialog: {
				visible: false,
				user: null,
				primaryAction: {
					label: 'Löschen',
					callback: () => this.confirmDeleteUser()
				}
			},

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
		...mapGetters('socket', ['connectionStatus']),

		isLoading() {
			return this.loading
		},

		isCreating() {
			return this.creating
		},

		socketStatus() {
			return this.connectionStatus
		},

		socketStatusText() {
			const statusMap = {
				connected: 'Verbunden',
				connecting: 'Verbinde...',
				disconnected: 'Getrennt'
			}
			return statusMap[this.socketStatus] || 'Unbekannt'
		}
	},

	methods: {
		...mapActions('users', ['fetchUsers', 'createUser', 'clearError']),
		...mapActions('connection', ['checkHealth', 'reconnect']),

		// Notification methods using the new system
		async loadUsers() {
			try {
				await this.$store.dispatch('notifications/info', 'Lade Benutzer...')
				const result = await this.fetchUsers()
				this.lastUpdated = new Date().toLocaleTimeString()

				if (!this.hasUsers) {
					await this.$store.dispatch('notifications/info', '📋 Keine Benutzer in der Datenbank gefunden')
				} else {
					await this.$store.dispatch('notifications/success', `📋 ${result.count} Benutzer erfolgreich geladen`)
				}
			} catch (error) {
				console.error('Error loading users:', error)
				await this.$store.dispatch('notifications/error', `❌ Fehler beim Laden: ${error.message}`, {
					persistent: true
				})
			}
		},

		async handleCreateUser(protoMessage, formData) {
			if (!this.isConnected) {
				await this.$store.dispatch('notifications/warning', '⚠️ Backend-Verbindung nicht verfügbar')
				return
			}

			try {
				const response = await this.createUser(formData)
				this.lastUpdated = new Date().toLocaleTimeString()

				await this.$store.dispatch('notifications/success', `🎉 Benutzer "${formData.name}" erfolgreich erstellt!`)

				// Auto-refresh user list
				setTimeout(() => this.loadUsers(), 500)

			} catch (error) {
				console.error('Error creating user:', error)
				await this.$store.dispatch('notifications/error', `❌ Fehler beim Erstellen: ${error.message}`, {
					persistent: true
				})
			}
		},

		async handleFormError(error) {
			console.error('Form error:', error)
			await this.$store.dispatch('notifications/error', `❌ Formular-Fehler: ${error.message}`)
		},

		// Dialog Methods
		showDeleteDialog(user) {
			this.deleteDialog.user = user
			this.deleteDialog.visible = true
		},

		closeDeleteDialog() {
			this.deleteDialog.visible = false
			this.deleteDialog.user = null
		},

		async confirmDeleteUser() {
			const user = this.deleteDialog.user
			if (!user) return

			if (!this.isConnected) {
				await this.$store.dispatch('notifications/warning', '⚠️ Backend-Verbindung nicht verfügbar')
				this.closeDeleteDialog()
				return
			}

			try {
				const result = await this.$store.dispatch('users/deleteUser', user.id)
				this.lastUpdated = new Date().toLocaleTimeString()

				await this.$store.dispatch('notifications/success', `🗑️ Benutzer "${user.name}" erfolgreich gelöscht!`)
				this.closeDeleteDialog()

			} catch (error) {
				console.error('Error deleting user:', error)
				await this.$store.dispatch('notifications/error', `❌ Fehler beim Löschen: ${error.message}`, {
					persistent: true
				})
			}
		},

		// Test Methods for Notifications
		async testSuccessNotification() {
			await this.$store.dispatch('notifications/success', '✅ Das ist eine Erfolgs-Nachricht!')
		},

		async testErrorNotification() {
			await this.$store.dispatch('notifications/error', '❌ Das ist eine Fehler-Nachricht!')
		},

		async testWarningNotification() {
			await this.$store.dispatch('notifications/warning', '⚠️ Das ist eine Warnung!')
		},

		async testInfoNotification() {
			await this.$store.dispatch('notifications/info', 'ℹ️ Das ist eine Info-Nachricht!')
		},

		async testPersistentNotification() {
			await this.$store.dispatch('notifications/info', '💾 Das ist eine persistente Nachricht!', {
				persistent: true
			})
		},

		async testCustomSocketEvent() {
			// Send custom event via socket
			await this.$store.dispatch('socket/emit', {
				event: 'custom_test_event',
				data: {
					message: 'Das ist ein benutzerdefiniertes Socket Event',
					timestamp: new Date().toISOString(),
					userId: 1
				}
			})

			await this.$store.dispatch('notifications/info', '🔌 Custom Socket Event gesendet!')
		},

		async simulateBackendNotification() {
			// Simulate a backend notification by manually triggering the handler
			this.$store.dispatch('notifications/handleSocketNotification', {
				id: Date.now(),
				message: '📨 Simulierte Backend-Nachricht',
				type: 'info',
				persistent: false,
				createdAt: new Date().toISOString()
			})
		},

		async clearAllToasts() {
			this.$store.dispatch('notifications/clearAllToasts')
		}
	},

	async mounted() {
		// Setup custom socket event listeners for this component
		this.$store.dispatch('socket/on', {
			event: 'user_created',
			callback: (data) => {
				console.log('User created event received:', data)
				this.$store.dispatch('notifications/success', `👤 Neuer Benutzer: ${data.name}`)
			}
		})

		this.$store.dispatch('socket/on', {
			event: 'user_deleted',
			callback: (data) => {
				console.log('User deleted event received:', data)
				this.$store.dispatch('notifications/warning', `🗑️ Benutzer gelöscht: ${data.name}`)
			}
		})

		this.$store.dispatch('socket/on', {
			event: 'custom_test_event',
			callback: (data) => {
				console.log('Custom test event received:', data)
				this.$store.dispatch('notifications/info', `🎯 Custom Event: ${data.message}`)
			}
		})

		// Wait for app initialization
		setTimeout(async () => {
			if (this.$store.getters['connection/isConnected']) {
				try {
					await this.$store.dispatch('notifications/loadPersistentNotifications')
				} catch (error) {
					console.warn('Could not load initial notifications:', error.message)
				}
			}
		}, 1000)
	},

	beforeUnmount() {
		// Clean up socket listeners if needed
		// Note: In a real app, you might want to remove specific listeners
		console.log('Component unmounting - socket listeners will be cleaned up by socket disconnect')
	},

	watch: {
		'$store.getters["connection/isConnected"]': {
			handler(newConnected, oldConnected) {
				if (newConnected && !oldConnected) {
					// Connection restored
					this.$store.dispatch('notifications/success', '🔌 Verbindung wiederhergestellt!')
					this.$store.dispatch('notifications/loadPersistentNotifications').catch(error => {
						console.warn('Could not reload notifications:', error.message)
					})
				} else if (!newConnected && oldConnected) {
					// Connection lost
					this.$store.dispatch('notifications/warning', '⚠️ Verbindung verloren!')
				}
			},
			immediate: false
		},

		'$store.getters["socket/connectionStatus"]': {
			handler(newStatus, oldStatus) {
				if (newStatus === 'connected' && oldStatus !== 'connected') {
					this.$store.dispatch('notifications/success', '🔌 Socket verbunden!')
				} else if (newStatus === 'disconnected' && oldStatus === 'connected') {
					this.$store.dispatch('notifications/warning', '🔌 Socket getrennt!')
				}
			},
			immediate: false
		}
	}
}
</script>
