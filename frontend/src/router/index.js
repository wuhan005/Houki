import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue')
    },
    {
        path: '/modules',
        name: 'Modules',
        component: () => import('@/views/Modules.vue')
    },
    {
        path: '/store',
        name: 'Store',
        component: () => import('@/views/Store.vue')
    },
    {
        path: '/logs',
        name: 'Logs',
        component: () => import('@/views/Logs.vue')
    },
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

export default router
