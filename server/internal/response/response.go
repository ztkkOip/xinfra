package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Details string      `json:"details,omitempty"`
}

type PaginatedData struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

const (
	CodeSuccess = 0

	CodeLDAPAuthFailed     = 10001
	CodeInvalidCredentials = 10002
	CodeTokenExpired       = 10003
	CodeSSONotConfigured   = 10004
	CodeSSOGenerateFailed  = 10005

	CodeTaskCreateFailed = 20001
	CodeTaskExecTimeout  = 20002
)

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

func Message(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: Message(CodeSuccess),
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPaginated(c *gin.Context, total int64, items interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: Message(CodeSuccess),
		Data: PaginatedData{
			Total: total,
			Items: items,
		},
	})
}

func Error(c *gin.Context, httpCode, bizCode int, details string) {
	c.JSON(httpCode, Response{
		Code:    bizCode,
		Message: Message(bizCode),
		Details: details,
	})
}

func ErrorWithMessage(c *gin.Context, httpCode, bizCode int, message, details string) {
	c.JSON(httpCode, Response{
		Code:    bizCode,
		Message: message,
		Details: details,
	})
}

func BadRequest(c *gin.Context, details string) {
	Error(c, http.StatusBadRequest, CodeInvalidCredentials, details)
}

func Unauthorized(c *gin.Context, details string) {
	Error(c, http.StatusUnauthorized, CodeTokenExpired, details)
}

func Forbidden(c *gin.Context, details string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: "无权限访问",
		Details: details,
	})
}

func NotFound(c *gin.Context, details string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: "资源不存在",
		Details: details,
	})
}

func InternalError(c *gin.Context, details string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
		Details: details,
	})
}
