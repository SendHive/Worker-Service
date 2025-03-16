package job

import (
	"context"
	"log"
	"time"

	"github.com/SendHive/Infra-Common/queue"
	"github.com/SendHive/worker-service/models"
	pb "github.com/SendHive/worker-service/proto"
	"github.com/rabbitmq/amqp091-go"
)

type IJobService interface {
	ConsumeMessage() (*models.QueueResponse, error)
	GetJobStatus(jobId string) error
}

type JobService struct {
	client pb.TaskServiceClient
	Qu     amqp091.Queue
	Iq     queue.IQueueService
	QConn  *amqp091.Connection
}

func NewJobServiceRequest(c pb.TaskServiceClient, qu amqp091.Queue, iq queue.IQueueService, qConn *amqp091.Connection) (IJobService, error) {
	ser := &JobService{}
	ser.client = c
	ser.Qu = qu
	ser.Iq = iq
	ser.QConn = qConn
	return ser, nil
}

func (j *JobService) StartJob(req *models.QueueResponse) error {
	log.Println("from the startJob: ", req.TaskId)
	resp, err := j.client.StartJob(context.Background(), &pb.StartJobRequest{
		JobId:   req.TaskId.String(),
		JobName: req.Name,
	})
	if err != nil {
		return err
	}
	log.Println(resp)
	return nil
}

func (j *JobService) GetJobStatus(jobId string) error {
	stream, err := j.client.GetJobStatus(context.Background(), &pb.GetJobStatusRequest{
		JobId: jobId,
	})
	if err != nil {
		log.Println("error while listening to the streaming: ", err)
		return err
	}

	//Start stream listening
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Println("Stream closed.")
			break
		}
		log.Printf("Job Status: %s\n", resp.Status)

		// Stop when job is completed
		if resp.Status == "COMPLETED" {
			log.Println("Job processing completed.")
			break
		}

		time.Sleep(1 * time.Second)
	}

	return &models.ServiceResponse{
		Code:    200,
		Message: "Client Finished !",
	}
}
