package business

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type Token struct {
	Token string `json:"token"`
}
type UserUpdate struct {
	ID        int    `gorm:"primaryKey;column:id"`
	Username  string `gorm:"column:username"`
	Password  string `gorm:"column:password"`
	Email     string `gorm:"column:email"`
	CreatedAt string `gorm:"column:created_at"`
}
