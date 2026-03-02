type FileManifest struct {
    Version     string
    FileName    string
    FileSize    int64 
    ChunkSize   int 
    ChunkHashes []string
    ReplicaCnt  int     
    CreatedAt   int64   
}
func (m *Manifest) RootHash() string {
    data, _ := json.Marshal(m)
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}
