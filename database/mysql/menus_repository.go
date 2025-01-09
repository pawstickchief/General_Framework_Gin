package mysql

import "General_Framework_Gin/schemas/business"

func GetMenusByRole(role string) ([]business.Menu, error) {
	var menus []business.Menu

	// 数据库查询：联表查询 role_menu 和 menus 表
	err := DB.Table("menus").
		Select("menus.id, menus.title, menus.icon, menus.pathname,menus.type,menus.position,menus.parent_id").
		Joins("JOIN role_menu ON menus.id = role_menu.menu_id").
		Joins("JOIN roles ON role_menu.role_id = roles.id").
		Where("roles.name = ?", role).
		Order("menus.position ASC"). // 根据菜单顺序排列
		Scan(&menus).Error

	return menus, err
}
