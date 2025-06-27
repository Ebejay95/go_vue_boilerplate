<template>
  <div class="user-list">
    <h2>Benutzer Verwaltung (gRPC)</h2>

    <div class="add-user">
      <h3>Neuen Benutzer hinzufügen</h3>
      <form @submit.prevent="addUser">
        <input v-model="newUser.name" type="text" placeholder="Name" required />
        <input v-model="newUser.email" type="email" placeholder="Email" required />
        <button type="submit" :disabled="isLoading">
          {{ isLoading ? 'Hinzufügen...' : 'Hinzufügen' }}
        </button>
      </form>
    </div>

    <div class="users">
      <h3>Alle Benutzer</h3>
      <div v-if="loading" class="loading">Lade Benutzer über gRPC...</div>
      <div v-else-if="error" class="error">
        Fehler: {{ error }}
        <button @click="fetchUsers" class="retry-btn">Erneut versuchen</button>
      </div>
      <div v-else-if="users.length === 0" class="no-users">Keine Benutzer gefunden</div>
      <div v-else>
        <div v-for="user in users" :key="user.id" class="user-card">
          <div class="user-info">
            <h4>{{ user.name }}</h4>
            <p>{{ user.email }}</p>
            <small>Erstellt: {{ formatDate(user.created_at) }}</small>
          </div>
          <div class="user-actions">
            <button @click="editUser(user)" class="edit-btn">Bearbeiten</button>
            <button @click="deleteUser(user.id)" class="delete-btn" :disabled="isLoading">
              Löschen
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit User Modal -->
    <div v-if="editingUser" class="modal-overlay" @click="cancelEdit">
      <div class="modal" @click.stop>
        <h3>Benutzer bearbeiten</h3>
        <form @submit.prevent="updateUser">
          <input v-model="editingUser.name" type="text" placeholder="Name" required />
          <input v-model="editingUser.email" type="email" placeholder="Email" required />
          <div class="modal-actions">
            <button type="submit" :disabled="isLoading">Speichern</button>
            <button type="button" @click="cancelEdit">Abbrechen</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import userService from '../services/userService.js'

export default {
  name: 'UserList',
  data() {
    return {
      users: [],
      loading: false,
      isLoading: false,
      error: null,
      newUser: {
        name: '',
        email: ''
      },
      editingUser: null
    }
  },
  async mounted() {
    await this.fetchUsers()
  },
  methods: {
    async fetchUsers() {
      this.loading = true
      this.error = null
      try {
        this.users = await userService.getUsers()
      } catch (error) {
        console.error('Fehler beim Laden der Benutzer:', error)
        this.error = error.message || 'Unbekannter Fehler beim Laden der Benutzer'
      } finally {
        this.loading = false
      }
    },

    async addUser() {
      if (this.isLoading) return

      this.isLoading = true
      try {
        await userService.createUser(this.newUser.name, this.newUser.email)
        this.newUser = { name: '', email: '' }
        await this.fetchUsers()
        this.showSuccess('Benutzer erfolgreich hinzugefügt!')
      } catch (error) {
        console.error('Fehler beim Hinzufügen des Benutzers:', error)
        this.showError('Fehler beim Hinzufügen des Benutzers: ' + error.message)
      } finally {
        this.isLoading = false
      }
    },

    editUser(user) {
      this.editingUser = { ...user }
    },

    cancelEdit() {
      this.editingUser = null
    },

    async updateUser() {
      if (this.isLoading) return

      this.isLoading = true
      try {
        await userService.updateUser(
          this.editingUser.id,
          this.editingUser.name,
          this.editingUser.email
        )
        this.editingUser = null
        await this.fetchUsers()
        this.showSuccess('Benutzer erfolgreich aktualisiert!')
      } catch (error) {
        console.error('Fehler beim Aktualisieren des Benutzers:', error)
        this.showError('Fehler beim Aktualisieren des Benutzers: ' + error.message)
      } finally {
        this.isLoading = false
      }
    },

    async deleteUser(id) {
      if (this.isLoading) return

      if (confirm('Sind Sie sicher, dass Sie diesen Benutzer löschen möchten?')) {
        this.isLoading = true
        try {
          await userService.deleteUser(id)
          await this.fetchUsers()
          this.showSuccess('Benutzer erfolgreich gelöscht!')
        } catch (error) {
          console.error('Fehler beim Löschen des Benutzers:', error)
          this.showError('Fehler beim Löschen des Benutzers: ' + error.message)
        } finally {
          this.isLoading = false
        }
      }
    },

    formatDate(date) {
      if (!date) return 'N/A'
      return new Date(date).toLocaleDateString('de-DE')
    },

    showSuccess(message) {
      alert(message) // In production, use a proper toast notification
    },

    showError(message) {
      alert(message) // In production, use a proper toast notification
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

.add-user button:disabled {
  background: #6c757d;
  cursor: not-allowed;
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

.user-actions {
  display: flex;
  gap: 0.5rem;
}

.edit-btn {
  padding: 0.25rem 0.75rem;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.delete-btn {
  padding: 0.25rem 0.75rem;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.delete-btn:disabled, .edit-btn:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.loading, .no-users {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.error {
  text-align: center;
  padding: 2rem;
  color: #dc3545;
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.retry-btn {
  margin-top: 1rem;
  padding: 0.5rem 1rem;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

/* Modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
}

.modal h3 {
  margin-top: 0;
  margin-bottom: 1rem;
}

.modal input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-bottom: 1rem;
  box-sizing: border-box;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.modal-actions button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.modal-actions button[type="submit"] {
  background: #007bff;
  color: white;
}

.modal-actions button[type="button"] {
  background: #6c757d;
  color: white;
}