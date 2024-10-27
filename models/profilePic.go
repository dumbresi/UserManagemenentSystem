package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Image struct {
    ID         string     `gorm:"primaryKey;type:uuid" json:"id"`
    FileName   string     `gorm:"type:varchar(255);not null" json:"file_name"`
    URL        string     `gorm:"type:varchar(255);not null" json:"url"`
    UploadDate time.Time  `gorm:"not null" json:"upload_date"`
	UserID     string     `gorm:"type:uuid;not null;index" json:"user_id"`
}

func (image *Image) BeforeCreate(tx *gorm.DB) (err error) {
	if(image.ID == ""){
		image.ID = uuid.New().String()
	}
    return
}