<template>
  <div class="user-list">
    <h2>Benutzer Verwaltung</h2>

    <div class="add-user">
      <h3>Neuen Benutzer hinzufügen</h3>
      <form @submit.prevent="addUser">
        <input v-model="newUser.name" type="text" placeholder="Name" required />
        <input v-model="newUser.email" type="email" placeholder="Email" required />
        <button type="submit">Hinzufügen</button>
      </form>
    </div>

    <div class="users">
      <h3>Alle Benutzer</h3>
      <div v-if="loading" class="loading">Lade Benutzer...</div>
      <div v-else-if="users.length === 0" class="no-users">Keine Benutzer gefunden</div>
      <div v-else>
        <div v-for="user in users" :key="user.id" class="user-card">
          <div class="user-info">
            <h4>{{ user.name }}</h4>
            <p>{{ user.email }}</p>
            <small>Erstellt: {{ formatDate(user.created_at) }}</small>
          </div>
          <div class="user-actions">
            <button @click="deleteUser(user.id)" class="delete-btn">Löschen</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/api.js'

export default {
  name: 'UserList',
  data() {
    return {
      users: [],
      loading: false,
      newUser: {
        name: '',
        email: ''
      }
    }
  },
  async mounted() {
    await this.fetchUsers()
  },
  methods: {
    async fetchUsers() {
      this.loading = true
      try {
        const response = await api.get('/users')
        this.users = response.data || []
      } catch (error) {
        console.error('Fehler beim Laden der Benutzer:', error)
        alert('Fehler beim Laden der Benutzer')
      } finally {
        this.loading = false
      }
    },
    async addUser() {
      try {
        await api.post('/users', this.newUser)
        this.newUser = { name: '', email: '' }
        await this.fetchUsers()
        alert('Benutzer erfolgreich hinzugefügt!')
      } catch (error) {
        console.error('Fehler beim Hinzufügen des Benutzers:', error)
        alert('Fehler beim Hinzufügen des Benutzers')
      }
    },
    async deleteUser(id) {
      if (confirm('Sind Sie sicher, dass Sie diesen Benutzer löschen möchten?')) {
        try {
          await api.delete(`/users/${id}`)
          await this.fetchUsers()
          alert('Benutzer erfolgreich gelöscht!')
        } catch (error) {
          console.error('Fehler beim Löschen des Benutzers:', error)
          alert('Fehler beim Löschen des Benutzers')
        }
      }
    },
    formatDate(dateString) {
      return new Date(dateString).toLocaleDateString('de-DE')
    }
  }
}
</script>

<style scoped>
.user-list {
  max-width: 800px;
  margin: 0 auto;
}

.add-user {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
}

.add-user form {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.add-user input {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.add-user button {
  padding: 0.5rem 1rem;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.user-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.user-info h4 {
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
}

.user-info p {
  margin: 0 0 0.25rem 0;
  color: #666;
}

.delete-btn {
  padding: 0.25rem 0.75rem;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.loading, .no-users {
  text-align: center;
  padding: 2rem;
  color: #666;
}
</style>