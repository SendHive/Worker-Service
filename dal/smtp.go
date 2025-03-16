package dal

import (
	"github.com/SendHive/worker-service/external"
	"github.com/SendHive/worker-service/models"

	"errors"

	"github.com/google/uuid"
)

type SmtpDal struct{}

type ISmtpDal interface {
	Create(value *models.DBSMTPDetails) error
	GetAll(userId uuid.UUID) (response []*models.DBSMTPDetails, err error)
	Update(id uuid.UUID, value *models.DBSMTPDetails) error
	FindBy(conditions *models.DBSMTPDetails) (*models.DBSMTPDetails, error)
}

func NewSmtpDalRequest() (ISmtpDal, error) {
	return &SmtpDal{}, nil
}

func (smtp *SmtpDal) Create(value *models.DBSMTPDetails) error {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	customerAddition := transaction.Create(&value)
	if customerAddition.Error != nil {
		return customerAddition.Error
	}
	transaction.Commit()
	return nil
}

func (smtp *SmtpDal) GetAll(userId uuid.UUID) (response []*models.DBSMTPDetails, err error) {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return nil, err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()

	customerDetails := transaction.Find(&response, &models.DBSMTPDetails{
		UserId: userId,
	})
	if customerDetails.Error != nil {
		return nil, customerDetails.Error
	}
	return response, nil
}

func (smtp *SmtpDal) Update(id uuid.UUID, value *models.DBSMTPDetails) error {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	uerr := dbConn.Model(&models.DBSMTPDetails{}).Where("id = ?", id).Updates(&value)
	if uerr.Error != nil {
		return uerr.Error
	}
	transaction.Commit()
	return nil
}

func (smtp *SmtpDal) FindBy(conditions *models.DBSMTPDetails) (*models.DBSMTPDetails, error) {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return nil, err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var resp *models.DBSMTPDetails
	smtpDetails := transaction.Find(&resp, &conditions)
	if smtpDetails.Error != nil {
		return nil, errors.New("error while finding the smtp entry: " + smtpDetails.Error.Error())
	}
	return resp, nil
}
