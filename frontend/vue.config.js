const { defineConfig } = require('@vue/cli-service')
const { DefinePlugin } = require('webpack')
const path = require('path')

module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: process.env.BASE_URL || '/',

  productionSourceMap: process.env.NODE_ENV !== 'production',

  devServer: {
    host: '0.0.0.0',
    port: process.env.FRONTEND_PORT,
    hot: true,
    liveReload: true,

    watchFiles: {
      paths: ['src/**/*', 'public/**/*'],
      options: {
        usePolling: process.env.CHOKIDAR_USEPOLLING === 'true',
        interval: 1000,
        ignored: /node_modules/
      }
    },

    client: {
      webSocketURL: 'auto://0.0.0.0:0/ws',
      overlay: {
        errors: true,
        warnings: false
      },
      progress: true,
      reconnect: 5
    },

    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
      'Access-Control-Allow-Headers': 'X-Requested-With, content-type, Authorization'
    }
  },

  configureWebpack: config => {
    config.plugins.push(
      new DefinePlugin({
        __VUE_OPTIONS_API__: JSON.stringify(true),
        __VUE_PROD_DEVTOOLS__: JSON.stringify(false),
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: JSON.stringify(false)
      })
    )

    if (process.env.NODE_ENV === 'development') {
      config.watchOptions = {
        poll: process.env.WATCHPACK_POLLING === 'true' ? 1000 : false,
        aggregateTimeout: 300,
        ignored: /node_modules/
      }

      config.cache = {
        type: 'filesystem',
        cacheDirectory: path.resolve(__dirname, 'node_modules/.cache/webpack')
      }

      config.devtool = 'eval-cheap-module-source-map'
    }

    if (process.env.NODE_ENV === 'production') {
      config.optimization = {
        ...config.optimization,
        splitChunks: {
          chunks: 'all',
          cacheGroups: {
            vendor: {
              test: /[\\/]node_modules[\\/]/,
              name: 'vendors',
              chunks: 'all'
            }
          }
        }
      }
    }
  },

  css: {
    extract: process.env.NODE_ENV === 'production',
    sourceMap: process.env.NODE_ENV === 'development',
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
  },

  chainWebpack: config => {
    if (process.env.NODE_ENV === 'development') {
      config.plugin('html').tap(args => {
        args[0].template = './public/index.html'
        return args
      })
      config.plugins.delete('watch-ignore')
    }
  }
})
