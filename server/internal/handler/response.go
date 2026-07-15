package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`             // 业务错误码，0 表示成功
	Message string      `json:"message"`          // 响应消息
	Data    interface{} `json:"data,omitempty"`    // 响应数据，成功时返回
	Details string      `json:"details,omitempty"` // 错误详情，失败时返回
}

// 分页响应结构
type PaginatedData struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

// 业务错误码定义
const (
	// 成功
	CodeSuccess = 0

	// 认证相关错误 10001-10099
	CodeLDAPAuthFailed     = 10001 // LDAP 认证失败
	CodeInvalidCredentials = 10002 // 用户名或密码错误
	CodeTokenExpired       = 10003 // Token 已过期
	CodeSSONotConfigured   = 10004 // 子系统未配置 SSO
	CodeSSOGenerateFailed  = 10005 // SSO 跳转生成失败

	// 任务相关错误 20001-20099
	CodeTaskCreateFailed  = 20001 // Ansible 任务创建失败
	CodeTaskExecTimeout   = 20002 // Ansible 任务执行超时
)

// 错误码对应的默认消息
var codeMessages = map[int]string{
	CodeSuccess:            "success",
	CodeLDAPAuthFailed:     "LDAP 认证失败",
	CodeInvalidCredentials: "用户名或密码错误",
	CodeTokenExpired:       "Token 已过期",
	CodeSSONotConfigured:   "子系统未配置 SSO",
	CodeSSOGenerateFailed:  "SSO 跳转生成失败",
	CodeTaskCreateFailed:   "Ansible 任务创建失败",
	CodeTaskExecTimeout:    "Ansible 任务执行超时",
}

// GetMessage 根据错误码获取默认消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: GetMessage(CodeSuccess),
		Data:    data,
	})
}

// SuccessWithMessage 返回自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPaginated 返回分页成功响应
func SuccessWithPaginated(c *gin.Context, total int64, items interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: GetMessage(CodeSuccess),
		Data: PaginatedData{
			Total: total,
			Items: items,
		},
	})
}

// Error 返回错误响应
func Error(c *gin.Context, httpCode int, bizCode int, details string) {
	c.JSON(httpCode, Response{
		Code:    bizCode,
		Message: GetMessage(bizCode),
		Details: details,
	})
}

// ErrorWithMessage 返回自定义消息的错误响应
func ErrorWithMessage(c *gin.Context, httpCode int, bizCode int, message string, details string) {
	c.JSON(httpCode, Response{
		Code:    bizCode,
		Message: message,
		Details: details,
	})
}

// BadRequest 返回 400 错误
func BadRequest(c *gin.Context, details string) {
	Error(c, http.StatusBadRequest, CodeInvalidCredentials, details)
}

// Unauthorized 返回 401 错误
func Unauthorized(c *gin.Context, details string) {
	Error(c, http.StatusUnauthorized, CodeTokenExpired, details)
}

// Forbidden 返回 403 错误
func Forbidden(c *gin.Context, details string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: "无权限访问",
		Details: details,
	})
}

// NotFound 返回 404 错误
func NotFound(c *gin.Context, details string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: "资源不存在",
		Details: details,
	})
}

// InternalError 返回 500 错误
func InternalError(c *gin.Context, details string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
		Details: details,
	})
}
