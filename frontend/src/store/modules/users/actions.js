import axios from 'axios'

const API_BASE = process.env.NODE_ENV === 'production' ? '' : 'http://localhost:3000'

export default {
	async fetchUsers({ commit }) {
		commit('SET_LOADING', true)
		commit('SET_ERROR', null)

		try {
			const response = await axios.get(`${API_BASE}/api/users`)
			commit('SET_USERS', response.data)
			return response.data
		} catch (error) {
			const errorMessage = error.response?.data?.error || error.message
			commit('SET_ERROR', `Fehler beim Laden der Users: ${errorMessage}`)
			throw error
		} finally {
			commit('SET_LOADING', false)
		}
	},

	async createUser({ commit, dispatch }, userData) {
		commit('SET_CREATING', true)
		commit('SET_ERROR', null)

		try {
			const response = await axios.post(`${API_BASE}/api/users`, {
				name: userData.name,
				email: userData.email,
				age: parseInt(userData.age),
				role: userData.role || 'user'
			})

			// Nach dem Erstellen die Liste neu laden
			await dispatch('fetchUsers')
			return response.data
		} catch (error) {
			const errorMessage = error.response?.data?.error || error.message
			commit('SET_ERROR', `Fehler beim Erstellen des Users: ${errorMessage}`)
			throw error
		} finally {
			commit('SET_CREATING', false)
		}
	},

	async getUserById({ commit }, userId) {
		commit('SET_LOADING', true)
		commit('SET_ERROR', null)

		try {
			const response = await axios.get(`${API_BASE}/api/users/${userId}`)
			return response.data
		} catch (error) {
			const errorMessage = error.response?.data?.error || error.message
			commit('SET_ERROR', `Fehler beim Laden des Users: ${errorMessage}`)
			throw error
		} finally {
			commit('SET_LOADING', false)
		}
	},

	clearError({ commit }) {
		commit('SET_ERROR', null)
	},

	clearUsers({ commit }) {
		commit('SET_USERS', [])
	}
}