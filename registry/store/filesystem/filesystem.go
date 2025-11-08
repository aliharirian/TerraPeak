package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliharirian/TerraPeak/config"
	"github.com/aliharirian/TerraPeak/logger"
)

// Storage implements local filesystem storage backend
type Storage struct {
	basePath string
}

// New creates a new filesystem storage instance
func New(cfg *config.Config) (*Storage, error) {
	basePath := cfg.Storage.File.Path
	if basePath == "" {
		basePath = "./storage"
	}

	logger.Infof("Initializing FileSystem storage at: %s", basePath)

	// Ensure base directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		logger.Errorf("Failed to create base directory %s: %v", basePath, err)
		return nil, err
	}

	logger.Infof("FileSystem storage initialized successfully")
	return &Storage{basePath: basePath}, nil
}

// Exists checks if file exists in filesystem
func (s *Storage) Exists(filePath string) bool {
	fullPath := filepath.Join(s.basePath, filePath)
	_, err := os.Stat(fullPath)
	if err != nil {
		logger.Debugf("File %s not found in filesystem: %v", filePath, err)
		return false
	}
	logger.Debugf("File %s exists in filesystem", filePath)
	return true
}

// Read reads file from filesystem
func (s *Storage) Read(filePath string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, filePath)
	logger.Debugf("Reading file %s from filesystem", fullPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		logger.Errorf("Failed to read file from filesystem: %v", err)
		return nil, err
	}

	logger.Infof("Successfully read file %s from filesystem (%d bytes)", filePath, len(data))
	return data, nil
}

// Write writes file to filesystem
func (s *Storage) Write(filePath string, data []byte) error {
	fullPath := filepath.Join(s.basePath, filePath)
	logger.Debugf("Writing %s to filesystem at %s", filePath, fullPath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		logger.Errorf("Failed to create directories for %s: %v", fullPath, err)
		return err
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		logger.Errorf("Failed to write file %s: %v", fullPath, err)
		return err
	}

	logger.Infof("Successfully wrote %s to filesystem (%d bytes)", filePath, len(data))
	return nil
}

// StreamWrite streams data to filesystem
func (s *Storage) StreamWrite(filePath string, reader io.Reader, size int64) error {
	fullPath := filepath.Join(s.basePath, filePath)
	logger.Debugf("Streaming %s to filesystem at %s", filePath, fullPath)

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

	// Stream to file
	bytesWritten, err := io.Copy(file, reader)
	if err != nil {
		logger.Errorf("Failed to stream to file %s: %v", fullPath, err)
		return err
	}

	logger.Infof("Successfully streamed %s to filesystem (%d bytes)", filePath, bytesWritten)
	return nil
}

// StreamRead streams data from filesystem
func (s *Storage) StreamRead(filePath string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, filePath)
	logger.Debugf("Opening stream for file %s from filesystem", fullPath)

	file, err := os.Open(fullPath)
	if err != nil {
		logger.Errorf("Failed to open file %s: %v", fullPath, err)
		return nil, err
	}

	logger.Debugf("Successfully opened stream for file %s", filePath)
	return file, nil
}

// SaveMetadata saves metadata to filesystem
func (s *Storage) SaveMetadata(filePath, md5Sum, sha256Sum string, size int64) error {
	fullPath := filepath.Join(s.basePath, filePath)

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
// Move renames (moves) a file or directory to a new path
func (s *Storage) Move(oldPath, newPath string) error {
	fullOldPath := filepath.Join(s.basePath, oldPath)
	fullNewPath := filepath.Join(s.basePath, newPath)
	logger.Debugf("Moving from %s to %s", fullOldPath, fullNewPath)

	// Rename (move) the file or directory
	err := os.Rename(fullOldPath, fullNewPath)
	if err != nil {
		logger.Errorf("Failed to move %s to %s: %v", fullOldPath, fullNewPath, err)
		return err
	}

	logger.Infof("Successfully moved %s to %s", oldPath, newPath)
	return nil
}

// Delete removes a file or directory
func (s *Storage) Delete(filePath string) error {
	fullPath := filepath.Join(s.basePath, filePath)
	logger.Debugf("Deleting file or directory: %s", fullPath)

	// Remove the file or directory
	err := os.RemoveAll(fullPath)
	if err != nil {
		logger.Errorf("Failed to delete %s: %v", fullPath, err)
		return err
	}

	logger.Infof("Successfully deleted %s", filePath)
	return nil
}

// Copy copies a file or directory to a new location
func (s *Storage) Copy(srcPath, destPath string) error {
	fullSrcPath := filepath.Join(s.basePath, srcPath)
	fullDestPath := filepath.Join(s.basePath, destPath)
	logger.Debugf("Copying from %s to %s", fullSrcPath, fullDestPath)

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(fullDestPath), 0755); err != nil {
		logger.Errorf("Failed to create directories for %s: %v", fullDestPath, err)
		return err
	}

	// Copy the file or directory
	err := os.Rename(fullSrcPath, fullDestPath)
	if err != nil {
		logger.Errorf("Failed to copy %s to %s: %v", fullSrcPath, fullDestPath, err)
		return err
	}

	logger.Infof("Successfully copied %s to %s", srcPath, destPath)
	return nil
}

// Clear clears the storage by deleting all files and directories
func (s *Storage) Clear() error {
	logger.Infof("Clearing storage at: %s", s.basePath)

	// Remove all files and directories
	err := os.RemoveAll(s.basePath)
	if err != nil {
		logger.Errorf("Failed to clear storage: %v", err)
		return err
	}

	// Recreate the base directory
	err = os.MkdirAll(s.basePath, 0755)
	if err != nil {
		logger.Errorf("Failed to recreate base directory %s: %v", s.basePath, err)
		return err
	}

	logger.Infof("Storage cleared successfully")
	return nil
}

// PathJoin joins base path with given file path
func (s *Storage) PathJoin(filePath string) string {
	return filepath.Join(s.basePath, filePath)
}

// IsSuccessTag checks if the file is a success tag
func (s *Storage) IsSuccessTag(filePath string) bool {
	return strings.HasSuffix(filePath, ".success")
}

// IsMetadataFile checks if the file is a metadata file
func (s *Storage) IsMetadataFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".metadata.json")
}
