package storage

import (
	"sora_landing_be/pkg/logger"

	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"sora_landing_be/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/tags"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Object struct {
	client *minio.Client
}

type Storage interface {
	CreateObjectTemporary(ctx context.Context, key string, payload *multipart.FileHeader) (string, error)
	DeleteObjects(ctx context.Context, keys ...string) error
	RemoveObjectTags(ctx context.Context, fileKey ...string) error
	SetTags(ctx context.Context, fileKey string, tag map[string]string) error
	PresignURL(ctx context.Context, fileKey string) string
}

func (o *Object) CreateObjectTemporary(ctx context.Context, key string, payload *multipart.FileHeader) (string, error) {
	blob, err := payload.Open()
	if err != nil {
		return "", err
	}

	defer blob.Close()

	object, err := o.client.PutObject(ctx, config.LoadConfig().ObjectStorage.Bucket, key, blob, payload.Size, minio.PutObjectOptions{
		ContentType: payload.Header["Content-Type"][0],
		UserTags: map[string]string{
			"temporary": "true",
		},
	})
	if err != nil {
		return "", err
	}

	logger.Log.Info("Successfully upload file to object storage",
		zap.String("key", key),
		zap.String("object", object.Key),
		zap.Int64("size", payload.Size),
	)

	return object.Key, err
}

func (o *Object) DeleteObjects(ctx context.Context, keys ...string) error {
	g, asyncCtx := errgroup.WithContext(ctx)
	for _, key := range keys {
		if key == "" {
			continue
		}
		keyCopy := key
		g.Go(func() error {
			err := o.client.RemoveObject(asyncCtx, config.LoadConfig().ObjectStorage.Bucket, keyCopy, minio.RemoveObjectOptions{})
			return err
		})
	}

	return g.Wait()

}

func (o *Object) RemoveObjectTags(ctx context.Context, fileKey ...string) error {
	g, asyncCtx := errgroup.WithContext(ctx)
	for _, key := range fileKey {
		if key == "" {
			continue
		}
		keyCopy := key
		g.Go(func() error {
			err := o.client.RemoveObjectTagging(asyncCtx, config.LoadConfig().ObjectStorage.Bucket, keyCopy, minio.RemoveObjectTaggingOptions{})
			return err
		})
	}

	return g.Wait()
}

func (o *Object) SetTags(ctx context.Context, fileKey string, tag map[string]string) error {
	var minioTags tags.Tags
	for key, value := range tag {
		err := minioTags.Set(key, value)
		if err != nil {
			logger.Log.Error("Error setting tag", zap.String("key", key), zap.String("value", value), zap.Error(err))
		}
	}

	err := o.client.PutObjectTagging(ctx, config.LoadConfig().ObjectStorage.Bucket, fileKey, &minioTags, minio.PutObjectTaggingOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (o *Object) PresignURL(ctx context.Context, fileKey string) string {
	if fileKey == "" {
		return ""
	}

	reqParams := make(url.Values)
	reqParams.Set("content-disposition", fmt.Sprintf("inline; filename=\"%s\"", fileKey))

	presignedURL, err := o.client.PresignedGetObject(ctx, config.LoadConfig().ObjectStorage.Bucket, fileKey, config.LoadConfig().ObjectStorage.PresignExpiration, reqParams)
	if err != nil {
		logger.Log.Error("Error presigning url", zap.String("url", fileKey), zap.Error(err))
	}
	return presignedURL.String()
}
