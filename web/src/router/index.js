import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/front/Home.vue')
  },
  {
    path: '/detail/:id',
    name: 'Detail',
    component: () => import('../views/front/Detail.vue')
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('../views/admin/Login.vue')
  },
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('../views/admin/Layout.vue'),
    children: [
      {
        path: '',
        redirect: '/admin/documents'
      },
      {
        path: 'documents',
        name: 'AdminDocuments',
        component: () => import('../views/admin/Document.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
