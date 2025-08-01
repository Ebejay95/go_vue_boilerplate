import FormField from './FormField.vue'

describe('FormField Component', () => {
  it('renders text input correctly', () => {
    cy.mount(FormField, {
      props: {
        field: {
          key: 'name',
          label: 'Full Name',
          type: 'string',
          placeholder: 'Enter your name'
        },
        modelValue: ''
      }
    })

    cy.get('label').should('contain', 'Full Name')
    cy.get('input[type="text"]').should('have.attr', 'placeholder', 'Enter your name')
  })

  it('renders email input for email fields', () => {
    cy.mount(FormField, {
      props: {
        field: {
          key: 'email',
          label: 'Email Address'
        },
        modelValue: ''
      }
    })

    cy.get('input[type="email"]').should('exist')
  })

  it('renders select field with options', () => {
    cy.mount(FormField, {
      props: {
        field: {
          key: 'role',
          label: 'Role',
          inputType: 'select',
          options: [
            { value: 'user', label: 'User' },
            { value: 'admin', label: 'Administrator' }
          ]
        },
        modelValue: ''
      }
    })

    cy.get('select').should('exist')
    cy.get('option').should('have.length', 2)
    cy.get('option').first().should('contain', 'User')
    cy.get('option').last().should('contain', 'Administrator')
  })

  it('shows required indicator', () => {
    cy.mount(FormField, {
      props: {
        field: {
          key: 'name',
          label: 'Name',
          required: true
        },
        modelValue: ''
      }
    })

    cy.get('.text-red-500').should('contain', '*')
  })

  it('displays error message', () => {
    cy.mount(FormField, {
      props: {
        field: {
          key: 'email',
          label: 'Email'
        },
        modelValue: 'invalid-email',
        error: 'Please enter a valid email',
        touched: true
      }
    })

    cy.get('.text-red-600').should('contain', 'Please enter a valid email')
  })

  it('emits update events on input', () => {
    const onUpdate = cy.stub()

    cy.mount(FormField, {
      props: {
        field: {
          key: 'name',
          label: 'Name'
        },
        modelValue: '',
        'onUpdate:modelValue': onUpdate
      }
    })

    cy.get('input').type('John Doe')
    cy.then(() => {
      expect(onUpdate).to.have.been.calledWith('John Doe')
    })
  })

  it('handles range input correctly', () => {
    cy.mount(FormField, {
      props: {
        field: {
          key: 'age',
          label: 'Age',
          inputType: 'range',
          min: 18,
          max: 100,
          step: 1
        },
        modelValue: 25
      }
    })

    cy.get('input[type="range"]').should('have.attr', 'min', '18')
    cy.get('input[type="range"]').should('have.attr', 'max', '100')
    cy.get('input[type="range"]').should('have.value', '25')

    // Test range value display
    cy.get('.font-medium').should('contain', '25')
  })
})
