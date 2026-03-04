package storage

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
)
type Cloud struct {
    baseDir     string
    maxSize     int64
    currentSize int64
    mu          sync.RWMutex
    
    manifests   map[string]*cachedManifest
    manifestsMu sync.RWMutex
}


type cachedManifest struct {
    manifest *FileManifest
    loadedAt time.Time
    accessCount int
}

type Config struct {
    BaseDir     string
    MaxSize     int64
    ChunkSize   int
    ReplicaCnt  int
}

func NewCloud(cfg *Config) (*Cloud, error) {
    if cfg.BaseDir == "" {
		return nil, fmt.Errorf("base directory cannot be empty")
	}
    err := os.MkdirAll(filepath.Join(cfg.BaseDir, "chunks"), 0755)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(filepath.Join(cfg.BaseDir, "manifests"), 0755)
	if err != nil {
		return nil, err
	}
	cloud := Cloud{
		baseDir: 	 cfg.BaseDir,
		maxSize: 	 cfg.MaxSize,
		currentSize: 0,
		mu: 		 sync.RWMutex{},
		manifests: 	 make(map[string]*cachedManifest),
		manifestsMu: sync.RWMutex{},
	}
	return &cloud, nil
}
func (c *Cloud) SaveChunk(data []byte) ([32]byte, error) {
	if c.maxSize > 0 {
	    c.mu.RLock()
	    current := c.currentSize
	    c.mu.RUnlock()

	    if current+int64(len(data)) > c.maxSize {
	        return [32]byte{}, fmt.Errorf("storage size limit exceeded")
	    }
	}
	chunk := NewChunk(data)
	err := chunk.Save(c.baseDir)
	if err != nil {
	    return [32]byte{}, fmt.Errorf("failed to save chunk: %w", err)
	}
	c.mu.Lock()
	c.currentSize += int64(len(data))
	c.mu.Unlock()
	return chunk.Hash, nil
}
func (c *Cloud) GetChunk(hash [32]byte) ([]byte, error) {
	chunk, err := LoadChunk(c.baseDir, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to load chunk: %w", err)
	}
	return chunk.Data, nil
}
