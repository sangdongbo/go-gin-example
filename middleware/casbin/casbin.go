package casbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	casbinPkg "github.com/EDDYCJY/go-gin-example/pkg/casbin"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

// Casbin 权限校验中间件
func Casbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户标识（从 JWT token 中获取用户ID）
		userID, err := getUserIDFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_AUTH,
				"msg":  "未授权：" + err.Error(),
				"data": nil,
			})
			c.Abort()
			return
		}

		// 2. 构建 subject（主体）
		subject := fmt.Sprintf("user:%d", userID)

		// 3. 获取请求的资源和方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 4. 使用 Casbin 检查权限
		ok, err := casbinPkg.Enforce(subject, obj, act)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": e.ERROR,
				"msg":  "权限检查失败：" + err.Error(),
				"data": nil,
			})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code": e.ERROR_AUTH,
				"msg":  "无权限访问此资源",
				"data": map[string]interface{}{
					"user":   subject,
					"path":   obj,
					"method": act,
				},
			})
			c.Abort()
			return
		}

		// 权限验证通过，继续处理请求
		c.Next()
	}
}

// getUserIDFromToken 从 JWT token 中获取用户ID
func getUserIDFromToken(c *gin.Context) (int, error) {
	// 优先从 Authorization Header 获取 token
	token := c.GetHeader("Authorization")
	if token != "" {
		// 检查是否是 Bearer 格式
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:] // 去掉 "Bearer " 前缀
		}
	} else {
		// 如果 Header 中没有，则从 query 参数获取（向后兼容）
		token = c.Query("token")
	}

	if token == "" {
		return 0, fmt.Errorf("token 为空")
	}

	// 解析 token
	claims, err := util.ParseToken(token)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			return 0, fmt.Errorf("token 已过期")
		default:
			return 0, fmt.Errorf("token 无效")
		}
	}

	// 从 claims 中获取用户ID
	// 注意：这里假设你的 JWT claims 中有 user_id 字段
	// 你可能需要根据实际情况调整
	if userID, ok := claims["user_id"].(float64); ok {
		return int(userID), nil
	}

	// 如果没有 user_id，尝试使用 name 作为用户标识
	if name, ok := claims["name"].(string); ok {
		// 这里简单处理，实际应该从数据库查询用户ID
		// 或者你的 JWT token 中应该包含用户ID
		return parseUserIDFromName(name), nil
	}

	return 0, fmt.Errorf("无法从 token 中获取用户标识")
}

// parseUserIDFromName 从用户名解析用户ID（示例实现）
func parseUserIDFromName(name string) int {
	// 这只是一个示例，实际应该从数据库查询
	// 或者在生成 JWT 时就包含用户ID
	if strings.Contains(name, "admin") {
		return 1
	}
	return 2
}

// CasbinWithRoles 带角色检查的中间件（可选）
func CasbinWithRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserIDFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_AUTH,
				"msg":  "未授权：" + err.Error(),
				"data": nil,
			})
			c.Abort()
			return
		}

		subject := fmt.Sprintf("user:%d", userID)

		// 获取用户的所有角色
		roles, err := casbinPkg.GetRolesForUser(subject)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": e.ERROR,
				"msg":  "获取用户角色失败：" + err.Error(),
				"data": nil,
			})
			c.Abort()
			return
		}

		// 检查用户是否拥有所需角色
		hasRole := false
		for _, userRole := range roles {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"code": e.ERROR_AUTH,
				"msg":  fmt.Sprintf("需要以下角色之一: %v", requiredRoles),
				"data": map[string]interface{}{
					"user_roles":     roles,
					"required_roles": requiredRoles,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
