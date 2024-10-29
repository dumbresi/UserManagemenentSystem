package models

import (
	"fmt"
	"log"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             string    `gorm:"primaryKey;type:uuid" json:"id"`
	Email          string    `gorm:"unique" json:"email"`
	Password       string    `json:"password"`
	FirstName      string    `gorm:"size:255" json:"first_name"`
	LastName       string    `gorm:"size:255" json:"last_name"`
	AccountCreated time.Time `gorm:"autoCreateTime" json:"account_created"`
	AccountUpdated time.Time `gorm:"autoUpdateTime" json:"account_updated"`
}

type UserResponse struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email          string    `gorm:"unique" json:"email"`
	FirstName      string    `gorm:"size:255" json:"first_name"`
	LastName       string    `gorm:"size:255" json:"last_name"`
	AccountCreated time.Time `gorm:"autoCreateTime" json:"account_created"`
	AccountUpdated time.Time `gorm:"autoUpdateTime" json:"account_updated"`
}

func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (user *User) String() string {
	return fmt.Sprintf("FirstName:%s LastName:%s Email: %s", user.FirstName, user.LastName, user.Email)
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New().String()
    return
}