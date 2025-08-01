describe('Notification System', () => {
	beforeEach(() => {
	  // Mock WebSocket connection
	  cy.window().then((win) => {
		win.WebSocket = class MockWebSocket {
		  constructor(url) {
			this.url = url
			this.readyState = 1
			setTimeout(() => this.onopen?.(), 100)
		  }
		  send() {}
		  close() {}
		}
	  })

	  cy.intercept('POST', '**/notification.NotificationService/ListNotifications', {
		fixture: 'notifications-list.json'
	  }).as('listNotifications')

	  cy.visit('/')
	})

	it('displays notification bell with badge', () => {
	  cy.get('[data-testid="notification-bell"]').should('be.visible')
	  cy.get('.bg-red-500').should('contain', '3') // Badge count
	})

	it('opens notification panel', () => {
	  cy.get('[data-testid="notification-bell"]').click()
	  cy.get('[data-testid="notification-panel"]').should('be.visible')
	  cy.get('h3').should('contain', 'Benachrichtigungen')
	})

	it('displays notifications in panel', () => {
	  cy.get('[data-testid="notification-bell"]').click()
	  cy.wait('@listNotifications')

	  cy.get('[data-testid="notification-item"]').should('have.length.greaterThan', 0)
	  cy.get('[data-testid="notification-item"]').first().should('contain', 'Test Notification')
	})

	it('marks notification as read', () => {
	  cy.intercept('POST', '**/notification.NotificationService/MarkNotificationAsRead', {
		body: { success: true }
	  }).as('markAsRead')

	  cy.get('[data-testid="notification-bell"]').click()
	  cy.get('[data-testid="notification-item"]').first().click()

	  cy.wait('@markAsRead')

	  // Check that notification is marked as read (visual change)
	  cy.get('[data-testid="notification-item"]').first()
		.should('have.class', 'bg-gray-50')
	})

	it('deletes notification', () => {
	  cy.intercept('POST', '**/notification.NotificationService/DeleteNotification', {
		body: { success: true }
	  }).as('deleteNotification')

	  cy.get('[data-testid="notification-bell"]').click()
	  cy.get('[data-testid="delete-notification"]').first().click()

	  cy.wait('@deleteNotification')

	  // Verify notification is removed from list
	  cy.get('[data-testid="notification-item"]').should('have.length', 2)
	})

	it('shows toast notifications', () => {
	  // Trigger a success notification
	  cy.get('button').contains('Benutzer erstellen').click()

	  // Check toast appears
	  cy.get('.toast').should('be.visible')
	  cy.get('.toast').should('contain', 'erfolgreich')

	  // Check toast disappears after timeout
	  cy.get('.toast', { timeout: 6000 }).should('not.exist')
	})

	it('handles toast actions', () => {
	  // Create a notification with action
	  cy.window().its('store').invoke('dispatch', 'notifications/error', {
		message: 'Connection failed',
		actions: [{ label: 'Retry', handler: cy.stub().as('retryHandler') }]
	  })

	  cy.get('.toast').should('be.visible')
	  cy.get('button').contains('Retry').click()

	  cy.get('@retryHandler').should('have.been.called')
	})
})
