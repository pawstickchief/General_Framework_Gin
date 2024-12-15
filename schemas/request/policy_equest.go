package request

type PolicyOption struct {
	Option      string `json:"option" binding:"required"` // 操作类型，如"add", "remove", "check"
	Role        string `json:"role" binding:"required"`   // 角色
	Resource    string `json:"resource" `                 // 资源
	Action      string `json:"action" `                   // 动作
	NewRole     string `json:"new_role"`                  // 新角色（可选）
	NewResource string `json:"new_resource"`              // 新资源（可选）
	NewAction   string `json:"new_action"`
	Operator    string `json:"operator"`
}
type RoleInheritanceOption struct {
	UserRole      string `json:"user_role"`
	InheritedRole string `json:"inherited_role"`
}
