package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/aliharirian/TerraPeak/config"
	"github.com/aliharirian/TerraPeak/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Storage implements S3/MinIO storage backend
type Storage struct {
	client *minio.Client
	bucket string
	config *config.Config
}

// New creates a new S3 storage instance
func New(cfg *config.Config) (*Storage, error) {
	s3Config := cfg.Storage.S3

	if s3Config.Endpoint == "" || s3Config.AccessKey == "" ||
		s3Config.SecretKey == "" || s3Config.Bucket == "" {
		return nil, fmt.Errorf("S3 configuration is incomplete: endpoint, access key, secret key, and bucket must be set")
	}

	// Parse endpoint to extract host:port and determine if SSL
	endpoint, useSSL, err := parseEndpoint(s3Config.Endpoint, s3Config.SkipSSL)
	if err != nil {
		logger.Errorf("Invalid S3 endpoint: %v", err)
		return nil, err
	}

	// Use default region if not specified
	region := s3Config.Region
	if region == "" {
		region = "us-east-1"
	}

	logger.Debugf("Connecting to S3: endpoint=%s, ssl=%v", endpoint, useSSL)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Config.AccessKey, s3Config.SecretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		logger.Errorf("Error initializing S3 client: %s", err)
		return nil, err
	}

	// Check and create bucket if needed
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, s3Config.Bucket)
	if err != nil {
		logger.Errorf("Failed to check if bucket %s exists: %v", s3Config.Bucket, err)
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, s3Config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Errorf("Failed to create bucket %s: %v", s3Config.Bucket, err)
			return nil, err
		}
		logger.Infof("Bucket %s created", s3Config.Bucket)
	}

	logger.Infof("S3 storage initialized successfully")
	return &Storage{
		client: client,
		bucket: s3Config.Bucket,
		config: cfg,
	}, nil
}

// Exists checks if file exists in S3
func (s *Storage) Exists(filePath string) bool {
	ctx := context.Background()
	_, err := s.client.StatObject(ctx, s.bucket, filePath, minio.StatObjectOptions{})
	if err != nil {
		logger.Debugf("File %s not found in S3: %v", filePath, err)
		return false
	}
	logger.Debugf("File %s exists in S3", filePath)
	return true
}

// Read reads file from S3
func (s *Storage) Read(filePath string) ([]byte, error) {
	ctx := context.Background()
	logger.Debugf("Reading file %s from S3", filePath)

	object, err := s.client.GetObject(ctx, s.bucket, filePath, minio.GetObjectOptions{})
	if err != nil {
		logger.Errorf("Failed to get object from S3: %v", err)
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		logger.Errorf("Failed to read object data: %v", err)
		return nil, err
	}

	logger.Infof("Successfully read file %s from S3 (%d bytes)", filePath, len(data))
	return data, nil
}

// Write writes file to S3
func (s *Storage) Write(filePath string, data []byte) error {
	ctx := context.Background()
	logger.Debugf("Writing %s to S3 (%d bytes)", filePath, len(data))

	_, err := s.client.PutObject(ctx, s.bucket, filePath, bytes.NewReader(data),
		int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		logger.Errorf("Failed to put object %s to S3: %v", filePath, err)
		return err
	}

	logger.Infof("Successfully wrote %s to S3", filePath)
	return nil
}

// StreamWrite streams data to S3
func (s *Storage) StreamWrite(filePath string, reader io.Reader, size int64) error {
	ctx := context.Background()
	logger.Debugf("Streaming %s to S3 (size: %d bytes)", filePath, size)

	_, err := s.client.PutObject(ctx, s.bucket, filePath, reader, size, minio.PutObjectOptions{})
	if err != nil {
		logger.Errorf("Failed to stream to S3: %v", err)
		return err
	}

	logger.Infof("Successfully streamed %s to S3", filePath)
	return nil
}

// StreamRead streams data from S3
func (s *Storage) StreamRead(filePath string) (io.ReadCloser, error) {
	ctx := context.Background()
	logger.Debugf("Opening stream for file %s from S3", filePath)

	object, err := s.client.GetObject(ctx, s.bucket, filePath, minio.GetObjectOptions{})
	if err != nil {
		logger.Errorf("Failed to get object from S3: %v", err)
		return nil, err
	}

	logger.Debugf("Successfully opened stream for file %s", filePath)
	return object, nil
}

// SaveMetadata saves metadata to S3
func (s *Storage) SaveMetadata(filePath, md5Sum, sha256Sum string, size int64) error {
	ctx := context.Background()

	metadata := fmt.Sprintf(`{
  "file": "%s",
  "timestamp": "%s",
  "size": %d,
  "md5": "%s",
  "sha256": "%s",
  "status": "success"
}`, filePath, time.Now().Format(time.RFC3339), size, md5Sum, sha256Sum)

	metadataPath := filePath + ".metadata.json"
	_, err := s.client.PutObject(ctx, s.bucket, metadataPath,
		bytes.NewReader([]byte(metadata)), int64(len(metadata)), minio.PutObjectOptions{})

	if err != nil {
		return err
	}

	logger.Debugf("Metadata saved to S3: %s", metadataPath)
	return nil
}

// parseEndpoint parses S3 endpoint URL and returns host:port and SSL flag
func parseEndpoint(endpointURL string, skipSSL bool) (string, bool, error) {
	parsedURL, err := url.Parse(endpointURL)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse endpoint URL: %v", err)
	}

	host := parsedURL.Host
	if host == "" {
		host = endpointURL
	}

	useSSL := false
	switch parsedURL.Scheme {
	case "https":
		useSSL = true
	case "http":
		useSSL = false
	default:
		useSSL = !skipSSL
	}

	if skipSSL {
		useSSL = false
	}

	logger.Debugf("Parsed S3 endpoint: original=%s, host=%s, ssl=%v", endpointURL, host, useSSL)
	return host, useSSL, nil
}
