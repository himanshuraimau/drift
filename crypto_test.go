package main

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()
	
	assert.NotEqual(t, id1, id2, "Generated IDs should be unique")
	assert.Greater(t, len(id1), 0, "ID should not be empty")
}

func TestHashKey(t *testing.T) {
	key := "test_key"
	hash1 := hashKey(key)
	hash2 := hashKey(key)
	
	assert.Equal(t, hash1, hash2, "Hash should be consistent")
	assert.Equal(t, 32, len(hash1), "MD5 hash should be 32 characters")
}

func TestNewEncryptionKey(t *testing.T) {
	key1 := newEncryptionKey()
	key2 := newEncryptionKey()
	
	assert.NotEqual(t, key1, key2, "Encryption keys should be unique")
	assert.Equal(t, 32, len(key1), "Encryption key should be 32 bytes")
	assert.Equal(t, 32, len(key2), "Encryption key should be 32 bytes")
}

func TestValidateKey(t *testing.T) {
	validKey := make([]byte, 32)
	rand.Read(validKey)
	
	invalidKey := make([]byte, 16)
	rand.Read(invalidKey)
	
	assert.NoError(t, validateKey(validKey), "Valid key should pass validation")
	assert.Error(t, validateKey(invalidKey), "Invalid key should fail validation")
}

func TestCopyEncryptDecrypt(t *testing.T) {
	key := newEncryptionKey()
	plaintext := []byte("This is a test message for encryption and decryption")
	
	// Test encryption
	var encrypted bytes.Buffer
	src := bytes.NewReader(plaintext)
	
	n, err := copyEncrypt(key, src, &encrypted)
	assert.NoError(t, err, "Encryption should not error")
	assert.Greater(t, n, len(plaintext), "Encrypted data should be larger due to IV")
	
	// Test decryption
	var decrypted bytes.Buffer
	encryptedReader := bytes.NewReader(encrypted.Bytes())
	
	n, err = copyDecrypt(key, encryptedReader, &decrypted)
	assert.NoError(t, err, "Decryption should not error")
	assert.Equal(t, plaintext, decrypted.Bytes(), "Decrypted data should match original")
}

func TestEncryptDecryptLargeData(t *testing.T) {
	key := newEncryptionKey()
	
	// Create large test data
	largeData := make([]byte, 1024*1024) // 1MB
	rand.Read(largeData)
	
	// Encrypt
	var encrypted bytes.Buffer
	src := bytes.NewReader(largeData)
	
	_, err := copyEncrypt(key, src, &encrypted)
	assert.NoError(t, err, "Encryption of large data should not error")
	
	// Decrypt
	var decrypted bytes.Buffer
	encryptedReader := bytes.NewReader(encrypted.Bytes())
	
	_, err = copyDecrypt(key, encryptedReader, &decrypted)
	assert.NoError(t, err, "Decryption of large data should not error")
	assert.Equal(t, largeData, decrypted.Bytes(), "Decrypted large data should match original")
}

func TestEncryptDecryptWithDifferentKeys(t *testing.T) {
	key1 := newEncryptionKey()
	key2 := newEncryptionKey()
	plaintext := []byte("This is a test message")
	
	// Encrypt with key1
	var encrypted bytes.Buffer
	src := bytes.NewReader(plaintext)
	
	_, err := copyEncrypt(key1, src, &encrypted)
	assert.NoError(t, err, "Encryption should not error")
	
	// Try to decrypt with key2 (should produce garbage)
	var decrypted bytes.Buffer
	encryptedReader := bytes.NewReader(encrypted.Bytes())
	
	_, err = copyDecrypt(key2, encryptedReader, &decrypted)
	assert.NoError(t, err, "Decryption should not error technically")
	assert.NotEqual(t, plaintext, decrypted.Bytes(), "Decrypted data should not match original with wrong key")
}

func TestHashSHA1(t *testing.T) {
	data := []byte("test data")
	hash1 := hashSHA1(data)
	hash2 := hashSHA1(data)
	
	assert.Equal(t, hash1, hash2, "SHA1 hash should be consistent")
	assert.Equal(t, 40, len(hash1), "SHA1 hash should be 40 characters")
	
	// Test with different data
	differentData := []byte("different test data")
	hash3 := hashSHA1(differentData)
	assert.NotEqual(t, hash1, hash3, "Different data should produce different hashes")
}