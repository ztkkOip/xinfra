import request from './request'

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

// Mock 数据
const mockUser: UserInfo = {
  id: 1,
  username: 'admin',
  display_name: '王晓东',
  email: 'admin@qiniu.com',
  business_line: 'kodo',
}

export const authApi = {
  // 登录
  login(data: LoginRequest): Promise<any> {
    // MVP 阶段使用 Mock 数据
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 0,
          message: 'success',
          data: {
            token: 'mock-jwt-token-' + Date.now(),
            expires_in: 3600,
            user: mockUser,
          },
        })
      }, 500)
    })
  },

  // 登出
  logout(): Promise<any> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({ code: 0, message: 'success' })
      }, 200)
    })
  },

  // 获取用户信息
  getUserInfo(): Promise<any> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 0,
          message: 'success',
          data: mockUser,
        })
      }, 200)
    })
  },
}
