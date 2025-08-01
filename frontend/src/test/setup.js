import { expect, vi } from 'vitest'
import '@testing-library/jest-dom'

// Global Mocks für gRPC-Web
vi.mock('@/proto/user_grpc_web_pb', () => ({
  UserServiceClient: vi.fn(() => ({
    listUsers: vi.fn(),
    createUser: vi.fn(),
    getUser: vi.fn(),
    deleteUser: vi.fn()
  }))
}))

vi.mock('@/proto/notification_grpc_web_pb', () => ({
  NotificationServiceClient: vi.fn(() => ({
    listNotifications: vi.fn(),
    createNotification: vi.fn(),
    deleteNotification: vi.fn(),
    markNotificationAsRead: vi.fn()
  }))
}))

// Proto Message Mocks
vi.mock('@/proto/user_pb', () => ({
  CreateUserRequest: vi.fn(() => ({
    setName: vi.fn(),
    setEmail: vi.fn(),
    setAge: vi.fn(),
    setRole: vi.fn()
  })),
  ListUsersRequest: vi.fn(),
  GetUserRequest: vi.fn(() => ({
    setId: vi.fn()
  })),
  DeleteUserRequest: vi.fn(() => ({
    setId: vi.fn()
  }))
}))

vi.mock('@/proto/notification_pb', () => ({
  CreateNotificationRequest: vi.fn(() => ({
    setMessage: vi.fn(),
    setType: vi.fn(),
    setPersistent: vi.fn(),
    setUserId: vi.fn()
  })),
  ListNotificationsRequest: vi.fn(),
  DeleteNotificationRequest: vi.fn(() => ({
    setId: vi.fn()
  })),
  MarkNotificationAsReadRequest: vi.fn(() => ({
    setId: vi.fn(),
    setUserId: vi.fn()
  }))
}))

// Environment Variables für Tests
Object.assign(process.env, {
  VUE_APP_GRPC_WEB_URL: 'http://localhost:8081',
  VUE_APP_WS_URL: 'ws://localhost:8082',
  NODE_ENV: 'test'
})

// WebSocket Mock für Tests
global.WebSocket = vi.fn(() => ({
  onopen: vi.fn(),
  onclose: vi.fn(),
  onmessage: vi.fn(),
  onerror: vi.fn(),
  send: vi.fn(),
  close: vi.fn(),
  readyState: 1
}))

// Router Mock
export const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
  go: vi.fn(),
  back: vi.fn(),
  forward: vi.fn(),
  currentRoute: {
    value: {
      path: '/',
      params: {},
      query: {},
      meta: {}
    }
  }
}

// Store Mock Helpers
export const createMockStore = (modules = {}) => {
  const store = {
    state: {},
    dispatch: vi.fn(),
    commit: vi.fn(),
    getters: {},
    modules: {},
    _actions: {},
    _mutations: {}
  }

  // Setup modules with proper action structure
  Object.keys(modules).forEach(moduleName => {
    const module = modules[moduleName]

    // Setup state
    if (module.state) {
      store.state[moduleName] = typeof module.state === 'function' ? module.state() : module.state
    }

    // Setup actions with Vuex internal structure
    if (module.actions) {
      Object.keys(module.actions).forEach(actionName => {
        const fullActionName = module.namespaced ? `${moduleName}/${actionName}` : actionName
        // Vuex stores actions as arrays of functions
        store._actions[fullActionName] = [module.actions[actionName]]
      })
    }

    // Setup mutations
    if (module.mutations) {
      Object.keys(module.mutations).forEach(mutationName => {
        const fullMutationName = module.namespaced ? `${moduleName}/${mutationName}` : mutationName
        store._mutations[fullMutationName] = [module.mutations[mutationName]]
      })
    }

    // Setup getters
    if (module.getters) {
      Object.keys(module.getters).forEach(getterName => {
        const fullGetterName = module.namespaced ? `${moduleName}/${getterName}` : getterName
        store.getters[fullGetterName] = module.getters[getterName]
      })
    }
  })

  return store
}

// Test Utilities
export const flushPromises = () => new Promise(resolve => setTimeout(resolve, 0))

// Custom Matchers
expect.extend({
  toHaveBeenCalledWithPayload(received, expectedPayload) {
    const calls = received.mock.calls
    const hasPayload = calls.some(call =>
      JSON.stringify(call[0]) === JSON.stringify(expectedPayload)
    )

    return {
      message: () => `expected function to ${this.isNot ? 'not ' : ''}have been called with payload ${JSON.stringify(expectedPayload)}`,
      pass: hasPayload
    }
  }
})

// Mock console methods to reduce noise in tests
global.console = {
  ...console,
  log: vi.fn(),
  warn: vi.fn(),
  error: vi.fn()
}

// Cleanup nach jedem Test
afterEach(() => {
  vi.clearAllMocks()
  vi.clearAllTimers()
})
