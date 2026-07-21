import { getToken } from '@/utils/auth'

export interface BusinessLine {
  id: number
  name: string
  created_at: string
  updated_at: string
  permission?: 0 | 1
}

export interface WayneNamespace {
  id: number
  name: string
  kubeNamespace: string
}

export const businessLineApi = {
  async listMine(): Promise<BusinessLine[]> {
    const token = getToken()
    const response = await fetch('/auth/api/v1/business-lines', {
      headers: {
        Accept: 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      throw new Error(data.error || `HTTP ${response.status}`)
    }
    return Array.isArray(data.items) ? data.items : []
  },

  async listAll(): Promise<BusinessLine[]> {
    const data = await request('/auth/api/v1/business-lines/all')
    return Array.isArray(data.items) ? data.items : []
  },

  async create(name: string): Promise<BusinessLine> {
    return request('/auth/api/v1/business-lines', {
      method: 'POST',
      body: JSON.stringify({ name }),
    })
  },

  async update(id: number, name: string): Promise<BusinessLine> {
    return request(`/auth/api/v1/business-lines/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name }),
    })
  },

  async remove(id: number): Promise<void> {
    await request(`/auth/api/v1/business-lines/${id}`, {
      method: 'DELETE',
    })
  },

  async grant(payload: {
    business_line_id: number
    target_user_id: number
    target_business_line_id: number
  }): Promise<void> {
    await request('/auth/api/v1/business-lines/authorizations', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  async listWayneNamespaces(): Promise<WayneNamespace[]> {
    const data = await request('/auth/api/v1/wayne/namespaces')
    const items = Array.isArray(data.data?.list) ? data.data.list : []
    return items.map((item: any) => ({
      id: Number(item.id),
      name: item.name || '',
      kubeNamespace: item.kubeNamespace || item.kube_namespace || '',
    }))
  },

  async listMappedWayneNamespaces(businessLineId: number): Promise<WayneNamespace[]> {
    const data = await request(`/auth/api/v1/business-lines/${businessLineId}/wayne-namespaces`)
    return Array.isArray(data.items) ? data.items : []
  },

  async replaceMappedWayneNamespaces(businessLineId: number, namespaces: WayneNamespace[]): Promise<void> {
    await request(`/auth/api/v1/business-lines/${businessLineId}/wayne-namespaces`, {
      method: 'PUT',
      body: JSON.stringify({ namespaces }),
    })
  },
}

async function request(path: string, init: RequestInit = {}) {
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
  const data = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new Error(data.error || `HTTP ${response.status}`)
  }
  return data
}
