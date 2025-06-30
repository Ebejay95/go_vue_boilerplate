import { createApp } from 'vue'
import App from './App.vue'
import './assets/styles/index.css'
import './assets/styles/base/typography.css'
import './assets/styles/base/variables.css'
import './assets/styles/base/reset.css'
import './assets/styles/components/layout.css'
import './assets/styles/components/buttons.css'
import './assets/styles/components/cards.css'
import './assets/styles/components/forms.css'
import './assets/styles/components/navigation.css'
import './assets/styles/utilities/animations.css'
import './assets/styles/utilities/spacing.css'
import './assets/styles/utilities/helpers.css'
import './assets/styles/vendor/overrides.css'
import FriendContact from './components/FriendContact.vue'

const app = createApp(App)
app.component('friend-contact', FriendContact)
app.mount('#app')
