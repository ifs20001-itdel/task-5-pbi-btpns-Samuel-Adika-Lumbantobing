package models

import (
	"time"
)

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(300);not null" json:"user_name"`
	Email     string    `gorm:"type:varchar(300);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(300);not null;min:6" json:"password"`
	Photo     Photo     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
