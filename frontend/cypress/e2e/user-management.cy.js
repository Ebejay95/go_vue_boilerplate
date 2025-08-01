describe('User Management', () => {
	beforeEach(() => {
	  // Intercept gRPC-Web calls
	  cy.intercept('POST', '**/user.UserService/ListUsers', {
		fixture: 'users-list.json'
	  }).as('listUsers')

	  cy.intercept('POST', '**/user.UserService/CreateUser', {
		fixture: 'user-created.json'
	  }).as('createUser')

	  cy.intercept('POST', '**/user.UserService/DeleteUser', {
		body: { success: true, message: 'User deleted successfully' }
	  }).as('deleteUser')

	  // Visit the main page
	  cy.visit('/')
	})

	it('displays user creation form', () => {
	  cy.get('h2').should('contain', 'Neuen Benutzer erstellen')
	  cy.get('form').should('be.visible')

	  // Check form fields
	  cy.get('input[type="text"]').should('be.visible') // Name field
	  cy.get('input[type="email"]').should('be.visible') // Email field
	  cy.get('input[type="range"]').should('be.visible') // Age field
	  cy.get('select').should('be.visible') // Role field
	})

	it('creates a new user successfully', () => {
	  // Fill out the form
	  cy.get('input[type="text"]').type('John Doe')
	  cy.get('input[type="email"]').type('john.doe@example.com')
	  cy.get('input[type="range"]').invoke('val', 30).trigger('input')
	  cy.get('select').select('user')

	  // Submit the form
	  cy.get('button[type="submit"]').click()

	  // Wait for API call
	  cy.wait('@createUser')

	  // Check success notification
	  cy.get('.toast').should('contain', 'erfolgreich erstellt')

	  // Verify user appears in list
	  cy.wait('@listUsers')
	  cy.get('table tbody tr').should('contain', 'John Doe')
	})

	it('validates form fields', () => {
	  // Try to submit empty form
	  cy.get('button[type="submit"]').click()

	  // Check validation errors
	  cy.get('.text-red-600').should('be.visible')
	  cy.get('.text-red-600').should('contain', 'erforderlich')
	})

	it('loads and displays users', () => {
	  cy.wait('@listUsers')

	  // Check table headers
	  cy.get('table thead th').should('contain', 'Benutzer')
	  cy.get('table thead th').should('contain', 'Email')
	  cy.get('table thead th').should('contain', 'Alter')
	  cy.get('table thead th').should('contain', 'Rolle')

	  // Check if users are displayed
	  cy.get('table tbody tr').should('have.length.greaterThan', 0)
	})

	it('deletes a user', () => {
	  cy.wait('@listUsers')

	  // Click delete button on first user
	  cy.get('table tbody tr').first().find('button').contains('Löschen').click()

	  // Confirm deletion in dialog
	  cy.get('[data-testid="delete-dialog"]').should('be.visible')
	  cy.get('button').contains('Löschen').click()

	  // Wait for delete API call
	  cy.wait('@deleteUser')

	  // Check success notification
	  cy.get('.toast').should('contain', 'erfolgreich gelöscht')
	})

	it('handles connection errors gracefully', () => {
	  // Simulate network error
	  cy.intercept('POST', '**/user.UserService/ListUsers', {
		forceNetworkError: true
	  }).as('listUsersError')

	  cy.get('button').contains('Aktualisieren').click()
	  cy.wait('@listUsersError')

	  // Check error notification
	  cy.get('.toast').should('contain', 'Fehler')
	})
})
