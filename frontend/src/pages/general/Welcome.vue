<!-- src/pages/general/Welcome.vue - Updated with Dynamic Form -->
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
		<base-card>
			<h3>Persistent Notifications</h3>
			<ul>
				<li v-for="pNote in persistentNotifications" :key="pNote.id">
				{{ pNote.message }} - {{ pNote.type }}
				</li>
			</ul>
		</base-card>
	  </div>

	</base-section>
<div class="mt-4 space-x-2">
  <button @click="testPersistentNotification" class="btn primary">
    Test Persistent
  </button>
  <button @click="testToastNotification" class="btn secondary">
    Test Toast
  </button>
  <button @click="debugNotifications" class="btn flat">
    Debug Notifications
  </button>
</div>
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
	  ...mapGetters('notifications', ['persistentNotifications']),

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
		'dispatchInfo'
	  ]),

	  debugNotifications() {
    console.log('🔍 Debug Notifications:')
    console.log('Persistent:', this.persistentNotifications)
    console.log('Unread Count:', this.unreadCount)
    console.log('Has Unread:', this.hasUnreadPersistent)
  },

  testPersistentNotification() {
    console.log('🧪 Testing persistent notification...')
    this.dispatchSuccess('Test Persistent Notification', {
      persistent: true,
      actions: [
        {
          label: 'Test Action',
          handler: () => console.log('Action clicked!')
        }
      ]
    })

    // Debug nach 100ms
    setTimeout(() => {
      this.debugNotifications()
    }, 100)
  },

  testToastNotification() {
    console.log('🧪 Testing toast notification...')
    this.dispatchSuccess('Test Toast Notification')
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

	  // Handle dynamic form submission
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

	  // Handle form errors
	  handleFormError(error) {
		console.error('❌ Form error:', error)
		this.dispatchError(`❌ Formular-Fehler: ${error.message}`)
	  },
	}
  }
  </script>