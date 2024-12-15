package business

import "time"

type Policy struct {
	ID        int       `gorm:"primary_key"`
	Role      string    `gorm:"not null"`
	Resource  string    `gorm:"not null"`
	Action    string    `gorm:"not null"`
	Remark    string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Operator  string    `gorm:"type:varchar(255)"`
}
