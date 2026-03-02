package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "os"
    "path/filepath"
	"bytes"
)
type Chunk struct {
    Hash [32]byte
    Data []byte
    Size int
}
func NewChunk(data []byte) *Chunk {
    hash := sha256.Sum256(data)
    size := len(data)
    
    chunk := &Chunk{
        Hash: hash,
        Data: data,
        Size: size,
    }
    return chunk
}
func (c *Chunk) HashHex() string {
	return hex.EncodeToString(c.Hash[:])
}
func (c *Chunk) Validate() bool {
    hash := sha256.Sum256(c.Data)
	return hash == c.Hash
}
func (c *Chunk) Save(baseDir string) error {
    hashHex := c.HashHex()
    subDir := hashHex[:2]
    chunkPath := filepath.Join(baseDir, "chunks", subDir, hashHex)
    chunkDir := filepath.Join(baseDir, "chunks", subDir)
    err := os.MkdirAll(chunkDir, 0755)
    if err != nil {
        return fmt.Errorf("failed to create chunk directory: %w", err)
    }
    file, err := os.Create(chunkPath)
    if err != nil {
        return fmt.Errorf("failed to create chunk file: %w", err)
    }
    defer file.Close()
    _, err = file.Write(c.Data)
    if err != nil {
        return fmt.Errorf("failed to write chunk data: %w", err)
    }
    
    return nil
}

func LoadChunk(baseDir string, hash [32]byte) (*Chunk, error) {
    hashHex := hex.EncodeToString(hash[:])
    subDir := hashHex[:2]
    chunkPath := filepath.Join(baseDir, "chunks", subDir, hashHex)
    data, err := os.ReadFile(chunkPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read chunk file: %w", err)
    }
    chunk := NewChunk(data)
    if chunk.Hash != hash {
        return nil, fmt.Errorf("chunk corrupted: hash mismatch")
    }
    
    return chunk, nil
}

func ChunkExists(baseDir string, hash [32]byte) bool {
    hashHex := hex.EncodeToString(hash[:])
    subDir := hashHex[:2]
    chunkPath := filepath.Join(baseDir, "chunks", subDir, hashHex)
    _, err := os.Stat(chunkPath)
    if err == nil {
        return true
    }
    if os.IsNotExist(err) {
        return false 
    }
    return false
}