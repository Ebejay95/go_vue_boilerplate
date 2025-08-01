import { mount } from 'cypress/vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createStore } from 'vuex'

// Global component registration
import Button from '../../src/components/base/Button.vue'
import Card from '../../src/components/base/Card.vue'
import FormField from '../../src/components/base/FormField.vue'

Cypress.Commands.add('mount', (component, options = {}) => {
  // Default router for component tests
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div>Home</div>' } },
      { path: '/test', component: { template: '<div>Test</div>' } }
    ]
  })

  // Default store for component tests
  const store = createStore({
    modules: {
      users: {
        namespaced: true,
        state: { users: [], loading: false },
        actions: { fetchUsers: cy.stub().resolves([]) }
      },
      notifications: {
        namespaced: true,
        state: { toasts: [], persistentNotifications: [] },
        actions: {
          success: cy.stub(),
          error: cy.stub(),
          info: cy.stub()
        }
      }
    }
  })

  const defaultOptions = {
    global: {
      plugins: [router, store],
      components: {
        'base-button': Button,
        'base-card': Card,
        'base-form-field': FormField
      }
    }
  }

  const mergedOptions = {
    ...defaultOptions,
    ...options,
    global: {
      ...defaultOptions.global,
      ...options.global
    }
  }

  return mount(component, mergedOptions)
})
