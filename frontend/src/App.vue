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

  <style>
  #app {
	font-family: Avenir, Helvetica, Arial, sans-serif;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	min-height: 100vh;
	background-color: #f5f5f5;
  }

  .container {
	max-width: 1200px;
	margin: 0 auto;
	padding: 2rem 1rem;
  }

  h1 {
	color: #333;
	font-size: 2rem;
	margin-bottom: 2rem;
	text-align: center;
  }

  h2 {
	color: #333;
	font-size: 1.5rem;
	margin-bottom: 1rem;
  }

  .status-banner {
	background-color: #d4edda;
	border: 1px solid #c3e6cb;
	color: #155724;
	padding: 1rem;
	border-radius: 8px;
	margin-bottom: 2rem;
  }

  .form-section, .table-section {
	background: white;
	border-radius: 8px;
	padding: 2rem;
	box-shadow: 0 2px 4px rgba(0,0,0,0.1);
	margin-bottom: 2rem;
  }

  .form-group {
	margin-bottom: 1rem;
  }

  .form-group label {
	display: block;
	margin-bottom: 0.5rem;
	font-weight: 500;
	color: #333;
  }

  .form-group input {
	width: 100%;
	padding: 0.75rem;
	border: 1px solid #ddd;
	border-radius: 4px;
	font-size: 1rem;
	box-sizing: border-box;
  }

  .form-group input:focus {
	outline: none;
	border-color: #007bff;
	box-shadow: 0 0 0 2px rgba(0,123,255,0.25);
  }

  .btn {
	padding: 0.75rem 1.5rem;
	border: none;
	border-radius: 4px;
	font-size: 1rem;
	cursor: pointer;
	transition: background-color 0.2s;
  }

  .btn:disabled {
	opacity: 0.6;
	cursor: not-allowed;
  }

  .btn-primary {
	background-color: #007bff;
	color: white;
	width: 100%;
  }

  .btn-primary:hover:not(:disabled) {
	background-color: #0056b3;
  }

  .btn-secondary {
	background-color: #28a745;
	color: white;
  }

  .btn-secondary:hover:not(:disabled) {
	background-color: #1e7e34;
  }

  .table-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 1rem;
  }

  .table-container {
	overflow-x: auto;
  }

  .users-table {
	width: 100%;
	border-collapse: collapse;
	margin-top: 1rem;
  }

  .users-table th,
  .users-table td {
	padding: 1rem;
	text-align: left;
	border-bottom: 1px solid #ddd;
  }

  .users-table th {
	background-color: #f8f9fa;
	font-weight: 600;
	color: #333;
  }

  .users-table tr:hover {
	background-color: #f8f9fa;
  }

  .no-data {
	text-align: center;
	font-style: italic;
	color: #666;
  }

  .loading {
	display: flex;
	justify-content: center;
	align-items: center;
	padding: 2rem;
  }

  .spinner {
	width: 32px;
	height: 32px;
	border: 3px solid #f3f3f3;
	border-top: 3px solid #007bff;
	border-radius: 50%;
	animation: spin 1s linear infinite;
  }

  @keyframes spin {
	0% { transform: rotate(0deg); }
	100% { transform: rotate(360deg); }
  }

  .error-message {
	background-color: #f8d7da;
	border: 1px solid #f5c6cb;
	color: #721c24;
	padding: 1rem;
	border-radius: 4px;
  }

  .toast {
	position: fixed;
	bottom: 1rem;
	right: 1rem;
	padding: 1rem;
	border-radius: 4px;
	box-shadow: 0 4px 6px rgba(0,0,0,0.1);
	z-index: 1000;
	max-width: 300px;
  }

  .toast.success {
	background-color: #d4edda;
	border: 1px solid #c3e6cb;
	color: #155724;
  }

  .toast.error {
	background-color: #f8d7da;
	border: 1px solid #f5c6cb;
	color: #721c24;
  }

  @media (max-width: 768px) {
	.container {
	  padding: 1rem;
	}

	.table-header {
	  flex-direction: column;
	  gap: 1rem;
	}

	.btn-secondary {
	  width: 100%;
	}
  }
  </style>
