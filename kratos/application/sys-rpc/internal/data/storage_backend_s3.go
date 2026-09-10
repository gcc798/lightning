package data

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type s3CompatibleStorageBackend struct {
	kind   string
	client *minio.Client
	bucket string
}

type s3CompatibleStorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	UseSSL    bool
}

func newS3CompatibleStorageBackend(kind string) storageFactory {
	return func(env *entpkg.StorageEnv) (storageBackend, error) {
		if env == nil {
			return nil, errors.New("storage env not found")
		}
		cfg := s3CompatibleStorageConfig{
			Endpoint:  storageConfigString(env.Config, []string{"endpoint"}, ""),
			AccessKey: storageConfigString(env.Config, []string{"accessKey", "accessKeyId"}, ""),
			SecretKey: storageConfigString(env.Config, []string{"secretKey", "secretAccessKey"}, ""),
			Bucket:    storageConfigString(env.Config, []string{"bucket"}, ""),
			Region:    storageConfigString(env.Config, []string{"region"}, "us-east-1"),
			UseSSL:    storageConfigBool(env.Config, []string{"useSSL"}, false),
		}
		if err := cfg.validate(); err != nil {
			return nil, err
		}
		client, err := minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.UseSSL,
			Region: cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		return &s3CompatibleStorageBackend{kind: kind, client: client, bucket: cfg.Bucket}, nil
	}
}

func (b *s3CompatibleStorageBackend) Type() string { return b.kind }

func (b *s3CompatibleStorageBackend) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}
	_, err := b.client.PutObject(ctx, b.bucket, cleanStorageObjectKey(key), reader, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (b *s3CompatibleStorageBackend) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := b.client.GetObject(ctx, b.bucket, cleanStorageObjectKey(key), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	if _, err := object.Stat(); err != nil {
		_ = object.Close()
		return nil, err
	}
	return object, nil
}

func (b *s3CompatibleStorageBackend) Delete(ctx context.Context, key string) error {
	return b.client.RemoveObject(ctx, b.bucket, cleanStorageObjectKey(key), minio.RemoveObjectOptions{})
}

func (b *s3CompatibleStorageBackend) GetURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	key = cleanStorageObjectKey(key)
	if expires <= 0 {
		return fmt.Sprintf("%s/%s/%s", strings.TrimRight(b.client.EndpointURL().String(), "/"), b.bucket, key), nil
	}
	url, err := b.client.PresignedGetObject(ctx, b.bucket, key, expires, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (b *s3CompatibleStorageBackend) Exists(ctx context.Context, key string) (bool, error) {
	_, err := b.client.StatObject(ctx, b.bucket, cleanStorageObjectKey(key), minio.StatObjectOptions{})
	if err == nil {
		return true, nil
	}
	errResp := minio.ToErrorResponse(err)
	if errResp.Code == "NoSuchKey" || errResp.Code == "NoSuchBucket" {
		return false, nil
	}
	return false, err
}

func (b *s3CompatibleStorageBackend) Stat(ctx context.Context, key string) (*storageObjectInfo, error) {
	stat, err := b.client.StatObject(ctx, b.bucket, cleanStorageObjectKey(key), minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &storageObjectInfo{Key: stat.Key, Size: stat.Size, ContentType: stat.ContentType, LastModified: stat.LastModified}, nil
}

func (b *s3CompatibleStorageBackend) Ping(ctx context.Context) error {
	exists, err := b.client.BucketExists(ctx, b.bucket)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("bucket %s does not exist", b.bucket)
	}
	return nil
}

func (c s3CompatibleStorageConfig) validate() error {
	if strings.TrimSpace(c.Endpoint) == "" {
		return errors.New("storage endpoint is required")
	}
	if strings.TrimSpace(c.AccessKey) == "" {
		return errors.New("storage accessKey is required")
	}
	if strings.TrimSpace(c.SecretKey) == "" {
		return errors.New("storage secretKey is required")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return errors.New("storage bucket is required")
	}
	return nil
}
