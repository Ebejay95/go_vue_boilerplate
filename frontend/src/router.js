import { createRouter, createWebHistory } from 'vue-router'
import Login from './pages/authentication/Login.vue'
import Signup from './pages/authentication/Signup.vue'
import Welcome from './pages/general/Welcome.vue'
import Dashboard from './pages/general/Dashboard.vue'
import NotFound from './pages/NotFound.vue'

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/', component: Welcome },
		{ path: '/dashboard', component: Dashboard },
		{ path: '/signup', component: Signup },
		{ path: '/login', component: Login },
		//{ path: '/', redirect: null },
		// { path: '/coaches', component: null },
		// {
		// 	path: '/coaches/:id', component: null, children: [
		// 		{ path: 'contact', component: null },
		// 	]
		// },
		{ path: '/:notFound(.*)', component: NotFound },
	]
})

export default router;
