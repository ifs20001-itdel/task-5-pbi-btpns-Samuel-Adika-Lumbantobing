package models

type Photo struct {
	ID     int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID int64  `gorm:"not null" json:"-"`
	URL    string `gorm:"type:varchar(300);not null" json:"url"`
}
