import request from './request'

export interface Subsystem {
  id: number
  name: string
  description: string
  icon: string
  url: string
  status: 'integrated' | 'integrating'
  sso_enabled: boolean
}

// Mock 数据
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
  // 获取子系统列表
  getSubsystems(): Promise<any> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 0,
          message: 'success',
          data: mockSubsystems,
        })
      }, 300)
    })
  },

  // 获取 SSO 跳转 URL
  getSSOUrl(id: number): Promise<any> {
    return new Promise((resolve) => {
      setTimeout(() => {
        const subsystem = mockSubsystems.find((s) => s.id === id)
        resolve({
          code: 0,
          message: 'success',
          data: {
            sso_url: subsystem?.url + '/sso/callback?code=mock-code&state=mock-state',
            expires_in: 300,
          },
        })
      }, 200)
    })
  },
}
