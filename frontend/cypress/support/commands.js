Cypress.Commands.add('login', (email = 'test@example.com', password = 'password') => {
	cy.session([email, password], () => {
	  cy.visit('/login')
	  cy.get('input[type="email"]').type(email)
	  cy.get('input[type="password"]').type(password)
	  cy.get('button[type="submit"]').click()
	  cy.url().should('not.include', '/login')
	})
  })

  // Notification commands
  Cypress.Commands.add('waitForNotification', (type = 'success', message = null) => {
	cy.get('.toast', { timeout: 10000 }).should('be.visible')
	if (type) {
	  cy.get('.toast').should('have.class', `border-${type === 'success' ? 'green' : type === 'error' ? 'red' : 'blue'}-400`)
	}
	if (message) {
	  cy.get('.toast').should('contain', message)
	}
  })

  Cypress.Commands.add('dismissNotification', () => {
	cy.get('.toast button[aria-label="Schließen"]').click()
	cy.get('.toast').should('not.exist')
  })

  // User management commands
  Cypress.Commands.add('createUser', (userData = {}) => {
	const defaultUser = {
	  name: 'Test User',
	  email: 'test@example.com',
	  age: 25,
	  role: 'user'
	}

	const user = { ...defaultUser, ...userData }

	cy.get('input[type="text"]').clear().type(user.name)
	cy.get('input[type="email"]').clear().type(user.email)
	cy.get('input[type="range"]').invoke('val', user.age).trigger('input')
	cy.get('select').select(user.role)
	cy.get('button[type="submit"]').click()

	return cy.wrap(user)
  })

  Cypress.Commands.add('deleteUser', (userName) => {
	cy.get('table tbody tr').contains(userName).parent('tr').within(() => {
	  cy.get('button').contains('Löschen').click()
	})
	cy.get('[data-testid="delete-dialog"]').within(() => {
	  cy.get('button').contains('Löschen').click()
	})
  })

  // Table commands
  Cypress.Commands.add('getTableRow', (rowIndex) => {
	return cy.get('table tbody tr').eq(rowIndex)
  })

  Cypress.Commands.add('sortTableBy', (columnName) => {
	cy.get('table thead th').contains(columnName).click()
  })

  // Form commands
  Cypress.Commands.add('fillForm', (formData) => {
	Object.entries(formData).forEach(([field, value]) => {
	  if (typeof value === 'string') {
		cy.get(`[name="${field}"], #${field}`).clear().type(value)
	  } else if (typeof value === 'number') {
		cy.get(`[name="${field}"], #${field}`).clear().type(value.toString())
	  } else if (typeof value === 'boolean') {
		if (value) {
		  cy.get(`[name="${field}"], #${field}`).check()
		} else {
		  cy.get(`[name="${field}"], #${field}`).uncheck()
		}
	  }
	})
  })

  // Accessibility commands
  Cypress.Commands.add('tab', { prevSubject: 'optional' }, (subject) => {
	if (subject) {
	  cy.wrap(subject).trigger('keydown', { key: 'Tab' })
	} else {
	  cy.get('body').trigger('keydown', { key: 'Tab' })
	}
  })

  Cypress.Commands.add('checkFocusOrder', (selectors) => {
	selectors.forEach((selector, index) => {
	  if (index === 0) {
		cy.get(selector).focus()
	  } else {
		cy.tab()
	  }
	  cy.focused().should('match', selector)
	})
  })

  // API mocking commands
  Cypress.Commands.add('mockUserAPI', (scenario = 'success') => {
	const scenarios = {
	  success: {
		listUsers: { fixture: 'users-list.json' },
		createUser: { fixture: 'user-created.json' },
		deleteUser: { body: { success: true, message: 'User deleted' } }
	  },
	  error: {
		listUsers: { statusCode: 500, body: { error: 'Server error' } },
		createUser: { statusCode: 400, body: { error: 'Validation failed' } },
		deleteUser: { statusCode: 404, body: { error: 'User not found' } }
	  },
	  slow: {
		listUsers: { fixture: 'users-list.json', delay: 3000 },
		createUser: { fixture: 'user-created.json', delay: 2000 },
		deleteUser: { body: { success: true }, delay: 1000 }
	  }
	}

	const config = scenarios[scenario]

	cy.intercept('POST', '**/user.UserService/ListUsers', config.listUsers).as('listUsers')
	cy.intercept('POST', '**/user.UserService/CreateUser', config.createUser).as('createUser')
	cy.intercept('POST', '**/user.UserService/DeleteUser', config.deleteUser).as('deleteUser')
  })

  // Visual regression commands
  Cypress.Commands.add('percySnapshot', (name, options = {}) => {
	if (Cypress.env('PERCY_TOKEN')) {
	  cy.percySnapshot(name, options)
	} else {
	  cy.log('Percy snapshot skipped - no PERCY_TOKEN found')
	}
  })
