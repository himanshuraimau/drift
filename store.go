package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "driftnetwork"

// PathKey represents a path key for content-addressable storage
type PathKey struct {
	PathName string
	Filename string
}

// FirstPathName returns the first path component
func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// FullPath returns the full path including filename
func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

// PathTransformFunc defines how to transform a key into a path
type PathTransformFunc func(string) PathKey

// CASPathTransformFunc creates a content-addressable path from a key
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

// DefaultPathTransformFunc is the default path transformation function
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		Filename: key,
	}
}

// StoreOpts contains options for the store
type StoreOpts struct {
	// Root is the folder name containing all files
	Root              string
	PathTransformFunc PathTransformFunc
}

// Store represents the file storage system
type Store struct {
	StoreOpts
}

// NewStore creates a new store instance
func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

// Has checks if a file exists in the store
func (s *Store) Has(id string, key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

// Clear removes all files from the store
func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

// Delete removes a file from the store
func (s *Store) Delete(id string, key string) error {
	pathKey := s.PathTransformFunc(key)
	
	// Check if file exists before deleting
	if !s.Has(id, key) {
		return fmt.Errorf("file with key %s does not exist", key)
	}

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()

	firstPathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FirstPathName())

	return os.RemoveAll(firstPathNameWithRoot)
}

// Write writes data to the store
func (s *Store) Write(id string, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)
}

// WriteDecrypt writes encrypted data to the store with decryption
func (s *Store) WriteDecrypt(encKey []byte, id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n, err := copyDecrypt(encKey, r, f)
	return int64(n), err
}

// openFileForWriting opens a file for writing, creating directories as needed
func (s *Store) openFileForWriting(id string, key string) (*os.File, error) {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.PathName)
	
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil, err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	return os.Create(fullPathWithRoot)
}

// writeStream writes data to a file stream
func (s *Store) writeStream(id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return io.Copy(f, r)
}

// Read reads data from the store
func (s *Store) Read(id string, key string) (int64, io.ReadCloser, error) {
	return s.readStream(id, key)
}

// readStream reads data from a file stream
func (s *Store) readStream(id string, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return 0, nil, err
	}

	return fi.Size(), file, nil
}

// Size returns the size of a file in the store
func (s *Store) Size(id string, key string) (int64, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	fi, err := os.Stat(fullPathWithRoot)
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}