package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	DOB       time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// IsUsernameExists checks if a username already exists in the database
func (u *User) IsUsernameExists(db *gorm.DB) bool {
	var exists User
	result := db.Where("username = ?", u.Username).First(&exists)
	return result.Error == nil
}

// IsEmailExists checks if an email already exists in the database
func (u *User) IsEmailExists(db *gorm.DB) bool {
	var exists User
	result := db.Where("email = ?", u.Email).First(&exists)
	return result.Error == nil
}

// Validate checks if all required fields are present and valid
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	if u.LastName == "" {
		return errors.New("last name is required")
	}
	if u.DOB.IsZero() {
		return errors.New("date of birth is required")
	}
	return nil
}
