// userService.js - gRPC-Web client for user operations

import { UserServiceClient } from './proto/user_grpc_web_pb.js'
import {
	GetUsersRequest,
	GetUserRequest,
	CreateUserRequest,
	UpdateUserRequest,
	DeleteUserRequest
} from './proto/user_pb.js'

class UserService {
	constructor() {
		// Initialize gRPC-Web client
		this.client = new UserServiceClient('http://localhost:8080', null, null)
	}

	async getUsers() {
		return new Promise((resolve, reject) => {
			const request = new GetUsersRequest()

			this.client.getUsers(request, {}, (err, response) => {
				if (err) {
					console.error('Error fetching users:', err)
					reject(err)
					return
				}

				const users = response.getUsersList().map(user => ({
					id: user.getId(),
					name: user.getName(),
					email: user.getEmail(),
					created_at: user.getCreatedAt()?.toDate(),
					updated_at: user.getUpdatedAt()?.toDate()
				}))

				resolve(users)
			})
		})
	}

	async getUser(id) {
		return new Promise((resolve, reject) => {
			const request = new GetUserRequest()
			request.setId(id)

			this.client.getUser(request, {}, (err, response) => {
				if (err) {
					console.error('Error fetching user:', err)
					reject(err)
					return
				}

				const user = response.getUser()
				resolve({
					id: user.getId(),
					name: user.getName(),
					email: user.getEmail(),
					created_at: user.getCreatedAt()?.toDate(),
					updated_at: user.getUpdatedAt()?.toDate()
				})
			})
		})
	}

	async createUser(name, email) {
		return new Promise((resolve, reject) => {
			const request = new CreateUserRequest()
			request.setName(name)
			request.setEmail(email)

			this.client.createUser(request, {}, (err, response) => {
				if (err) {
					console.error('Error creating user:', err)
					reject(err)
					return
				}

				const user = response.getUser()
				resolve({
					id: user.getId(),
					name: user.getName(),
					email: user.getEmail(),
					created_at: user.getCreatedAt()?.toDate(),
					updated_at: user.getUpdatedAt()?.toDate()
				})
			})
		})
	}

	async updateUser(id, name, email) {
		return new Promise((resolve, reject) => {
			const request = new UpdateUserRequest()
			request.setId(id)
			request.setName(name)
			request.setEmail(email)

			this.client.updateUser(request, {}, (err, response) => {
				if (err) {
					console.error('Error updating user:', err)
					reject(err)
					return
				}

				const user = response.getUser()
				resolve({
					id: user.getId(),
					name: user.getName(),
					email: user.getEmail(),
					created_at: user.getCreatedAt()?.toDate(),
					updated_at: user.getUpdatedAt()?.toDate()
				})
			})
		})
	}

	async deleteUser(id) {
		return new Promise((resolve, reject) => {
			const request = new DeleteUserRequest()
			request.setId(id)

			this.client.deleteUser(request, {}, (err, response) => {
				if (err) {
					console.error('Error deleting user:', err)
					reject(err)
					return
				}

				resolve({
					message: response.getMessage()
				})
			})
		})
	}
}

export default new UserService()