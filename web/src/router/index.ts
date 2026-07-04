import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { setAPIRouter } from '@/api/client'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      component: () => import('@/components/layout/AppLayout.vue'),
      children: [
        { path: '', name: 'dashboard', component: () => import('@/views/Dashboard.vue') },
        { path: 'sites', name: 'sites', component: () => import('@/views/Placeholder.vue'), props: { title: '网站管理' } },
        { path: 'environment', name: 'environment', component: () => import('@/views/Placeholder.vue'), props: { title: '环境管理' } },
        { path: 'database', name: 'database', component: () => import('@/views/Placeholder.vue'), props: { title: '数据库' } },
        { path: 'files', name: 'files', component: () => import('@/views/Placeholder.vue'), props: { title: '文件管理' } },
        { path: 'cron', name: 'cron', component: () => import('@/views/Placeholder.vue'), props: { title: '计划任务' } },
        { path: 'logs', name: 'logs', component: () => import('@/views/Placeholder.vue'), props: { title: '日志' } },
        { path: 'settings', name: 'settings', component: () => import('@/views/Placeholder.vue'), props: { title: '设置' } },
      ],
    },
  ],
})

setAPIRouter(router)

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.isLoggedIn) {
    return { name: 'dashboard' }
  }
})

export default router
