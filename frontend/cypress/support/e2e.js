import './commands'
import 'cypress-axe'
import '@percy/cypress'

// Global error handling
Cypress.on('uncaught:exception', (err, runnable) => {
  // Don't fail tests on unhandled promise rejections
  if (err.message.includes('ResizeObserver loop limit exceeded')) {
    return false
  }
  return true
})

// Global before hook
beforeEach(() => {
  // Set up common interceptors
  cy.intercept('GET', '**/health', { statusCode: 200, body: { status: 'healthy' } })

  // Mock WebSocket connection
  cy.window().then((win) => {
    if (!win.WebSocket.isMocked) {
      const OriginalWebSocket = win.WebSocket
      win.WebSocket = class MockWebSocket extends OriginalWebSocket {
        constructor(url) {
          super(url)
          this.isMocked = true
          // Auto-trigger onopen for tests
          setTimeout(() => {
            if (this.onopen) this.onopen(new Event('open'))
          }, 100)
        }
      }
      win.WebSocket.isMocked = true
    }
  })
})
