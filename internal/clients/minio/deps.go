package minio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type externalClient interface {
	IsOnline() bool
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	PutObject(
		ctx context.Context,
		bucketName, objectName string,
		reader io.Reader,
		objectSize int64,
		opts minio.PutObjectOptions,
	) (info minio.UploadInfo, err error)
}
