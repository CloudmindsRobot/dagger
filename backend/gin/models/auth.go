package models

import "time"

type User struct {
	ID          int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	IsSuperuser bool      `gorm:"type:tinyint;column:is_superuser;not null;default:0;" json:"is_superuser"`
	IsActive    bool      `gorm:"type:tinyint;column:is_active;not null;default:1;" json:"is_active"`
	Username    string    `gorm:"type:varchar(128);column:username;unique_index;not null;" json:"username"`
	Password    string    `gorm:"type:varchar(256);column:password;not null;" json:"password"`
	Email       string    `gorm:"type:varchar(128);column:email;unique_index" json:"email"`
	CreateAt    time.Time `gorm:"type:datetime;column:create_at;not null;" json:"create_at"`
	LastLoginAt time.Time `gorm:"type:datetime;column:last_login_at;null;" json:"last_login_at"`
}

func (User) TableName() string {
	return "auth_user"
}
