import type { ApiResponse } from '@/types/api'
import { getToken } from '@/utils/auth'

export interface LoginRequest {
  username: string
  password?: string
}

export interface UserInfo {
  id: number
  username: string
  display_name: string
  email: string
  business_line: string
}

export interface LoginResponse {
  token: string
  expires_in: number
  user: UserInfo
}

export interface AuthConfig {
  sso_enabled: boolean
}

export const authApi = {
  async getConfig(): Promise<ApiResponse<AuthConfig>> {
    const response = await fetch('/auth/api/v1/config', {
      headers: {
        Accept: 'application/json',
      },
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      throw new Error(data.error || `HTTP ${response.status}`)
    }
    return {
      code: 0,
      message: 'success',
      data: {
        sso_enabled: data.sso_enabled !== false,
      },
    }
  },

  async login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    const response = await fetch('/auth/api/v1/login', {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username: data.username }),
    })
    const result = await response.json().catch(() => ({}))
    if (!response.ok) {
      throw new Error(result.error || `HTTP ${response.status}`)
    }
    return {
      code: 0,
      message: 'success',
      data: result,
    }
  },

  async logout(): Promise<ApiResponse<null>> {
    const token = getToken()
    const response = await fetch('/auth/api/v1/logout', {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      credentials: 'include',
    })
    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.error || `HTTP ${response.status}`)
    }
    return { code: 0, message: 'success', data: null }
  },

  async getUserInfo(): Promise<ApiResponse<UserInfo>> {
    const token = getToken()
    const response = await fetch('/auth/api/v1/users/me', {
      headers: {
        Accept: 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      throw new Error(data.error || `HTTP ${response.status}`)
    }
    return {
      code: 0,
      message: 'success',
      data: {
        id: data.id,
        username: data.username,
        display_name: data.display_name || data.username,
        email: data.email || '',
        business_line: data.business_line || '',
      },
    }
  },
}
