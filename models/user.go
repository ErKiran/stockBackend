package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email    string `gorm:"size:255;not null;unique" json:"email,omitempty"`
	UserName string `gorm:"size:255;null;" json:"username,omitempty"`
	Password string `gorm:"size:255;not null;" json:"password,omitempty"`
	IsAdmin  bool   `gorm:"not null;default:false" json:"is_admin,omitempty"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func isHashed(password string) bool {
	return len(password) == 60
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	if !isHashed(u.Password) {
		hashedPassword, err := Hash(u.Password)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

func (u *User) Prepare() {
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.IsAdmin = false
}

func (u *User) Validate() error {
	if u.UserName == "" {
		return errors.New("Required UserName")
	}
	if u.Password == "" {
		return errors.New("Required Password")
	}
	if u.Email == "" {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("Invalid Email")
	}
	return nil
}

func (u *User) Save(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (data *User) FindByEmail(db *gorm.DB, email string) (*User, error) {
	var err error
	err = db.Debug().Model(&User{}).Where("email=?", email).Take(data).Error
	if err != nil {
		return &User{}, err
	}
	return data, nil
}
