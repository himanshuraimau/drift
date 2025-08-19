package main

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCASPathTransformFunc(t *testing.T) {
	key := "test_key"
	pathKey := CASPathTransformFunc(key)
	
	assert.NotEmpty(t, pathKey.PathName, "Path name should not be empty")
	assert.NotEmpty(t, pathKey.Filename, "Filename should not be empty")
	assert.Equal(t, 40, len(pathKey.Filename), "Filename should be SHA1 hash (40 chars)")
	
	// Test consistency
	pathKey2 := CASPathTransformFunc(key)
	assert.Equal(t, pathKey, pathKey2, "Path transformation should be consistent")
}

func TestDefaultPathTransformFunc(t *testing.T) {
	key := "test_key"
	pathKey := DefaultPathTransformFunc(key)
	
	assert.Equal(t, key, pathKey.PathName, "Default path name should equal key")
	assert.Equal(t, key, pathKey.Filename, "Default filename should equal key")
}

func TestPathKeyMethods(t *testing.T) {
	pathKey := PathKey{
		PathName: "dir1/dir2/dir3",
		Filename: "testfile.txt",
	}
	
	assert.Equal(t, "dir1", pathKey.FirstPathName(), "First path name should be correct")
	assert.Equal(t, "dir1/dir2/dir3/testfile.txt", pathKey.FullPath(), "Full path should be correct")
	
	// Test empty path
	emptyPathKey := PathKey{
		PathName: "",
		Filename: "testfile.txt",
	}
	assert.Equal(t, "", emptyPathKey.FirstPathName(), "Empty path should return empty string")
}

func TestNewStore(t *testing.T) {
	// Test with default options
	store1 := NewStore(StoreOpts{})
	assert.Equal(t, defaultRootFolderName, store1.Root, "Default root should be used")
	assert.NotNil(t, store1.PathTransformFunc, "Default path transform func should be set")
	
	// Test with custom options
	customRoot := "custom_root"
	store2 := NewStore(StoreOpts{
		Root:              customRoot,
		PathTransformFunc: CASPathTransformFunc,
	})
	assert.Equal(t, customRoot, store2.Root, "Custom root should be used")
}

func TestStoreWriteAndRead(t *testing.T) {
	store := NewStore(StoreOpts{
		Root:              "test_store",
		PathTransformFunc: CASPathTransformFunc,
	})
	
	// Clean up after test
	defer func() {
		store.Clear()
	}()
	
	id := "test_id"
	key := "test_key"
	data := []byte("This is test data for the store")
	
	// Test write
	n, err := store.Write(id, key, bytes.NewReader(data))
	assert.NoError(t, err, "Write should not error")
	assert.Equal(t, int64(len(data)), n, "Written bytes should match data length")
	
	// Test has
	assert.True(t, store.Has(id, key), "Store should have the written file")
	
	// Test read
	size, reader, err := store.Read(id, key)
	assert.NoError(t, err, "Read should not error")
	assert.Equal(t, int64(len(data)), size, "Size should match written data")
	
	readData, err := io.ReadAll(reader)
	assert.NoError(t, err, "Reading data should not error")
	assert.Equal(t, data, readData, "Read data should match written data")
	
	reader.Close()
}

func TestStoreWriteDecryptAndRead(t *testing.T) {
	store := NewStore(StoreOpts{
		Root:              "test_store_decrypt",
		PathTransformFunc: CASPathTransformFunc,
	})
	
	// Clean up after test
	defer func() {
		store.Clear()
	}()
	
	id := "test_id"
	key := "test_key"
	encKey := newEncryptionKey()
	data := []byte("This is test data for encryption")
	
	// First encrypt the data
	var encrypted bytes.Buffer
	_, err := copyEncrypt(encKey, bytes.NewReader(data), &encrypted)
	assert.NoError(t, err, "Encryption should not error")
	
	// Test write decrypt
	n, err := store.WriteDecrypt(encKey, id, key, bytes.NewReader(encrypted.Bytes()))
	assert.NoError(t, err, "Write decrypt should not error")
	assert.Greater(t, n, int64(0), "Should write some bytes")
	
	// Test read the decrypted data
	_, reader, err := store.Read(id, key)
	assert.NoError(t, err, "Read should not error")
	
	readData, err := io.ReadAll(reader)
	assert.NoError(t, err, "Reading data should not error")
	assert.Equal(t, data, readData, "Read data should match original data")
	
	reader.Close()
}

func TestStoreDelete(t *testing.T) {
	store := NewStore(StoreOpts{
		Root:              "test_store_delete",
		PathTransformFunc: CASPathTransformFunc,
	})
	
	// Clean up after test
	defer func() {
		store.Clear()
	}()
	
	id := "test_id"
	key := "test_key"
	data := []byte("This is test data for deletion")
	
	// Write data
	_, err := store.Write(id, key, bytes.NewReader(data))
	assert.NoError(t, err, "Write should not error")
	assert.True(t, store.Has(id, key), "Store should have the file")
	
	// Delete data
	err = store.Delete(id, key)
	assert.NoError(t, err, "Delete should not error")
	assert.False(t, store.Has(id, key), "Store should not have the file after deletion")
}

func TestStoreSize(t *testing.T) {
	store := NewStore(StoreOpts{
		Root:              "test_store_size",
		PathTransformFunc: CASPathTransformFunc,
	})
	
	// Clean up after test
	defer func() {
		store.Clear()
	}()
	
	id := "test_id"
	key := "test_key"
	data := []byte("This is test data for size testing")
	
	// Write data
	_, err := store.Write(id, key, bytes.NewReader(data))
	assert.NoError(t, err, "Write should not error")
	
	// Test size
	size, err := store.Size(id, key)
	assert.NoError(t, err, "Size should not error")
	assert.Equal(t, int64(len(data)), size, "Size should match data length")
}

func TestStoreNonExistentFile(t *testing.T) {
	store := NewStore(StoreOpts{
		Root:              "test_store_nonexistent",
		PathTransformFunc: CASPathTransformFunc,
	})
	
	// Clean up after test
	defer func() {
		store.Clear()
	}()
	
	id := "test_id"
	key := "nonexistent_key"
	
	// Test has on non-existent file
	assert.False(t, store.Has(id, key), "Store should not have non-existent file")
	
	// Test read on non-existent file
	_, _, err := store.Read(id, key)
	assert.Error(t, err, "Read should error for non-existent file")
	
	// Test size on non-existent file
	_, err = store.Size(id, key)
	assert.Error(t, err, "Size should error for non-existent file")
	
	// Test delete on non-existent file
	err = store.Delete(id, key)
	assert.Error(t, err, "Delete should error for non-existent file")
}

func TestStoreLargeFile(t *testing.T) {
	store := NewStore(StoreOpts{
		Root:              "test_store_large",
		PathTransformFunc: CASPathTransformFunc,
	})
	
	// Clean up after test
	defer func() {
		store.Clear()
	}()
	
	id := "test_id"
	key := "large_file_key"
	
	// Create large data (1MB)
	largeData := make([]byte, 1024*1024)
	rand.Read(largeData)
	
	// Write large data
	n, err := store.Write(id, key, bytes.NewReader(largeData))
	assert.NoError(t, err, "Write should not error for large file")
	assert.Equal(t, int64(len(largeData)), n, "Written bytes should match data length")
	
	// Read large data
	_, reader, err := store.Read(id, key)
	assert.NoError(t, err, "Read should not error for large file")
	
	readData, err := io.ReadAll(reader)
	assert.NoError(t, err, "Reading large data should not error")
	assert.Equal(t, largeData, readData, "Read data should match written data")
	
	reader.Close()
}