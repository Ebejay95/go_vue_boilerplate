import { createApp } from 'vue'
import App from './App.vue'
import './assets/styles/index.css'
import FriendContact from './components/FriendContact.vue'

const app = createApp(App)
app.component('friend-contact', FriendContact)
app.mount('#app')
