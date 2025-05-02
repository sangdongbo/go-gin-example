package mq

import "time"

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
