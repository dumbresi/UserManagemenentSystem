package models

import (
	"log"
	"time"
	"gorm.io/gorm"
)

type User struct{
	ID string					`gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email string				`gorm:"unique" 				json:"email"`
	Password string				`json:"-"`
	FirstName string			`gorm:"size:255" 			json:"first_name"`
	LastName string				`gorm:"size:255" 			json:"last_name"`
	AccountCreated time.Time   	`gorm:"autoCreateTime"		json:"account_created"`
	AccountUpdated time.Time  	`gorm:"autoUpdateTime"		json:"updated_at"` 
}

func MigrateUser(db *gorm.DB) error{
    err:=db.AutoMigrate(&User{})
	if(err!=nil){
		log.Fatalln(err)
		return err
	}
	return nil
}