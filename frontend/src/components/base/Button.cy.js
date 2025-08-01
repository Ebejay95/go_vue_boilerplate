import Button from './Button.vue'

describe('Button Component', () => {
  it('renders as button by default', () => {
    cy.mount(Button, {
      slots: {
        default: 'Click Me'
      }
    })

    cy.get('button').should('contain', 'Click Me')
    cy.get('button').should('have.class', 'btn')
  })

  it('renders as router-link when link prop is true', () => {
    cy.mount(Button, {
      props: {
        link: true,
        to: '/test'
      },
      slots: {
        default: 'Navigate'
      }
    })

    cy.get('router-link').should('exist')
    cy.get('router-link').should('contain', 'Navigate')
  })

  it('applies correct CSS classes for different modes', () => {
    const modes = ['primary', 'secondary', 'flat']

    modes.forEach(mode => {
      cy.mount(Button, {
        props: { mode },
        slots: { default: `${mode} Button` }
      })

      cy.get('button').should('have.class', mode)
    })
  })

  it('handles click events', () => {
    const onClick = cy.stub()

    cy.mount(Button, {
      props: {
        onClick
      },
      slots: {
        default: 'Click Me'
      }
    })

    cy.get('button').click()
    cy.then(() => {
      expect(onClick).to.have.been.called
    })
  })

  it('is accessible', () => {
    cy.mount(Button, {
      slots: {
        default: 'Accessible Button'
      }
    })

    cy.get('button').should('be.visible')
    cy.get('button').focus()
    cy.get('button').should('have.focus')

    // Test keyboard interaction
    cy.get('button').type('{enter}')
  })
})
