describe('Responsive Design', () => {
	const viewports = [
	  { name: 'mobile', width: 375, height: 667 },
	  { name: 'tablet', width: 768, height: 1024 },
	  { name: 'desktop', width: 1280, height: 720 }
	]

	viewports.forEach(({ name, width, height }) => {
	  describe(`${name} viewport`, () => {
		beforeEach(() => {
		  cy.viewport(width, height)
		  cy.visit('/')
		})

		it('displays correctly on ' + name, () => {
		  cy.get('header').should('be.visible')
		  cy.get('main').should('be.visible')

		  if (name === 'mobile') {
			// Mobile-specific checks
			cy.get('.space-y-2').should('exist') // Stacked layout
		  } else {
			// Desktop/tablet checks
			cy.get('.grid-cols-2').should('exist') // Side-by-side layout
		  }
		})

		it('navigation works on ' + name, () => {
		  cy.get('nav').should('be.visible')
		  cy.get('nav a').should('have.length.greaterThan', 0)

		  // Test navigation
		  cy.get('nav a').contains('Login').click()
		  cy.url().should('include', '/login')
		})

		it('forms are usable on ' + name, () => {
		  cy.get('form').should('be.visible')
		  cy.get('input, select').each($el => {
			cy.wrap($el).should('be.visible')
		  })

		  // Test form interaction
		  cy.get('input[type="text"]').type('Test Input')
		  cy.get('input[type="text"]').should('have.value', 'Test Input')
		})
	  })
	})
})
