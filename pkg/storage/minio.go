package storage

import (
	"context"
	"sora_landing_be/pkg/config"
	"sora_landing_be/pkg/logger"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	once        = &sync.Once{}
	MinioClient *Object
)

func InitMinioStorage(config config.ObjectStorage) {
	once.Do(func() {

		minioClientObj, err := minio.New(config.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
			Secure: config.UseSSL,
		})
		if err != nil {
			logger.Log.Fatal(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		isExist, err := minioClientObj.BucketExists(ctx, config.Bucket)
		if err != nil {
			logger.Log.Fatal(err.Error(), zap.String("bucket", config.Bucket))
		}

		if !isExist {
			logger.Log.Error("Bucket is non exists: " + config.Bucket)
		}

		MinioClient = &Object{
			client: minioClientObj,
		}
	})

}
