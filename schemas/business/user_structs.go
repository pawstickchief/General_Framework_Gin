package business

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"` // 主键，自增
	Username  string    `json:"username" gorm:"size:50;not null"`   // 用户名，最多 50 个字符
	Password  string    `json:"password" gorm:"size:255;not null"`  // 密码，隐藏在 JSON 输出中
	Email     string    `json:"email,omitempty" gorm:"size:100"`    // 邮箱，可为空
	Role      string    `json:"role" gorm:"size:50;default:user"`   // 用户角色，默认为 user
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`   // 创建时间
}

type Token struct {
	Token string `json:"token"`
}

type UserLoginLog struct {
	ID          int    `gorm:"primaryKey"`
	UserID      uint   `gorm:"index"`
	Username    string `gorm:"size:100"`
	Role        string `gorm:"size:50"`
	LoginTime   time.Time
	IPAddress   string `gorm:"size:100"`
	TokenExpire time.Time
	Remark      string `gorm:"size:255"`
}
