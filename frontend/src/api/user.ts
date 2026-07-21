import { getToken } from '@/utils/auth'

export interface UserOption {
  uid: number
  username: string
}

export const userApi = {
  async list(params: { businessLineId?: number } = {}): Promise<UserOption[]> {
    const token = getToken()
    const query = new URLSearchParams()
    if (params.businessLineId) {
      query.set('business_line_id', String(params.businessLineId))
    }
    const path = query.toString() ? `/auth/api/v1/users?${query}` : '/auth/api/v1/users'
    const response = await fetch(path, {
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
}
