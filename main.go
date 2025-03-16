package main

import (
	"context"
	"log"

	"github.com/SendHive/worker-service/client"
	"github.com/SendHive/worker-service/external"
	"github.com/SendHive/worker-service/job"
	pb "github.com/SendHive/worker-service/proto"
)

func main() {

	qConn, Iq, err := external.SetupQueue()

	if err != nil {
		log.Println(err)
		return
	}

	qu, err := external.DeclareQueue(qConn, Iq)
	if err != nil {
		log.Println(err)
		return
	}

	clientInstance := client.InitClient()
	resp, err := clientInstance.HealthCheck(context.Background(), &pb.NoParams{})
	if err != nil {
		log.Println("error while getting health: ", err)
		return
	}
	log.Println(resp)

	// Creating the jobInstance
	Ijob, err := job.NewJobServiceRequest(clientInstance, qu, Iq, qConn)
	if err != nil {
		log.Println("error while getting the job instance: ", err)
		return
	}
	Jresp, err := Ijob.ConsumeMessage()
	if err != nil {
		log.Println("error while consuming message: ", err)
		return
	}
	log.Println(Jresp)

}
