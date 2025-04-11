package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	minioDb "github.com/SendHive/Infra-Common/minio"
	"github.com/SendHive/worker-service/dal"
	"github.com/SendHive/worker-service/external"
	"github.com/SendHive/worker-service/job"
	"github.com/SendHive/worker-service/models"
	pb "github.com/SendHive/worker-service/proto"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
)

// Server dependencies
type ServerDependencies struct {
	JobRepo      dal.IJob
	SmtpRepo     dal.ISmtpDal
	FileRepo     dal.IFile
	UserRepo     dal.IUser
	MinioClient  *minio.Client
	MinioService minioDb.IMinioService
}

// GRPC Server implementation
type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	deps ServerDependencies
}

func NewTaskServer(deps ServerDependencies) *TaskServer {
	return &TaskServer{deps: deps}
}

// HealthCheck implements a simple health check endpoint
func (s *TaskServer) HealthCheck(ctx context.Context, _ *pb.NoParams) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: "Connection Successful"}, nil
}

// StartJob handles job initiation
func (s *TaskServer) StartJob(ctx context.Context, req *pb.StartJobRequest) (*pb.StartJobResponse, error) {
	if err := validateStartJobRequest(req); err != nil {
		return nil, err
	}

	jobDetails, err := s.getJobDetails(req.JobId)
	if err != nil {
		return nil, err
	}

	log.Println("the taskId: ", jobDetails.TaskId)

	if err := s.updateJobStatus(jobDetails.TaskId, "IN_PROGRESS"); err != nil {
		return nil, err
	}

	if err := s.processJob(jobDetails); err != nil {
		return nil, err
	}

	return &pb.StartJobResponse{Status: "IN_PROGRESS"}, nil
}

// GetJobStatus streams job status updates
func (s *TaskServer) GetJobStatus(req *pb.GetJobStatusRequest, stream pb.TaskService_GetJobStatusServer) error {
	if req.JobId == "" {
		return &models.ServiceResponse{
			Code:    404,
			Message: "Job ID is required",
		}
	}

	jobID := uuid.MustParse(req.JobId)
	return s.streamJobStatus(stream, jobID)
}

// Helper methods
func validateStartJobRequest(req *pb.StartJobRequest) error {
	if req.JobId == "" || req.JobName == "" {
		return &models.ServiceResponse{
			Code:    504,
			Message: "Request Body should have the name and jobId",
		}
	}
	return nil
}

func (s *TaskServer) getJobDetails(jobID string) (*models.DBJobDetails, error) {
	jobDetails, err := s.deps.JobRepo.FindBy(&models.DBJobDetails{
		TaskId: uuid.MustParse(jobID),
	})
	if err != nil {
		return nil, &models.ServiceResponse{
			Code:    500,
			Message: "error whie finding the the job details: " + err.Error(),
		}
	}
	return jobDetails, nil
}

func (s *TaskServer) updateJobStatus(jobID uuid.UUID, status string) error {
	if err := s.deps.JobRepo.UpdateStatus(jobID, status); err != nil {
		return &models.ServiceResponse{
			Code:    500,
			Message: "error whie updating the job details: " + err.Error(),
		}
	}
	return nil
}

func (s *TaskServer) processJob(jobDetails *models.DBJobDetails) error {
	userDetails, err := s.deps.UserRepo.FindByConditions(&models.DBUserDetails{
		UserId: jobDetails.UserId,
	})
	if err != nil {
		return &models.ServiceResponse{
			Code:    500,
			Message: "error whie finding the the user details: " + err.Error(),
		}
	}

	fileDetails, err := s.deps.FileRepo.FindBy(&models.DbFileDetails{
		Name: jobDetails.ObjectName,
	})
	if err != nil {
		return &models.ServiceResponse{
			Code:    500,
			Message: "error whie finding the the file details: " + err.Error(),
		}
	}

	obj, err := external.GetObject(s.deps.MinioClient, s.deps.MinioService, userDetails.Name, fileDetails.Name)
	if err != nil {
		return err
	}

	rec, err := job.ReadCSV(obj)
	if err != nil {
		return &models.ServiceResponse{
			Code:    500,
			Message: "error whie reading the csv: " + err.Error(),
		}
	}

	log.Println("the records: ", rec)

	smtpDetails, err := s.deps.SmtpRepo.FindBy(&models.DBSMTPDetails{
		UserId: jobDetails.UserId,
	})
	if err != nil {
		return &models.ServiceResponse{
			Code:    500,
			Message: "error whie finding the the smtp details: " + err.Error(),
		}
	}

	log.Printf("SMTP details found: %+v", smtpDetails)
	return nil
}

func (s *TaskServer) streamJobStatus(stream pb.TaskService_GetJobStatusServer, jobID uuid.UUID) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			jobDetails, err := s.deps.JobRepo.FindBy(&models.DBJobDetails{TaskId: jobID})
			if err != nil {
				return err
			}

			if err := stream.Send(&pb.GetJobStatusResponse{Status: jobDetails.Status}); err != nil {
				log.Printf("Error sending status update: %v", err)
				return err
			}

			if jobDetails.Status == "COMPLETED" {
				log.Printf("Job %s completed", jobID)
				return nil
			}
		}
	}
}

// Initialization and main
func initializeDependencies() (*ServerDependencies, error) {
	jobRepo, err := dal.NewJobDalRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize job DAL: %w", err)
	}

	smtpRepo, err := dal.NewSmtpDalRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SMTP DAL: %w", err)
	}

	fileRepo, err := dal.NewFileDalRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file DAL: %w", err)
	}

	userRepo, err := dal.NewUserDalRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user DAL: %w", err)
	}

	minioClient, minioService, err := external.ConnectMinio()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	return &ServerDependencies{
		JobRepo:      jobRepo,
		SmtpRepo:     smtpRepo,
		FileRepo:     fileRepo,
		UserRepo:     userRepo,
		MinioClient:  minioClient,
		MinioService: minioService,
	}, nil
}

func main() {
	deps, err := initializeDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	server := NewTaskServer(*deps)

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, server)

	log.Println("Server running at port 3000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
