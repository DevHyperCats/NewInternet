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
