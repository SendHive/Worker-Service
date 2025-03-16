package external

import (
	minioDb "github.com/SendHive/Infra-Common/minio"
	"github.com/minio/minio-go/v7"
)

func ConnectMinio() (*minio.Client, minioDb.IMinioService, error) {
	dbI, err := minioDb.NewMinioRequest()
	if err != nil {
		return nil, nil, err
	}
	conn, err := dbI.MinioConnect()
	if err != nil {
		return nil, nil, err
	}
	return conn, dbI, nil
}
