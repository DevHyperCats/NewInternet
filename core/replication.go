package core

import (
    "time"
    "p2p-lib/storage"
    "p2p-lib/core/routing"
)

type ReplicationManager struct {
    node        *Node
    storage     *storage.Storage
    routing     *routing.Table
    stopChan    chan struct{}
    interval    time.Duration
}

func NewReplicationManager(node *Node) *ReplicationManager {
    return &ReplicationManager{
        node:     node,
        storage:  node.storage,
        routing:  node.routing,
        interval: 5 * time.Minute,
    }
}

func (r *ReplicationManager) Start() {
    r.stopChan = make(chan struct{})
    go r.loop()
}

func (r *ReplicationManager) Stop() {
    close(r.stopChan)
}

func (r *ReplicationManager) loop() {
    ticker := time.NewTicker(r.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            r.checkAllChunks()
        case <-r.stopChan:
            return
        }
    }
}
и
func (r *ReplicationManager) checkAllChunks() {
    localChunks := r.storage.ListChunks()
    
    for _, chunkHash := range localChunks {
        r.checkChunkReplication(chunkHash)
    }
}
func (r *ReplicationManager) checkChunkReplication(chunkHash string) {
    manifests := r.storage.FindManifestsByChunk(chunkHash)
    if len(manifests) == 0 {
        return
    }

    targetReplicas := 0
    for _, m := range manifests {
        if m.ReplicaCnt > targetReplicas {
            targetReplicas = m.ReplicaCnt
        }
    }

    peers, err := r.findPeersWithChunk(chunkHash)
    if err != nil {
        return
    }
    
    currentReplicas := len(peers)

    if currentReplicas < targetReplicas {
        r.replicateChunk(chunkHash, peers, targetReplicas-currentReplicas)
    }

    if currentReplicas > targetReplicas {
        r.pruneChunk(chunkHash, peers, currentReplicas-targetReplicas)
    }
}

func (r *ReplicationManager) replicateChunk(chunkHash string, existingPeers []Peer, need int) {
    candidates := r.routing.FindNearestWithoutChunk(chunkHash, need, existingPeers)

    data := r.storage.GetChunk(chunkHash)
    if data == nil {
        return
    }

    for _, peer := range candidates {
        go func(p Peer) {
            err := r.node.network.SendChunk(p, chunkHash, data)
            if err == nil {
                r.node.log.Printf("Реплицирован чанк %s на пир %s", chunkHash, p.ID)
            }
        }(peer)
    }
}
func (r *ReplicationManager) pruneChunk(chunkHash string, existingPeers []Peer, extra int) {
    sorted := r.sortPeersByDistance(chunkHash, existingPeers)

    toDelete := sorted[len(sorted)-extra:]
    
    for _, peer := range toDelete {
        go func(p Peer) {
            err := r.node.network.RequestDelete(p, chunkHash)
            if err == nil {
                r.node.log.Printf("Удалена лишняя копия чанка %s с пира %s", chunkHash, p.ID)
            }
        }(peer)
    }
}

func (r *ReplicationManager) findPeersWithChunk(chunkHash string) ([]Peer, error) {
    closest := r.routing.FindNearest(chunkHash, 20)
    
    var result []Peer
    for _, peer := range closest {
        has, err := r.node.network.HasChunk(peer, chunkHash)
        if err == nil && has {
            result = append(result, peer)
        }
    }
    return result, nil
}