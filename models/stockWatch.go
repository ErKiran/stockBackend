package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type StockWatch struct {
	gorm.Model
	Script                  string `gorm:"not null; size:255" json:"script"`
	DesiredPrice            int64  `gorm:"null;" json:"desiredPrice"`
	DesiredPercentageChange int64  `gorm:"null;" json:"desiredPercentageChange"`
	Change                  string `gorm:"null;" json:"change"`
	User                    *User  `gorm:"foreignkey:UserID" json:"user"`
	UserID                  int64  `gorm:"not null" json:"userId"`
}

func (s *StockWatch) Prepare() {
	s.Script = html.EscapeString(strings.TrimSpace(s.Script))
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *StockWatch) Validate() error {
	if s.Script == "" {
		return errors.New("Script is Required")
	}
	if s.UserID != 0 {
		return errors.New("Shouldn't include UserID")
	}
	if s.DesiredPrice == 0 && s.DesiredPercentageChange == 0 {
		return errors.New("Either Desired Price or Desired Percentage Change is Required")
	}
	return nil
}

func (s *StockWatch) Save(db *gorm.DB) (*StockWatch, error) {
	var err error
	err = db.Debug().Create(&s).Error
	if err != nil {
		return &StockWatch{}, err
	}
	return s, nil
}

func (s *StockWatch) GetStockWatchByUserID(db *gorm.DB, userId int64) ([]*StockWatch, error) {
	var stocks []*StockWatch
	err := db.Debug().Model(&StockWatch{}).Where("user_id=?", userId).Find(&stocks).Error
	if err != nil {
		return []*StockWatch{}, err
	}
	return stocks, nil
}
