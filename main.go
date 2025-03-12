package main

import (
	"fmt"
	"log"

	"github.com/SendHive/worker-service/external"
)

func main() {
	dbConn, err := external.GetDbConn()
	if err != nil {
		log.Println("error while connecting to the database: ", err)
		return
	}
	fmt.Println(dbConn)

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

	qerr := external.ConsumeMessage(qu, qConn, false)
	if qerr != nil {
		log.Println(err)
		return
	}
}
