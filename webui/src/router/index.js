import { createRouter, createWebHistory } from 'vue-router';
import AppLayout from '@/layout/AppLayout.vue';
// import process from "eslint-plugin-vue/lib/configs/base.js";


const router = createRouter({
  // history: createWebHistory(process.env.VITE_BASE || ''),
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: AppLayout,
      children: [
        {
          path: '/',
          name: 'home',
          component: () => import('@/views/HomeView.vue')
        },
        {
          path: '/about',
          name: 'about',
          component: () => import('@/views/AboutView.vue')
        },
        // {
        //   path: '/uikit/menu',
        //   component: () => import('@/views/uikit/Menu.vue'),
        //   children: [
        //     {
        //       path: '/uikit/menu',
        //       component: () => import('@/views/uikit/menu/PersonalDemo.vue')
        //     },
        //     {
        //       path: '/uikit/menu/seat',
        //       component: () => import('@/views/uikit/menu/SeatDemo.vue')
        //     },
        //     {
        //       path: '/uikit/menu/payment',
        //       component: () => import('@/views/uikit/menu/PaymentDemo.vue')
        //     },
        //     {
        //       path: '/uikit/menu/confirmation',
        //       component: () => import('@/views/uikit/menu/ConfirmationDemo.vue')
        //     }
        //   ]
        // },
      ],
    },
    {
        path: '/login',
        name: 'login',
        component: () => import('@/views/LoginView.vue')
    },
    // {
    //   path: '/landing',
    //   name: 'landing',
    //   component: () => import('@/views/pages/Landing.vue')
    // },
    // {
    //   path: '/pages/notfound',
    //   name: 'notfound',
    //   component: () => import('@/views/pages/NotFound.vue')
    // },
    //
    // {
    //   path: '/auth/login',
    //   name: 'login',
    //   component: () => import('@/views/pages/auth/Login.vue')
    // },
    // {
    //   path: '/auth/access',
    //   name: 'accessDenied',
    //   component: () => import('@/views/pages/auth/Access.vue')
    // },
    // {
    //   path: '/auth/error',
    //   name: 'error',
    //   component: () => import('@/views/pages/auth/Error.vue')
    // }
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/404.vue')
    },

    // { path: "*", component: {        template: '<p>Page Not Found</p>'      }
    
  ]
});

export default router;