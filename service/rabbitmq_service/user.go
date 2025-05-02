package rabbitmq_service

import (
	"github.com/EDDYCJY/go-gin-example/models/mq"
)

type User struct {
	Username string
	Email    string
}

type BaseRabbitMQUserForm struct {
	Username string `json:"username" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
}

type AddRabbitMQUserForm struct {
	BaseRabbitMQUserForm
}

type EditRabbitMQUserForm struct {
	Id int `json:"id" binding:"required"`
	BaseRabbitMQUserForm
}

func ConvertAddFormToUser(form AddRabbitMQUserForm) User {
	user := User{
		Username: form.Username,
		Email:    form.Email,
	}
	return user
}

func toModelMQUser(o *User) *mq.MQUser {
	model := &mq.MQUser{
		Username: o.Username,
		Email:    o.Email,
	}
	return model
}
func (o *User) Add() (int, error) {
	return mq.AddUser(toModelMQUser(o))
}
