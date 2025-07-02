import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/styles/index.css'
import Header from './components/navigation/Header.vue'
import Footer from './components/navigation/Footer.vue'
import Navigation from './components/navigation/Navigation.vue'

const app = createApp(App)

app.component('header-component', Header)
app.component('navigation-component', Navigation)
app.component('footer-component', Footer)
app.use(router)
app.mount('#app')
