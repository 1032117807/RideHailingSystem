package model

import "time"

const (
	RolePassenger = "passenger"
	RoleDriver    = "driver"
	RoleAdmin     = "admin"

	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
	UserStatusFrozen   = "frozen"
)

type User struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Phone            string     `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	PasswordHash     string     `gorm:"column:password_hash;size:255;not null" json:"-"`
	Nickname         string     `gorm:"size:50;not null" json:"nickname"`
	Role             string     `gorm:"size:20;not null;default:passenger" json:"role"`
	DefaultRole      string     `gorm:"column:default_role;size:20;not null;default:passenger" json:"defaultRole"`
	RealName         string     `gorm:"column:real_name;size:50" json:"realName"`
	IDCard           string     `gorm:"column:id_card;size:32" json:"idCard"`
	RealNameVerified bool       `gorm:"column:real_name_verified;not null;default:false" json:"realNameVerified"`
	Avatar           string     `gorm:"size:255" json:"avatar"`
	Email            string     `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Gender           string     `gorm:"size:20" json:"gender"`
	Birthday         *time.Time `json:"birthday,omitempty"`
	Status           string     `gorm:"size:20;not null;default:active" json:"status"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
