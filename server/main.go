package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/SendHive/worker-service/dal"
	"github.com/SendHive/worker-service/external"
	"github.com/SendHive/worker-service/job"
	"github.com/SendHive/worker-service/models"
	pb "github.com/SendHive/worker-service/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedTaskServiceServer
	JobRepo  dal.IJob
	SmtpRepo dal.ISmtpDal
	Job      job.IJobService
}

func IntilizeDataBase() *gorm.DB {
	dbConn, err := external.GetDbConn()
	if err != nil {
		log.Println("error while connecting to the database: ", err)
		return nil
	}
	fmt.Println(dbConn)
	return dbConn
}

func (s *Server) HealthCheck(ctx context.Context, in *pb.NoParams) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: "Connection Successfull",
	}, nil
}

func (s *Server) StartJob(ctx context.Context, in *pb.StartJobRequest) (*pb.StartJobResponse, error) {
	if in.JobId == "" || in.JobName == "" {
		return nil, &models.ServiceResponse{
			Code:    404,
			Message: "Name or id is messing in the grpc request",
			Data:    nil,
		}
	}

	fmt.Println(" the requestId: ", in.JobId)

	Jresp, err := s.JobRepo.FindBy(&models.DBJobDetails{
		TaskId: uuid.MustParse(in.JobId),
	})
	if err != nil {
		return nil, &models.ServiceResponse{
			Code:    500,
			Message: "error while finding the task service: " + err.Error(),
		}
	}

	log.Println("The details of job with this id: ", Jresp)

	//Updating the JobStatus
	err = s.JobRepo.UpdateStatus(Jresp.TaskId)
	if err != nil {
		return nil, &models.ServiceResponse{
			Code:    500,
			Message: "error while updating the task service: " + err.Error(),
		}
	}

	// SMTP Sending email
	resp, err := s.SmtpRepo.FindBy(&models.DBSMTPDetails{
		UserId: Jresp.UserId,
	})
	if err != nil {
		return nil, &models.ServiceResponse{
			Code:    500,
			Message: "error while updating the task service: " + err.Error(),
		}
	}
	log.Println("The smtp details: ", resp)
	return &pb.StartJobResponse{
		Status: "IN_PROGRESS",
	}, nil
}

func (s *Server) GetJobStatus(in *pb.GetJobStatusRequest, strem grpc.ServerStreamingServer[pb.GetJobStatusResponse]) error {
	if in.JobId == "" {
		return &models.ServiceResponse{
			Code:    404,
			Message: "Bad Request",
		}
	}

	jobId := uuid.MustParse(in.JobId)

	// pushing message to the stream
	for {
		resp, err := s.JobRepo.FindBy(&models.DBJobDetails{
			TaskId: jobId,
		})
		if err != nil {
			return &models.ServiceResponse{}
		}
		err = strem.Send(&pb.GetJobStatusResponse{
			Status: resp.Status,
		})
		if err != nil {
			log.Printf("Error sending stream: %v", err)
			return err
		}

		if resp.Status == "COMPLETED" {
			log.Printf("Job %s completed, stopping stream.", in.JobId)
			break
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

func main() {
	var ser Server
	// Connect to the Job Dal
	IJob, err := dal.NewJobDalRequest()
	if err != nil {
		log.Println("error while connecting to dal(job): ", err)
		return
	}

	// Connect to the SMTP Dal
	ISmtp, err := dal.NewSmtpDalRequest()
	if err != nil {
		log.Println("error while connecting to dal(smtp): ", err)
		return
	}

	conn, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal("error while starting the GRPC server: ", err)
	}
	s := grpc.NewServer()
	ser.JobRepo = IJob
	ser.SmtpRepo = ISmtp
	pb.RegisterTaskServiceServer(s, &ser)
	log.Println("Server running at port 3000")
	s.Serve(conn)
}
