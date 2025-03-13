package node

import (
	"context"
	"sync"

	"github.com/greg/hashchain/pkg/blockchain"
	pb "github.com/greg/hashchain/proto"
)

type NodeService struct {
	pb.UnimplementedNodeServiceServer
	chain    *blockchain.Blockchain
	peers    map[string]pb.NodeServiceClient
	peersMu  sync.RWMutex
}

func NewNodeService() *NodeService {
	return &NodeService{
		chain: blockchain.NewBlockchain(),
		peers: make(map[string]pb.NodeServiceClient),
	}
}

func (s *NodeService) BroadcastBlock(ctx context.Context, block *pb.Block) (*pb.BroadcastResponse, error) {
	// Convert protobuf block to internal block type
	internalBlock := &blockchain.Block{Block: block}
	
	if err := s.chain.AddBlock(internalBlock); err != nil {
		return &pb.BroadcastResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Broadcast to peers
	s.broadcastToPeers(ctx, block)

	return &pb.BroadcastResponse{
		Success: true,
		Message: "Block accepted",
	}, nil
}

func (s *NodeService) BroadcastTransaction(ctx context.Context, tx *pb.Transaction) (*pb.BroadcastResponse, error) {
	if err := s.chain.AddTransaction(tx); err != nil {
		return &pb.BroadcastResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Broadcast to peers
	s.broadcastTransactionToPeers(ctx, tx)

	return &pb.BroadcastResponse{
		Success: true,
		Message: "Transaction accepted",
	}, nil
}

func (s *NodeService) GetBlock(ctx context.Context, req *pb.GetBlockRequest) (*pb.Block, error) {
	block := s.chain.GetBlock(req.BlockHash)
	if block == nil {
		return nil, nil
	}
	return block.Block, nil
}

func (s *NodeService) GetLatestBlock(ctx context.Context, empty *pb.Empty) (*pb.Block, error) {
	block := s.chain.GetLatestBlock()
	if block == nil {
		return &pb.Block{}, nil
	}
	return block.Block, nil
}

func (s *NodeService) AddPeer(address string, client pb.NodeServiceClient) {
	s.peersMu.Lock()
	defer s.peersMu.Unlock()
	s.peers[address] = client
}

func (s *NodeService) RemovePeer(address string) {
	s.peersMu.Lock()
	defer s.peersMu.Unlock()
	delete(s.peers, address)
}

func (s *NodeService) broadcastToPeers(ctx context.Context, block *pb.Block) {
	s.peersMu.RLock()
	defer s.peersMu.RUnlock()

	for _, peer := range s.peers {
		go func(p pb.NodeServiceClient) {
			_, _ = p.BroadcastBlock(ctx, block)
		}(peer)
	}
}

func (s *NodeService) broadcastTransactionToPeers(ctx context.Context, tx *pb.Transaction) {
	s.peersMu.RLock()
	defer s.peersMu.RUnlock()

	for _, peer := range s.peers {
		go func(p pb.NodeServiceClient) {
			_, _ = p.BroadcastTransaction(ctx, tx)
		}(peer)
	}
}
