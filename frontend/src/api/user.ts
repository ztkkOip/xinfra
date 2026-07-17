import { getToken } from '@/utils/auth'

export interface UserOption {
  uid: number
  username: string
}

export const userApi = {
  async list(): Promise<UserOption[]> {
    const token = getToken()
    const response = await fetch('/auth/api/v1/users', {
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
