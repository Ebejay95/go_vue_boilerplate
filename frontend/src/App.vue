<template>
  <div id="app">
    <header>
      <h1>My Fullstack App</h1>
    </header>
    <main>
      <UserList />
    </main>
  </div>
</template>

<script>
import UserList from './components/UserList.vue'

export default {
  name: 'App',
  components: {
    UserList
  }
}
</script>

<style>
#app {
  font-family: Arial, sans-serif;
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

header {
  text-align: center;
  margin-bottom: 2rem;
}

h1 {
  color: #2c3e50;
}
</style>