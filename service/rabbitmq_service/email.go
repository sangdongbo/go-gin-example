package rabbitmq_service

import (
	"errors"
	"github.com/EDDYCJY/go-gin-example/models/mq"
)

type Email struct {
	ID      int
	UserID  int
	Subject string
	Body    string
	Status  string
}

type BaseEmailForm struct {
	UserID  int    `gorm:"column:user_id" json:"user_id"`
	Subject string `gorm:"column:subject" json:"subject"`
	Body    string `gorm:"column:body" json:"body"`
	Status  string `gorm:"column:status;default:pending" json:"status"` // enum: pending, sent, failed
}

type AddBaseEmailForm struct {
	BaseEmailForm
}

type UpdateBaseEmailForm struct {
	ID int `gorm:"column:id" json:"id"`
	BaseEmailForm
}

func ConvertAddFormToUEmail(form AddBaseEmailForm) Email {
	email := Email{
		UserID:  form.UserID,
		Subject: form.Subject,
		Body:    form.Body,
		Status:  form.Status,
	}
	return email
}

func ConvertEditFormToUEmail(form UpdateBaseEmailForm) Email {
	email := Email{
		ID:      form.ID,
		UserID:  form.UserID,
		Subject: form.Subject,
		Body:    form.Body,
		Status:  form.Status,
	}
	return email
}

func toModelEmail(o *Email) *mq.MQEmail {
	model := &mq.MQEmail{
		ID:      o.ID,
		UserID:  o.UserID,
		Subject: o.Subject,
		Body:    o.Body,
		Status:  o.Status,
	}
	return model
}

func (o *Email) Add() (int, error) {
	return mq.AddEmail(toModelEmail(o))
}

func (o *Email) Edit() error {
	id := o.ID
	email, _ := mq.GetEmailById(id)
	if email == nil {
		return errors.New("email not found")
	}
	return mq.EditEmail(toModelEmail(o))
}
