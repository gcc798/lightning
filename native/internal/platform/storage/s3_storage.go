package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Storage S3 存储实现（支持 MinIO/AWS S3）
type S3Storage struct {
	client *minio.Client
	bucket string
}

// NewS3Storage 创建 S3 存储实例
func NewS3Storage(config Config) (*S3Storage, error) {
	// 初始化 MinIO 客户端
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化 S3 客户端失败: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("检查 S3 Bucket 失败: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, config.Bucket, minio.MakeBucketOptions{Region: config.Region}); err != nil {
			exists, checkErr := client.BucketExists(ctx, config.Bucket)
			if checkErr != nil || !exists {
				return nil, fmt.Errorf("创建 S3 Bucket 失败: %w", errors.Join(err, checkErr))
			}
		}
	}

	return &S3Storage{
		client: client,
		bucket: config.Bucket,
	}, nil
}

// Upload 上传文件
func (s *S3Storage) Upload(ctx context.Context, key string, reader io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}
	return nil
}

// Download 下载文件
func (s *S3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %w", err)
	}
	return object, nil
}

// Delete 删除文件
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// GetURL 获取访问 URL
func (s *S3Storage) GetURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if expires == 0 {
		// 永久 URL（需要桶为公开访问）
		return fmt.Sprintf("%s/%s/%s", s.client.EndpointURL(), s.bucket, key), nil
	}

	// 临时 URL
	url, err := s.client.PresignedGetObject(ctx, s.bucket, key, expires, nil)
	if err != nil {
		return "", fmt.Errorf("生成访问 URL 失败: %w", err)
	}
	return url.String(), nil
}

// Exists 检查文件是否存在
func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("检查文件是否存在失败: %w", err)
	}
	return true, nil
}

// List 列出文件
func (s *S3Storage) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	objectCh := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("列出文件失败: %w", object.Err)
		}
		keys = append(keys, object.Key)
	}

	return keys, nil
}

// GetInfo 获取文件信息
func (s *S3Storage) GetInfo(ctx context.Context, key string) (*FileInfo, error) {
	stat, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	return &FileInfo{
		Key:          stat.Key,
		Size:         stat.Size,
		LastModified: stat.LastModified,
		ContentType:  stat.ContentType,
		ETag:         stat.ETag,
	}, nil
}
