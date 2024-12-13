package mysql

import (
	"General_Framework_Gin/schemas/business"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GetUserByUsername(name *string) (user business.User, err error) {
	query := `SELECT id, username, password, role FROM users WHERE username = ? LIMIT 1`

	// 使用 GORM 的 Raw 方法执行查询
	result := DB.Raw(query, *name).Scan(&user)
	if result.Error != nil {
		log.Printf("查询用户失败: %v", result.Error)
		return business.User{}, result.Error
	}

	return user, nil
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
