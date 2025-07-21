// src/store/modules/user/actions.js
// 100% gRPC-Web - KEIN AXIOS mehr!

import {
	GetUserRequest,
	CreateUserRequest,
	ListUsersRequest,
	DeleteUserRequest
} from '../../../proto/user_pb'

export default {
	async fetchUsers({ commit, dispatch, rootGetters }) {
		commit('SET_LOADING', true)
		commit('SET_ERROR', null)

		try {
			// Check if connection is available
			if (!rootGetters['connection/isConnected']) {
				throw new Error('No gRPC connection available. Please check your connection.')
			}

			const request = new ListUsersRequest()
			const response = await dispatch('connection/promisifyGrpcCall', {
				method: rootGetters['connection/grpcClient'].listUsers,
				request: request
			}, { root: true })

			const users = response.getUsersList().map(user => ({
				id: user.getId(),
				name: user.getName(),
				email: user.getEmail(),
				age: user.getAge(),
				role: user.getRole()
			}))

			commit('SET_USERS', users)

			return { users, count: users.length }

		} catch (error) {
			console.error('Error fetching users:', error)
			const errorMessage = error.message || 'Failed to fetch users'
			commit('SET_ERROR', `Error loading users: ${errorMessage}`)
			throw error
		} finally {
			commit('SET_LOADING', false)
		}
	},

	async createUser({ commit, dispatch, rootGetters }, userData) {
		commit('SET_CREATING', true)
		commit('SET_ERROR', null)

		try {
			// Check if connection is available
			if (!rootGetters['connection/isConnected']) {
				throw new Error('No gRPC connection available. Please check your connection.')
			}

			const request = new CreateUserRequest()
			request.setName(userData.name)
			request.setEmail(userData.email)
			request.setAge(parseInt(userData.age))
			request.setRole(userData.role || 'user')

			const response = await dispatch('connection/promisifyGrpcCall', {
				method: rootGetters['connection/grpcClient'].createUser,
				request: request
			}, { root: true })

			const user = response.getUser()
			const createdUser = {
				id: user.getId(),
				name: user.getName(),
				email: user.getEmail(),
				age: user.getAge(),
				role: user.getRole()
			}

			// Add to local state
			commit('ADD_USER', createdUser)

			return { user: createdUser }

		} catch (error) {
			console.error('Error creating user:', error)
			const errorMessage = error.message || 'Failed to create user'
			commit('SET_ERROR', `Error creating user: ${errorMessage}`)
			throw error
		} finally {
			commit('SET_CREATING', false)
		}
	},

	async deleteUser({ commit, dispatch, rootGetters }, userId) {
		commit('SET_LOADING', true)
		commit('SET_ERROR', null)

		try {
			// Check if connection is available
			if (!rootGetters['connection/isConnected']) {
				throw new Error('No gRPC connection available. Please check your connection.')
			}

			const request = new DeleteUserRequest()
			request.setId(parseInt(userId))

			const response = await dispatch('connection/promisifyGrpcCall', {
				method: rootGetters['connection/grpcClient'].deleteUser,
				request: request
			}, { root: true })

			// Check if deletion was successful
			if (response.getSuccess()) {
				// Remove from local state
				commit('REMOVE_USER', parseInt(userId))
				return { success: true, message: response.getMessage() }
			} else {
				throw new Error(response.getMessage() || 'Delete failed')
			}

		} catch (error) {
			console.error('Error deleting user:', error)
			const errorMessage = error.message || 'Failed to delete user'
			commit('SET_ERROR', `Error deleting user: ${errorMessage}`)
			throw error
		} finally {
			commit('SET_LOADING', false)
		}
	},
	
	async getUserById({ commit, dispatch, rootGetters }, userId) {
		commit('SET_LOADING', true)
		commit('SET_ERROR', null)

		try {
			// Check if connection is available
			if (!rootGetters['connection/isConnected']) {
				throw new Error('No gRPC connection available. Please check your connection.')
			}

			const request = new GetUserRequest()
			request.setId(parseInt(userId))

			const response = await dispatch('connection/promisifyGrpcCall', {
				method: rootGetters['connection/grpcClient'].getUser,
				request: request
			}, { root: true })

			const user = response.getUser()
			const userData = {
				id: user.getId(),
				name: user.getName(),
				email: user.getEmail(),
				age: user.getAge(),
				role: user.getRole()
			}

			return { user: userData }

		} catch (error) {
			console.error('Error getting user:', error)
			const errorMessage = error.message || 'Failed to get user'
			commit('SET_ERROR', `Error loading user: ${errorMessage}`)
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
