package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID             int    `gorm:"primary_key"`
	Nickname       string `gorm:"type:text"`
	Email          string `gorm:"uniqueIndex"`
	Password       string `gorm:"-"`
	HashedPassword string `gorm:"type:varchar(100)"`
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
		if err != nil {
			return err
		}
		u.HashedPassword = string(hashedPassword)
	}

	return nil
}
