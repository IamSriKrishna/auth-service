package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type SupportStatus string

const (
	SupportStatusOpen       SupportStatus = "open"
	SupportStatusInProgress SupportStatus = "in_progress"
	SupportStatusResolved   SupportStatus = "resolved"
	SupportStatusClosed     SupportStatus = "closed"
)

const SupportTableName = "support"

type Support struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null"`
	Phone       *string        `json:"phone,omitempty" gorm:"type:varchar(15)"`
	Email       *string        `json:"email,omitempty" gorm:"type:varchar(100)"`
	IssueType   *string        `json:"issue_type,omitempty" gorm:"type:varchar(100)"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	Status      SupportStatus  `json:"status" gorm:"type:varchar(50);default:open"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Support) TableName() string {
	return SupportTableName
}

func (s *Support) BeforeCreate(tx *gorm.DB) (err error) {
	if (s.Phone == nil || *s.Phone == "") && (s.Email == nil || *s.Email == "") {
		return errors.New("either phone or email is required")
	}
	return
}
