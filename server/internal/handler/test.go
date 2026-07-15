package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/response"

	"github.com/gin-gonic/gin"
)

// Ping 健康检查
// @Summary      健康检查
// @Description  检查服务器是否正常运行
// @Tags         系统
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response{data=object{time=string}}
// @Router       /ping [get]
func Ping(c *gin.Context) {
	response.Success(c, gin.H{"time": time.Now().Format(time.RFC3339)})
}

// TestSuccess 测试成功响应
// @Summary      测试成功响应
// @Description  测试返回成功的统一响应格式
// @Tags         测试
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response{data=object{username=string,role=string}}
// @Router       /test/success [get]
func TestSuccess(c *gin.Context) {
	response.Success(c, gin.H{
		"username": "testuser",
		"role":     "admin",
	})
}

// TestError 测试业务错误响应
// @Summary      测试业务错误响应
// @Description  根据错误码返回对应的业务错误响应
// @Tags         测试
// @Accept       json
// @Produce      json
// @Param        code   path      int  true  "业务错误码 (10001-10099: 认证相关, 20001-20099: 任务相关)"
// @Success      200    {object}  Response{details=string}
// @Failure      400    {object}  Response
// @Router       /test/error/{code} [get]
func TestError(c *gin.Context) {
	codeStr := c.Param("code")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		response.BadRequest(c, "错误码必须是数字")
		return
	}
	details := fmt.Sprintf("这是错误码 %d 的测试响应", code)
	response.Error(c, http.StatusOK, code, details)
}

// Test500 测试 500 错误
// @Summary      测试 500 错误
// @Description  模拟返回 HTTP 500 服务器内部错误
// @Tags         测试
// @Accept       json
// @Produce      json
// @Success      500  {object}  Response
// @Router       /test/500 [get]
func Test500(c *gin.Context) {
	response.InternalError(c, "模拟服务器内部错误")
}

// Test401 测试 401 错误
// @Summary      测试 401 错误
// @Description  模拟返回 HTTP 401 未授权错误
// @Tags         测试
// @Accept       json
// @Produce      json
// @Success      401  {object}  Response
// @Router       /test/401 [get]
func Test401(c *gin.Context) {
	response.Unauthorized(c, "模拟 Token 过期")
}

// Test403 测试 403 错误
// @Summary      测试 403 错误
// @Description  模拟返回 HTTP 403 禁止访问错误
// @Tags         测试
// @Accept       json
// @Produce      json
// @Success      403  {object}  Response
// @Router       /test/403 [get]
func Test403(c *gin.Context) {
	response.Forbidden(c, "模拟无权限")
}

// TestTimeout 测试超时响应
// @Summary      测试超时响应
// @Description  模拟请求超时（15秒延迟，超过前端 10s timeout）
// @Tags         测试
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response
// @Failure      408  {object}  Response
// @Router       /test/timeout [get]
func TestTimeout(c *gin.Context) {
	time.Sleep(15 * time.Second)
	response.Success(c, "这条消息不应该出现")
}

// TestPaginated 测试分页响应
// @Summary      测试分页响应
// @Description  测试返回分页数据的统一响应格式
// @Tags         测试
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response{data=PaginatedData{items=[]object}}
// @Router       /test/paginated [get]
func TestPaginated(c *gin.Context) {
	items := []gin.H{
		{"id": 1, "name": "item-1"},
		{"id": 2, "name": "item-2"},
		{"id": 3, "name": "item-3"},
	}
	response.SuccessWithPaginated(c, 100, items)
}
