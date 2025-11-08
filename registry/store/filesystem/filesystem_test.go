package filesystem

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aliharirian/TerraPeak/config"
)

// setupTestStorage creates a temporary storage instance for testing
func setupTestStorage(t *testing.T) (*Storage, string, func()) {
	t.Helper()

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "filesystem-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create config
	cfg := &config.Config{}
	cfg.Storage.File.Path = tempDir

	// Create storage instance
	storage, err := New(cfg)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return storage, tempDir, cleanup
}

func TestNew(t *testing.T) {
	t.Run("success_with_custom_path", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "filesystem-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		cfg := &config.Config{}
		cfg.Storage.File.Path = tempDir

		storage, err := New(cfg)
		if err != nil {
			t.Errorf("New() error = %v, want nil", err)
		}
		if storage == nil {
			t.Error("New() returned nil storage")
		}
		if storage.basePath != tempDir {
			t.Errorf("basePath = %s, want %s", storage.basePath, tempDir)
		}
	})

	t.Run("success_with_default_path", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Storage.File.Path = ""

		storage, err := New(cfg)
		if err != nil {
			t.Errorf("New() error = %v, want nil", err)
		}
		if storage == nil {
			t.Error("New() returned nil storage")
		}
		if storage.basePath != "./storage" {
			t.Errorf("basePath = %s, want ./storage", storage.basePath)
		}

		// Cleanup
		os.RemoveAll("./storage")
	})

	t.Run("creates_directory_if_not_exists", func(t *testing.T) {
		tempBase, err := os.MkdirTemp("", "filesystem-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempBase)

		newPath := filepath.Join(tempBase, "nonexistent", "nested", "path")

		cfg := &config.Config{}
		cfg.Storage.File.Path = newPath

		storage, err := New(cfg)
		if err != nil {
			t.Errorf("New() error = %v, want nil", err)
		}
		if storage == nil {
			t.Error("New() returned nil storage")
		}

		// Verify directory was created
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			t.Error("Directory was not created")
		}
	})
}

func TestExists(t *testing.T) {
	storage, _, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("file_exists", func(t *testing.T) {
		// Create a test file
		testPath := "test/file.txt"
		data := []byte("test content")
		if err := storage.Write(testPath, data); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		if !storage.Exists(testPath) {
			t.Error("Exists() = false, want true for existing file")
		}
	})

	t.Run("file_not_exists", func(t *testing.T) {
		if storage.Exists("nonexistent/file.txt") {
			t.Error("Exists() = true, want false for nonexistent file")
		}
	})
}

func TestWrite(t *testing.T) {
	storage, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		testPath := "test/write.txt"
		testData := []byte("test write content")

		err := storage.Write(testPath, testData)
		if err != nil {
			t.Errorf("Write() error = %v, want nil", err)
		}

		// Verify file was written
		fullPath := filepath.Join(tempDir, testPath)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read written file: %v", err)
		}
		if !bytes.Equal(data, testData) {
			t.Errorf("Written data = %s, want %s", data, testData)
		}
	})

	t.Run("creates_nested_directories", func(t *testing.T) {
		testPath := "deep/nested/path/file.txt"
		testData := []byte("nested content")

		err := storage.Write(testPath, testData)
		if err != nil {
			t.Errorf("Write() error = %v, want nil", err)
		}

		// Verify file exists
		fullPath := filepath.Join(tempDir, testPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Error("Nested file was not created")
		}
	})

	t.Run("overwrites_existing_file", func(t *testing.T) {
		testPath := "test/overwrite.txt"
		originalData := []byte("original")
		newData := []byte("updated")

		// Write original
		if err := storage.Write(testPath, originalData); err != nil {
			t.Fatalf("Failed to write original: %v", err)
		}

		// Overwrite
		if err := storage.Write(testPath, newData); err != nil {
			t.Errorf("Write() error = %v, want nil", err)
		}

		// Verify overwrite
		fullPath := filepath.Join(tempDir, testPath)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)
		}
		if !bytes.Equal(data, newData) {
			t.Errorf("Data = %s, want %s", data, newData)
		}
	})
}

func TestRead(t *testing.T) {
	storage, _, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		testPath := "test/read.txt"
		testData := []byte("test read content")

		// Write file first
		if err := storage.Write(testPath, testData); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Read it back
		data, err := storage.Read(testPath)
		if err != nil {
			t.Errorf("Read() error = %v, want nil", err)
		}
		if !bytes.Equal(data, testData) {
			t.Errorf("Read() data = %s, want %s", data, testData)
		}
	})

	t.Run("file_not_found", func(t *testing.T) {
		_, err := storage.Read("nonexistent/file.txt")
		if err == nil {
			t.Error("Read() error = nil, want error for nonexistent file")
		}
	})
}

func TestStreamWrite(t *testing.T) {
	storage, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		testPath := "test/stream-write.txt"
		testData := []byte("stream write content")
		reader := bytes.NewReader(testData)

		err := storage.StreamWrite(testPath, reader, int64(len(testData)))
		if err != nil {
			t.Errorf("StreamWrite() error = %v, want nil", err)
		}

		// Verify file was written
		fullPath := filepath.Join(tempDir, testPath)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read streamed file: %v", err)
		}
		if !bytes.Equal(data, testData) {
			t.Errorf("Streamed data = %s, want %s", data, testData)
		}
	})

	t.Run("large_stream", func(t *testing.T) {
		testPath := "test/large-stream.bin"
		// Create 1MB of data
		testData := bytes.Repeat([]byte("A"), 1024*1024)
		reader := bytes.NewReader(testData)

		err := storage.StreamWrite(testPath, reader, int64(len(testData)))
		if err != nil {
			t.Errorf("StreamWrite() error = %v, want nil", err)
		}

		// Verify size
		fullPath := filepath.Join(tempDir, testPath)
		stat, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("Failed to stat file: %v", err)
		}
		if stat.Size() != int64(len(testData)) {
			t.Errorf("File size = %d, want %d", stat.Size(), len(testData))
		}
	})
}

func TestStreamRead(t *testing.T) {
	storage, _, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		testPath := "test/stream-read.txt"
		testData := []byte("stream read content")

		// Write file first
		if err := storage.Write(testPath, testData); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Stream read
		reader, err := storage.StreamRead(testPath)
		if err != nil {
			t.Errorf("StreamRead() error = %v, want nil", err)
		}
		defer reader.Close()

		// Read all data
		data, err := io.ReadAll(reader)
		if err != nil {
			t.Errorf("Failed to read stream: %v", err)
		}
		if !bytes.Equal(data, testData) {
			t.Errorf("Streamed data = %s, want %s", data, testData)
		}
	})

	t.Run("file_not_found", func(t *testing.T) {
		_, err := storage.StreamRead("nonexistent/file.txt")
		if err == nil {
			t.Error("StreamRead() error = nil, want error for nonexistent file")
		}
	})

	t.Run("reads_large_file", func(t *testing.T) {
		testPath := "test/large-read.bin"
		// Create 2MB of data
		testData := bytes.Repeat([]byte("B"), 2*1024*1024)

		if err := storage.Write(testPath, testData); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		reader, err := storage.StreamRead(testPath)
		if err != nil {
			t.Errorf("StreamRead() error = %v, want nil", err)
		}
		defer reader.Close()

		// Read in chunks to simulate streaming
		buffer := make([]byte, 8192)
		totalRead := 0
		for {
			n, err := reader.Read(buffer)
			totalRead += n
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("Error reading stream: %v", err)
				break
			}
		}

		if totalRead != len(testData) {
			t.Errorf("Total read = %d, want %d", totalRead, len(testData))
		}
	})
}

func TestSaveMetadata(t *testing.T) {
	storage, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		testPath := "test/file-with-metadata.txt"
		testData := []byte("test content")
		md5Sum := "abc123"
		sha256Sum := "def456"
		size := int64(len(testData))

		// Write file first
		if err := storage.Write(testPath, testData); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Save metadata
		err := storage.SaveMetadata(testPath, md5Sum, sha256Sum, size)
		if err != nil {
			t.Errorf("SaveMetadata() error = %v, want nil", err)
		}

		// Verify metadata file exists
		metadataPath := filepath.Join(tempDir, testPath+".metadata.json")
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			t.Error("Metadata file was not created")
		}

		// Verify success tag exists
		successPath := filepath.Join(tempDir, testPath+".success")
		if _, err := os.Stat(successPath); os.IsNotExist(err) {
			t.Error("Success tag file was not created")
		}

		// Read and verify metadata content
		metadataContent, err := os.ReadFile(metadataPath)
		if err != nil {
			t.Errorf("Failed to read metadata: %v", err)
		}
		metadataStr := string(metadataContent)
		if !strings.Contains(metadataStr, md5Sum) {
			t.Error("Metadata doesn't contain MD5 sum")
		}
		if !strings.Contains(metadataStr, sha256Sum) {
			t.Error("Metadata doesn't contain SHA256 sum")
		}
		if !strings.Contains(metadataStr, "success") {
			t.Error("Metadata doesn't contain success status")
		}

		// Read and verify success tag content
		successContent, err := os.ReadFile(successPath)
		if err != nil {
			t.Errorf("Failed to read success tag: %v", err)
		}
		successStr := string(successContent)
		if !strings.Contains(successStr, md5Sum) {
			t.Error("Success tag doesn't contain MD5 sum")
		}
		if !strings.Contains(successStr, sha256Sum) {
			t.Error("Success tag doesn't contain SHA256 sum")
		}
	})

	t.Run("creates_nested_metadata", func(t *testing.T) {
		testPath := "deep/nested/metadata/file.txt"
		md5Sum := "nested-md5"
		sha256Sum := "nested-sha256"
		size := int64(100)

		// Write file first
		if err := storage.Write(testPath, []byte("test")); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		err := storage.SaveMetadata(testPath, md5Sum, sha256Sum, size)
		if err != nil {
			t.Errorf("SaveMetadata() error = %v, want nil", err)
		}

		// Verify metadata file exists in nested path
		metadataPath := filepath.Join(tempDir, testPath+".metadata.json")
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			t.Error("Nested metadata file was not created")
		}
	})
}

func TestIntegration(t *testing.T) {
	storage, _, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("full_workflow", func(t *testing.T) {
		testPath := "integration/test-file.txt"
		testData := []byte("integration test content")

		// 1. Verify file doesn't exist
		if storage.Exists(testPath) {
			t.Error("File should not exist initially")
		}

		// 2. Write file
		if err := storage.Write(testPath, testData); err != nil {
			t.Fatalf("Failed to write: %v", err)
		}

		// 3. Verify file exists
		if !storage.Exists(testPath) {
			t.Error("File should exist after write")
		}

		// 4. Read file back
		readData, err := storage.Read(testPath)
		if err != nil {
			t.Fatalf("Failed to read: %v", err)
		}
		if !bytes.Equal(readData, testData) {
			t.Error("Read data doesn't match written data")
		}

		// 5. Save metadata
		md5Sum := "integration-md5"
		sha256Sum := "integration-sha256"
		if err := storage.SaveMetadata(testPath, md5Sum, sha256Sum, int64(len(testData))); err != nil {
			t.Fatalf("Failed to save metadata: %v", err)
		}

		// 6. Stream read
		reader, err := storage.StreamRead(testPath)
		if err != nil {
			t.Fatalf("Failed to stream read: %v", err)
		}
		defer reader.Close()

		streamData, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("Failed to read stream: %v", err)
		}
		if !bytes.Equal(streamData, testData) {
			t.Error("Stream read data doesn't match written data")
		}
	})
}
