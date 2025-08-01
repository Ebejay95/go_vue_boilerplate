describe('Accessibility', () => {
	beforeEach(() => {
	  cy.visit('/')
	  cy.injectAxe() // Requires cypress-axe
	})

	it('has no accessibility violations', () => {
	  cy.checkA11y()
	})

	it('keyboard navigation works', () => {
	  // Tab through interactive elements
	  cy.get('body').tab()
	  cy.focused().should('have.prop', 'tagName', 'A') // First focusable element

	  cy.tab()
	  cy.focused().should('have.prop', 'tagName', 'BUTTON')

	  // Test form navigation
	  cy.get('input[type="text"]').focus().type('Test')
	  cy.tab()
	  cy.focused().should('have.attr', 'type', 'email')
	})

	it('has proper ARIA labels', () => {
	  cy.get('button').should('have.attr', 'aria-label')
	  cy.get('input').should('have.attr', 'aria-label')

	  // Check for required field indicators
	  cy.get('input[required]').should('have.attr', 'aria-required', 'true')
	})

	it('provides feedback for screen readers', () => {
	  // Test error messages have proper ARIA attributes
	  cy.get('input[type="email"]').type('invalid-email')
	  cy.get('button[type="submit"]').click()

	  cy.get('[role="alert"]').should('exist')
	  cy.get('.text-red-600').should('have.attr', 'aria-live', 'polite')
	})

	it('supports high contrast mode', () => {
	  // Test color contrast
	  cy.get('button.primary').should('have.css', 'background-color')
	  cy.get('.text-gray-900').should('have.css', 'color')

	  // Verify focus indicators are visible
	  cy.get('button').focus()
	  cy.get('button:focus').should('have.css', 'outline')
	})
})
