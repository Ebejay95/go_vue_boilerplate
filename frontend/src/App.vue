<template>
	<div class="max-w-6xl mx-auto p-6">
	  <!-- Header -->
	  <div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-900 mb-2">🚀 Direct gRPC-Web Testing</h1>
		<p class="text-gray-600">Test direkte gRPC-Web Verbindung (ohne REST Proxy)</p>
	  </div>

	  <!-- Connection Method Selection -->
	  <div class="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
		<h3 class="text-lg font-medium text-blue-900 mb-3">Verbindungsmethode wählen:</h3>
		<div class="space-y-2">
		  <label class="inline-flex items-center">
			<input
			  type="radio"
			  v-model="connectionMethod"
			  value="grpc-web"
			  class="text-blue-600 focus:ring-blue-500"
			>
			<span class="ml-2">⚡ Direct gRPC-Web (Empfohlen)</span>
		  </label>
		  <label class="inline-flex items-center">
			<input
			  type="radio"
			  v-model="connectionMethod"
			  value="rest"
			  class="text-blue-600 focus:ring-blue-500"
			>
			<span class="ml-2">🌐 REST API (Express Proxy - Fallback)</span>
		  </label>
		</div>
	  </div>

	  <!-- Status Display -->
	  <div class="mb-6 p-4 rounded-lg" :class="connectionStatus.healthy ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'">
		<div class="flex items-center justify-between">
		  <div class="flex items-center">
			<div class="flex-shrink-0">
			  <div class="w-3 h-3 rounded-full mr-3" :class="connectionStatus.healthy ? 'bg-green-500 animate-pulse' : 'bg-red-500'"></div>
			</div>
			<div>
			  <p class="text-sm font-medium" :class="connectionStatus.healthy ? 'text-green-800' : 'text-red-800'">
				{{ connectionStatus.healthy ? `✅ ${getCurrentMethodName()} Connected` : `❌ ${getCurrentMethodName()} Disconnected` }}
			  </p>
			  <p class="text-xs" :class="connectionStatus.healthy ? 'text-green-600' : 'text-red-600'">
				{{ connectionStatus.message }}
			  </p>
			</div>
		  </div>
		  <button
			@click="checkHealth"
			class="text-xs px-3 py-1 rounded border bg-white"
			:class="connectionStatus.healthy ? 'border-green-300 text-green-700 hover:bg-green-50' : 'border-red-300 text-red-700 hover:bg-red-50'"
		  >
			Test Connection
		  </button>
		</div>
	  </div>

	  <!-- Create User Form -->
	  <div class="bg-white shadow-lg rounded-lg p-6 mb-8 border">
		<h2 class="text-xl font-semibold text-gray-900 mb-4 flex items-center">
		  <span class="text-2xl mr-2">➕</span>
		  User erstellen via {{ getCurrentMethodName() }}
		</h2>

		<form @submit.prevent="submitCreateUser" class="space-y-4">
		  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<div>
			  <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
			  <input
				type="text"
				id="name"
				v-model="newUser.name"
				required
				placeholder="Name eingeben"
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			  />
			</div>

			<div>
			  <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
			  <input
				type="email"
				id="email"
				v-model="newUser.email"
				required
				placeholder="email@example.com"
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			  />
			</div>

			<div>
			  <label for="age" class="block text-sm font-medium text-gray-700 mb-1">Alter</label>
			  <input
				type="number"
				id="age"
				v-model.number="newUser.age"
				required
				min="1"
				max="150"
				placeholder="25"
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			  />
			</div>

			<div>
			  <label for="role" class="block text-sm font-medium text-gray-700 mb-1">Rolle</label>
			  <select
				id="role"
				v-model="newUser.role"
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			  >
				<option value="user">User</option>
				<option value="admin">Admin</option>
			  </select>
			</div>
		  </div>

		  <button
			type="submit"
			:disabled="isCreating || !connectionStatus.healthy"
			class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
		  >
			<svg v-if="isCreating" class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
			  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
			  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			{{ isCreating ? 'Creating...' : `User erstellen (${getCurrentMethodName()})` }}
		  </button>
		</form>
	  </div>

	  <!-- Users List -->
	  <div class="bg-white shadow-lg rounded-lg border">
		<div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center bg-gray-50">
		  <h2 class="text-xl font-semibold text-gray-900 flex items-center">
			<span class="text-2xl mr-2">📋</span>
			Users via {{ getCurrentMethodName() }}
		  </h2>
		  <button
			@click="loadUsers"
			:disabled="isLoading || !connectionStatus.healthy"
			class="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
		  >
			<svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
			  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
			  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			{{ isLoading ? 'Loading...' : 'Refresh' }}
		  </button>
		</div>

		<div class="p-6">
		  <!-- Loading State -->
		  <div v-if="isLoading && !hasUsers" class="text-center py-8">
			<div class="inline-flex items-center">
			  <svg class="animate-spin -ml-1 mr-3 h-6 w-6 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			  </svg>
			  <span class="text-gray-600">Loading users via {{ getCurrentMethodName() }}...</span>
			</div>
		  </div>

		  <!-- Empty State -->
		  <div v-else-if="!hasUsers && !isLoading" class="text-center py-8">
			<div class="text-6xl mb-4">👥</div>
			<h3 class="text-lg font-medium text-gray-900 mb-2">Keine Users gefunden</h3>
			<p class="text-gray-500 mb-4">Erstellen Sie den ersten User.</p>
		  </div>

		  <!-- Users Table -->
		  <div v-else class="overflow-x-auto">
			<table class="min-w-full divide-y divide-gray-200">
			  <thead class="bg-gray-50">
				<tr>
				  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">👤 User</th>
				  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">📧 Email</th>
				  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">🎂 Alter</th>
				  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">🏷️ Rolle</th>
				</tr>
			  </thead>
			  <tbody class="bg-white divide-y divide-gray-200">
				<tr v-for="user in users" :key="user.id" class="hover:bg-gray-50">
				  <td class="px-6 py-4 whitespace-nowrap">
					<div class="flex items-center">
					  <div class="flex-shrink-0 h-8 w-8">
						<div class="h-8 w-8 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center">
						  <span class="text-xs font-bold text-white">{{ user.name.charAt(0).toUpperCase() }}</span>
						</div>
					  </div>
					  <div class="ml-4">
						<div class="text-sm font-medium text-gray-900">{{ user.name }}</div>
						<div class="text-xs text-gray-500">ID: {{ user.id }}</div>
					  </div>
					</div>
				  </td>
				  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.email }}</td>
				  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.age }}</td>
				  <td class="px-6 py-4 whitespace-nowrap">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full"
						  :class="user.role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-green-100 text-green-800'">
					  {{ user.role === 'admin' ? '👑 Admin' : '👤 User' }}
					</span>
				  </td>
				</tr>
			  </tbody>
			</table>
		  </div>
		</div>
	  </div>

	  <!-- Success/Error Messages -->
	  <div v-if="message" class="fixed bottom-4 right-4 max-w-sm z-50">
		<div class="p-4 rounded-lg shadow-lg border" :class="message.type === 'success' ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'">
		  <div class="flex items-center">
			<span class="text-lg mr-2">{{ message.type === 'success' ? '✅' : '❌' }}</span>
			<p class="text-sm" :class="message.type === 'success' ? 'text-green-800' : 'text-red-800'">
			  {{ message.text }}
			</p>
		  </div>
		</div>
	  </div>

	  <!-- Debug Info -->
	  <div class="mt-8 p-4 bg-gray-100 rounded-lg border">
		<h3 class="text-sm font-medium text-gray-900 mb-2">🔧 Debug Info</h3>
		<div class="text-xs text-gray-600 space-y-1">
		  <div>Connection Method: {{ connectionMethod }}</div>
		  <div>Connection Type: {{ connectionStatus.connection || 'unknown' }}</div>
		  <div>Backend Status: {{ connectionStatus.status || 'unknown' }}</div>
		  <div>Users Count: {{ users.length }}</div>
		  <div>Last Updated: {{ lastUpdated || 'never' }}</div>
		  <div>gRPC-Web URL: {{ grpcWebUrl }}</div>
		  <div>REST API URL: {{ restApiUrl }}</div>
		</div>
	  </div>
	</div>
  </template>

  <script>
  import { simpleGrpcClient, restLikeGrpcClient } from './services/grpcClient.js'

  export default {
	name: 'DirectGrpcWebTest',
	data() {
	  return {
		connectionMethod: 'grpc-web', // Start with gRPC-Web as default
		users: [],
		isLoading: false,
		isCreating: false,
		message: null,
		lastUpdated: null,
		connectionStatus: {
		  healthy: false,
		  message: 'Checking connection...',
		  status: 'unknown',
		  connection: 'unknown'
		},
		newUser: {
		  name: '',
		  email: '',
		  age: '',
		  role: 'user'
		},
		// Environment URLs (Vue.js macht diese automatisch verfügbar)
		grpcWebUrl: import.meta.env?.VUE_APP_GRPC_WEB_URL || process.env?.VUE_APP_GRPC_WEB_URL || 'http://localhost:8081',
		restApiUrl: import.meta.env?.VUE_APP_API_BASE_URL || process.env?.VUE_APP_API_BASE_URL || 'http://localhost:3000/api'
	  }
	},

	computed: {
	  hasUsers() {
		return this.users && this.users.length > 0
	  },

	  currentClient() {
		return this.connectionMethod === 'grpc-web' ? simpleGrpcClient : restLikeGrpcClient
	  }
	},

	methods: {
	  getCurrentMethodName() {
		return this.connectionMethod === 'grpc-web' ? 'gRPC-Web' : 'REST API'
	  },

	  async checkHealth() {
		console.log(`🔍 Checking health via ${this.getCurrentMethodName()}...`)
		try {
		  const response = await this.currentClient.checkHealth()
		  this.connectionStatus = {
			healthy: response.status === 'healthy',
			message: response.status === 'healthy'
			  ? `✅ Connected via ${this.getCurrentMethodName()} with ${response.users || 0} users`
			  : `❌ ${response.error || 'Connection failed'}`,
			status: response.status,
			connection: response.connection || this.connectionMethod
		  }

		  if (response.status === 'healthy') {
			this.showMessage(`${this.getCurrentMethodName()} connection successful!`, 'success')
		  } else {
			this.showMessage(`${this.getCurrentMethodName()} connection failed`, 'error')
		  }
		} catch (error) {
		  console.error(`❌ Health check failed:`, error)
		  this.connectionStatus = {
			healthy: false,
			message: `❌ ${this.getCurrentMethodName()} server not reachable: ${error.message}`,
			status: 'disconnected',
			connection: this.connectionMethod
		  }
		  this.showMessage(`${this.getCurrentMethodName()} connection failed: ${error.message}`, 'error')
		}
	  },

	  async loadUsers() {
		if (!this.connectionStatus.healthy) {
		  this.showMessage('⚠️ Backend connection not available', 'error')
		  return
		}

		this.isLoading = true
		console.log(`📋 Loading users via ${this.getCurrentMethodName()}...`)

		try {
		  const response = await this.currentClient.listUsers()
		  this.users = response.users || []
		  this.lastUpdated = new Date().toLocaleTimeString()

		  console.log(`✅ Loaded ${this.users.length} users via ${this.getCurrentMethodName()}`)

		  if (this.users.length === 0) {
			this.showMessage('📋 No users found in database', 'success')
		  }
		} catch (error) {
		  console.error(`❌ Error loading users:`, error)
		  this.showMessage(`❌ Error: ${error.message}`, 'error')

		  // Reset connection status on error
		  this.connectionStatus.healthy = false
		  this.connectionStatus.message = `❌ ${error.message}`
		} finally {
		  this.isLoading = false
		}
	  },

	  async submitCreateUser() {
		if (!this.connectionStatus.healthy) {
		  this.showMessage('⚠️ Backend connection not available', 'error')
		  return
		}

		this.isCreating = true
		console.log(`➕ Creating user via ${this.getCurrentMethodName()}: ${this.newUser.name}`)

		try {
		  const response = await this.currentClient.createUser(this.newUser)
		  console.log(`✅ User created:`, response.user)

		  // Reset form
		  this.newUser = {
			name: '',
			email: '',
			age: '',
			role: 'user'
		  }

		  // Reload users list
		  await this.loadUsers()

		  this.showMessage(`🎉 User "${response.user.name}" created via ${this.getCurrentMethodName()}!`, 'success')
		} catch (error) {
		  console.error(`❌ Error creating user:`, error)
		  this.showMessage(`❌ Error: ${error.message}`, 'error')
		} finally {
		  this.isCreating = false
		}
	  },

	  showMessage(text, type = 'success') {
		this.message = { text, type }
		setTimeout(() => {
		  this.message = null
		}, 5000)
	  }
	},

	watch: {
	  // When connection method changes, re-check health
	  connectionMethod(newMethod) {
		console.log(`🔄 Switching to ${this.getCurrentMethodName()}`)
		this.connectionStatus = {
		  healthy: false,
		  message: `Switching to ${this.getCurrentMethodName()}...`,
		  status: 'switching',
		  connection: newMethod
		}

		// Small delay to let UI update
		setTimeout(async () => {
		  await this.checkHealth()
		  if (this.connectionStatus.healthy) {
			await this.loadUsers()
		  }
		}, 500)
	  }
	},

	async mounted() {
	  console.log(`🚀 Direct gRPC-Web Test Component mounted`)
	  console.log(`🔗 Starting with: ${this.getCurrentMethodName()}`)
	  console.log(`🌐 gRPC-Web URL: ${this.grpcWebUrl}`)
	  console.log(`🌐 REST API URL: ${this.restApiUrl}`)

	  await this.checkHealth()

	  if (this.connectionStatus.healthy) {
		await this.loadUsers()
	  } else {
		console.warn('⚠️ Skipping initial load due to connection issues')
	  }
	}
  }
  </script>
