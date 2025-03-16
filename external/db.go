package external

import (
	"log"

	infraDb "github.com/SendHive/Infra-Common/db"
	"gorm.io/gorm"
)


func GetDbConn() (*gorm.DB, error) {
	IdbConn, err := infraDb.NewDbRequest()
	if err != nil {
		log.Println("error while connecting to the database: ", err)
		return nil, err
	}

	dbConn, err := IdbConn.InitDB()
	if err != nil {
		log.Println("error while getting the database instance: ", err)
		return nil, err
	}
	return dbConn, nil
}
