<template>
	<div id="app">
	  <div class="container">

		<!-- User Creation Form -->
		<div class="form-section">
		  <h2>Neuen asdasd</h2>

		  <form @submit.prevent="createUser" class="user-form">
			<div class="form-group">
			  <label for="name">Name</label>
			  <input
				type="text"
				id="name"
				v-model="newUser.name"
				required
				placeholder="Name eingeben"
			  />
			</div>

			<div class="form-group">
			  <label for="email">Email</label>
			  <input
				type="email"
				id="email"
				v-model="newUser.email"
				required
				placeholder="email@example.com"
			  />
			</div>

			<div class="form-group">
			  <label for="age">Alter</label>
			  <input
				type="number"
				id="age"
				v-model.number="newUser.age"
				required
				min="1"
				max="150"
				placeholder="25"
			  />
			</div>

			<button
			  type="submit"
			  :disabled="isCreating"
			  class="btn btn-primary"
			>
			  <span v-if="isCreating">Erstelle User...</span>
			  <span v-else>User erstellen</span>
			</button>
		  </form>
		</div>

		<!-- Users Table -->
		<div class="table-section">
		  <div class="table-header">
			<h2>Alle Users</h2>
			<button
			  @click="loadUsers"
			  :disabled="isLoading"
			  class="btn btn-secondary"
			>
			  <span v-if="isLoading">Laden...</span>
			  <span v-else>Users laden</span>
			</button>
		  </div>

		  <!-- Loading State -->
		  <div v-if="isLoading" class="loading">
			<div class="spinner"></div>
		  </div>

		  <!-- Error State -->
		  <div v-else-if="error" class="error-message">
			<p>{{ error }}</p>
		  </div>

		  <!-- Users Table -->
		  <div v-else class="table-container">
			<table class="users-table">
			  <thead>
				<tr>
				  <th>ID</th>
				  <th>Name</th>
				  <th>Email</th>
				  <th>Alter</th>
				  <th>Role</th>
				</tr>
			  </thead>
			  <tbody>
				<tr v-if="users.length === 0">
				  <td colspan="4" class="no-data">
					Keine Users gefunden
				  </td>
				</tr>
				<tr v-for="user in users" :key="user.id">
				  <td>{{ user.id }}</td>
				  <td>{{ user.name }}</td>
				  <td>{{ user.email }}</td>
				  <td>{{ user.age }}</td>
				  <td>{{ user.role }}</td>
				</tr>
			  </tbody>
			</table>
		  </div>
		</div>

		<!-- Success/Error Messages -->
		<div v-if="message" class="toast" :class="messageType">
		  <p>{{ message }}</p>
		</div>
	  </div>
	</div>
  </template>

  <script>
  import { ref, reactive, onMounted } from 'vue'
  import axios from 'axios'

  export default {
	name: 'App',
	setup() {
	  const users = ref([])
	  const isLoading = ref(false)
	  const isCreating = ref(false)
	  const error = ref('')
	  const message = ref('')
	  const messageType = ref('success')

	  const newUser = reactive({
		name: '',
		email: '',
		age: ''
	  })

	  const serverStatus = reactive({
		port: process.env.NODE_ENV === 'production' ? 3000 : 8080,
		grpcServer: 'localhost:50051'
	  })

	  const API_BASE = process.env.NODE_ENV === 'production' ? '' : 'http://localhost:3000'

	  const showMessage = (msg, type = 'success') => {
		message.value = msg
		messageType.value = type
		setTimeout(() => {
		  message.value = ''
		}, 3000)
	  }

	  const loadUsers = async () => {
		isLoading.value = true
		error.value = ''

		try {
		  const response = await axios.get(`${API_BASE}/api/users`)
		  users.value = response.data
		} catch (err) {
		  error.value = `Fehler beim Laden der Users: ${err.response?.data?.error || err.message}`
		  console.error('Error loading users:', err)
		} finally {
		  isLoading.value = false
		}
	  }

	  const createUser = async () => {
		isCreating.value = true

		try {
		  const response = await axios.post(`${API_BASE}/api/users`, {
			name: newUser.name,
			email: newUser.email,
			age: parseInt(newUser.age)
		  })

		  // Reset form
		  newUser.name = ''
		  newUser.email = ''
		  newUser.age = ''

		  // Reload users and show success message
		  await loadUsers()
		  showMessage('User erfolgreich erstellt!', 'success')
		} catch (err) {
		  const errorMsg = err.response?.data?.error || err.message
		  showMessage(`Fehler beim Erstellen des Users: ${errorMsg}`, 'error')
		  console.error('Error creating user:', err)
		} finally {
		  isCreating.value = false
		}
	  }

	  const loadServerStatus = async () => {
		try {
		  const response = await axios.get(`${API_BASE}/api/health`)
		  serverStatus.port = response.data.port
		  serverStatus.grpcServer = response.data.grpcServer
		} catch (err) {
		  console.warn('Could not load server status:', err)
		}
	  }

	  onMounted(() => {
		loadUsers()
		loadServerStatus()
	  })

	  return {
		users,
		isLoading,
		isCreating,
		error,
		message,
		messageType,
		newUser,
		serverStatus,
		loadUsers,
		createUser
	  }
	}
  }
  </script>
