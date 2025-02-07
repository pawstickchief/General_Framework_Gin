package mysql

import (
	"General_Framework_Gin/schemas/business"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GetUsers(params business.PaginationParams) ([]business.User, int64, error) {
	var users []business.User
	var total int64

	// 计算偏移量
	offset := (params.Page - 1) * params.Limit

	// 创建查询对象
	query := DB.Model(&business.User{})

	// 处理搜索条件
	for key, value := range params.Filters {
		if value != "" {
			query = query.Where(fmt.Sprintf("%s LIKE ?", key), "%"+value+"%")
		}
	}

	// 查询符合条件的总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting users: %w", err)
	}

	// 查询当前页的数据
	err := query.Select("id, username, email, role, created_at").
		Offset(offset).
		Limit(params.Limit).
		Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching users: %w", err)
	}

	return users, total, nil
}
func GetUserByUsername(name *string) (user business.User, err error) {
	query := `SELECT id, username, password ,role,email,created_at FROM users WHERE username = ? LIMIT 1`

	// 使用 GORM 的 Raw 方法执行查询
	result := DB.Raw(query, *name).Scan(&user)
	if result.Error != nil {
		log.Printf("查询用户失败: %v", result.Error)
		return business.User{}, result.Error
	}

	return user, nil
}
func DeleteUserByID(id int) error {
	query := `DELETE FROM users WHERE id = ?`

	// 执行删除 SQL
	result := DB.Exec(query, id)
	if result.Error != nil {
		log.Printf("删除用户失败: %v", result.Error)
		return result.Error
	}

	// 检查是否有行被影响
	if result.RowsAffected == 0 {
		log.Printf("删除失败: 用户 ID %d 不存在", id)
		return fmt.Errorf("用户 ID %d 不存在", id)
	}

	return nil
}
func AddUser(username, password, email, role string) error {
	// 检查用户是否已存在
	var count int64
	DB.Model(&business.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return fmt.Errorf("用户名已存在")
	}

	// 哈希加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户实例
	user := business.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Role:     role,
	}

	// 存入数据库
	if err := DB.Create(&user).Error; err != nil {
		return fmt.Errorf("添加用户失败: %v", err)
	}

	return nil
}
func UpdateUserByUsernameAndEmail(username, email, newPassword string) (bool, error) {
	if username == "" || email == "" {
		return false, fmt.Errorf("用户名和邮箱不能为空")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		return false, fmt.Errorf("密码加密失败: %w", err)
	}

	// 检查用户是否存在
	var count int64
	err = DB.Model(&business.User{}).
		Where("username = ? AND email = ?", username, email).
		Count(&count).Error

	if err != nil {
		log.Printf("查询用户失败: %v", err)
		return false, fmt.Errorf("数据库查询失败: %w", err)
	}

	if count == 0 {
		log.Printf("未找到匹配的用户: %s", username)
		return false, fmt.Errorf("未找到匹配的用户")
	}

	// 执行更新操作
	result := DB.Model(&business.User{}).
		Where("username = ? AND email = ?", username, email).
		Update("password", string(hashedPassword))

	if result.Error != nil {
		log.Printf("更新用户密码失败: %v", result.Error)
		return false, fmt.Errorf("更新失败: %w", result.Error)
	}

	// 检查是否有行被更新
	if result.RowsAffected == 0 {
		log.Printf("未更新任何记录: %s", username)
		return false, fmt.Errorf("未更新任何记录")
	}

	log.Printf("成功更新用户密码: %s", username)
	return true, nil
}
