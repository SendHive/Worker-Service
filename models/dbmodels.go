package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DBSMTPDetails struct {
	Id        uuid.UUID `gorm:"primaryKey,column:id"`
	UserId    uuid.UUID `gorm:"column:user_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not_null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not_null"`
	Server    string    `gorm:"column:server;type:varchar(100);not null"`
	Port      string    `gorm:"column:port;type:varchar(100);not null"`
	Username  string    `gorm:"column:username;type:varchar(100);not null"`
	Password  string    `gorm:"column:password;type:varchar(100);not null"`
}

func (DBSMTPDetails) TableName() string {
	return "smtp_tbl"
}

func (*DBSMTPDetails) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBJobDetails struct {
	Id         uuid.UUID `gorm:"primaryKey,column:id"`
	Name       string    `gorm:"column:name;not null"`
	UserId     uuid.UUID `gorm:"column:user_id;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not_null"`
	TaskId     uuid.UUID `gorm:"column:task_id;not null"`
	Type       string    `gorm:"column:type;not null"`
	ObjectName string    `gorm:"column:object_name;not null"`
	Status     string    `gorm:"column:status;not null"`
}

func (DBJobDetails) TableName() string {
	return "job_tbl"
}

func (*DBJobDetails) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBUserDetails struct {
	Id        uuid.UUID `gorm:"primaryKey,column:id"`
	UserId    uuid.UUID `gorm:"column:user_id;not null"`
	Name      string    `gorm:"column:name;not null"`
	SecretKey string    `gorm:"column:secret_key"`
}

func (DBUserDetails) TableName() string {
	return "user_tbl"
}

func (*DBUserDetails) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBSecretsDetails struct {
	Id        uuid.UUID `gorm:"primaryKey,column:id"`
	UserId    uuid.UUID `gorm:"column:user_id;not null"`
	SecretKey string    `gorm:"column:secret_key;not null"`
}

func (DBSecretsDetails) TableName() string {
	return "secret_tbl"
}

func (*DBSecretsDetails) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DbFileDetails struct {
	Id     uuid.UUID `gorm:"primaryKey,column:id"`
	Name   string    `gorm:"column:name;not null"`
	UserId uuid.UUID `gorm:"column:user_id;not null"`
}

func (DbFileDetails) TableName() string {
	return "file_tbl"
}

func (*DbFileDetails) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}