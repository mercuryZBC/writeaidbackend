package models

import (
	"time"
)

// 用户模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`           // 主键，自动增长
	Email        string    `json:"email" binding:"required,email" gorm:"unique"` // 唯一约束
	Nickname     string    `json:"nickname" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	RegisteredAt time.Time `json:"registered_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
	ExpiryAt     time.Time `json:"expiry_at"` // 会员到期时间
}
