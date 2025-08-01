describe('Performance', () => {
	it('loads page within acceptable time', () => {
	  const start = Date.now()

	  cy.visit('/')

	  cy.window().then(() => {
		const loadTime = Date.now() - start
		expect(loadTime).to.be.lessThan(3000) // 3 seconds max
	  })
	})

	it('lazy loads components efficiently', () => {
	  cy.visit('/')

	  // Check that heavy components aren't loaded initially
	  cy.get('[data-testid="heavy-component"]').should('not.exist')

	  // Trigger lazy loading
	  cy.get('button').contains('Load More').click()

	  // Verify component loads
	  cy.get('[data-testid="heavy-component"]', { timeout: 5000 }).should('exist')
	})

	it('handles large datasets efficiently', () => {
	  // Mock large dataset
	  cy.intercept('POST', '**/user.UserService/ListUsers', {
		fixture: 'large-users-list.json'
	  })

	  cy.visit('/')

	  // Verify virtual scrolling or pagination works
	  cy.get('table tbody tr').should('have.length.lessThan', 50) // Limited rendering
	  cy.get('[data-testid="pagination"]').should('exist')
	})

	it('optimizes network requests', () => {
	  cy.visit('/')

	  // Count network requests
	  let requestCount = 0
	  cy.intercept('**', () => {
		requestCount++
	  })

	  cy.get('button').contains('Aktualisieren').click()

	  cy.then(() => {
		expect(requestCount).to.be.lessThan(10) // Reasonable limit
	  })
	})
})
