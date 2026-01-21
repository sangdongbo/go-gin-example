package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	casbinPkg "github.com/EDDYCJY/go-gin-example/pkg/casbin"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

// ===== Casbin 权限管理接口 =====

// @Summary 为用户添加角色
// @Tags 权限管理
// @Produce json
// @Param body body object true "用户和角色"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/add-role [post]
func AddRoleForUser(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	subject := getUserSubject(req.UserID)
	ok, err := casbinPkg.AddRoleForUser(subject, req.Role)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if !ok {
		appG.Response(http.StatusOK, e.SUCCESS, "角色已存在")
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "添加角色成功")
}

// @Summary 删除用户角色
// @Tags 权限管理
// @Produce json
// @Param body body object true "用户和角色"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/delete-role [delete]
func DeleteRoleForUser(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	subject := getUserSubject(req.UserID)
	ok, err := casbinPkg.DeleteRoleForUser(subject, req.Role)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if !ok {
		appG.Response(http.StatusOK, e.SUCCESS, "角色不存在")
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "删除角色成功")
}

// @Summary 获取用户的所有角色
// @Tags 权限管理
// @Produce json
// @Param user_id query int true "用户ID"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/user-roles [get]
func GetRolesForUser(c *gin.Context) {
	appG := app.Gin{C: c}

	userID := c.Query("user_id")
	if userID == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	subject := "user:" + userID
	roles, err := casbinPkg.GetRolesForUser(subject)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, gin.H{
		"user_id": userID,
		"roles":   roles,
	})
}

// @Summary 添加权限策略
// @Tags 权限管理
// @Produce json
// @Param body body object true "权限策略"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/add-policy [post]
func AddPolicy(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		Role   string `json:"role" binding:"required"`
		Path   string `json:"path" binding:"required"`
		Method string `json:"method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	ok, err := casbinPkg.AddPolicy(req.Role, req.Path, req.Method)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if !ok {
		appG.Response(http.StatusOK, e.SUCCESS, "策略已存在")
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "添加策略成功")
}

// @Summary 删除权限策略
// @Tags 权限管理
// @Produce json
// @Param body body object true "权限策略"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/delete-policy [delete]
func DeletePolicy(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		Role   string `json:"role" binding:"required"`
		Path   string `json:"path" binding:"required"`
		Method string `json:"method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	ok, err := casbinPkg.RemovePolicy(req.Role, req.Path, req.Method)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if !ok {
		appG.Response(http.StatusOK, e.SUCCESS, "策略不存在")
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "删除策略成功")
}

// @Summary 检查权限
// @Tags 权限管理
// @Produce json
// @Param user_id query int true "用户ID"
// @Param path query string true "路径"
// @Param method query string true "方法"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/check-permission [get]
func CheckPermission(c *gin.Context) {
	appG := app.Gin{C: c}

	userID := c.Query("user_id")
	path := c.Query("path")
	method := c.Query("method")

	if userID == "" || path == "" || method == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	subject := "user:" + userID
	ok, err := casbinPkg.Enforce(subject, path, method)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, gin.H{
		"user_id":    userID,
		"path":       path,
		"method":     method,
		"has_access": ok,
	})
}

// @Summary 创建新角色
// @Tags 权限管理
// @Produce json
// @Param body body object true "角色信息"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/create-role [post]
func CreateRole(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		Role        string `json:"role" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 创建角色（通过添加一个基础权限来创建）
	// 这里添加一个空路径的权限，表示角色已创建但暂无实际权限
	ok, err := casbinPkg.CreateRole(req.Role)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if !ok {
		appG.Response(http.StatusOK, e.SUCCESS, "角色已存在")
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, gin.H{
		"role":        req.Role,
		"description": req.Description,
		"message":     "角色创建成功",
	})
}

// @Summary 删除角色
// @Tags 权限管理
// @Produce json
// @Param role query string true "角色名称"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/delete-role [delete]
func DeleteRole(c *gin.Context) {
	appG := app.Gin{C: c}

	role := c.Query("role")
	if role == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	err := casbinPkg.DeleteRole(role)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "角色删除成功")
}

// @Summary 获取所有角色列表
// @Tags 权限管理
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/roles [get]
func GetAllRoles(c *gin.Context) {
	appG := app.Gin{C: c}

	roles, err := casbinPkg.GetAllRoles()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, gin.H{
		"roles": roles,
		"count": len(roles),
	})
}

// @Summary 获取角色的所有权限
// @Tags 权限管理
// @Produce json
// @Param role query string true "角色名称"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/role-permissions [get]
func GetPermissionsForRole(c *gin.Context) {
	appG := app.Gin{C: c}

	role := c.Query("role")
	if role == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	permissions, err := casbinPkg.GetPermissionsForRole(role)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, gin.H{
		"role":        role,
		"permissions": permissions,
		"count":       len(permissions),
	})
}

// @Summary 获取拥有某个角色的所有用户
// @Tags 权限管理
// @Produce json
// @Param role query string true "角色名称"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/casbin/role-users [get]
func GetUsersForRole(c *gin.Context) {
	appG := app.Gin{C: c}

	role := c.Query("role")
	if role == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	users, err := casbinPkg.GetUsersForRole(role)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, gin.H{
		"role":  role,
		"users": users,
		"count": len(users),
	})
}

// getUserSubject 构建用户标识
func getUserSubject(userID int) string {
	return "user:" + string(rune(userID+'0'))
}
