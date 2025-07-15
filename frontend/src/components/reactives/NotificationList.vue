<template>
	<!-- Notification List Component -->
	<div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm">
	  <transition-group
		name="notification"
		tag="div"
		class="space-y-2"
	  >
		<div
		  v-for="notification in visibleNotifications"
		  :key="notification.id"
		  :class="[
			'transform transition-all duration-300 ease-in-out',
			'bg-white rounded-lg shadow-lg border-l-4 p-4',
			'flex items-start space-x-3 min-w-80',
			getTypeClasses(notification.type)
		  ]"
		  @click="markAsRead(notification.id)"
		>
		  <!-- Icon -->
		  <div class="flex-shrink-0">
			<component
			  :is="getIcon(notification.type)"
			  :class="[
				'h-5 w-5',
				getIconColor(notification.type)
			  ]"
			/>
		  </div>

		  <!-- Content -->
		  <div class="flex-1 min-w-0">
			<p :class="[
			  'text-sm font-medium',
			  notification.read ? 'text-gray-600' : 'text-gray-900'
			]">
			  {{ notification.message }}
			</p>

			<!-- Timestamp -->
			<p class="text-xs text-gray-500 mt-1">
			  {{ formatTime(notification.createdAt) }}
			</p>

			<!-- Actions -->
			<div v-if="notification.actions && notification.actions.length"
				 class="mt-2 space-x-2">
			  <button
				v-for="action in notification.actions"
				:key="action.label"
				@click.stop="handleAction(action, notification)"
				class="text-xs px-2 py-1 bg-gray-100 hover:bg-gray-200 rounded transition-colors"
			  >
				{{ action.label }}
			  </button>
			</div>
		  </div>

		  <!-- Close Button -->
		  <button
			@click.stop="removeNotification(notification.id)"
			class="flex-shrink-0 text-gray-400 hover:text-gray-600 transition-colors"
		  >
			<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
			</svg>
		  </button>
		</div>
	  </transition-group>
	</div>

	<!-- Notification Badge (für Header) -->
	<div
	  v-if="showBadge && hasUnreadNotifications"
	  class="relative inline-block cursor-pointer"
	  @click="toggleNotificationPanel"
	>
	  <!-- Bell Icon -->
	  <svg class="h-6 w-6 text-gray-600 hover:text-gray-800 transition-colors"
		   fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
			  d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
	  </svg>

	  <!-- Badge -->
	  <span
		v-if="unreadCount > 0"
		class="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center font-medium"
	  >
		{{ unreadCount > 9 ? '9+' : unreadCount }}
	  </span>
	</div>

	<!-- Notification Panel (optional dropdown) -->
	<div
	  v-if="showPanel && showNotificationPanel"
	  class="fixed top-16 right-4 z-50 bg-white rounded-lg shadow-xl border w-80 max-h-96 overflow-hidden"
	>
	  <!-- Header -->
	  <div class="px-4 py-3 border-b bg-gray-50">
		<div class="flex items-center justify-between">
		  <h3 class="text-sm font-medium text-gray-900">Benachrichtigungen</h3>
		  <div class="flex space-x-2">
			<button
			  v-if="hasUnreadNotifications"
			  @click="markAllAsRead"
			  class="text-xs text-blue-600 hover:text-blue-800"
			>
			  Alle als gelesen markieren
			</button>
			<button
			  @click="clearAllNotifications"
			  class="text-xs text-gray-500 hover:text-gray-700"
			>
			  Alle löschen
			</button>
		  </div>
		</div>
	  </div>

	  <!-- Notifications List -->
	  <div class="max-h-64 overflow-y-auto">
		<div v-if="!hasNotifications" class="p-4 text-center text-gray-500 text-sm">
		  Keine Benachrichtigungen
		</div>

		<div
		  v-for="notification in allNotifications"
		  :key="notification.id"
		  :class="[
			'p-3 border-b last:border-b-0 cursor-pointer hover:bg-gray-50 transition-colors',
			!notification.read ? 'bg-blue-50' : ''
		  ]"
		  @click="markAsRead(notification.id)"
		>
		  <div class="flex items-start space-x-2">
			<component
			  :is="getIcon(notification.type)"
			  :class="['h-4 w-4 mt-0.5 flex-shrink-0', getIconColor(notification.type)]"
			/>
			<div class="flex-1 min-w-0">
			  <p :class="[
				'text-sm',
				notification.read ? 'text-gray-600' : 'text-gray-900 font-medium'
			  ]">
				{{ notification.message }}
			  </p>
			  <p class="text-xs text-gray-500 mt-1">
				{{ formatTime(notification.createdAt) }}
			  </p>
			</div>
			<div v-if="!notification.read" class="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0 mt-2"></div>
		  </div>
		</div>
	  </div>
	</div>
  </template>

  <script>
  import { mapState, mapGetters, mapActions } from 'vuex'

  export default {
	name: 'NotificationList',

	props: {
	  showBadge: {
		type: Boolean,
		default: true
	  },
	  showPanel: {
		type: Boolean,
		default: true
	  },
	  maxVisible: {
		type: Number,
		default: 5
	  }
	},

	data() {
	  return {
		showNotificationPanel: false
	  }
	},

	computed: {
	  ...mapState('notifications', ['notifications']),
	  ...mapGetters('notifications', [
		'allNotifications',
		'unreadNotifications',
		'unreadCount',
		'hasNotifications',
		'hasUnreadNotifications',
		'latestNotifications'
	  ]),

	  visibleNotifications() {
		return this.latestNotifications.slice(0, this.maxVisible)
	  }
	},

	methods: {
	  ...mapActions('notifications', [
		'removeNotification',
		'markAsRead',
		'markAllAsRead',
		'clearAllNotifications'
	  ]),

	  getTypeClasses(type) {
		const classes = {
		  success: 'border-green-400 bg-green-50',
		  error: 'border-red-400 bg-red-50',
		  warning: 'border-yellow-400 bg-yellow-50',
		  info: 'border-blue-400 bg-blue-50'
		}
		return classes[type] || classes.info
	  },

	  getIconColor(type) {
		const colors = {
		  success: 'text-green-500',
		  error: 'text-red-500',
		  warning: 'text-yellow-500',
		  info: 'text-blue-500'
		}
		return colors[type] || colors.info
	  },

	  getIcon(type) {
		// Return inline SVG components
		const icons = {
		  success: {
			template: `
			  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
					  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
			  </svg>
			`
		  },
		  error: {
			template: `
			  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
					  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
			  </svg>
			`
		  },
		  warning: {
			template: `
			  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
					  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.268 16.5c-.77.833.192 2.5 1.732 2.5z" />
			  </svg>
			`
		  },
		  info: {
			template: `
			  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
					  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
			  </svg>
			`
		  }
		}
		return icons[type] || icons.info
	  },

	  formatTime(timestamp) {
		const now = new Date()
		const time = new Date(timestamp)
		const diffMs = now - time
		const diffMins = Math.floor(diffMs / 60000)

		if (diffMins < 1) return 'Jetzt'
		if (diffMins < 60) return `vor ${diffMins}m`

		const diffHours = Math.floor(diffMins / 60)
		if (diffHours < 24) return `vor ${diffHours}h`

		const diffDays = Math.floor(diffHours / 24)
		return `vor ${diffDays}d`
	  },

	  toggleNotificationPanel() {
		this.showNotificationPanel = !this.showNotificationPanel
	  },

	  handleAction(action, notification) {
		if (typeof action.handler === 'function') {
		  action.handler(notification)
		}
	  },

	  // Click outside to close panel
	  handleClickOutside(event) {
		if (this.showNotificationPanel && !this.$el.contains(event.target)) {
		  this.showNotificationPanel = false
		}
	  }
	},

	mounted() {
	  document.addEventListener('click', this.handleClickOutside)
	},

	beforeUnmount() {
	  document.removeEventListener('click', this.handleClickOutside)
	}
  }
  </script>

  <style scoped>
  /* Transition animations */
  .notification-enter-active,
  .notification-leave-active {
	transition: all 0.3s ease;
  }

  .notification-enter-from {
	opacity: 0;
	transform: translateX(100%);
  }

  .notification-leave-to {
	opacity: 0;
	transform: translateX(100%);
  }

  .notification-move {
	transition: transform 0.3s ease;
  }
  </style>
