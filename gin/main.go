package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SendHive/worker-service/client"
	pb "github.com/SendHive/worker-service/proto"
	"github.com/gin-gonic/gin"
)

func getJobStatusHandler(c *gin.Context) {
	jobID := c.Param("jobId")
	Iclient := client.InitClient()

	stream, err := Iclient.GetJobStatus(context.Background(), &pb.GetJobStatusRequest{
		JobId: jobID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message":"error while getting the grpc client: "+err.Error()})
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
	c.JSON(http.StatusOK, gin.H{"Message":"Status Completed"})
}

func main() {
	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Message":"GIN WORKING SUCCESSFULLY"})
	})
	r.GET("/job/:jobId/status", getJobStatusHandler)
	r.Run(":8000")
}
