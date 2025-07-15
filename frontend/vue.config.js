// vue.config.js - Enhanced Development Configuration
const { defineConfig } = require('@vue/cli-service')
const { DefinePlugin } = require('webpack')
const path = require('path')

module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: process.env.BASE_URL || '/',

  // Production Build Optimierungen
  productionSourceMap: process.env.NODE_ENV !== 'production',

  // Development Server Configuration
  devServer: {
    host: '0.0.0.0',
    port: process.env.FRONTEND_PORT,
    hot: true,
    liveReload: true,

    // Enhanced file watching
    watchFiles: {
      paths: ['src/**/*', 'public/**/*'],
      options: {
        usePolling: process.env.CHOKIDAR_USEPOLLING === 'true',
        interval: 1000,
        ignored: /node_modules/
      }
    },

    // Client configuration for better HMR
    client: {
      webSocketURL: 'auto://0.0.0.0:0/ws',
      overlay: {
        errors: true,
        warnings: false
      },
      progress: true,
      reconnect: 5
    },

    // Headers for development
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

    // Development optimizations
    if (process.env.NODE_ENV === 'development') {
      // Enhanced file watching
      config.watchOptions = {
        poll: process.env.WATCHPACK_POLLING === 'true' ? 1000 : false,
        aggregateTimeout: 300,
        ignored: /node_modules/
      }

      // FIXED: Use absolute path for webpack cache directory
      config.cache = {
        type: 'filesystem',
        cacheDirectory: path.resolve(__dirname, 'node_modules/.cache/webpack')
      }

      // Source maps for debugging
      config.devtool = 'eval-cheap-module-source-map'
    }

    // Production optimizations
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

  // Chain webpack configuration for additional HMR optimizations
  chainWebpack: config => {
    // Ensure proper HMR for Vue files
    if (process.env.NODE_ENV === 'development') {
      config.plugin('html').tap(args => {
        args[0].template = './public/index.html'
        return args
      })

      // Remove problematic watch ignore plugin that might interfere
      config.plugins.delete('watch-ignore')
    }
  }
})
