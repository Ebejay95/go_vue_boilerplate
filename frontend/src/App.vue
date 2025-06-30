<template>
	<section>
	  <h2>My friends</h2>
	  <friend-contact
		v-for="friend in friends"
		:key="friend.id"
		:friend="friend">
	  </friend-contact>
	</section>
  </template>

  <script>
  export default {
	data() {
	  return {
		friends: [
		  {
			id: 'manuel',
			name: 'Manuel',
			phone: '123-456-789',
			email: 'manuel@example.com'
		  },
		  {
			id: 'julie',
			name: 'Julie',
			phone: '987-654-321',
			email: 'julie@example.com'
		  }
		]
	  }
	}
  }
  </script>
