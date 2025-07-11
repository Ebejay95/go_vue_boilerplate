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
		port: process.env.PORT,
		host: '0.0.0.0',
		allowedHosts: 'all',

		client: {
			overlay: {
				errors: true,
				warnings: false,
				runtimeErrors: false
			}
		},

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
				interval: 1000,
				ignored: ['**/node_modules/**', '**/.git/**']
			}
		},

		headers: {
			'Access-Control-Allow-Origin': '*',
			'Access-Control-Allow-Methods': '*',
			'Access-Control-Allow-Headers': '*'
		},


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
