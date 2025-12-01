package models

import "time"

type File struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	Name      string
	Path      string
	Size      int64
	Type      string
	CreatedAt time.Time
}