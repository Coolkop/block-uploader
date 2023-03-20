package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	accessKeyID     = "minioadmin"
	secretAccessKey = "minioadmin"
	location        = "us-east-1"
	bucketName      = "storage"
)

type Client struct {
	clientByStorageIDs map[int]externalClient
}

func New(externalClients map[int]externalClient) *Client {
	return &Client{
		clientByStorageIDs: externalClients,
	}
}

func NewExternalClients(cfg Config) (map[int]externalClient, error) {
	clientByStorageIDs := make(map[int]externalClient, len(cfg.Hosts))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for id, host := range cfg.Hosts {
		minioClient, err := minio.New(host, &minio.Options{
			Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		})
		if err != nil {
			return nil, fmt.Errorf("init minio client with id '%d' and host '%s'", id, host)
		}

		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
			if errBucketExists != nil {
				return nil, fmt.Errorf("check bucket existence for storageID '%d': %w", id, errBucketExists)
			}

			if !exists {
				return nil, fmt.Errorf("create bucket for storageID '%d': %w", id, err)
			}
		}

		clientByStorageIDs[id] = minioClient
	}

	return clientByStorageIDs, nil
}

func (c *Client) Ping() error {
	for storageID, externalClient := range c.clientByStorageIDs {
		if externalClient.IsOnline() {
			return fmt.Errorf("external client with id '%d' not online", storageID)
		}
	}

	return nil
}

func (c *Client) Put(ctx context.Context, storageID int, fileID uuid.UUID, data []byte) error {
	externalClient, ok := c.clientByStorageIDs[storageID]
	if !ok {
		return fmt.Errorf("no client with ID '%d'", storageID)
	}

	_, err := externalClient.PutObject(
		ctx,
		bucketName,
		fileID.String(),
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{DisableMultipart: true},
	)
	if err != nil {
		return fmt.Errorf("put object in external client: %w", err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, storageID int, fileID uuid.UUID) ([]byte, error) {
	externalClient, ok := c.clientByStorageIDs[storageID]
	if !ok {
		return nil, fmt.Errorf("no client with ID '%d'", storageID)
	}

	obj, err := externalClient.GetObject(ctx, bucketName, fileID.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object from external client: %w", err)
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("obtain stat from returned object: %w", err)
	}

	result := make([]byte, stat.Size)

	_, err = obj.Read(result)
	obj.Close()

	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read data from object: %w", err)
	}

	return result, nil
}
