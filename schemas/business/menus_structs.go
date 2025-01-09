package business

type Menu struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`                     // 主键 ID
	Title    string `json:"title" gorm:"type:varchar(100)"`                         // 菜单标题
	Icon     string `json:"icon" gorm:"type:varchar(50)"`                           // 菜单图标
	Pathname string `json:"pathname" gorm:"type:varchar(255)"`                      // 路径名称
	Type     string `json:"type" gorm:"type:enum('item','divider');default:'item'"` // 菜单类型
	Position int    `json:"position" gorm:"type:int;default:0"`                     // 菜单顺序
	ParentID *int   `json:"parent_id" gorm:"type:int"`                              // 父菜单 ID（允许为空）
	Children []Menu `json:"children" gorm:"-"`                                      // 子菜单，用于递归嵌套
}
