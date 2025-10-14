import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      redirect: (to) => {
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        return isLoggedIn ? '/dashboard' : '/login'
      },
    },
    {
      path: '/tenaga/edit/:id',
      name: 'edit-worker',
      component: () => import('../views/EditWorkerView.vue'),
      props: true,
      beforeEnter: (to, from, next) => {
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        if (isLoggedIn) next()
        else next('/login')
      },
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
    },
    {
      path: '/tenaga/tambah',
      name: 'add-worker',
      component: () => import('../views/AddWorkerView.vue'),
    },
    {
      path: '/tenaga/form-lanjutan/:id?',
      name: 'registration-form',
      component: () => import('../views/RegistrationFormView.vue'),
      props: true,
      beforeEnter: (to, from, next) => {
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        if (isLoggedIn) next()
        else next('/login')
      },
    },
    {
      path: '/tenaga/upload',
      name: 'upload-worker',
      component: () => import('../views/UploadWorkerView.vue'),
      beforeEnter: (to, from, next) => {
        // Check if user is logged in
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        if (isLoggedIn) {
          next()
        } else {
          next('/login')
        }
      },
    },
    {
      path: '/tenaga/upload-na',
      name: 'upload-tk-na',
      component: () => import('../views/UploadTKNAView.vue'),
      beforeEnter: (to, from, next) => {
        // Check if user is logged in
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        if (isLoggedIn) {
          next()
        } else {
          next('/login')
        }
      },
    },
    {
      path: '/tenaga/upload-upah',
      name: 'upload-upah',
      component: () => import('../views/UploadUpahView.vue'),
      beforeEnter: (to, from, next) => {
        // Check if user is logged in
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        if (isLoggedIn) {
          next()
        } else {
          next('/login')
        }
      },
    },
    {
      path: '/tenaga/koreksi-data',
      name: 'koreksi-data',
      component: () => import('../views/KoreksiDataView.vue'),
      beforeEnter: (to, from, next) => {
        // Check if user is logged in
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        if (isLoggedIn) {
          next()
        } else {
          next('/login')
        }
      },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
      beforeEnter: (to, from, next) => {
        // Check if user is logged in
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        if (isLoggedIn) {
          next()
        } else {
          next('/login')
        }
      },
    },
    {
      path: '/edit/:data',
      name: 'edit',
      component: () => import('../views/EditView.vue'),
      beforeEnter: (to, from, next) => {
        // Check if user is logged in
        const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
        if (isLoggedIn) {
          next()
        } else {
          next('/login')
        }
      },
    },
  ],
})

export default router
