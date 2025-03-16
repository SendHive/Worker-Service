package dal

import (
	"log"

	"github.com/SendHive/worker-service/external"
	"github.com/SendHive/worker-service/models"
	"github.com/google/uuid"
)

type IJob interface {
	Create(value *models.DBJobDetails) error
	FindBy(conditions *models.DBJobDetails) (*models.DBJobDetails, error)
	UpdateStatus(id uuid.UUID) error
}

type Job struct{}

func NewJobDalRequest() (IJob, error) {
	return &Job{}, nil
}

func (j *Job) Create(value *models.DBJobDetails) error {
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

func (j *Job) FindBy(conditions *models.DBJobDetails) (*models.DBJobDetails, error) {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return nil, err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var resp *models.DBJobDetails
	ferr := transaction.Find(&resp, &conditions)
	if ferr.Error != nil {
		log.Println("the error while finding the job:", ferr.Error)
		return nil, ferr.Error
	}
	return resp, nil
}

func (j *Job) UpdateStatus(id uuid.UUID) error {
	dbConn, err := external.GetDbConn()
	if err != nil {
		return err
	}
	transaction := dbConn.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	update := transaction.Model(models.DBJobDetails{}).Where("task_id = ?", id).Update("status", "IN_PROGRESS")
	if update.Error != nil {
		log.Println("error whhile ")
		return err
	}
	transaction.Commit()
	return nil
}
