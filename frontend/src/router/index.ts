import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/auth'
import { consumeSSOToken, redirectToSSO } from '@/utils/sso'
import { useAuthStore } from '@/stores/auth'
import { useBusinessLineStore } from '@/stores/businessLine'
import { authApi } from '@/api/auth'

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
          path: 'subsystem/detail/:name',
          name: 'SubsystemDetail',
          component: () => import('@/views/subsystem/SubsystemDetail.vue'),
          meta: { title: '子系统详情' },
        },
        {
          path: 'subsystem/authorization',
          name: 'SubsystemAuthorization',
          component: () => import('@/views/subsystem/Authorization.vue'),
          meta: { title: '子系统赋权' },
        },
        {
          path: 'business-line/permissions',
          redirect: '/business-line/manage',
        },
        {
          path: 'business-line/manage',
          name: 'BusinessLineManage',
          component: () => import('@/views/businessLine/Manage.vue'),
          meta: { title: '业务线管理' },
        },
        {
          path: 'business-line/assignment',
          name: 'BusinessLineAssignment',
          component: () => import('@/views/businessLine/Assignment.vue'),
          meta: { title: '业务线分配' },
        },
        {
          path: 'audit/login',
          name: 'LoginAudit',
          component: () => import('@/views/audit/LoginAudit.vue'),
          meta: { title: '登录审计' },
        },
        {
          path: 'audit/operation',
          name: 'OperationAudit',
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
          path: 'machine/management',
          name: 'MachineManagement',
          component: () => import('@/views/resource/Management.vue'),
          meta: { title: '机器管理' },
        },
        {
          path: 'infrastructure/cluster',
          name: 'InfrastructureCluster',
          component: () => import('@/views/cluster/ClusterList.vue'),
          meta: { title: '集群与容器' },
        },
        {
          path: 'business/quota',
          name: 'BusinessQuota',
          component: () => import('@/views/resource/Tenants.vue'),
          meta: { title: '业务配额' },
        },
        {
          path: 'service/catalog',
          name: 'ServiceCatalog',
          component: () => import('@/views/service/Catalog.vue'),
          meta: { title: '基础服务' },
        },
        {
          path: 'service/management',
          name: 'ServiceManagement',
          component: () => import('@/views/service/Management.vue'),
          meta: { title: '服务管理' },
        },
        {
          path: 'task/log',
          name: 'TaskLog',
          component: () => import('@/views/task/TaskCenter.vue'),
          meta: { title: '任务日志' },
        },
        {
          path: 'config',
          name: 'ConfigCenter',
          component: () => import('@/views/config/ConfigCenter.vue'),
          meta: { title: '配置中心管理' },
        },
        {
          path: 'observable/monitoring',
          name: 'MonitoringDashboard',
          component: () => import('@/views/monitor/Monitor.vue'),
          meta: { title: '监控看板' },
        },
        {
          path: 'observable/alert',
          name: 'AlertInformation',
          component: () => import('@/views/monitor/Alert.vue'),
          meta: { title: '告警信息' },
        },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()
  const ssoToken = consumeSSOToken()
  if (ssoToken) {
    authStore.setSessionToken(ssoToken)
  }
  const token = authStore.token || getToken()

  if (to.meta.requiresAuth !== false && !token) {
    const { data } = await authApi.getConfig()
    if (!data.sso_enabled) {
      return { path: '/login', query: to.fullPath === '/' ? {} : { redirect: to.fullPath } }
    }
    redirectToSSO()
    return false
  }
  if (token && (!authStore.user || ssoToken)) {
    try {
      await authStore.refreshUser()
      await useBusinessLineStore().loadMine()
    } catch {
      authStore.clearAuth()
      const { data } = await authApi.getConfig()
      if (!data.sso_enabled) {
        return { path: '/login', query: to.fullPath === '/' ? {} : { redirect: to.fullPath } }
      }
      redirectToSSO()
      return false
    }
  } else if (token && !useBusinessLineStore().businessLines.length) {
    await useBusinessLineStore().loadMine().catch(() => {})
  }
  if (to.path === '/login' && token) {
    return '/'
  }
  return true
})

export default router
