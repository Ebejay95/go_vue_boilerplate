<template>
	<div v-if="isVisible" class="dialog-overlay" @click="handleOverlayClick">
	  <div class="dialog" @click.stop>
		<div class="dialog-content">
		  <slot></slot>
		</div>

		<div v-if="mode === 'message'" class="dialog-actions message-actions">
		  <button class="close-button" @click="close" aria-label="Schließen">
			×
		  </button>
		</div>

		<div v-else-if="mode === 'dialog'" class="dialog-actions dialog-buttons">
		  <base-button
			v-if="primaryAction"
			mode="primary"
			@click="handlePrimaryAction"
		  >
			{{ primaryAction.label }}
		  </base-button>
		  <base-button mode="flat" @click="close">
			{{ closeLabel }}
		  </base-button>
		</div>
	  </div>
	</div>
  </template>

  <script>
  export default {
	name: 'Dialog',
	props: {
	  mode: {
		type: String,
		default: 'message', // 'message' oder 'dialog'
		validator: (value) => ['message', 'dialog'].includes(value)
	  },
	  visible: {
		type: Boolean,
		default: false
	  },
	  primaryAction: {
		type: Object,
		default: null,
		// Erwartete Struktur: { label: 'OK', callback: function() {} }
	  },
	  closeLabel: {
		type: String,
		default: 'Abbrechen'
	  },
	  closeOnOverlay: {
		type: Boolean,
		default: true
	  }
	},
	data() {
	  return {
		isVisible: this.visible
	  };
	},
	watch: {
	  visible(newValue) {
		this.isVisible = newValue;
	  }
	},
	methods: {
	  close() {
		this.isVisible = false;
		this.$emit('close');
	  },
	  handlePrimaryAction() {
		if (this.primaryAction && this.primaryAction.callback) {
		  this.primaryAction.callback();
		}
		this.$emit('primary-action');
		// Dialog automatisch schließen nach Primary Action
		this.close();
	  },
	  handleOverlayClick() {
		if (this.closeOnOverlay) {
		  this.close();
		}
	  }
	}
  };
  </script>

  <style scoped>
  .dialog-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 1000;
  }

  .dialog {
	background: white;
	border-radius: 8px;
	box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
	max-width: 500px;
	width: 90%;
	max-height: 80vh;
	overflow-y: auto;
	position: relative;
  }

  .dialog-content {
	padding: 20px;
  }

  .dialog-actions {
	padding: 16px 20px;
	border-top: 1px solid #e0e0e0;
  }

  .message-actions {
	display: flex;
	justify-content: flex-end;
  }

  .close-button {
	background: none;
	border: none;
	font-size: 24px;
	font-weight: bold;
	color: #666;
	cursor: pointer;
	padding: 4px 8px;
	border-radius: 4px;
	transition: background-color 0.2s;
  }

  .close-button:hover {
	background-color: #f0f0f0;
	color: #333;
  }

  .dialog-buttons {
	display: flex;
	gap: 12px;
	justify-content: flex-end;
  }

  .btn {
	padding: 8px 16px;
	border: none;
	border-radius: 4px;
	cursor: pointer;
	font-size: 14px;
	transition: background-color 0.2s;
  }

  .btn-primary {
	background-color: #007bff;
	color: white;
  }

  .btn-primary:hover {
	background-color: #0056b3;
  }

  .btn-secondary {
	background-color: #6c757d;
	color: white;
  }

  .btn-secondary:hover {
	background-color: #545b62;
  }
  </style>
