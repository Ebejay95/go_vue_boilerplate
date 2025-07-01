export default {
	users(state) {
		return state.users
	},

	hasUsers(state) {
		return state.users && state.users.length > 0
	},

	userCount(state) {
		return state.users ? state.users.length : 0
	},

	isLoading(state) {
		return state.loading
	},

	isCreating(state) {
		return state.creating
	},

	error(state) {
		return state.error
	},

	getUserById: (state) => (id) => {
		return state.users.find(user => user.id === id)
	},

	getUsersByRole: (state) => (role) => {
		return state.users.filter(user => user.role === role)
	},

	adminUsers(state) {
		return state.users.filter(user => user.role === 'admin')
	},

	regularUsers(state) {
		return state.users.filter(user => user.role === 'user')
	}
}