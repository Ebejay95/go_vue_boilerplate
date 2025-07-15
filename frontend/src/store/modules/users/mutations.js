export default {
	SET_USERS(state, users) {
		state.users = users || []
	},

	SET_LOADING(state, loading) {
		state.loading = loading
	},

	SET_CREATING(state, creating) {
		state.creating = creating
	},

	SET_ERROR(state, error) {
		state.error = error
	},

	ADD_USER(state, user) {
		state.users.push(user)
	},

	UPDATE_USER(state, updatedUser) {
		const index = state.users.findIndex(user => user.id === updatedUser.id)
		if (index !== -1) {
			state.users.splice(index, 1, updatedUser)
		}
	},

	REMOVE_USER(state, userId) {
		state.users = state.users.filter(user => user.id !== userId)
	}
}
