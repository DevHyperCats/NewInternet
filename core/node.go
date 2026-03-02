package main
import (
    "context"
    "sync"
    "time"
    
    "p2p-lib/crypto"
    "p2p-lib/routing"
    "p2p-lib/network"
    "p2p-lib/storage"
    "p2p-lib/messaging"
)

type PeerID [20]byte

type Message struct {
    Type    network.MessageType
    From    PeerID
    To      PeerID
    Payload []byte
    Relay   bool
    TTL     int
}

type MessageHandler func(msg *Message) error

type Node struct {
    ID        PeerID
    Cloud     *CloudAPI
    Messaging *MessagingAPI
    
    privateKey crypto.PrivateKey
    publicKey  crypto.PublicKey

    routing     *routing.Table
    network     *network.Manager
    storage     *storage.Storage
    replication *replicationManager
    relay       *relayManager
    
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
    
    inboundChan chan *Message

    handlers    map[network.MessageType]MessageHandler
    handlersMu  sync.RWMutex
    
    config      *Config
}

type CloudAPI struct {
    node *Node
}

type CloudOptions struct {
    ReplicaCnt int
    ChunkSize  int
    Progress   func(done, total int64)
}

type FileInfo struct {
    ID        string
    Name      string
    Size      int64
    CreatedAt time.Time
    Replicas  int
}

func (c *CloudAPI) splitIntoChunks(filePath string, chunkSize int64, progress func(done, total int64)) ([][32]byte, int64, error) {
    // Открываем файл
    file, err := os.Open(filePath)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // Получаем информацию о файле
    fileInfo, err := file.Stat()
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get file info: %w", err)
    }
    totalSize := fileInfo.Size()

    // Буфер для чтения
    buffer := make([]byte, chunkSize)
    var chunkHashes [][32]byte
    var bytesProcessed int64 = 0

    for {
        bytesRead, err := file.Read(buffer)
        if err != nil && err != io.EOF {
            return nil, 0, fmt.Errorf("failed to read file: %w", err)
        }

        if bytesRead == 0 {
            break
        }

        chunkData := buffer[:bytesRead]
        chunkHash := sha256.Sum256(chunkData)

        chunk := &storage.Chunk{
            Hash: chunkHash,
            Data: chunkData,
            Size: bytesRead,
        }

        // TODO: сохранить через storage (когда будет реализован)
        // err = c.node.storage.SaveChunk(chunk)
        // if err != nil {
        //     return nil, 0, fmt.Errorf("failed to save chunk: %w", err)
        // }

        // Добавляем хеш в список
        chunkHashes = append(chunkHashes, chunkHash)
        bytesProcessed += int64(bytesRead)
        if progress != nil {
            progress(bytesProcessed, totalSize)
		}
        select {
        case <-c.node.ctx.Done():
            return nil, 0, fmt.Errorf("operation cancelled")
        default:
        }
        if err == io.EOF {
            break
        }
    }

    return chunkHashes, totalSize, nil
}

func (cloudApi *CloudAPI) StoreFile(p path, opt)