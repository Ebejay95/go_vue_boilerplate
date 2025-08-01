module.exports = {
	root: true,
	env: {
	  node: true,
	  browser: true,
	  es2022: true
	},
	extends: [
	  'plugin:vue/vue3-essential',
	  'eslint:recommended',
	  '@vue/eslint-config-prettier'
	],
	plugins: ['vue'],
	parserOptions: {
	  ecmaVersion: 2022,
	  sourceType: 'module'
	},
	rules: {
	  'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
	  'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
	  'vue/multi-word-component-names': 'off',
	  'vue/no-unused-vars': 'error',
	  'vue/no-unused-components': 'warn'
	},
	overrides: [
	  {
		files: ['**/*.test.js', '**/*.spec.js', '**/*.cy.js'],
		env: {
		  'cypress/globals': true,
		  jest: true
		},
		extends: ['plugin:cypress/recommended'],
		rules: {
		  'no-unused-expressions': 'off'
		}
	  }
	]
}
