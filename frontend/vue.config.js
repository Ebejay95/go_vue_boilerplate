const { defineConfig } = require('@vue/cli-service')
const { DefinePlugin } = require('webpack')
const packageJson = require('./package.json')

const baseUrl = process.env.VUE_APP_BASE_URL || process.env.BASE_URL || '/'

module.exports = defineConfig({
	transpileDependencies: true,

	publicPath: baseUrl,

	chainWebpack: config => {
		config
			.plugin('html')
			.tap(args => {
				args[0].title = packageJson.name || 'MyApp'
				args[0].templateParameters = {
					...args[0].templateParameters,
					APP_NAME: packageJson.name || 'MyApp',
					BASE_URL: baseUrl
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

		if (process.env.NODE_ENV === 'development') {
			config.watchOptions = {
				poll: 500,
				aggregateTimeout: 200,
				ignored: /node_modules/
			}
			config.devtool = 'eval-cheap-module-source-map'
		}
	},

	devServer: {
		port: 8080,
		host: '0.0.0.0',
		allowedHosts: 'all',
		hot: 'only',
		liveReload: false,
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
				interval: 500,
				ignored: ['**/node_modules/**', '**/.git/**']
			}
		},

		client: {
			webSocketURL: {
				hostname: '0.0.0.0',
				pathname: '/ws',
				port: 8080,
				protocol: 'ws'
			},
			overlay: {
				errors: true,
				warnings: false,
				runtimeErrors: false
			},
			progress: false,
			reconnect: 3
		},

		headers: {
			'Access-Control-Allow-Origin': '*',
			'Access-Control-Allow-Methods': '*',
			'Access-Control-Allow-Headers': '*'
		},

		proxy: {
			'/api': {
				target: 'http://localhost:3000',
				changeOrigin: true
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