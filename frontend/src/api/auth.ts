import type { ApiResponse } from '@/types/api'
import { getToken } from '@/utils/auth'

export interface LoginRequest {
  username: string
  password: string
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

export const authApi = {
  login(_data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return Promise.reject(new Error('password login is disabled'))
  },

  logout(): Promise<ApiResponse<null>> {
    return Promise.resolve({ code: 0, message: 'success', data: null })
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
