// vue.config.js - Production optimiert
const { defineConfig } = require('@vue/cli-service')
const { DefinePlugin } = require('webpack')

module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: process.env.BASE_URL || '/',

  // Production Build Optimierungen
  productionSourceMap: false,  // Keine source maps in production

  configureWebpack: config => {
    config.plugins.push(
      new DefinePlugin({
        __VUE_OPTIONS_API__: JSON.stringify(true),
        __VUE_PROD_DEVTOOLS__: JSON.stringify(false),
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: JSON.stringify(false)
      })
    )

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
    extract: true,  // CSS in separate files
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
