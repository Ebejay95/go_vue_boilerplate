// src/services/grpcClient.js
import { UserServiceClient } from '../proto/user_grpc_web_pb'
import {
  GetUserRequest,
  CreateUserRequest,
  ListUsersRequest
} from '../proto/user_pb'

// ==========================================
// ECHTER gRPC-Web CLIENT (direkt zu gRPC Server)
// ==========================================
class DirectGrpcWebClient {
  constructor() {
    // gRPC-Web Client - muss vom Browser erreichbar sein!
    // Nicht backend-grpc-server:8081 (container-intern), sondern localhost:8081 (host)
    const grpcWebUrl = process.env.VUE_APP_GRPC_WEB_URL || 'http://localhost:8081'

    console.log(`🔗 Direct gRPC-Web Client connecting to: ${grpcWebUrl}`)
    console.log(`🌐 Running in browser - using host-accessible URL`)

    this.client = new UserServiceClient(grpcWebUrl, null, null)
    this.grpcWebUrl = grpcWebUrl
  }

  // Promise wrapper für gRPC Calls mit verbessertem Error Handling
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

          // Bessere Error Messages für häufige Probleme
          let userFriendlyMessage = err.message || 'gRPC-Web connection failed'

          if (err.code === 14 || err.message?.includes('UNAVAILABLE')) {
            userFriendlyMessage = `gRPC server not reachable at ${this.grpcWebUrl}. Is the backend running?`
          } else if (err.code === 12 || err.message?.includes('UNIMPLEMENTED')) {
            userFriendlyMessage = 'gRPC method not implemented on server'
          } else if (err.message?.includes('fetch')) {
            userFriendlyMessage = `Network error: Cannot connect to ${this.grpcWebUrl}. Check if port 8081 is accessible.`
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

  // Health check using ListUsers
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
// REST-LIKE CLIENT (über Express Proxy)
// ==========================================
class RestLikeGrpcClient {
  constructor() {
    this.baseUrl = process.env.VUE_APP_API_BASE_URL || 'http://localhost:3000/api'
    console.log(`🌐 REST-like Client connecting to: ${this.baseUrl}`)
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      ...options
    }

    console.log(`📡 REST API Call: ${config.method || 'GET'} ${url}`)

    try {
      const response = await fetch(url, config)

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`HTTP ${response.status}: ${errorText}`)
      }

      const data = await response.json()
      console.log(`✅ REST API Response:`, data)
      return data
    } catch (error) {
      console.error(`❌ REST API Error:`, error)
      throw error
    }
  }

  async checkHealth() {
    try {
      const response = await this.request('/health')
      return {
        status: 'healthy',
        connection: 'rest-api',
        url: this.baseUrl,
        ...response
      }
    } catch (error) {
      return {
        status: 'unhealthy',
        connection: 'rest-api',
        error: error.message || 'REST API connection failed',
        url: this.baseUrl
      }
    }
  }

  async listUsers() {
    return await this.request('/users')
  }

  async getUser(id) {
    return await this.request(`/users/${id}`)
  }

  async createUser(userData) {
    return await this.request('/users', {
      method: 'POST',
      body: JSON.stringify(userData)
    })
  }
}

// ==========================================
// EXPORTS
// ==========================================

// Export instances
export const simpleGrpcClient = new DirectGrpcWebClient()
export const restLikeGrpcClient = new RestLikeGrpcClient()

// Default export für backwards compatibility
export default simpleGrpcClient

// Debug info
console.log('🔧 gRPC Clients initialized:')
console.log(`  - simpleGrpcClient: Direct gRPC-Web to ${simpleGrpcClient.grpcWebUrl}`)
console.log(`  - restLikeGrpcClient: REST API Proxy to ${restLikeGrpcClient.baseUrl}`)
console.log('🌐 Note: gRPC-Web URL must be accessible from browser (host), not container network!')
