package store

import (
	"io"
)

// Storage defines the interface for storage backends (S3, FileSystem, etc.)
type Storage interface {
	// Check if file exists
	Exists(filePath string) bool

	// Read file data
	Read(filePath string) ([]byte, error)

	// Write file data
	Write(filePath string, data []byte) error

	// Stream write (for large files)
	StreamWrite(filePath string, reader io.Reader, size int64) error

	// Stream read (for large files)
	StreamRead(filePath string) (io.ReadCloser, error)

	// Save metadata (checksums, etc.)
	SaveMetadata(filePath, md5Sum, sha256Sum string, size int64) error
}
