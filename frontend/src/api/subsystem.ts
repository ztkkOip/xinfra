import type { ApiResponse } from '@/types/api'
import { getToken, removeToken } from '@/utils/auth'
import { redirectToSSO } from '@/utils/sso'

export interface Subsystem {
  id: number
  name: string
  label: string
  description: string
  icon: string
  url: string
  domain: string
  category: string
  status: 'active' | 'integrated' | 'integrating'
  sso_enabled: boolean
}

export interface SSOResponse {
  sso_url: string
  expires_in: number
}

const mockSubsystems: Subsystem[] = [
  {
    id: 1,
    name: 'Wayne',
    label: '多集群发布平台',
    description: '业务容器发布、命名空间与配额管理，复用 Wayne 原生多租户能力。',
    icon: 'W',
    url: 'https://wayne.qiniu.com',
    domain: 'wayne.xinfra.internal',
    category: '容器管理',
    status: 'active',
    sso_enabled: true,
  },
  {
    id: 2,
    name: 'CloudDM',
    label: 'SQL 审核平台',
    description: '数据库 SQL 上线统一审核，支持 LDAP 用户组到角色的自动映射。',
    icon: 'DM',
    url: 'https://clouddm.qiniu.com',
    domain: 'clouddm.xinfra.internal',
    category: '数据库',
    status: 'active',
    sso_enabled: true,
  },
  {
    id: 3,
    name: 'CacheCloud',
    label: 'Redis 管理平台',
    description: 'Redis 实例全生命周期管理，登录与监控运维统一入口。',
    icon: 'CC',
    url: 'https://cachecloud.qiniu.com',
    domain: 'cachecloud.xinfra.internal',
    category: '缓存',
    status: 'integrating',
    sso_enabled: true,
  },
  {
    id: 4,
    name: 'Apollo',
    label: '配置中心',
    description: '统一配置管理，按机房 Cluster 隔离，支持灰度发布、版本回滚与变更审计。',
    icon: 'AP',
    url: 'https://apollo.xinfra.internal',
    domain: 'apollo.xinfra.internal',
    category: '配置中心',
    status: 'integrating',
    sso_enabled: true,
  },
  {
    id: 5,
    name: 'qpass',
    label: '七牛统一运维平台',
    description: '七牛云内部系统。',
    icon: 'QP',
    url: 'https://qpass.xinfra.internal',
    domain: 'qpass.xinfra.internal',
    category: '内部系统',
    status: 'integrating',
    sso_enabled: true,
  },
  {
    id: 6,
    name: 'Grafana',
    label: '指标可视化平台',
    description: '对接 VictoriaMetrics 数据源，K8s / 容器 / 业务指标统一仪表盘展示。',
    icon: 'GF',
    url: 'https://grafana.xinfra.internal',
    domain: 'grafana.xinfra.internal',
    category: '可观测性',
    status: 'integrating',
    sso_enabled: true,
  },
  {
    id: 7,
    name: 'Superset',
    label: '日志数据可视化',
    description: '日志数据探索与可视化分析平台，支持多数据源接入和丰富的图表展示。',
    icon: 'SS',
    url: 'https://superset.xinfra.internal',
    domain: 'superset.xinfra.internal',
    category: '可观测性',
    status: 'integrating',
    sso_enabled: true,
  },
]

export const subsystemApi = {
  getSubsystems(): Promise<ApiResponse<Subsystem[]>> {
    return Promise.resolve({
      code: 0,
      message: 'success',
      data: mockSubsystems,
    })
  },

  getSubsystemByName(name: string): Subsystem | undefined {
    return mockSubsystems.find((s) => s.name.toLowerCase() === name.toLowerCase())
  },

  async getSSOUrl(id: number): Promise<ApiResponse<SSOResponse>> {
    const subsystem = mockSubsystems.find((s) => s.id === id)
    if (!subsystem) {
      throw new Error('subsystem not found')
    }
    if (subsystem.name !== 'Wayne' && subsystem.name !== 'CloudDM') {
      return {
        code: 0,
        message: 'success',
        data: {
          sso_url: subsystem.url,
          expires_in: 300,
        },
      }
    }

    const token = getToken()
    const openApp = subsystem.name === 'CloudDM' ? 'clouddm' : 'wayne'
    const path = subsystem.name === 'CloudDM' ? '/auth/api/v1/clouddm/login' : '/auth/api/v1/wayen/login'
    const response = await fetch(path, {
      headers: {
        Accept: 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      credentials: 'include',
    })
    const data = await response.json().catch(() => ({}))
    if (response.status === 401) {
      removeToken()
      const configResponse = await fetch('/auth/api/v1/config', {
        headers: { Accept: 'application/json' },
      }).catch(() => null)
      const config = await configResponse?.json().catch(() => ({}))
      if (config?.sso_enabled === false) {
        window.location.assign(`/login?redirect=${encodeURIComponent(`/subsystem?open_app=${openApp}`)}`)
      } else {
        redirectToSSO(openApp)
      }
      throw new Error('unauthorized')
    }
    if (!response.ok) {
      throw new Error(data.error || `HTTP ${response.status}`)
    }
    if (!data.target_url) {
      throw new Error(`${subsystem.name} 跳转地址为空`)
    }
    return {
      code: 0,
      message: 'success',
      data: {
        sso_url: data.target_url,
        expires_in: 300,
      },
    }
  },
}
