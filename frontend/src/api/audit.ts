import type { ApiResponse, PaginatedData } from '@/types/api'

export interface LoginAudit {
  id: number
  user_id: number
  username: string
  login_time: string
  source_ip: string
  target_system: string
  status: 'success' | 'failed'
}

export interface OpsAudit {
  id: number
  user_id: number
  username: string
  operation_type: string
  operation: string
  target: string
  status: 'success' | 'failed' | 'running'
  created_at: string
}

// Mock 数据
const mockLoginAudits: LoginAudit[] = [
  { id: 1, user_id: 1, username: 'zhangsan', login_time: '2024-01-15T10:30:00Z', source_ip: '10.0.0.1', target_system: 'Wayne', status: 'success' },
  { id: 2, user_id: 2, username: 'lisi', login_time: '2024-01-15T10:25:00Z', source_ip: '10.0.0.2', target_system: 'CloudDM', status: 'success' },
  { id: 3, user_id: 3, username: 'wangwu', login_time: '2024-01-15T10:20:00Z', source_ip: '10.0.0.3', target_system: 'CacheCloud', status: 'failed' },
  { id: 4, user_id: 1, username: 'zhangsan', login_time: '2024-01-15T09:15:00Z', source_ip: '10.0.0.1', target_system: 'Grafana', status: 'success' },
  { id: 5, user_id: 4, username: 'zhaoliu', login_time: '2024-01-15T09:00:00Z', source_ip: '10.0.0.4', target_system: 'Apollo', status: 'success' },
]

const mockOpsAudits: OpsAudit[] = [
  { id: 1, user_id: 1, username: 'zhangsan', operation_type: 'ansible', operation: '执行 playbook: node-join', target: 'node-10.0.0.5', status: 'success', created_at: '2024-01-15T10:35:00Z' },
  { id: 2, user_id: 2, username: 'lisi', operation_type: 'ansible', operation: '执行 playbook: mysql-deploy', target: 'mysql-kodo-01', status: 'running', created_at: '2024-01-15T10:30:00Z' },
  { id: 3, user_id: 1, username: 'zhangsan', operation_type: 'ansible', operation: '执行 playbook: redis-deploy', target: 'redis-las-cluster', status: 'success', created_at: '2024-01-15T10:25:00Z' },
  { id: 4, user_id: 3, username: 'wangwu', operation_type: 'ansible', operation: '执行 playbook: openresty-deploy', target: 'openresty-gateway', status: 'failed', created_at: '2024-01-15T10:20:00Z' },
]

export const auditApi = {
  // 获取登录审计
  getLoginAudit(_params?: any): Promise<ApiResponse<PaginatedData<LoginAudit>>> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 0,
          message: 'success',
          data: {
            total: mockLoginAudits.length,
            items: mockLoginAudits,
          },
        })
      }, 300)
    })
  },

  // 获取运维操作审计
  getOpsAudit(_params?: any): Promise<ApiResponse<PaginatedData<OpsAudit>>> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 0,
          message: 'success',
          data: {
            total: mockOpsAudits.length,
            items: mockOpsAudits,
          },
        })
      }, 300)
    })
  },
}
