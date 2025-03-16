package models

import (
	"fmt"

	"github.com/google/uuid"
)

type ServiceResponse struct {
	Code    int
	Message string
	Data    interface{}
}

func (e *ServiceResponse) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Data: %+v", e.Code, e.Message, e.Data)
}

type QueueResponse struct {
	TaskId uuid.UUID
	Name   string
}
