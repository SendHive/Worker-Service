package job

//Comman function file

import (
	"encoding/json"
	"log"
	"time"

	"github.com/SendHive/worker-service/models"
)

func (j *JobService) ConsumeMessage() (*models.QueueResponse, error) {
	ch, err := j.QConn.Channel()
	if err != nil {
		log.Println("Error while creating a channel:", err)
		return nil, err
	}
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Println("Error setting QoS:", err)
		return nil, err
	}

	msgs, err := ch.Consume(
		j.Qu.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error while consuming the msgs:", err)
		return nil, err
	}
	var resp *models.QueueResponse

	errChan := make(chan error)

	for d := range msgs {
		d := d
		go func() {
			log.Printf("Processing message: %s", d.Body)
			resp, err = processMessage(d.Body)
			if err != nil {
				log.Println("error while processing the error: ", err)
				errChan <- err
				return
			}
			err = j.StartJob(resp)
			if err != nil {
				log.Println("error while starting the job: ", err)
				errChan <- err
				return
			}
			err := d.Ack(false)
			if err != nil {
				log.Println("Failed to acknowledge message:", err)
				errChan <- err
			}
		}()
	}

	return resp, nil
}

func processMessage(req []byte) (*models.QueueResponse, error) {
	var res models.QueueResponse
	err := json.Unmarshal(req, &res)
	if err != nil {
		return nil, err
	}
	time.Sleep(10 * time.Second)
	log.Println("the process: ", res)
	return &res, nil
}
