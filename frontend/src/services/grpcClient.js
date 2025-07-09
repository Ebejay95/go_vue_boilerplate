// src/services/grpcClient.js
import { UserServiceClient } from '../proto/user_grpc_web_pb'
import {
  GetUserRequest,
  CreateUserRequest,
  ListUsersRequest
} from '../proto/user_pb'

// ==========================================
// DIREKTER gRPC-Web CLIENT
// ==========================================
class DirectGrpcWebClient {
  constructor() {
    const grpcWebUrl = process.env.VUE_APP_GRPC_WEB_URL
    this.client = new UserServiceClient(grpcWebUrl, null, null)
    this.grpcWebUrl = grpcWebUrl
  }

  // Promise wrapper für gRPC Calls
  promisify(method, request) {
    return new Promise((resolve, reject) => {
      const call = method.call(this.client, request, {}, (err, response) => {
        if (err) {
          console.error('❌ gRPC-Web Error Details:', {
            code: err.code,
            message: err.message,
            details: err.details,
            metadata: err.metadata
          })

          // Benutzerfreundliche Error Messages
          let userFriendlyMessage = err.message || 'gRPC-Web connection failed'

          if (err.code === 14 || err.message?.includes('UNAVAILABLE')) {
            userFriendlyMessage = `gRPC server not reachable at ${this.grpcWebUrl}. Is the backend running?`
          } else if (err.code === 12 || err.message?.includes('UNIMPLEMENTED')) {
            userFriendlyMessage = 'gRPC method not implemented on server'
          }

          reject(new Error(userFriendlyMessage))
        } else {
          resolve(response)
        }
      })

      // Timeout nach 10 Sekunden
      setTimeout(() => {
        if (call && typeof call.cancel === 'function') {
          call.cancel()
        }
        reject(new Error(`gRPC call timeout after 10 seconds to ${this.grpcWebUrl}`))
      }, 10000)
    })
  }

  // Health check
  async checkHealth() {
    try {
      console.log(`🏥 gRPC-Web Health Check to ${this.grpcWebUrl}...`)
      const response = await this.listUsers()
      return {
        status: 'healthy',
        connection: 'grpc-web',
        users: response.users.length,
        server: 'direct-grpc-web',
        url: this.grpcWebUrl
      }
    } catch (error) {
      console.error('❌ gRPC-Web Health check failed:', error)
      return {
        status: 'unhealthy',
        connection: 'grpc-web',
        error: error.message || 'gRPC-Web connection failed',
        url: this.grpcWebUrl
      }
    }
  }

  // List all users
  async listUsers() {
    console.log('📋 gRPC-Web: ListUsers called')
    try {
      const request = new ListUsersRequest()
      const response = await this.promisify(this.client.listUsers, request)

      const users = response.getUsersList().map(user => ({
        id: user.getId(),
        name: user.getName(),
        email: user.getEmail(),
        age: user.getAge(),
        role: user.getRole()
      }))

      console.log(`✅ gRPC-Web: Found ${users.length} users`)
      return { users, count: users.length }
    } catch (error) {
      console.error('❌ gRPC-Web ListUsers failed:', error)
      throw error
    }
  }

  // Get specific user
  async getUser(id) {
    console.log(`👤 gRPC-Web: GetUser called for ID ${id}`)
    try {
      const request = new GetUserRequest()
      request.setId(id)

      const response = await this.promisify(this.client.getUser, request)
      const user = response.getUser()

      const userData = {
        id: user.getId(),
        name: user.getName(),
        email: user.getEmail(),
        age: user.getAge(),
        role: user.getRole()
      }

      console.log(`✅ gRPC-Web: Found user ${userData.name}`)
      return { user: userData }
    } catch (error) {
      console.error('❌ gRPC-Web GetUser failed:', error)
      throw error
    }
  }

  // Create new user
  async createUser(userData) {
    console.log(`➕ gRPC-Web: CreateUser called for ${userData.name}`)
    try {
      const request = new CreateUserRequest()
      request.setName(userData.name)
      request.setEmail(userData.email)
      request.setAge(userData.age)
      request.setRole(userData.role || 'user')

      const response = await this.promisify(this.client.createUser, request)
      const user = response.getUser()

      const createdUser = {
        id: user.getId(),
        name: user.getName(),
        email: user.getEmail(),
        age: user.getAge(),
        role: user.getRole()
      }

      console.log(`✅ gRPC-Web: Created user ${createdUser.name} with ID ${createdUser.id}`)
      return { user: createdUser }
    } catch (error) {
      console.error('❌ gRPC-Web CreateUser failed:', error)
      throw error
    }
  }
}

// ==========================================
// EXPORTS
// ==========================================

// Hauptclient exportieren
export const grpcClient = new DirectGrpcWebClient()

// Für backwards compatibility
export const simpleGrpcClient = grpcClient

// Default export
export default grpcClient
