package main

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/SendHive/worker-service/proto"
	"github.com/gin-gonic/gin"
)

var grpcClient pb.TaskServiceClient

func getJobStatusHandler(c *gin.Context) {
	jobID := c.Param("jobId")

	// Call gRPC streaming method
	stream, err := grpcClient.GetJobStatus(context.Background(), &pb.GetJobStatusRequest{JobId: jobID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start job status stream"})
		return
	}

	// Set up SSE (Server-Sent Events) for streaming HTTP response
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		fmt.Fprintf(c.Writer, "data: %s\n\n", resp.Status)
		c.Writer.Flush()

		if resp.Status == "COMPLETED" {
			break
		}
	}
}

func main() {
	r := gin.Default()

	// Define API endpoint
	r.GET("/job/:jobId/status", getJobStatusHandler)

	// Start the API server
	r.Run(":8080")
}
