import type { ApiResponse } from '@/types/api'
import { getToken, removeToken } from '@/utils/auth'
import { redirectToSSO } from '@/utils/sso'

export interface Subsystem {
  id: number
  name: string
  description: string
  icon: string
  url: string
  status: 'integrated' | 'integrating'
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
    description: '多集群容器管理平台',
    icon: 'W',
    url: 'https://wayne.qiniu.com',
    status: 'integrated',
    sso_enabled: true,
  },
  {
    id: 2,
    name: 'CloudDM',
    description: '数据库管理与SQL审核',
    icon: 'DM',
    url: 'https://clouddm.qiniu.com',
    status: 'integrated',
    sso_enabled: true,
  },
  {
    id: 3,
    name: 'CacheCloud',
    description: 'Redis 云管理平台',
    icon: 'CC',
    url: 'https://cachecloud.qiniu.com',
    status: 'integrated',
    sso_enabled: true,
  },
  {
    id: 4,
    name: 'Apollo',
    description: '配置中心',
    icon: 'AP',
    url: 'https://apollo.xinfra.internal',
    status: 'integrated',
    sso_enabled: true,
  },
  {
    id: 5,
    name: 'qpass',
    description: '密码管理平台',
    icon: 'QP',
    url: 'https://qpass.xinfra.internal',
    status: 'integrated',
    sso_enabled: true,
  },
  {
    id: 6,
    name: 'Grafana',
    description: '监控可视化平台',
    icon: 'GF',
    url: 'https://grafana.xinfra.internal',
    status: 'integrated',
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
