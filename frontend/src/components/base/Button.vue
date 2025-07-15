<template>
	<button v-if="!link" class="btn" :class="mode">
		<slot></slot>
	</button>
	<router-link v-else :to="to" class="btn" :class="mode">
		<slot></slot>
	</router-link>
</template>

<script>
	export default {
	props: {
			mode: {
				type: String,
				required: false,
				default: ''
			},
			link: {
				type: Boolean,
				required: false,
				default: false
			},
			to: {
				type: String,
				required: false,
				default: '/'
			},
		}
	}
</script>
