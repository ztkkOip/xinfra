import { getToken } from '@/utils/auth'

export interface SubsystemAuthSystem {
  key: string
  name: string
  enabled: boolean
}

export interface WayneRole {
  id: number
  name: string
  comment?: string
  type: number
}

export interface WaynePermission {
  create: boolean
  update: boolean
  delete: boolean
}

export interface WayneBusinessLineNamespace {
  id: number
  name: string
  kubeNamespace: string
  permissions?: WaynePermission
  can_bind?: boolean
  can_unbind?: boolean
  permission_error?: string
}

export interface WayneRoleBindingPayload {
  groupIds?: number[]
  replace?: boolean
  requestId?: string
  reason?: string
  dryRun?: boolean
}

export interface WayneUserRoles {
  userId?: number
  userName?: string
  namespaces?: Array<{
    namespace?: {
      id: number
      name: string
    }
    groups?: Array<{
      id: number
      name: string
    }>
  }>
  apps?: unknown[]
}

export const subsystemAuthApi = {
  async listSystems(): Promise<SubsystemAuthSystem[]> {
    const data = await authRequest('/auth/api/v1/subsystem-auth/systems')
    return Array.isArray(data.items) ? data.items : []
  },

  async listWayneRoles(): Promise<WayneRole[]> {
    const data = await authRequest('/auth/api/v1/subsystem-auth/wayne/roles')
    return Array.isArray(data.items) ? data.items : []
  },

  async listWayneNamespaces(businessLineId: number): Promise<WayneBusinessLineNamespace[]> {
    const data = await authRequest(`/auth/api/v1/subsystem-auth/wayne/business-lines/${businessLineId}/namespaces`)
    return Array.isArray(data.items) ? data.items : []
  },

  async getWayneUserRoles(username: string): Promise<WayneUserRoles> {
    const data = await authRequest(`/auth/api/v1/subsystem-auth/wayne/users/${encodeURIComponent(username)}/roles`)
    return data.data || data
  },

  async bindWayneNamespaceRoles(
    businessLineId: number,
    namespaceId: number,
    username: string,
    payload: WayneRoleBindingPayload,
  ): Promise<unknown> {
    return authRequest(
      `/auth/api/v1/subsystem-auth/wayne/business-lines/${businessLineId}/namespaces/${namespaceId}/users/${encodeURIComponent(username)}/roles`,
      {
        method: 'PUT',
        body: JSON.stringify(payload),
      },
    )
  },

  async unbindWayneNamespaceRoles(
    businessLineId: number,
    namespaceId: number,
    username: string,
    payload: WayneRoleBindingPayload = {},
  ): Promise<unknown> {
    return authRequest(
      `/auth/api/v1/subsystem-auth/wayne/business-lines/${businessLineId}/namespaces/${namespaceId}/users/${encodeURIComponent(username)}/roles`,
      {
        method: 'DELETE',
        body: JSON.stringify(payload),
      },
    )
  },

  async initWayneBusinessLineUser(businessLineId: number, userId: number): Promise<unknown> {
    return authRequest(`/auth/api/v1/subsystem-auth/wayne/business-lines/${businessLineId}/users/${userId}/init`, {
      method: 'POST',
    })
  },
}

async function authRequest(path: string, init: RequestInit = {}) {
  const token = getToken()
  const response = await fetch(path, {
    ...init,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...init.headers,
    },
  })
  const text = await response.text()
  const data = parseResponseBody(text)
  if (!response.ok) {
    const message = data?.error || data?.message || text || `HTTP ${response.status}`
    throw new Error(message)
  }
  return data || {}
}

function parseResponseBody(text: string) {
  if (!text.trim()) {
    return {}
  }
  try {
    return JSON.parse(text)
  } catch {
    return { error: text }
  }
}
