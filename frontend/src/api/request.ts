import axios, { type AxiosResponse, type AxiosError } from 'axios'
import { getToken, removeToken } from '@/utils/auth'
import { redirectToSSO } from '@/utils/sso'
import { ElMessage } from 'element-plus'
import type { ApiResponse } from '@/types/api'
import { codeMessageMap } from '@/types/api'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

/** HTTP 状态码对应的用户友好提示 */
const httpMessageMap: Record<number, string> = {
  400: '请求参数错误',
  401: '未授权，请重新登录',
  403: '没有权限执行此操作',
  404: '请求的资源不存在',
  408: '请求超时，请稍后重试',
  500: '服务器内部错误',
  502: '网关错误',
  503: '服务暂时不可用，请稍后重试',
}

/** 根据后端错误码获取用户友好提示 */
function getBusinessMessage(code: number, fallback?: string): string {
  return codeMessageMap[code] || fallback || '请求失败'
}

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { data } = response
    if (data.code !== 0) {
      const message = getBusinessMessage(data.code, data.message)
      ElMessage.error(message)
      return Promise.reject(new Error(message))
    }
    return data as any
  },
  (error: AxiosError<ApiResponse>) => {
    // 优先使用后端返回的错误信息
    const backendMessage = error.response?.data?.message

    if (error.response) {
      const status = error.response.status

      if (status === 401) {
        removeToken()
        redirectToSSO()
        ElMessage.error(backendMessage || '未授权，请重新登录')
      } else {
        const message = backendMessage || httpMessageMap[status] || `请求失败 (${status})`
        ElMessage.error(message)
      }
    } else if (error.code === 'ECONNABORTED') {
      ElMessage.error('请求超时，请检查网络后重试')
    } else {
      ElMessage.error('网络连接异常，请检查网络')
    }

    return Promise.reject(error)
  }
)

export default request
