import { createRouter, createWebHistory } from 'vue-router'
import MainAppView from '@/views/MainAppView.vue'
import { useUserStore } from '@/stores/user.js'
import { ref } from 'vue'
// import process from "eslint-plugin-vue/lib/configs/base.js";

let firstLoad = true

const router = createRouter({
  // history: createWebHistory(),
  history: createWebHistory('/spa'),
  routes: [
    {
      path: '/',
      name: 'landing',
      meta: {
        // requiresAuth: true
      },
      component: () => import('@/views/LandingPage.vue')
    },
    {
      path: '/app',
      name: 'home',
      meta: {
        requiresAuth: true
      },
      component: () => import('@/views/MainAppView.vue')
    },
    {
      path: '/login',
      name: 'login',
      meta: {
        hideFromAuth: true
      },
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

    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/404.vue')
    }

    // { path: "*", component: {        template: '<p>Page Not Found</p>'      }
  ]
})

// this checks for metadata in the router and redirects to login page if the user is not logged in
// the same happens if the user is logged in he is redirected to the entry away from the login page
// this relies on the user store
// based on: https://stackoverflow.com/questions/52653337/vuejs-redirect-from-login-register-to-home-if-already-loggedin-redirect-from
router.beforeEach((to, from, next) => {
  const user = useUserStore()

  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (!user.isLoggedIn) {
      next({ name: 'login' })
    } else {
      next() // go to wherever I'm going
    }
  } else {
    next() // does not require auth, make sure to always call next()!
  }

  if (to.matched.some((record) => record.meta.hideFromAuth)) {
    if (user.isLoggedIn) {
      next({ name: 'login' })
    } else {
      next()
    }
  } else {
    next()
  }
})

router.beforeEach((to, from, next) => {
  const user = useUserStore()
  if (firstLoad) {
    firstLoad = false
    user.checkState()
  }
  next()
})

export default router
