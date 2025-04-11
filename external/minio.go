package external

import (
	minioDb "github.com/SendHive/Infra-Common/minio"
	"github.com/SendHive/worker-service/models"
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

func GetObject(mc *minio.Client, mI minioDb.IMinioService, bucketName string, objectName string) (object *minio.Object, err error) {
	obj, err := mI.GetObject(mc, bucketName, objectName)
	if err != nil {
		return nil,  &models.ServiceResponse{
			Code: 500,
			Message: "error while listing the obejct : "+err.Error(),
		}
	}
	return obj, nil
}