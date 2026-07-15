/**
 * 后端统一响应格式类型定义
 * 对应 server/internal/handler/response.go
 */

/** 统一 JSON 响应信封 */
export interface ApiResponse<T = any> {
  /** 业务状态码，0 表示成功 */
  code: number
  /** 提示信息 */
  message: string
  /** 响应数据 */
  data: T
  /** 参数校验错误详情（可选） */
  details?: Record<string, string[]>
}

/** 分页数据结构 */
export interface PaginatedData<T> {
  total: number
  items: T[]
}

/** 业务错误码常量，与后端一一对应 */
export const ApiCode = {
  /** 成功 */
  SUCCESS: 0,
  /** LDAP 认证失败 */
  LDAP_AUTH_FAILED: 10001,
  /** 用户名或密码错误 */
  INVALID_CREDENTIALS: 10002,
  /** Token 已过期 */
  TOKEN_EXPIRED: 10003,
  /** SSO 未配置 */
  SSO_NOT_CONFIGURED: 10004,
  /** SSO 跳转地址生成失败 */
  SSO_GENERATE_FAILED: 10005,
  /** 任务创建失败 */
  TASK_CREATE_FAILED: 20001,
  /** 任务执行超时 */
  TASK_EXEC_TIMEOUT: 20002,
} as const

/** 错误码到用户友好提示的映射 */
export const codeMessageMap: Record<number, string> = {
  [ApiCode.LDAP_AUTH_FAILED]: 'LDAP 认证服务异常，请稍后重试',
  [ApiCode.INVALID_CREDENTIALS]: '用户名或密码错误',
  [ApiCode.TOKEN_EXPIRED]: '登录已过期，请重新登录',
  [ApiCode.SSO_NOT_CONFIGURED]: '该子系统未配置 SSO，请联系管理员',
  [ApiCode.SSO_GENERATE_FAILED]: '生成 SSO 跳转链接失败，请稍后重试',
  [ApiCode.TASK_CREATE_FAILED]: '创建任务失败，请检查参数后重试',
  [ApiCode.TASK_EXEC_TIMEOUT]: '任务执行超时，请稍后查看结果',
}
