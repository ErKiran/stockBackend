package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	MobileNumber string `gorm:"size:255;null;unique" json:"mobileNumber,omitempty"`
	Email        string `gorm:"size:255;null;" json:"email,omitempty"`
	User         *User  `gorm:"foreignkey:UserID" json:"user"`
	UserID       int64  `gorm:"not null" json:"user_id"`
}

func (c *Contact) Prepare() {
	c.Email = html.EscapeString(strings.TrimSpace(c.Email))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Contact) Validate() error {
	if c.UserID != 0 {
		return errors.New("User Id isn't required")
	}
	if c.MobileNumber == "" && c.Email == "" {
		return errors.New("Either Mobile Numer or Email is required")
	}
	return nil
}

func (c *Contact) Create(db *gorm.DB) (*Contact, error) {
	err := db.Debug().Model(&Contact{}).Where("user_id=?", c.UserID).FirstOrCreate(&c).Error
	if err != nil {
		return &Contact{}, err
	}
	return c, nil
}

func (c *Contact) Find(db *gorm.DB, id int64) (*Contact, error) {
	err := db.Debug().Model(&Contact{}).Where("id=?", id).Find(&c).Error
	if err != nil {
		return &Contact{}, err
	}
	return c, nil
}

// func (c *Contact) FindContactByEmailOrNumber(db *gorm.DB, email, number string) (Contact, error) {
// 	contact := db.Debug().Model(&Contact{})

// 	if email != "" {
// 		contact = contact.Where("email=?", email)
// 	}

// 	if number != "" {
// 		contact = contact.Where("mobile_number=?", number)
// 	}

// 	// err := contact.Find(&).E
// }

func (c *Contact) Update(db *gorm.DB, id int64) (*Contact, error) {
	var contact = Contact{}

	if c.Email != "" {
		contact.Email = c.Email
	}

	if c.MobileNumber != "" {
		contact.MobileNumber = c.MobileNumber
	}

	err := db.Debug().Model(&Contact{}).Where("id=?", id).Updates(contact).Error

	if err != nil {
		return &Contact{}, err
	}
	return &contact, nil
}

func (c *Contact) FindContactByUserID(db *gorm.DB, userID int64) (*Contact, error) {
	err := db.Debug().Model(&Contact{}).Where("user_id=?", userID).Take(&c).Error
	if err != nil {
		return &Contact{}, err
	}
	return c, nil
}
