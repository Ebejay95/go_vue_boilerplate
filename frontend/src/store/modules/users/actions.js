// src/store/modules/users/actions.js
// 100% gRPC-Web - KEIN AXIOS mehr!

import {
	GetUserRequest,
	CreateUserRequest,
	ListUsersRequest
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

		console.log('📋 Fetching users via gRPC-Web...')

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
		console.log(`✅ Loaded ${users.length} users via gRPC-Web`)

		return { users, count: users.length }

	  } catch (error) {
		console.error('❌ Error fetching users:', error)
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

		console.log(`➕ Creating user via gRPC-Web: ${userData.name}`)

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

		console.log(`✅ User created: ${createdUser.name} with ID ${createdUser.id}`)

		return { user: createdUser }

	  } catch (error) {
		console.error('❌ Error creating user:', error)
		const errorMessage = error.message || 'Failed to create user'
		commit('SET_ERROR', `Error creating user: ${errorMessage}`)
		throw error
	  } finally {
		commit('SET_CREATING', false)
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

		console.log(`👤 Getting user by ID via gRPC-Web: ${userId}`)

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

		console.log(`✅ Found user: ${userData.name}`)

		return { user: userData }

	  } catch (error) {
		console.error('❌ Error getting user:', error)
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
