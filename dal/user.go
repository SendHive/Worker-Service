package dal

import (
	"log"

	"github.com/SendHive/worker-service/external"
	"github.com/SendHive/worker-service/models"

	"github.com/google/uuid"
)

type IUser interface {
	Create(value *models.DBUserDetails) error
	FindBy(userId uuid.UUID) (*models.DBUserDetails, error)
	FindByConditions(conditions *models.DBUserDetails) (*models.DBUserDetails, error)
}

type User struct{}

func NewUserDalRequest() (IUser, error) {
	return &User{}, nil
}

func (u *User) Create(value *models.DBUserDetails) error {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	uerr := transaction.Create(&value)
	if uerr.Error != nil {
		return uerr.Error
	}
	transaction.Commit()
	return nil
}

func (u *User) FindByConditions(conditions *models.DBUserDetails) (*models.DBUserDetails, error) {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return nil, err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	resp := &models.DBUserDetails{}
	ferr := transaction.Find(&resp, &conditions)
	if ferr.Error != nil {
		log.Println("the error : ", ferr)
		return nil, ferr.Error
	}
	return resp, nil
}

func (u *User) FindBy(userId uuid.UUID) (*models.DBUserDetails, error) {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return nil, err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	resp := &models.DBUserDetails{}
	ferr := transaction.Find(&resp, &models.DBUserDetails{
		UserId: userId,
	})
	if ferr.Error != nil {
		log.Println("the error : ", ferr)
		return nil, ferr.Error
	}
	return resp, nil
}
