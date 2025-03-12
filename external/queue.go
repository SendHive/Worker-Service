package external

import (
	"fmt"
	"log"
	"time"

	queue "github.com/SendHive/Infra-Common/queue"
	"github.com/rabbitmq/amqp091-go"
)

func SetupQueue() (*amqp091.Connection, queue.IQueueService, error) {
	qConn, err := queue.NewQueueRequest()
	if err != nil {
		log.Fatal("the error while creating the queue instance: ", err)
		return nil, nil, err
	}
	time.Sleep(3 * time.Second)
	qconn, err := qConn.Connect()
	if err != nil {
		return nil, nil, err
	}
	time.Sleep(3 * time.Second)
	return qconn, qConn, nil
}

func DeclareQueue(qConn *amqp091.Connection, Iq queue.IQueueService) (qu amqp091.Queue, err error) {
	queue, err := Iq.DeclareQueue(qConn)

	if err != nil {
		return amqp091.Queue{}, err
	}
	return queue, nil

}

func ConsumeMessage(qu amqp091.Queue, conn *amqp091.Connection, isTest bool) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Error while creating a channel:", err)
		return err
	}
	defer ch.Close()

	// Ensure prefetch is set to avoid overwhelming the consumer
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Println("Error setting QoS:", err)
		return err
	}

	msgs, err := ch.Consume(
		qu.Name,
		"",
		false, // Auto-Ack set to false to manually acknowledge after processing
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error while consuming the msgs:", err)
		return err
	}

	if isTest {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			_ = d.Ack(false) // Acknowledge message after processing
			break
		}
		fmt.Println("Consumed Message Successfully!")
	} else {
		for d := range msgs {
			d := d // Create a new variable to avoid loop variable reuse
			go func() {
				log.Printf("Processing message: %s", d.Body)
				processMessage()
				err := d.Ack(false) // Acknowledge message after processing
				if err != nil {
					log.Println("Failed to acknowledge message:", err)
				}
			}()
		}
	}

	return nil
}

func processMessage() {
	time.Sleep(10 * time.Second)
}
