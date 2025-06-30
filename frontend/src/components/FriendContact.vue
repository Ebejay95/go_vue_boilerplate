<template>
	<li>
	  <h2>{{ friend.name }}</h2>
	  <button @click="toggleDetails" class="btn">Show Details</button>
	  <ul v-if="detailsAreVisible">
		<li><strong>Phone:</strong> {{ friend.phone }}</li>
		<li><strong>Email:</strong> {{ friend.email }}</li>
	  </ul>
	</li>
  </template>

  <script>
  export default {
	props: ['friend'], // Props definieren
	data() {
	  return {
		detailsAreVisible: false
	  }
	},
	methods: {
	  toggleDetails() {
		this.detailsAreVisible = !this.detailsAreVisible
	  }
	}
  }
  </script>
