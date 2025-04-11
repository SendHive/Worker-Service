package dal

import (
	"log"

	"github.com/SendHive/worker-service/external"
	"github.com/SendHive/worker-service/models"

	"github.com/google/uuid"
)

type IFile interface {
	Create(value *models.DbFileDetails) error
	FindBy(conditions *models.DbFileDetails) (*models.DbFileDetails, error)
	FindAll(userId uuid.UUID) (response []*models.DbFileDetails, err error)
}

type File struct{}

func NewFileDalRequest() (IFile, error) {
	return &File{}, nil
}

func (f *File) Create(value *models.DbFileDetails) error {
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

func (f *File) FindBy(conditions *models.DbFileDetails) (*models.DbFileDetails, error) {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return nil, err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var resp *models.DbFileDetails
	ferr := transaction.Find(&resp, &conditions)
	if ferr.Error != nil {
		log.Println("the error while finding the job:", ferr.Error)
		return nil, ferr.Error
	}
	return resp, nil
}

func (f *File) FindAll(userId uuid.UUID) (response []*models.DbFileDetails, err error) {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return nil, err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	fileDetails := transaction.Find(&response, &models.DbFileDetails{
		UserId: userId,
	})
	if fileDetails.Error != nil {
		return nil, fileDetails.Error
	}
	return response, nil
}
