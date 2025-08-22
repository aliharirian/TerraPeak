package store

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/aliharirian/TerraPeak/config"
	"github.com/aliharirian/TerraPeak/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Store struct {
	config *config.Config
	client *minio.Client // MinIO client (nil if not using MinIO)
}

// New creates a new Store instance with the given config
func New(cfg *config.Config) (*Store, error) {
	store := &Store{
		config: cfg,
	}

	if err := store.init(); err != nil {
		return nil, err
	}

	return store, nil
}

// init initializes MinIO client if MinIO is enabled
func (s *Store) init() error {
	if !s.config.Storage.Minio.Enabled {
		logger.Infof("MinIO is disabled, using file storage")
		return nil
	}

	minioConfig := s.config.Storage.Minio
	logger.Debugf("Initializing MinIO client with endpoint %s, bucket %s", minioConfig.Endpoint, minioConfig.Bucket)

	if minioConfig.Endpoint == "" || minioConfig.AccessKey == "" || minioConfig.SecretKey == "" || minioConfig.Bucket == "" {
		return fmt.Errorf("MinIO configuration is incomplete: endpoint, access key, secret key, and bucket must be set")
	}

	// Parse endpoint to extract host:port and determine if SSL
	endpoint, useSSL, err := s.parseMinIOEndpoint(minioConfig.Endpoint, minioConfig.SkipSSL)
	if err != nil {
		logger.Errorf("Invalid MinIO endpoint: %v", err)
		return err
	}

	// Use default region if not specified
	region := minioConfig.Region
	if region == "" {
		region = "us-east-1"
	}

	logger.Debugf("Connecting to MinIO: endpoint=%s, ssl=%v", endpoint, useSSL)
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		logger.Errorf("Error initializing MinIO client: %s", err)
		return err
	}
	s.client = minioClient

	// Check and create bucket if needed
	ctx := context.Background()
	exists, err := s.client.BucketExists(ctx, minioConfig.Bucket)
	if err != nil {
		logger.Errorf("Failed to check if bucket %s exists: %v", minioConfig.Bucket, err)
		return err
	}

	if !exists {
		err = s.client.MakeBucket(ctx, minioConfig.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Errorf("Failed to create bucket %s: %v", minioConfig.Bucket, err)
			return err
		}
		logger.Infof("Bucket %s created", minioConfig.Bucket)
	}

	logger.Infof("MinIO client initialized successfully")
	return nil
}

// HandleRequest is the main function that handles all requests
func (s *Store) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Generate download URL and file path from request
	downloadURL, filePath := s.generateURL(r.URL.String())
	if downloadURL == "" || filePath == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Check if file already exists in storage (cache hit)
	if s.FileExists(filePath) {
		logger.Infof("Cache HIT: File %s found in storage", filePath)

		// Read from storage and serve
		data, err := s.ReadFromStorage(filePath)
		if err != nil {
			logger.Errorf("Failed to read from storage: %v", err)
			http.Error(w, "Storage read failed", http.StatusInternalServerError)
			return
		}

		// Send cached file to user
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Header().Set("X-Cache-Status", "HIT")
		w.WriteHeader(http.StatusOK)
		w.Write(data)

		logger.Infof("Served cached file %s (%d bytes)", filePath, len(data))
		return
	}

	// Cache miss - need to download and store
	logger.Infof("Cache MISS: Downloading file %s", filePath)

	// Get stream from source
	sourceStream, contentLength, err := s.getSourceStream(downloadURL)
	if err != nil {
		logger.Errorf("Failed to get source stream: %v", err)
		http.Error(w, "Download failed", http.StatusInternalServerError)
		return
	}
	defer sourceStream.Close()

	// Create buffer for user response and multi-writer for streaming
	var userBuffer bytes.Buffer

	// Create a TeeReader that writes to user buffer while we process
	streamReader := io.TeeReader(sourceStream, &userBuffer)

	// Stream to appropriate storage
	var storageErr error
	if s.config.Storage.Minio.Enabled {
		storageErr = s.streamToMinio(filePath, streamReader, contentLength)
	} else {
		storageErr = s.streamToFile(filePath, streamReader)
	}

	if storageErr != nil {
		logger.Errorf("Failed to save to storage: %v", storageErr)
		http.Error(w, "Save failed", http.StatusInternalServerError)
		return
	}

	// Send file data to user
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", userBuffer.Len()))
	w.Header().Set("X-Cache-Status", "MISS")
	w.WriteHeader(http.StatusOK)
	w.Write(userBuffer.Bytes())

	logger.Infof("Downloaded and cached file %s (%d bytes)", filePath, userBuffer.Len())
}

// getSourceStream opens a stream from the download URL
func (s *Store) getSourceStream(downloadURL string) (io.ReadCloser, int64, error) {
	logger.Debugf("Opening stream from %s", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil {
		logger.Errorf("Error creating HTTP request: %v", err)
		return nil, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	logger.Debugf("Stream opened successfully, content length: %d", resp.ContentLength)
	return resp.Body, resp.ContentLength, nil
}

// streamToMinio streams data directly to MinIO storage
func (s *Store) streamToMinio(filePath string, reader io.Reader, contentLength int64) error {
	minioConfig := s.config.Storage.Minio
	ctx := context.Background()

	logger.Debugf("Streaming %s to MinIO (size: %d bytes)", filePath, contentLength)

	// Create hash readers for checksum
	md5Hash := md5.New()
	sha256Hash := sha256.New()

	// Multi-writer for hashing while streaming
	hashReader := io.TeeReader(reader, io.MultiWriter(md5Hash, sha256Hash))

	// Stream to MinIO
	_, err := s.client.PutObject(ctx, minioConfig.Bucket, filePath, hashReader, contentLength, minio.PutObjectOptions{})
	if err != nil {
		logger.Errorf("Failed to stream to MinIO: %v", err)
		return err
	}

	// Save checksums as metadata files
	md5Sum := hex.EncodeToString(md5Hash.Sum(nil))
	sha256Sum := hex.EncodeToString(sha256Hash.Sum(nil))

	err = s.saveMetadataToMinio(filePath, md5Sum, sha256Sum)
	if err != nil {
		logger.Warnf("Failed to save metadata to MinIO: %v", err)
	}

	proto := "http"
	if !minioConfig.SkipSSL {
		proto = "https"
	}

	logger.Infof("Successfully streamed %s to MinIO at %s://%s/%s/%s", filePath, proto, minioConfig.Endpoint, minioConfig.Bucket, filePath)
	return nil
}

// streamToFile streams data directly to file system with checksum
func (s *Store) streamToFile(filePath string, reader io.Reader) error {
	basePath := s.config.Storage.File.Path
	if basePath == "" {
		basePath = "./storage"
	}

	fullPath := filepath.Join(basePath, filePath)
	logger.Debugf("Streaming %s to file system at %s", filePath, fullPath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		logger.Errorf("Failed to create directories for %s: %v", fullPath, err)
		return err
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		logger.Errorf("Failed to create file %s: %v", fullPath, err)
		return err
	}
	defer file.Close()

	// Create hash writers for checksum
	md5Hash := md5.New()
	sha256Hash := sha256.New()

	// Multi-writer for file + hashing
	multiWriter := io.MultiWriter(file, md5Hash, sha256Hash)

	// Stream to file while calculating hashes
	bytesWritten, err := io.Copy(multiWriter, reader)
	if err != nil {
		logger.Errorf("Failed to stream to file %s: %v", fullPath, err)
		return err
	}

	// Calculate checksums
	md5Sum := hex.EncodeToString(md5Hash.Sum(nil))
	sha256Sum := hex.EncodeToString(sha256Hash.Sum(nil))

	// Save metadata file
	err = s.saveMetadataToFile(fullPath, md5Sum, sha256Sum, bytesWritten)
	if err != nil {
		logger.Warnf("Failed to save metadata: %v", err)
	}

	logger.Infof("Successfully saved %s to file system (%d bytes)", filePath, bytesWritten)
	return nil
}

// saveMetadataToMinio saves metadata file to MinIO
func (s *Store) saveMetadataToMinio(filePath, md5Sum, sha256Sum string) error {
	minioConfig := s.config.Storage.Minio
	ctx := context.Background()

	metadata := fmt.Sprintf(`{
  "file": "%s",
  "timestamp": "%s",
  "md5": "%s",
  "sha256": "%s",
  "status": "success"
}`, filePath, time.Now().Format(time.RFC3339), md5Sum, sha256Sum)

	metadataPath := filePath + ".metadata.json"
	_, err := s.client.PutObject(ctx, minioConfig.Bucket, metadataPath,
		bytes.NewReader([]byte(metadata)), int64(len(metadata)), minio.PutObjectOptions{})

	if err != nil {
		return err
	}

	logger.Debugf("Metadata saved to MinIO: %s", metadataPath)
	return nil
}

// saveMetadataToFile saves metadata file to file system
func (s *Store) saveMetadataToFile(fullPath, md5Sum, sha256Sum string, size int64) error {
	metadata := fmt.Sprintf(`{
  "file": "%s",
  "timestamp": "%s",
  "size": %d,
  "md5": "%s",
  "sha256": "%s",
  "status": "success"
}`, filepath.Base(fullPath), time.Now().Format(time.RFC3339), size, md5Sum, sha256Sum)

	metadataPath := fullPath + ".metadata.json"
	err := os.WriteFile(metadataPath, []byte(metadata), 0644)
	if err != nil {
		return err
	}

	// Also create a simple success tag file
	successTagPath := fullPath + ".success"
	successTag := fmt.Sprintf("File successfully saved at %s\nMD5: %s\nSHA256: %s",
		time.Now().Format(time.RFC3339), md5Sum, sha256Sum)

	err = os.WriteFile(successTagPath, []byte(successTag), 0644)
	if err != nil {
		logger.Warnf("Failed to create success tag: %v", err)
	}

	logger.Debugf("Metadata and success tag saved: %s", metadataPath)
	return nil
}

// generateURL generates download URL and file path from request URL
func (s *Store) generateURL(requestURL string) (downloadURL string, filePath string) {
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		logger.Errorf("Failed to parse URL: %s", err)
		return "", ""
	}

	uri := parsedURL.Path
	httpsURL := "https:/" + uri

	// Validate the generated URL
	_, err = url.Parse(httpsURL)
	if err != nil {
		logger.Errorf("Error parsing generated URL: %v", err)
		return "", ""
	}

	logger.Debugf("Generated download URL: %s", httpsURL)

	// Remove leading slash for file path
	if len(uri) > 1 {
		filePath = uri[1:]
	} else {
		logger.Errorf("Invalid URI path: %s", uri)
		return "", ""
	}

	return httpsURL, filePath
}

// parseMinIOEndpoint parses MinIO endpoint URL and returns host:port and SSL flag
func (s *Store) parseMinIOEndpoint(endpointURL string, skipSSL bool) (string, bool, error) {
	// Parse the URL
	parsedURL, err := url.Parse(endpointURL)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse endpoint URL: %v", err)
	}

	// Extract host (includes port if specified)
	host := parsedURL.Host
	if host == "" {
		// If no protocol was provided, treat the whole thing as host
		host = endpointURL
	}

	// Determine SSL usage based on scheme or config
	useSSL := false
	if parsedURL.Scheme == "https" {
		useSSL = true
	} else if parsedURL.Scheme == "http" {
		useSSL = false
	} else {
		// No scheme provided, use config setting
		useSSL = !skipSSL
	}

	// Override with skipSSL config if explicitly set
	if skipSSL {
		useSSL = false
	}

	logger.Debugf("Parsed MinIO endpoint: original=%s, host=%s, ssl=%v", endpointURL, host, useSSL)
	return host, useSSL, nil
}

// FileExists checks if a file exists in storage
func (s *Store) FileExists(filePath string) bool {
	if s.config.Storage.Minio.Enabled {
		return s.fileExistsInMinio(filePath)
	}
	return s.fileExistsInFileSystem(filePath)
}

// fileExistsInMinio checks if file exists in MinIO
func (s *Store) fileExistsInMinio(filePath string) bool {
	minioConfig := s.config.Storage.Minio
	ctx := context.Background()

	_, err := s.client.StatObject(ctx, minioConfig.Bucket, filePath, minio.StatObjectOptions{})
	if err != nil {
		logger.Debugf("File %s not found in MinIO: %v", filePath, err)
		return false
	}

	logger.Debugf("File %s exists in MinIO", filePath)
	return true
}

// fileExistsInFileSystem checks if file exists in file system
func (s *Store) fileExistsInFileSystem(filePath string) bool {
	basePath := s.config.Storage.File.Path
	if basePath == "" {
		basePath = "./storage"
	}

	fullPath := filepath.Join(basePath, filePath)
	_, err := os.Stat(fullPath)
	if err != nil {
		logger.Debugf("File %s not found in file system: %v", filePath, err)
		return false
	}

	logger.Debugf("File %s exists in file system", filePath)
	return true
}

// ReadFromStorage reads file from storage and returns the data
func (s *Store) ReadFromStorage(filePath string) ([]byte, error) {
	if s.config.Storage.Minio.Enabled {
		return s.readFromMinio(filePath)
	}
	return s.readFromFileSystem(filePath)
}

// readFromMinio reads file from MinIO storage
func (s *Store) readFromMinio(filePath string) ([]byte, error) {
	minioConfig := s.config.Storage.Minio
	ctx := context.Background()

	logger.Debugf("Reading file %s from MinIO", filePath)

	object, err := s.client.GetObject(ctx, minioConfig.Bucket, filePath, minio.GetObjectOptions{})
	if err != nil {
		logger.Errorf("Failed to get object from MinIO: %v", err)
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		logger.Errorf("Failed to read object data: %v", err)
		return nil, err
	}

	logger.Infof("Successfully read file %s from MinIO (%d bytes)", filePath, len(data))
	return data, nil
}

// readFromFileSystem reads file from file system
func (s *Store) readFromFileSystem(filePath string) ([]byte, error) {
	basePath := s.config.Storage.File.Path
	if basePath == "" {
		basePath = "./storage"
	}

	fullPath := filepath.Join(basePath, filePath)
	logger.Debugf("Reading file %s from file system", fullPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		logger.Errorf("Failed to read file from file system: %v", err)
		return nil, err
	}

	logger.Infof("Successfully read file %s from file system (%d bytes)", filePath, len(data))
	return data, nil
}

// Save saves data to storage using the configured method
func (s *Store) Save(filename string, data []byte) error {
	if s.config.Storage.Minio.Enabled {
		return s.saveToMinio(filename, data)
	}
	return s.saveToFile(filename, data)
}

// saveToMinio saves data to MinIO storage
func (s *Store) saveToMinio(filename string, data []byte) error {
	minioConfig := s.config.Storage.Minio
	ctx := context.Background()

	_, err := s.client.PutObject(ctx, minioConfig.Bucket, filename, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		logger.Errorf("Failed to put object %s to MinIO: %v", filename, err)
		return err
	}

	proto := "http"
	if !minioConfig.SkipSSL {
		proto = "https"
	}
	logger.Infof("Successfully saved %s to MinIO at %s://%s/%s/%s", filename, proto, minioConfig.Endpoint, minioConfig.Bucket, filename)
	return nil
}

// saveToFile saves data to file system
func (s *Store) saveToFile(filename string, data []byte) error {
	basePath := s.config.Storage.File.Path
	if basePath == "" {
		basePath = "./storage"
	}

	fullPath := filepath.Join(basePath, filename)
	logger.Debugf("Saving %s to local path %s", filename, fullPath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		logger.Errorf("Failed to create directories for %s: %v", fullPath, err)
		return err
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		logger.Errorf("Failed to save local file %s: %v", fullPath, err)
		return err
	}

	logger.Infof("Successfully saved %s to local path %s", filename, fullPath)
	return nil
}
