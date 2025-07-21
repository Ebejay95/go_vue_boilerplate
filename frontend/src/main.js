import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/styles/index.css'

import Header from './components/navigation/Header.vue'
import Footer from './components/navigation/Footer.vue'
import Navigation from './components/navigation/Navigation.vue'

import { store } from './store/index'

import NotificationList from './components/reactives/NotificationList.vue'

import Card from './components/base/Card.vue'
import Button from './components/base/Button.vue'
import Section from './components/base/Section.vue'
import Form from './components/base/Form.vue'
import FormField from './components/base/FormField.vue'
import Dialog from './components/base/Dialog.vue'

const app = createApp(App)

app.component('header-component', Header)
app.component('navigation-component', Navigation)
app.component('notification-list', NotificationList)
app.component('footer-component', Footer)
app.component('base-button', Button)
app.component('base-card', Card)
app.component('base-section', Section)
app.component('base-form', Form)
app.component('base-form-field', FormField)
app.component('base-dialog', Dialog)

app.use(router)
app.use(store)

async function initializeAndMount() {
	try {
		await store.dispatch('initializeApp')
	} catch (error) {
		console.error('App initialization failed:', error)
	}
	app.mount('#app')
}

// Start the initialization process
initializeAndMount()
