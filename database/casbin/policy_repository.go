package casbin

import (
	"General_Framework_Gin/config"
	"General_Framework_Gin/database/mysql"
	"General_Framework_Gin/schemas/request"
	"fmt"
	"github.com/casbin/casbin/v2"
)

func AddPolicy(enforcer *casbin.Enforcer, policy request.PolicyOption) error {
	// 使用 Casbin 添加策略
	_, err := enforcer.AddPolicy(policy.Role, policy.Resource, policy.Action)
	if err != nil {
		return fmt.Errorf("添加策略失败: %v", err)
	}
	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %v", err)
	}
	err = mysql.UpdatePoliciesFromFile(mysql.DB, config.AppConfig.Casbin.PolicyFile, policy.Operator)
	if err != nil {
		return err
	}
	return err
}

// 删除策略
func RemovePolicy(enforcer *casbin.Enforcer, policy request.PolicyOption) error {

	// 使用 Casbin 删除策略
	_, err := enforcer.RemovePolicy(policy.Role, policy.Resource, policy.Action)
	if err != nil {
		return fmt.Errorf("删除策略失败: %v", err)
	}
	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %v", err)
	}
	err = mysql.UpdatePoliciesFromFile(mysql.DB, config.AppConfig.Casbin.PolicyFile, policy.Operator)
	if err != nil {
		return err
	}
	return err
}

// 查询角色的可访问资源
func GetRolePolicy(enforcer *casbin.Enforcer, policy request.PolicyOption) error {
	// 获取角色所有的权限策略（可以查询角色的所有资源和动作）
	policies, _ := enforcer.GetFilteredPolicy(0, policy.Role) // 0 表示过滤角色

	if len(policies) == 0 {
		// 如果没有策略，返回未找到数据
		return fmt.Errorf("没有找到该角色的权限策略")
	}

	// 返回角色的所有可访问资源和动作
	var rolePolicies []map[string]string
	for _, p := range policies {
		rolePolicies = append(rolePolicies, map[string]string{
			"role":     p[0], // 角色
			"resource": p[1], // 资源
			"action":   p[2], // 动作
		})
	}
	err := mysql.UpdatePoliciesFromFile(mysql.DB, config.AppConfig.Casbin.PolicyFile, policy.Operator)
	if err != nil {
		return err
	}
	return err
}

// 编辑策略
func EditPolicy(enforcer *casbin.Enforcer, policy request.PolicyOption) error {

	// 获取请求中的旧策略信息和新策略信息
	oldRole := policy.Role
	oldResource := policy.Resource // 假设我们用资源来标识旧策略
	oldAction := policy.Action

	// 检查旧的策略是否存在
	allowed, err := enforcer.Enforce(oldRole, oldResource, oldAction)
	if err != nil {
		return fmt.Errorf("策略检查失败: %v", err)
	}
	if !allowed {
		return fmt.Errorf("旧策略不存在，无法编辑")
	}

	// 删除旧的策略
	if _, err := enforcer.RemovePolicy(oldRole, oldResource, oldAction); err != nil {
		return fmt.Errorf("删除旧策略失败: %v", err)
	}

	// 新的策略信息
	newRole := policy.NewRole         // 新的角色（如果有更改）
	newResource := policy.NewResource // 新的资源（如果有更改）
	newAction := policy.NewAction     // 新的动作（如果有更改）

	// 添加新的策略
	if _, err := enforcer.AddPolicy(newRole, newResource, newAction); err != nil {
		return fmt.Errorf("添加新策略失败: %v", err)
	}

	// 保存策略更新
	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %v", err)
	}
	err = mysql.UpdatePoliciesFromFile(mysql.DB, config.AppConfig.Casbin.PolicyFile, policy.Operator)
	if err != nil {
		return err
	}
	return err
}
