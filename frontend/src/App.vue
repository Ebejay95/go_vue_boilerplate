<template>
	<div id="app" class="min-h-screen bg-gray-50">
	  <header class="bg-white shadow-sm border-b">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		  <div class="flex justify-between items-center h-16">
			<h1 class="text-xl font-semibold text-gray-900">
			  Direct gRPC-Web Test
			</h1>
			<div class="flex items-center space-x-4">
			  <!-- Connection Status -->
			  <div class="flex items-center space-x-2">
				<div
				  :class="[
					'w-3 h-3 rounded-full',
					connectionStatus.healthy ? 'bg-green-500' : 'bg-red-500'
				  ]"
				></div>
				<span class="text-sm text-gray-600">
				  {{ connectionStatus.message }}
				</span>
			  </div>
			</div>
		  </div>
		</div>
	  </header>
	  <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
		<!-- Status Message -->
		<div v-if="message"
			 :class="[
			   'mb-6 p-4 rounded-md',
			   message.type === 'success' ? 'bg-green-50 text-green-800 border border-green-200' : 'bg-red-50 text-red-800 border border-red-200'
			 ]">
		  {{ message.text }}
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
		  <!-- Create User Form -->
		  <div class="bg-white p-6 rounded-lg shadow">
			<h2 class="text-lg font-medium text-gray-900 mb-4">Create New User</h2>

			<form @submit.prevent="submitCreateUser" class="space-y-4">
			  <div>
				<label for="name" class="block text-sm font-medium text-gray-700">Name</label>
				<input
				  type="text"
				  id="name"
				  v-model="newUser.name"
				  required
				  placeholder="Enter name"
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
				<label for="age" class="block text-sm font-medium text-gray-700">Age</label>
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
				<label for="role" class="block text-sm font-medium text-gray-700">Role</label>
				<select
				  id="role"
				  v-model="newUser.role"
				  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
				>
				  <option value="user">User</option>
				  <option value="admin">Admin</option>
				</select>
			  </div>

			  <button
				type="submit"
				:disabled="isCreating || !connectionStatus.healthy"
				class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
			  >
				<svg v-if="isCreating" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				{{ isCreating ? 'Creating User...' : 'Create User' }}
			  </button>
			</form>
		  </div>

		  <!-- Users List -->
		  <div class="bg-white p-6 rounded-lg shadow">
			<div class="flex justify-between items-center mb-4">
			  <h2 class="text-lg font-medium text-gray-900">All Users</h2>
			  <button
				@click="loadUsers"
				:disabled="isLoading || !connectionStatus.healthy"
				class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
			  >
				<svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				{{ isLoading ? 'Loading...' : 'Refresh' }}
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
				  <span class="text-gray-500">Loading users...</span>
				</div>
			  </div>

			  <!-- Empty State -->
			  <div v-else-if="!hasUsers && !isLoading" class="text-center py-8">
				<svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
				  <path d="M34 40h10v-4a6 6 0 00-10.712-3.714M34 40H14m20 0v-4a9.971 9.971 0 00-.712-3.714M14 40H4v-4a6 6 0 0110.713-3.714M14 40v-4c0-1.313.253-2.566.713-3.714m0 0A10.003 10.003 0 0124 26c4.21 0 7.813 2.602 9.288 6.286M30 14a6 6 0 11-12 0 6 6 0 0112 0z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900">No users found</h3>
				<p class="mt-1 text-sm text-gray-500">Create the first user to get started.</p>
			  </div>

			  <!-- Users Table -->
			  <div v-else class="overflow-hidden">
				<table class="min-w-full divide-y divide-gray-200">
				  <thead class="bg-gray-50">
					<tr>
					  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
					  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
					  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Age</th>
					  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
					</tr>
				  </thead>
				  <tbody class="bg-white divide-y divide-gray-200">
					<tr v-for="user in users" :key="user.id" class="hover:bg-gray-50">
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
						  {{ user.role }}
						</span>
					  </td>
					</tr>
				  </tbody>
				</table>
			  </div>

			  <!-- Last Updated -->
			  <div v-if="lastUpdated" class="text-center text-xs text-gray-500 mt-4">
				Last updated: {{ lastUpdated }}
			  </div>
			</div>
		  </div>
		</div>

		<!-- Debug Info -->
		<div class="mt-8 bg-gray-100 p-4 rounded-lg">
		  <h3 class="text-sm font-medium text-gray-900 mb-2">Debug Information</h3>
		  <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs text-gray-600">
			<div>
			  <strong>gRPC-Web URL:</strong> {{ grpcWebUrl }}
			</div>
			<div>
			  <strong>Connection:</strong> {{ connectionStatus.connection }}
			</div>
			<div>
			  <strong>Status:</strong> {{ connectionStatus.status }}
			</div>
			<div>
			  <strong>Users Count:</strong> {{ users.length }}
			</div>
		  </div>
		</div>
	  </main>
	</div>
  </template>

  <script>
  import { grpcClient } from './services/grpcClient.js'

  export default {
	name: 'DirectGrpcWebTest',
	data() {
	  return {
		users: [],
		isLoading: false,
		isCreating: false,
		message: null,
		lastUpdated: null,
		connectionStatus: {
		  healthy: false,
		  message: 'Checking connection...',
		  status: 'unknown',
		  connection: 'grpc-web'
		},
		newUser: {
		  name: '',
		  email: '',
		  age: '',
		  role: 'user'
		},
		// ✅ Nur eine URL - gRPC-Web
		grpcWebUrl: process.env.VUE_APP_GRPC_WEB_URL
	  }
	},

	computed: {
	  hasUsers() {
		return this.users && this.users.length > 0
	  }
	},

	methods: {
	  async checkHealth() {
		console.log(`🔍 Checking gRPC-Web health...`)
		try {
		  const response = await grpcClient.checkHealth()
		  this.connectionStatus = {
			healthy: response.status === 'healthy',
			message: response.status === 'healthy'
			  ? `✅ Connected via gRPC-Web with ${response.users || 0} users`
			  : `❌ ${response.error || 'Connection failed'}`,
			status: response.status,
			connection: 'grpc-web'
		  }

		  if (response.status === 'healthy') {
			this.showMessage(`gRPC-Web connection successful!`, 'success')
		  } else {
			this.showMessage(`gRPC-Web connection failed`, 'error')
		  }
		} catch (error) {
		  console.error(`❌ Health check failed:`, error)
		  this.connectionStatus = {
			healthy: false,
			message: `❌ gRPC-Web server not reachable: ${error.message}`,
			status: 'disconnected',
			connection: 'grpc-web'
		  }
		  this.showMessage(`gRPC-Web connection failed: ${error.message}`, 'error')
		}
	  },

	  async loadUsers() {
		if (!this.connectionStatus.healthy) {
		  this.showMessage('⚠️ Backend connection not available', 'error')
		  return
		}

		this.isLoading = true
		console.log(`📋 Loading users via gRPC-Web...`)

		try {
		  const response = await grpcClient.listUsers()
		  this.users = response.users || []
		  this.lastUpdated = new Date().toLocaleTimeString()

		  console.log(`✅ Loaded ${this.users.length} users via gRPC-Web`)

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
		console.log(`➕ Creating user via gRPC-Web: ${this.newUser.name}`)

		try {
		  const response = await grpcClient.createUser(this.newUser)
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

		  this.showMessage(`🎉 User "${response.user.name}" created via gRPC-Web!`, 'success')
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

	async mounted() {
	  console.log(`🚀 Direct gRPC-Web Test Component mounted`)
	  console.log(`🌐 gRPC-Web URL: ${this.grpcWebUrl}`)

	  await this.checkHealth()

	  if (this.connectionStatus.healthy) {
		await this.loadUsers()
	  } else {
		console.warn('⚠️ Skipping initial load due to connection issues')
	  }
	}
  }
  </script>
