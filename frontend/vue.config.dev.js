const { defineConfig } = require('@vue/cli-service')
const { DefinePlugin } = require('webpack')
const packageJson = require('./package.json')

module.exports = defineConfig({
	transpileDependencies: true,
	publicPath: process.env.BASE_URL || '/',

	chainWebpack: config => {
		config
			.plugin('html')
			.tap(args => {
				args[0].title = packageJson.name || 'MyApp'
				args[0].templateParameters = {
					...args[0].templateParameters,
					APP_NAME: packageJson.name || 'MyApp',
					BASE_URL: process.env.BASE_URL || '/'
				}
				return args
			})
	},

	configureWebpack: config => {
		config.plugins.push(
			new DefinePlugin({
				__VUE_OPTIONS_API__: JSON.stringify(true),
				__VUE_PROD_DEVTOOLS__: JSON.stringify(false),
				__VUE_PROD_HYDRATION_MISMATCH_DETAILS__: JSON.stringify(false),
				__FEATURE_SUSPENSE__: JSON.stringify(true),
				__FEATURE_OPTIONS_API__: JSON.stringify(true),
				__FEATURE_PROD_DEVTOOLS__: JSON.stringify(false)
			})
		)

		// NUCLEAR HOT RELOAD OPTIONS
		config.watchOptions = {
			poll: 500,           // Poll every 500ms
			aggregateTimeout: 100,
			ignored: /node_modules/
		}

		// Disable all caching for development
		config.cache = false
		config.resolve.symlinks = false

		// Enable source maps for debugging
		config.devtool = 'eval-cheap-module-source-map'
	},

	devServer: {
		port: process.env.FRONTEND_PORT || 3000,
		host: '0.0.0.0',
		allowedHosts: 'all',

		// ENHANCED HOT RELOAD
		hot: true,
		liveReload: true,

		// WebSocket configuration for better hot reload
		client: {
			webSocketURL: {
				hostname: 'localhost',
				port: process.env.FRONTEND_PORT || 3000,
				protocol: 'ws'
			},
			overlay: {
				errors: true,
				warnings: false
			},
			progress: true,
			reconnect: true
		},

		// AGGRESSIVE FILE WATCHING
		watchFiles: {
			paths: [
				'src/**/*',
				'public/**/*',
				'*.js',
				'*.json',
				'*.vue'
			],
			options: {
				usePolling: true,
				interval: 500,    // Check every 500ms
				ignored: ['**/node_modules/**', '**/.git/**']
			}
		},

		// CORS headers
		headers: {
			'Access-Control-Allow-Origin': '*',
			'Access-Control-Allow-Methods': '*',
			'Access-Control-Allow-Headers': '*'
		},

		// Static file serving
		static: {
			directory: './public',
			watch: {
				usePolling: true,
				interval: 500
			}
		}
	},

	css: {
		loaderOptions: {
			postcss: {
				postcssOptions: {
					plugins: [
						require('postcss-import'),
						require('tailwindcss'),
						require('autoprefixer'),
					],
				},
			},
		},
	}
})
