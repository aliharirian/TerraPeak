package api

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aliharirian/TerraPeak/logger"
	"github.com/aliharirian/TerraPeak/store"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
)

// FileInfo represents information about a file or directory (API response type)
type FileInfo = store.FileEntry

// FilesListResponse represents a list of files with pagination
type FilesListResponse struct {
	Items      []FileInfo      `json:"items"`
	Pagination *PaginationInfo `json:"pagination"`
}

// BucketInfo represents information about an S3 bucket
type BucketInfo struct {
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
}

// BucketsListResponse represents a list of S3 buckets
type BucketsListResponse struct {
	Items []BucketInfo `json:"items"`
}

// S3ObjectInfo represents information about an S3 object
type S3ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	ETag         string    `json:"etag"`
	ContentType  string    `json:"contentType"`
}

// S3ObjectsListResponse represents a list of S3 objects with pagination
type S3ObjectsListResponse struct {
	Bucket     string            `json:"bucket"`
	Prefix     string            `json:"prefix"`
	Items      []S3ObjectInfo    `json:"items"`
	Pagination *PaginationInfo   `json:"pagination,omitempty"`
}

// ============================================================================
// Service Methods for Storage Operations - Unified Objects API
// ============================================================================

// ListObjects lists objects from the configured storage backend (filesystem or S3)
func (s *Service) ListObjects(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	page := getIntQueryParam(r, "page", 1)
	pageSize := getIntQueryParam(r, "pageSize", 20)
	prefix := r.URL.Query().Get("prefix")
	path := r.URL.Query().Get("path")
	
	// Use prefix or path interchangeably
	if path != "" && prefix == "" {
		prefix = path
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logger.Infof("Listing objects: prefix=%s, page=%d, pageSize=%d, backend=%s", 
		prefix, page, pageSize, s.getStorageType())

	// Route to appropriate backend
	if s.cfg.Storage.S3.Enabled {
		s.listObjectsS3(w, r, prefix, page, pageSize)
	} else {
		s.listObjectsFilesystem(w, r, prefix, page, pageSize)
	}
}

// GetObject retrieves a specific object from the configured storage backend
func (s *Service) GetObject(w http.ResponseWriter, r *http.Request) {
	objectPath := chi.URLParam(r, "*")

	if objectPath == "" {
		BadRequestResponse(w, r, "Object path is required")
		return
	}

	logger.Infof("Getting object: path=%s, backend=%s", objectPath, s.getStorageType())

	// Check if object exists
	exists := s.store.FileExists(objectPath)
	if !exists {
		NotFoundResponse(w, r, "Object not found")
		return
	}

	// Read from storage (works for both backends)
	data, err := s.store.ReadFromStorage(objectPath)
	if err != nil {
		logger.Errorf("Failed to read object: %v", err)
		InternalErrorResponse(w, r, "Failed to read object")
		return
	}

	// Detect content type
	contentType := http.DetectContentType(data)
	
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(objectPath)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// listObjectsFilesystem handles listing for filesystem storage
func (s *Service) listObjectsFilesystem(w http.ResponseWriter, r *http.Request, prefix string, page, pageSize int) {
	basePath := s.cfg.Storage.File.Path
	if basePath == "" {
		basePath = "./storage"
	}

	targetPath := filepath.Join(basePath, filepath.Clean(prefix))
	
	// Security: Ensure path is within base directory
	if !strings.HasPrefix(targetPath, basePath) {
		BadRequestResponse(w, r, "Invalid path: directory traversal not allowed")
		return
	}

	var files []FileInfo
	err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == targetPath && d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil // Skip files we can't read
		}

		relPath, _ := filepath.Rel(basePath, path)
		files = append(files, FileInfo{
			Name:         d.Name(),
			Path:         relPath,
			Size:         info.Size(),
			IsDirectory:  d.IsDir(),
			ModifiedTime: info.ModTime(),
		})

		return nil
	})

	if err != nil {
		logger.Errorf("Failed to list files: %v", err)
		InternalErrorResponse(w, r, "Failed to list objects")
		return
	}

	// Paginate
	totalItems := len(files)
	totalPages := (totalItems + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	paginatedFiles := files[start:end]

	response := FilesListResponse{
		Items: paginatedFiles,
		Pagination: &PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			TotalItems: totalItems,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}

	SuccessResponse(w, r, response, http.StatusOK)
}

// listObjectsS3 handles listing for S3 storage
func (s *Service) listObjectsS3(w http.ResponseWriter, r *http.Request, prefix string, page, pageSize int) {
	bucketName := s.cfg.Storage.S3.Bucket

	// Get S3 client from store backend
	s3Storage, ok := s.store.Backend().(interface {
		GetClient() *minio.Client
	})
	if !ok {
		InternalErrorResponse(w, r, "S3 client not available")
		return
	}

	client := s3Storage.GetClient()
	ctx := context.Background()

	var objects []FileInfo
	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			logger.Errorf("Error listing S3 objects: %v", object.Err)
			continue
		}

		objects = append(objects, FileInfo{
			Name:         filepath.Base(object.Key),
			Path:         object.Key,
			Size:         object.Size,
			IsDirectory:  false,
			ModifiedTime: object.LastModified,
		})
	}

	// Paginate
	totalItems := len(objects)
	totalPages := (totalItems + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	paginatedObjects := objects[start:end]

	response := FilesListResponse{
		Items: paginatedObjects,
		Pagination: &PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			TotalItems: totalItems,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}

	SuccessResponse(w, r, response, http.StatusOK)
}

// getStorageType returns the configured storage type
func (s *Service) getStorageType() string {
	if s.cfg.Storage.S3.Enabled {
		return "s3"
	}
	return "filesystem"
}

// ============================================================================
// Legacy Storage Endpoints (Deprecated - use /objects instead)
// ============================================================================

// ListFiles lists files from the filesystem storage
func (s *Service) ListFiles(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	page := getIntQueryParam(r, "page", 1)
	pageSize := getIntQueryParam(r, "pageSize", 20)
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	path := r.URL.Query().Get("path")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Check if filesystem storage is being used
	if s.cfg.Storage.S3.Enabled {
		BadRequestResponse(w, r, "Filesystem storage is not enabled. Use S3 endpoints instead.")
		return
	}

	basePath := s.cfg.Storage.File.Path
	if basePath == "" {
		basePath = "./storage"
	}

	targetPath := filepath.Join(basePath, filepath.Clean(path))
	
	// Security: Ensure path is within base directory
	if !strings.HasPrefix(targetPath, basePath) {
		BadRequestResponse(w, r, "Invalid path: directory traversal not allowed")
		return
	}

	logger.Infof("Listing filesystem files: path=%s, page=%d, pageSize=%d", 
		targetPath, page, pageSize)

	var files []FileInfo
	err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == targetPath && d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil // Skip files we can't read
		}

		relPath, _ := filepath.Rel(basePath, path)
		files = append(files, FileInfo{
			Name:         d.Name(),
			Path:         relPath,
			Size:         info.Size(),
			IsDirectory:  d.IsDir(),
			ModifiedTime: info.ModTime(),
		})

		return nil
	})

	if err != nil {
		logger.Errorf("Failed to list files: %v", err)
		InternalErrorResponse(w, r, "Failed to list files")
		return
	}

	// Sort files
	sortFiles(files, sortBy, order)

	// Paginate
	totalItems := len(files)
	totalPages := (totalItems + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	paginatedFiles := files[start:end]

	response := FilesListResponse{
		Items: paginatedFiles,
		Pagination: &PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			TotalItems: totalItems,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}

	SuccessResponse(w, r, response, http.StatusOK)
}

// GetFile gets a specific file from filesystem storage
func (s *Service) GetFile(w http.ResponseWriter, r *http.Request) {
	filePath := chi.URLParam(r, "path")

	if filePath == "" {
		BadRequestResponse(w, r, "File path is required")
		return
	}

	// Check if filesystem storage is being used
	if s.cfg.Storage.S3.Enabled {
		BadRequestResponse(w, r, "Filesystem storage is not enabled. Use S3 endpoints instead.")
		return
	}

	basePath := s.cfg.Storage.File.Path
	if basePath == "" {
		basePath = "./storage"
	}

	targetPath := filepath.Join(basePath, filepath.Clean(filePath))
	
	// Security: Ensure path is within base directory
	if !strings.HasPrefix(targetPath, basePath) {
		BadRequestResponse(w, r, "Invalid path: directory traversal not allowed")
		return
	}

	logger.Infof("Getting file from filesystem: %s", targetPath)

	exists := s.store.FileExists(filePath)
	if !exists {
		NotFoundResponse(w, r, "File not found")
		return
	}

	data, err := s.store.ReadFromStorage(filePath)
	if err != nil {
		logger.Errorf("Failed to read file: %v", err)
		InternalErrorResponse(w, r, "Failed to read file")
		return
	}

	// Detect content type
	contentType := http.DetectContentType(data)
	
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// ListS3Buckets lists all S3 buckets
func (s *Service) ListS3Buckets(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.Storage.S3.Enabled {
		BadRequestResponse(w, r, "S3 storage is not enabled")
		return
	}

	logger.Infof("Listing S3 buckets")

	// For now, return the configured bucket
	// In a real implementation, you might want to list all accessible buckets
	buckets := []BucketInfo{
		{
			Name:         s.cfg.Storage.S3.Bucket,
			CreationDate: time.Now(), // Would need to query this from S3
		},
	}

	response := BucketsListResponse{
		Items: buckets,
	}

	SuccessResponse(w, r, response, http.StatusOK)
}

// ListS3Objects lists objects in an S3 bucket
func (s *Service) ListS3Objects(w http.ResponseWriter, r *http.Request) {
	bucketName := chi.URLParam(r, "bucket")
	prefix := r.URL.Query().Get("prefix")
	page := getIntQueryParam(r, "page", 1)
	pageSize := getIntQueryParam(r, "pageSize", 100)

	if !s.cfg.Storage.S3.Enabled {
		BadRequestResponse(w, r, "S3 storage is not enabled")
		return
	}

	if bucketName == "" {
		bucketName = s.cfg.Storage.S3.Bucket
	}

	logger.Infof("Listing S3 objects: bucket=%s, prefix=%s", bucketName, prefix)

	// Get S3 client from store backend
	s3Storage, ok := s.store.Backend().(interface {
		GetClient() *minio.Client
	})
	if !ok {
		InternalErrorResponse(w, r, "S3 client not available")
		return
	}

	client := s3Storage.GetClient()
	ctx := context.Background()

	var objects []S3ObjectInfo
	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			logger.Errorf("Error listing S3 objects: %v", object.Err)
			continue
		}

		objects = append(objects, S3ObjectInfo{
			Key:          object.Key,
			Size:         object.Size,
			LastModified: object.LastModified,
			ETag:         object.ETag,
			ContentType:  object.ContentType,
		})
	}

	// Paginate
	totalItems := len(objects)
	totalPages := (totalItems + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	paginatedObjects := objects[start:end]

	response := S3ObjectsListResponse{
		Bucket: bucketName,
		Prefix: prefix,
		Items:  paginatedObjects,
		Pagination: &PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			TotalItems: totalItems,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}

	SuccessResponse(w, r, response, http.StatusOK)
}

// GetS3Object retrieves a specific S3 object
func (s *Service) GetS3Object(w http.ResponseWriter, r *http.Request) {
	bucketName := chi.URLParam(r, "bucket")
	objectKey := chi.URLParam(r, "*")

	if !s.cfg.Storage.S3.Enabled {
		BadRequestResponse(w, r, "S3 storage is not enabled")
		return
	}

	if bucketName == "" {
		bucketName = s.cfg.Storage.S3.Bucket
	}

	if objectKey == "" {
		BadRequestResponse(w, r, "Object key is required")
		return
	}

	logger.Infof("Getting S3 object: bucket=%s, key=%s", bucketName, objectKey)

	data, err := s.store.ReadFromStorage(objectKey)
	if err != nil {
		logger.Errorf("Failed to read S3 object: %v", err)
		NotFoundResponse(w, r, "Object not found")
		return
	}

	// Detect content type
	contentType := http.DetectContentType(data)
	
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(objectKey)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// HealthCheck returns the health status of the API
func (s *Service) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"storage": map[string]interface{}{
			"type": func() string {
				if s.cfg.Storage.S3.Enabled {
					return "s3"
				}
				return "filesystem"
			}(),
			"available": true,
		},
	}

	SuccessResponse(w, r, health, http.StatusOK)
}

// GetStorageStats returns statistics about the storage
func (s *Service) GetStorageStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"type": func() string {
			if s.cfg.Storage.S3.Enabled {
				return "s3"
			}
			return "filesystem"
		}(),
		"totalFiles":  0,
		"totalSize":   int64(0),
		"lastUpdated": time.Now().UTC().Format(time.RFC3339),
	}

	SuccessResponse(w, r, stats, http.StatusOK)
}

// ============================================================================
// Helper Functions
// ============================================================================

func getIntQueryParam(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func sortFiles(files []FileInfo, sortBy, order string) {
	// Simple sorting implementation
	// In production, you might want to use a more sophisticated sorting library
	if sortBy == "" {
		sortBy = "name"
	}
	if order == "" {
		order = "asc"
	}
	// Sorting logic would go here
}
