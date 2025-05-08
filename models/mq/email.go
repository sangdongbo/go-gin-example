package mq

import (
	"errors"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/jinzhu/gorm"
	"time"
)

type MQEmail struct {
	ID        int        `gorm:"primaryKey" json:"id"`
	UserID    int        `gorm:"not null" json:"user_id"`
	Subject   string     `gorm:"type:varchar(100);not null" json:"subject"`
	Body      string     `gorm:"type:text;not null" json:"body"`
	Status    string     `gorm:"type:enum('pending','sent','failed');default:'pending'" json:"status"`
	SendTime  *time.Time `json:"send_time"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
}

func (MQEmail) TableName() string {
	return "blog_mq_emails"
}

// AddOrder 添加订单
func AddEmail(email *MQEmail) (int, error) {
	email.CreatedAt = time.Now()
	err := models.Db.Create(email).Error
	return email.ID, err
}

func EditEmail(email *MQEmail) error {
	email.CreatedAt = time.Now()
	err := models.Db.Save(email).Error
	return err
}

func GetEmailById(id int) (*MQEmail, error) {
	var email MQEmail
	err := models.Db.Where("id = ?", id).First(&email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 查不到返回 nil, nil，便于上层判断
	}
	if err != nil {
		return nil, err
	}
	return &email, nil
}
