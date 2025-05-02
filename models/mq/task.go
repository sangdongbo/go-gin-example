package mq

import "time"

type MQTask struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	TaskType   string    `gorm:"type:varchar(50);not null" json:"task_type"`
	Payload    string    `gorm:"type:text;not null" json:"payload"`
	Status     string    `gorm:"type:enum('pending','processing','success','failed');default:'pending'" json:"status"`
	RetryCount int       `gorm:"default:0" json:"retry_count"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (MQTask) TableName() string {
	return "blog_mq_tasks"
}
