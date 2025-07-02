<template>
	<form @submit.prevent="submit_createUser">
		<div>
			<label for="name" >Name</label>
			<input
			type="text"
			id="name"
			v-model="newUser.name"
			required
			placeholder="Name eingeben"
			/>
		</div>

		<div>
			<label for="email" >Email</label>
			<input
			type="email"
			id="email"
			v-model="newUser.email"
			required
			placeholder="email@example.com"
			/>
		</div>

		<div>
			<label for="age" >Alter</label>
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

		<div>
			<label for="role" >Rolle</label>
			<select
			id="role"
			v-model="newUser.role"
			>
			<option value="user">User</option>
			<option value="admin">Admin</option>
			</select>
		</div>

		<button
			type="submit"
			:disabled="isCreating"
		>
			<svg v-if="isCreating" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
			<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
			<path  fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			{{ isCreating ? 'Erstelle User...' : 'User erstellen' }}
		</button>
	</form>

	<div>
		<div>
			<h2>Alle Users</h2>
			<button
			@click="loadUsers"
			:disabled="isLoading"

			>
			<svg v-if="isLoading" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle  cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			{{ isLoading ? 'Laden...' : 'Aktualisieren' }}
			</button>
		</div>

		<div>
			<div v-if="isLoading && !hasUsers">
			<div>
				<svg  xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				<span>Lade Users...</span>
			</div>
			</div>

			<div v-else-if="!hasUsers && !isLoading">
			<svg stroke="currentColor" fill="none" viewBox="0 0 48 48">
				<path d="M34 40h10v-4a6 6 0 00-10.712-3.714M34 40H14m20 0v-4a9.971 9.971 0 00-.712-3.714M14 40H4v-4a6 6 0 0110.713-3.714M14 40v-4c0-1.313.253-2.566.713-3.714m0 0A10.003 10.003 0 0124 26c4.21 0 7.813 2.602 9.288 6.286M30 14a6 6 0 11-12 0 6 6 0 0112 0z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
			</svg>
			<h3>Keine Users gefunden</h3>
			<p >Erstellen Sie den ersten User.</p>
			</div>
			<div v-else>
				<table>
					<thead>
					<tr>
						<th>User</th>
						<th>Email</th>
						<th>Alter</th>
						<th>Rolle</th>
					</tr>
					</thead>
					<tbody >
					<tr v-for="user in users" :key="user.id">
						<td >
						<div >
							<div >
							<div >
								<span>{{ user.name.charAt(0).toUpperCase() }}</span>
							</div>
							</div>
							<div >
							<div >{{ user.name }}</div>
							<div >ID: {{ user.id }}</div>
							</div>
						</div>
						</td>
						<td>{{ user.email }}</td>
						<td>{{ user.age }}</td>
						<td>
						<span>
							{{ user.role }}
						</span>
						</td>
					</tr>
					</tbody>
				</table>
			</div>
		</div>
	</div>
  </template>

  <script>

  export default {
	name: 'Welcome',
	data() {
		return {
		newUser: {
			name: '',
			email: '',
			age: '',
			role: 'user'
		}
		}
	},
	// computed: {
	// 	...mapGetters('users', [
	// 	'users',
	// 	'hasUsers',
	// 	'userCount',
	// 	'isLoading',
	// 	'isCreating',
	// 	'error',
	// 	'adminUsers',
	// 	'regularUsers'
	// 	])
	// },
	methods: {
		// ...mapActions('users', [
		// 'fetchUsers',
		// 'createUser',
		// 'clearError'
		// ]),

		async loadUsers() {
		try {
			//await this.fetchUsers()
		} catch (error) {
			console.error('Fehler beim Laden der Users:', error)
		}
		},

		async submit_createUser() {
			console.log("createUser")
			try {
			await this.createUser(this.newUser)

			this.newUser = {
			name: '',
			email: '',
			age: '',
			role: 'user'
			}

			this.$nextTick(() => {
			console.log('User erfolgreich erstellt!')
			})
		} catch (error) {
			console.error('Fehler beim Erstellen des Users:', error)
		}
		}
	},

	async mounted() {
		await this.loadUsers()
	}
  }
  </script>