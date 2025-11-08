package store

import "time"

// FileEntry represents a file or directory entry in storage
// This is used by the API layer for browsing storage
type FileEntry struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	IsDirectory  bool      `json:"isDirectory"`
	ModifiedTime time.Time `json:"modifiedTime"`
	ContentType  string    `json:"contentType,omitempty"`
}

// BrowseOptions contains options for listing/browsing storage
type BrowseOptions struct {
	Path      string
	Recursive bool
	MaxDepth  int
}

// StorageStats represents statistics about the storage backend
type StorageStats struct {
	Type        string    `json:"type"`
	TotalFiles  int       `json:"totalFiles"`
	TotalSize   int64     `json:"totalSize"`
	LastUpdated time.Time `json:"lastUpdated"`
}
