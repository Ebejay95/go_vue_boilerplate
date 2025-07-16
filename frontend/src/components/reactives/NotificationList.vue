<template>
	<!-- Toast Notifications (temporär, oben rechts, OHNE X-Button) -->
	<div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm pointer-events-none">
	  <transition-group
		name="toast"
		tag="div"
		class="space-y-2"
	  >
		<div
		  v-for="toast in visibleToasts"
		  :key="toast.id"
		  :class="[
			'transform transition-all duration-300 ease-in-out pointer-events-auto',
			'bg-white rounded-lg shadow-lg border-l-4 p-4',
			'flex items-start space-x-3 min-w-80 cursor-default',
			getTypeClasses(toast.type)
		  ]"
		>
		  <!-- Icon -->
		  <div class="flex-shrink-0">
			<component
			  :is="getIcon(toast.type)"
			  :class="['h-5 w-5', getIconColor(toast.type)]"
			/>
		  </div>

		  <!-- Content -->
		  <div class="flex-1 min-w-0">
			<p class="text-sm font-medium text-gray-900">
			  {{ toast.message }}
			</p>

			<!-- Actions für Toast (falls vorhanden) -->
			<div v-if="toast.actions && toast.actions.length"
				 class="mt-2 space-x-2">
			  <button
				v-for="action in toast.actions"
				:key="action.label"
				@click="handleToastAction(action, toast)"
				class="text-xs px-2 py-1 bg-gray-100 hover:bg-gray-200 rounded transition-colors"
			  >
				{{ action.label }}
			  </button>
			</div>
		  </div>

		  <!-- Progress Bar für verbleibende Zeit -->
		  <div class="absolute bottom-0 left-0 right-0 h-1 bg-gray-200 rounded-b-lg overflow-hidden">
			<div
			  class="h-full bg-current opacity-30 transition-all linear"
			  :class="getIconColor(toast.type)"
			  :style="{
				width: '100%',
				animation: `shrink ${toast.duration}ms linear forwards`
			  }"
			></div>
		  </div>
		</div>
	  </transition-group>
	</div>

	<!-- Notification Bell Badge -->
	<div
	  v-if="showBadge"
	  class="relative inline-block cursor-pointer"
	  @click="toggleNotificationPanel"
	>
	  <!-- Bell Icon -->
	  <svg class="h-6 w-6 text-gray-600 hover:text-gray-800 transition-colors"
		   fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
			  d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
	  </svg>

	  <!-- Badge für ungelesene persistente Nachrichten -->
	  <span
		v-if="unreadCount > 0"
		class="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center font-medium"
	  >
		{{ unreadCount > 9 ? '9+' : unreadCount }}
	  </span>
	</div>

	<!-- Notification Panel (persistente Liste) -->
	<div
	  v-if="showPanel && showNotificationPanel"
	  class="fixed top-16 right-4 z-50 bg-white rounded-lg shadow-xl border w-80 max-h-96 overflow-hidden"
	  @click.stop
	>
	  <!-- Header -->
	  <div class="px-4 py-3 border-b bg-gray-50">
		<div class="flex items-center justify-between">
		  <h3 class="text-sm font-medium text-gray-900">
			Benachrichtigungen ({{ persistentNotifications.length }})
		  </h3>
		  <div class="flex space-x-2">
			<button
			  v-if="hasUnreadPersistent"
			  @click="markAllAsRead"
			  class="text-xs text-blue-600 hover:text-blue-800"
			>
			  Alle gelesen
			</button>
			<button
			  @click="clearAllPersistent"
			  class="text-xs text-gray-500 hover:text-gray-700"
			>
			  Alle löschen
			</button>
			<button
			  @click="closePanel"
			  class="text-xs text-gray-400 hover:text-gray-600"
			>
			  ✕
			</button>
		  </div>
		</div>
	  </div>

	  <!-- Persistent Notifications List -->
	  <div class="max-h-64 overflow-y-auto">
		<div v-if="persistentNotifications.length === 0" class="p-4 text-center text-gray-500 text-sm">
		  <svg class="mx-auto h-8 w-8 text-gray-300 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
		  </svg>
		  Keine gespeicherten Benachrichtigungen
		</div>

		<div
		  v-for="notification in persistentNotifications"
		  :key="notification.id"
		  :class="[
			'p-3 border-b last:border-b-0 cursor-pointer hover:bg-gray-50 transition-colors',
			!notification.read ? 'bg-blue-50 border-l-2 border-l-blue-500' : ''
		  ]"
		  @click="markAsRead(notification.id)"
		>
		  <div class="flex items-start space-x-3">
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

			  <!-- Actions für persistente Nachrichten -->
			  <div v-if="notification.actions && notification.actions.length"
				   class="mt-2 space-x-2">
				<button
				  v-for="action in notification.actions"
				  :key="action.label"
				  @click.stop="handlePersistentAction(action, notification)"
				  class="text-xs px-2 py-1 bg-blue-100 hover:bg-blue-200 text-blue-800 rounded transition-colors"
				>
				  {{ action.label }}
				</button>
			  </div>
			</div>

			<div class="flex items-center space-x-2">
			  <!-- Unread indicator -->
			  <div v-if="!notification.read" class="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0"></div>

			  <!-- Delete button -->
			  <button
				@click.stop="removePersistentNotification(notification.id)"
				class="text-gray-400 hover:text-red-500 transition-colors"
			  >
				<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
				</svg>
			  </button>
			</div>
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
	  maxToasts: {
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
	  ...mapState('notifications', ['toasts', 'persistentNotifications']),
	  ...mapGetters('notifications', [
		'unreadCount',
		'hasUnreadPersistent'
	  ]),

	  // Nur die sichtbaren Toasts (begrenzt)
	  visibleToasts() {
		return this.toasts.slice(0, this.maxToasts)
	  }
	},

	methods: {
	  ...mapActions('notifications', [
		'dismissToast',
		'removePersistentNotification',
		'markAsRead',
		'markAllAsRead',
		'clearAllPersistent'
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

	  closePanel() {
		this.showNotificationPanel = false
	  },

	  handleToastAction(action, toast) {
		if (typeof action.handler === 'function') {
		  action.handler(toast)
		}
		// Toast nach Action automatisch schließen
		this.dismissToast(toast.id)
	  },

	  handlePersistentAction(action, notification) {
		if (typeof action.handler === 'function') {
		  action.handler(notification)
		}
		// Notification als gelesen markieren nach Action
		if (!notification.read) {
		  this.markAsRead(notification.id)
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
  /* Toast animations */
  .toast-enter-active,
  .toast-leave-active {
	transition: all 0.4s ease;
  }

  .toast-enter-from {
	opacity: 0;
	transform: translateX(100%) scale(0.9);
  }

  .toast-leave-to {
	opacity: 0;
	transform: translateX(100%) scale(0.9);
  }

  .toast-move {
	transition: transform 0.4s ease;
  }

  /* Progress bar animation */
  @keyframes shrink {
	from {
	  width: 100%;
	}
	to {
	  width: 0%;
	}
  }

  /* Smooth scrollbar */
  .max-h-64::-webkit-scrollbar {
	width: 4px;
  }

  .max-h-64::-webkit-scrollbar-track {
	background: #f1f5f9;
  }

  .max-h-64::-webkit-scrollbar-thumb {
	background: #cbd5e1;
	border-radius: 2px;
  }

  .max-h-64::-webkit-scrollbar-thumb:hover {
	background: #94a3b8;
  }
  </style>