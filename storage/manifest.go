package storage

import (
    "encoding/json"
    "crypto/sha256"
    "encoding/hex"
    "time"
    "fmt"
)

const ManifestVersion = "1.0"
type FileManifest struct {
    Version     string   `json:"version"`
    FileName    string   `json:"file_name"`
    FileSize    int64    `json:"file_size"`
    ChunkSize   int      `json:"chunk_size"`
    ChunkHashes []string `json:"chunk_hashes"`
    ReplicaCnt  int      `json:"replica_cnt"`
    CreatedAt   int64    `json:"created_at"`
}
func (m *FileManifest) RootHash() string {
    data, _ := json.Marshal(m)
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}
func NewManifest(fileName string, fileSize int64, chunkSize int, chunkHashes []string, replicaCnt int) *FileManifest {
    return &FileManifest{
        Version:     ManifestVersion,
        FileName:    fileName,
        FileSize:    fileSize,
        ChunkSize:   chunkSize,
        ChunkHashes: chunkHashes,
        ReplicaCnt:  replicaCnt,
        CreatedAt:   time.Now().Unix(),
    }
}
func (m *FileManifest) Marshal() ([]byte, error) {
    return json.Marshal(m)
}
func UnmarshalManifest(data []byte) (*FileManifest, error) {
    var man FileManifest
    err := json.Unmarshal(data, &man)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
    }
    return &man, nil
}
func (m *FileManifest) Validate() error {
    if m.Version != ManifestVersion {
        return fmt.Errorf("Incorrect version")
    }
    if m.FileName == "" {
        return fmt.Errorf("Empty name")
    }
    if m.FileSize <= 0 {
        return fmt.Errorf("Empty file or invalid size")
    }
    if m.ChunkSize <= 0 {
        return fmt.Errorf("Empty chunk or invalid size")
    }
    if len(m.ChunkHashes) == 0 {
        return fmt.Errorf("No hashes")
    }
    if m.ReplicaCnt < 1 {
        return fmt.Errorf("No replicas")
    }
    return nil
}