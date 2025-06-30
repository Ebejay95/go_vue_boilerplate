import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/', component: null },
		//{ path: '/', redirect: null },
		{ path: '/coaches', component: null },
		{
			path: '/coaches/:id', component: null, children: [
				{ path: 'contact', component: null },
		] },
		{ path: '/register', component: null },
		{ path: '/:notFound(.*)', component: null },
	]
})

export default router;
