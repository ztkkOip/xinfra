import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/auth/Login.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/dashboard',
        },
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/Index.vue'),
          meta: { title: '资源大盘' },
        },
        {
          path: 'subsystem',
          name: 'Subsystem',
          component: () => import('@/views/subsystem/Navigation.vue'),
          meta: { title: '子系统导航' },
        },
        {
          path: 'audit/login',
          name: 'LoginAudit',
          component: () => import('@/views/audit/LoginAudit.vue'),
          meta: { title: '登录审计' },
        },
        {
          path: 'audit/ops',
          name: 'OpsAudit',
          component: () => import('@/views/audit/OpsAudit.vue'),
          meta: { title: '运维操作审计' },
        },
        {
          path: 'resource/status',
          name: 'ResourceStatus',
          component: () => import('@/views/resource/StatusBoard.vue'),
          meta: { title: '资源状态看板' },
        },
        {
          path: 'resource/mgmt',
          name: 'ResourceMgmt',
          component: () => import('@/views/resource/Management.vue'),
          meta: { title: '资源管理' },
        },
        {
          path: 'cluster',
          name: 'Cluster',
          component: () => import('@/views/cluster/ClusterList.vue'),
          meta: { title: '集群 / 节点' },
        },
        {
          path: 'service/catalog',
          name: 'ServiceCatalog',
          component: () => import('@/views/service/Catalog.vue'),
          meta: { title: '基础服务目录' },
        },
        {
          path: 'service/mgmt',
          name: 'ServiceMgmt',
          component: () => import('@/views/service/Management.vue'),
          meta: { title: '服务管理' },
        },
        {
          path: 'tasks',
          name: 'Tasks',
          component: () => import('@/views/task/TaskCenter.vue'),
          meta: { title: '任务中心' },
        },
        {
          path: 'config',
          name: 'ConfigCenter',
          component: () => import('@/views/config/ConfigCenter.vue'),
          meta: { title: '配置中心管理' },
        },
        {
          path: 'monitor',
          name: 'Monitor',
          component: () => import('@/views/monitor/Monitor.vue'),
          meta: { title: '监控 / 日志 / 告警' },
        },
      ],
    },
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = getToken()

  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
