package store

import (
	"github.com/aliharirian/TerraPeak/config"
	"github.com/aliharirian/TerraPeak/store/filesystem"
	"github.com/aliharirian/TerraPeak/store/s3"
)

// Store handles file storage with automatic backend selection
type Store struct {
	config  *config.Config
	backend Storage
}

// New creates a new Store instance
// Automatically selects backend based on config:
// - If S3.Enabled = true, uses S3/MinIO
// - Otherwise uses FileSystem (default)
func New(cfg *config.Config) (*Store, error) {
	// Select backend based on config
	var backend Storage
	var err error
	if cfg.Storage.S3.Enabled {
		backend, err = s3.New(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		backend, err = filesystem.New(cfg)
		if err != nil {
			return nil, err
		}
	}
	return &Store{config: cfg, backend: backend}, nil
}

// Store now only provides storage operations; HTTP caching and proxy are handled in cache package.

// FileExists checks if a file exists in storage
func (s *Store) FileExists(filePath string) bool {
	return s.backend.Exists(filePath)
}

// ReadFromStorage reads file from storage and returns the data
func (s *Store) ReadFromStorage(filePath string) ([]byte, error) {
	return s.backend.Read(filePath)
}

// Save saves data to storage
func (s *Store) Save(filename string, data []byte) error {
	return s.backend.Write(filename, data)
}
