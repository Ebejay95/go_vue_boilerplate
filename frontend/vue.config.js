// frontend/vue.config.js - Garantiertes Hot Reload für Docker
const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,

  devServer: {
    port: 8080,
    host: '0.0.0.0',
    allowedHosts: 'all',

    // KRITISCH: Diese müssen alle true sein
    hot: 'only',  // Nur Hot Reload, kein vollständiges Neuladen
    liveReload: false, // Deaktivieren um Konflikte zu vermeiden

    // File Watching - AGGRESSIV
    watchFiles: {
      paths: [
        'src/**/*',
        'public/**/*',
        '*.js',
        '*.json',
        '*.vue'
      ],
      options: {
        usePolling: true,     // MUSS true sein für Docker
        interval: 500,        // Schnell genug
        ignored: ['**/node_modules/**', '**/.git/**']
      }
    },

    // WebSocket Konfiguration
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

    // Zusätzliche Headers
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': '*',
      'Access-Control-Allow-Headers': '*'
    },

    // Proxy
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true
      }
    }
  },

  // Webpack Konfiguration für File Watching
  configureWebpack: config => {
    if (process.env.NODE_ENV === 'development') {
      config.watchOptions = {
        poll: 500,
        aggregateTimeout: 200,
        ignored: /node_modules/
      }

      // Ensure HMR is enabled
      config.devtool = 'eval-cheap-module-source-map'
    }
  },

  css: {
    loaderOptions: {
      postcss: {
        postcssOptions: {
          plugins: [
            require('tailwindcss'),
            require('autoprefixer'),
          ],
        },
      },
    },
  }
})
