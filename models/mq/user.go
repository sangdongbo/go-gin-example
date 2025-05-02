package mq

import (
	"github.com/EDDYCJY/go-gin-example/models"
)

type MQUser struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(50);not null" json:"username"`
	Email    string `gorm:"type:varchar(100);not null" json:"email"`
}

func (MQUser) TableName() string {
	return "blog_mq_users"
}

// AddOrder 添加订单
func AddUser(user *MQUser) (int, error) {
	error := models.Db.Create(user).Error
	return user.ID, error
}
