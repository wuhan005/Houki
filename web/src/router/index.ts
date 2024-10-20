import { createRouter, createWebHistory } from 'vue-router';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

NProgress.configure({ showSpinner: false }); // NProgress Configuration

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '',
      name: 'landing',
      redirect: '/modules',
      component: () => import('@/layout/default-layout.vue'),
      children: [
        {
          path: '/modules',
          name: 'modules',
          component: () => import('@/views/modules.vue'),
        },
        {
          path: '/modules/new',
          name: 'new-module',
          component: () => import('@/views/module-editor.vue'),
        },
        {
          path: '/modules/:id',
          name: 'update-module',
          component: () => import('@/views/module-editor.vue'),
        },
        {
          path: '/certificate',
          name: 'certificate',
          component: () => import('@/views/certificate.vue'),
        },
        {
          path: '/proxy',
          name: 'proxy',
          component: () => import('@/views/proxy.vue'),
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'notFound',
      component: () => import('@/views/not-found.vue'),
    },
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});

export default router;
