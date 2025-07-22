<template>
	<!-- Toast Container -->
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
					'flex items-start space-x-3 min-w-80 cursor-default relative',
					getTypeClasses(toast.type)
				]"
			>
				<!-- Close Button -->
				<button
					@click="dismissToastManually(toast.id)"
					class="absolute top-2 right-2 text-gray-400 hover:text-gray-600 transition-colors"
				>
					<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>

				<!-- Icon -->
				<div class="flex-shrink-0">
					<span :class="['text-lg', getIconColor(toast.type)]">
						{{ getIcon(toast.type) }}
					</span>
				</div>

				<!-- Content -->
				<div class="flex-1 min-w-0 pr-6">
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
						:ref="`progress-${toast.id}`"
						class="h-full bg-current opacity-30 transition-all linear"
						:class="getIconColor(toast.type)"
						:style="{
							width: '100%',
							animation: `shrink ${toast.duration || 5000}ms linear forwards`
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
		@click.stop="toggleNotificationPanel"
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

	<!-- Slide-out Notification Panel -->
	<transition name="slide-out">
		<div
			v-if="showNotificationPanel"
			class="fixed inset-0 z-50"
			@click="handleOverlayClick"
		>
			<!-- Panel Content -->
			<div
				class="absolute top-0 right-0 h-full w-96 bg-white shadow-2xl flex flex-col"
				@click.stop
			>
				<!-- Header -->
				<div class="px-6 py-4 border-b bg-gray-50 flex-shrink-0">
					<div class="flex items-center justify-between">
						<h3 class="text-lg font-medium text-gray-900">
							Benachrichtigungen
						</h3>
						<button
							@click="closePanel"
							class="text-gray-400 hover:text-gray-600 transition-colors"
						>
							<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>

					<!-- Action Buttons -->
					<div v-if="persistentNotifications.length > 0" class="flex space-x-3 mt-3">
						<button
							v-if="hasUnreadPersistent"
							@click="markAllAsRead"
							:disabled="isLoading"
							class="text-sm text-blue-600 hover:text-blue-800 font-medium disabled:text-gray-400"
						>
							{{ isLoading ? 'Lädt...' : 'Alle gelesen' }}
						</button>
						<button
							@click="clearAllPersistent"
							:disabled="isLoading"
							class="text-sm text-gray-500 hover:text-gray-700 disabled:text-gray-400"
						>
							{{ isLoading ? 'Lädt...' : 'Alle löschen' }}
						</button>
					</div>
				</div>

				<!-- Loading State -->
				<div v-if="isLoading && persistentNotifications.length === 0"
					 class="flex-1 flex items-center justify-center">
					<div class="text-center">
						<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-2"></div>
						<p class="text-gray-500">Lade Benachrichtigungen...</p>
					</div>
				</div>

				<!-- Error State -->
				<div v-else-if="error" class="flex-1 flex items-center justify-center">
					<div class="text-center text-red-500">
						<p class="mb-2">Fehler beim Laden der Benachrichtigungen</p>
						<button
							@click="loadPersistentNotifications"
							class="text-sm bg-red-100 hover:bg-red-200 px-3 py-1 rounded"
						>
							Erneut versuchen
						</button>
					</div>
				</div>

				<!-- Content Area -->
				<div v-else class="flex-1 overflow-hidden">
					<div class="h-full">
						<!-- Empty State -->
						<div v-if="persistentNotifications.length === 0" class="text-center py-12 px-6">
							<svg class="mx-auto h-16 w-16 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
							</svg>
							<h4 class="text-lg font-medium text-gray-900 mb-2">Keine Benachrichtigungen</h4>
							<p class="text-gray-500">Sie haben keine persistenten Benachrichtigungen</p>
						</div>

						<!-- Notifications List -->
						<div v-else class="p-6 h-full overflow-y-auto">
							<div class="space-y-3">
								<div
									v-for="pNote in persistentNotifications"
									:key="pNote.id"
									:class="[
										'p-3 rounded-lg border transition-colors cursor-pointer',
										pNote.read
											? 'bg-gray-50 border-gray-200 hover:bg-gray-100'
											: 'bg-blue-50 border-blue-200 hover:bg-blue-100'
									]"
									@click="markNotificationAsRead(pNote)"
								>
									<div class="flex justify-between items-start">
										<div class="flex-1 min-w-0">
											<div class="flex items-start space-x-2">
												<!-- Unread indicator -->
												<span
													v-if="!pNote.read"
													class="w-2 h-2 bg-blue-500 rounded-full mt-2 flex-shrink-0"
												></span>

												<div class="flex-1">
													<p :class="[
														'text-sm break-words',
														pNote.read ? 'text-gray-700' : 'text-gray-900 font-medium'
													]">
														{{ pNote.message }}
													</p>
													<div class="flex items-center space-x-2 mt-1">
														<span :class="[
															'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
															pNote.type === 'success' ? 'bg-green-100 text-green-800' :
															pNote.type === 'error' ? 'bg-red-100 text-red-800' :
															pNote.type === 'warning' ? 'bg-yellow-100 text-yellow-800' :
															'bg-blue-100 text-blue-800'
														]">
															{{ pNote.type }}
														</span>
														<span class="text-xs text-gray-500">
															{{ formatTime(pNote.createdAt) }}
														</span>
														<span v-if="!pNote.persistent"
															  class="text-xs text-orange-500"
															  title="Nur lokal gespeichert">
															⚠️
														</span>
													</div>
												</div>
											</div>
										</div>

										<!-- Delete Button -->
										<button
											@click.stop="removePersistentNotification(pNote.id)"
											:disabled="isDeleting === pNote.id"
											class="text-gray-400 hover:text-red-500 transition-colors p-1 rounded-full hover:bg-gray-200 flex-shrink-0 disabled:opacity-50"
										>
											<svg v-if="isDeleting === pNote.id"
												 class="h-4 w-4 animate-spin"
												 fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
											</svg>
											<svg v-else
												 class="h-4 w-4"
												 fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
											</svg>
										</button>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</transition>
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
			showNotificationPanel: false,
			toastTimers: new Map(),
			isDeleting: null,
			processedToasts: new Set()
		}
	},

	computed: {
		...mapState('notifications', ['toasts', 'persistentNotifications']),
		...mapGetters('notifications', [
			'unreadCount',
			'hasUnreadPersistent',
			'isLoading',
			'error'
		]),

		visibleToasts() {
			return this.toasts.slice(0, this.maxToasts)
		}
	},

	watch: {
		toasts: {
			handler(newToasts) {
				newToasts.forEach(toast => {
					if (!this.processedToasts.has(toast.id)) {
						this.processedToasts.add(toast.id)
						this.startToastTimer(toast)
					}
				})

				const currentIds = new Set(newToasts.map(t => t.id))
				for (const processedId of this.processedToasts) {
					if (!currentIds.has(processedId)) {
						this.processedToasts.delete(processedId)
						this.clearToastTimer(processedId)
					}
				}
			},
			immediate: true,
			deep: true
		}
	},

	methods: {
		...mapActions('notifications', [
			'dismissToast',
			'removePersistentNotification',
			'markAsRead',
			'markAllAsRead',
			'clearAllPersistent',
			'loadPersistentNotifications'
		]),

		startToastTimer(toast) {
			if (this.toastTimers.has(toast.id)) {
				return
			}

			const duration = parseInt(toast.duration) || 5000

			if (duration <= 0) {
				console.warn(`⚠️ Invalid duration for toast ${toast.id}: ${duration}`)
				return
			}

			const timer = setTimeout(() => {
				this.dismissToastAutomatically(toast.id)
			}, duration)

			this.toastTimers.set(toast.id, timer)
		},

		clearToastTimer(toastId) {
			if (this.toastTimers.has(toastId)) {
				clearTimeout(this.toastTimers.get(toastId))
				this.toastTimers.delete(toastId)
			}
		},

		clearAllToastTimers() {
			this.toastTimers.forEach(timer => clearTimeout(timer))
			this.toastTimers.clear()
			this.processedToasts.clear()
		},

		dismissToastManually(toastId) {
			this.clearToastTimer(toastId)
			this.processedToasts.delete(toastId)
			this.dismissToast(toastId)
		},

		dismissToastAutomatically(toastId) {
			this.clearToastTimer(toastId)
			this.processedToasts.delete(toastId)

			const toastExists = this.toasts.some(t => t.id === toastId)
			if (toastExists) {
				this.dismissToast(toastId)
			} else {
				console.warn(`⚠️ Toast ${toastId} not found in store`)
			}
		},

		async removePersistentNotification(notificationId) {
			try {
				this.isDeleting = notificationId
				await this.$store.dispatch('notifications/removePersistentNotification', notificationId)
			} catch (error) {
				console.error('Error removing notification:', error)
				this.$store.dispatch('notifications/error', 'Fehler beim Löschen der Benachrichtigung')
			} finally {
				this.isDeleting = null
			}
		},

		async markNotificationAsRead(notification) {
			if (!notification.read) {
				try {
					await this.markAsRead(notification.id)
				} catch (error) {
					console.error('Error marking as read:', error)
				}
			}
		},

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
				success: '✅',
				error: '❌',
				warning: '⚠️',
				info: 'ℹ️'
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

			if (this.showNotificationPanel && this.persistentNotifications.length === 0) {
				this.loadPersistentNotifications()
			}
		},

		closePanel() {
			this.showNotificationPanel = false
		},

		handleOverlayClick() {
			this.closePanel()
		},

		handleToastAction(action, toast) {
			if (typeof action.handler === 'function') {
				action.handler(toast)
			}
			this.dismissToastManually(toast.id)
		}
	},

	mounted() {
		this.toasts.forEach(toast => {
			if (!this.processedToasts.has(toast.id)) {
				this.processedToasts.add(toast.id)
				this.startToastTimer(toast)
			}
		})
	},

	beforeUnmount() {
		this.clearAllToastTimers()
	}
}
</script>

<style scoped>
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

.slide-out-enter-active,
.slide-out-leave-active {
	transition: transform 0.3s ease-in-out;
}

.slide-out-enter-from,
.slide-out-leave-to {
	transform: translateX(100%);
}

@keyframes shrink {
	from {
		width: 100%;
	}
	to {
		width: 0%;
	}
}

.overflow-y-auto::-webkit-scrollbar {
	width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
	background: #f1f5f9;
	border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
	background: #cbd5e1;
	border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
	background: #94a3b8;
}
</style>
