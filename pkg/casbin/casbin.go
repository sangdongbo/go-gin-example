package casbin

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

// Setup 初始化 Casbin enforcer
func Setup() error {
	var err error
	once.Do(func() {
		// 使用文件适配器（将权限策略存储在文件中）
		policyPath := filepath.Join(setting.AppSetting.RuntimeRootPath, "../conf/rbac_policy.csv")
		adapter := fileadapter.NewAdapter(policyPath)

		// 加载模型配置文件
		modelPath := filepath.Join(setting.AppSetting.RuntimeRootPath, "../conf/rbac_model.conf")
		enforcer, err = casbin.NewEnforcer(modelPath, adapter)
		if err != nil {
			err = fmt.Errorf("failed to create casbin enforcer: %v", err)
			return
		}

		// 从文件加载策略
		if err = enforcer.LoadPolicy(); err != nil {
			err = fmt.Errorf("failed to load policy: %v", err)
			return
		}

		log.Println("Casbin enforcer initialized successfully")
	})

	return err
}

// GetEnforcer 获取 Casbin enforcer 实例
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// Enforce 检查权限
func Enforce(sub, obj, act string) (bool, error) {
	if enforcer == nil {
		return false, fmt.Errorf("casbin enforcer not initialized")
	}
	return enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加权限策略
// 示例: AddPolicy("admin", "/api/v1/articles", "POST")
func AddPolicy(role, path, method string) (bool, error) {
	if enforcer == nil {
		return false, fmt.Errorf("casbin enforcer not initialized")
	}
	return enforcer.AddPolicy(role, path, method)
}

// RemovePolicy 删除权限策略
func RemovePolicy(role, path, method string) (bool, error) {
	if enforcer == nil {
		return false, fmt.Errorf("casbin enforcer not initialized")
	}
	return enforcer.RemovePolicy(role, path, method)
}

// AddRoleForUser 为用户添加角色
// 示例: AddRoleForUser("user:1", "admin")
func AddRoleForUser(user, role string) (bool, error) {
	if enforcer == nil {
		return false, fmt.Errorf("casbin enforcer not initialized")
	}
	return enforcer.AddRoleForUser(user, role)
}

// DeleteRoleForUser 删除用户的角色
func DeleteRoleForUser(user, role string) (bool, error) {
	if enforcer == nil {
		return false, fmt.Errorf("casbin enforcer not initialized")
	}
	return enforcer.DeleteRoleForUser(user, role)
}

// GetRolesForUser 获取用户的所有角色
func GetRolesForUser(user string) ([]string, error) {
	if enforcer == nil {
		return nil, fmt.Errorf("casbin enforcer not initialized")
	}
	return enforcer.GetRolesForUser(user)
}

// GetUsersForRole 获取拥有某个角色的所有用户
func GetUsersForRole(role string) ([]string, error) {
	if enforcer == nil {
		return nil, fmt.Errorf("casbin enforcer not initialized")
	}
	return enforcer.GetUsersForRole(role)
}

// GetAllRoles 获取所有角色列表
func GetAllRoles() ([]string, error) {
	if enforcer == nil {
		return nil, fmt.Errorf("casbin enforcer not initialized")
	}

	allRoles := make(map[string]bool)

	// 1. 从权限策略(p)中获取角色（第一列是主体，可能是角色）
	allSubjects, err := enforcer.GetAllSubjects()
	if err != nil {
		return nil, fmt.Errorf("failed to get all subjects: %v", err)
	}
	for _, subject := range allSubjects {
		// 过滤掉 user: 开头的用户，保留角色
		if len(subject) < 5 || subject[:5] != "user:" {
			allRoles[subject] = true
		}
	}

	// 2. 从分组策略(g)中获取角色（第二列是角色）
	groupingPolicy, err := enforcer.GetGroupingPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to get grouping policy: %v", err)
	}
	for _, gp := range groupingPolicy {
		if len(gp) >= 2 {
			// gp[1] 是角色名
			allRoles[gp[1]] = true
		}
	}

	// 转换为切片并过滤掉占位符
	roles := make([]string, 0, len(allRoles))
	for role := range allRoles {
		// 过滤掉空字符串和占位符
		if role != "" && role != "__placeholder__" {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

// CreateRole 创建新角色
// 通过添加一个占位符权限来创建角色
func CreateRole(role string) (bool, error) {
	if enforcer == nil {
		return false, fmt.Errorf("casbin enforcer not initialized")
	}

	// 检查角色是否已存在
	roles, err := enforcer.GetAllRoles()
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		if r == role {
			return false, nil // 角色已存在
		}
	}

	// 创建角色（通过添加一个占位符权限）
	// 这样角色就会出现在角色列表中
	ok, err := enforcer.AddPolicy(role, "/__placeholder__", "NONE")
	if err != nil {
		return false, err
	}

	// 保存到文件
	if err := enforcer.SavePolicy(); err != nil {
		return false, fmt.Errorf("failed to save policy: %v", err)
	}

	return ok, nil
}

// DeleteRole 删除角色及其所有相关策略
func DeleteRole(role string) error {
	if enforcer == nil {
		return fmt.Errorf("casbin enforcer not initialized")
	}

	// 删除角色的所有权限策略
	_, err := enforcer.DeletePermissionsForUser(role)
	if err != nil {
		return fmt.Errorf("failed to delete permissions: %v", err)
	}

	// 删除所有用户与该角色的关联
	users, err := enforcer.GetUsersForRole(role)
	if err != nil {
		return fmt.Errorf("failed to get users for role: %v", err)
	}

	for _, user := range users {
		_, err := enforcer.DeleteRoleForUser(user, role)
		if err != nil {
			return fmt.Errorf("failed to delete role for user %s: %v", user, err)
		}
	}

	// 保存到文件
	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy: %v", err)
	}

	return nil
}

// GetPermissionsForRole 获取角色的所有权限
func GetPermissionsForRole(role string) ([][]string, error) {
	if enforcer == nil {
		return nil, fmt.Errorf("casbin enforcer not initialized")
	}
	return enforcer.GetPermissionsForUser(role)
}

// LoadPolicyFromFile 从 CSV 文件加载策略（用于初始化）
func LoadPolicyFromFile(policyPath string) error {
	if enforcer == nil {
		return fmt.Errorf("casbin enforcer not initialized")
	}

	// 读取 CSV 文件并添加策略
	adapter := enforcer.GetAdapter()
	if err := adapter.LoadPolicy(enforcer.GetModel()); err != nil {
		return err
	}

	return enforcer.SavePolicy()
}
