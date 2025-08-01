const { defineConfig } = require('cypress')

module.exports = defineConfig({
  // E2E Tests
  e2e: {
    baseUrl: 'http://localhost:3000',
    supportFile: 'cypress/support/e2e.js',
    specPattern: 'cypress/e2e/**/*.cy.{js,jsx,ts,tsx}',
    videosFolder: 'cypress/videos',
    screenshotsFolder: 'cypress/screenshots',

    // Test Settings
    viewportWidth: 1280,
    viewportHeight: 720,
    video: true,
    screenshotOnRunFailure: true,

    // Timeouts
    defaultCommandTimeout: 8000,
    requestTimeout: 10000,
    responseTimeout: 10000,
    pageLoadTimeout: 30000,

    // Retry Configuration
    retries: {
      runMode: 2,
      openMode: 0
    },

    setupNodeEvents(on, config) {
      // Code Coverage Plugin
      require('@cypress/code-coverage/task')(on, config)

      // Percy Plugin für Visual Testing
      on('task', {
        log(message) {
          console.log(message)
          return null
        }
      })

      return config
    },

    env: {
      // Test Environment Variables
      GRPC_WEB_URL: 'http://localhost:8081',
      WS_URL: 'ws://localhost:8082',
      API_BASE_URL: 'http://localhost:50051',
      COVERAGE: true
    }
  },

  // Component Tests
  component: {
    devServer: {
      framework: 'vue-cli',
      bundler: 'webpack'
    },
    supportFile: 'cypress/support/component.js',
    specPattern: 'src/**/*.cy.{js,jsx,ts,tsx,vue}',
    indexHtmlFile: 'cypress/support/component-index.html'
  },

  // Global Configuration
  chromeWebSecurity: false,
  experimentalStudio: true,
  experimentalWebKitSupport: true,

  // Weitere Einstellungen
  watchForFileChanges: true,
  numTestsKeptInMemory: 50,

  // Browser Configuration
  browsers: [
    {
      name: 'chrome',
      family: 'chromium',
      channel: 'stable',
      displayName: 'Chrome',
      version: '',
      path: '',
      minSupportedVersion: 64,
      majorVersion: ''
    }
  ]
})
