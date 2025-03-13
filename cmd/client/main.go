package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/greg/hashchain/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddr := flag.String("server", "localhost:50051", "The server address in the format of host:port")
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)
	ctx := context.Background()

	// Create a sample transaction
	tx := createSampleTransaction()
	log.Printf("Broadcasting transaction: %x", tx.Txid)
	
	resp, err := client.BroadcastTransaction(ctx, tx)
	if err != nil {
		log.Fatalf("Failed to broadcast transaction: %v", err)
	}
	log.Printf("Transaction broadcast response: %v", resp.Message)

	// Get the latest block to get its hash
	latestBlock, err := client.GetLatestBlock(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}

	// Create and mine a new block
	block := createSampleBlock([]*pb.Transaction{tx}, latestBlock)
	log.Printf("Mining block...")
	mineBlock(block)
	blockHash := calculateBlockHash(block)
	log.Printf("Broadcasting block with hash: %x", blockHash)
	
	resp, err = client.BroadcastBlock(ctx, block)
	if err != nil {
		log.Fatalf("Failed to broadcast block: %v", err)
	}
	log.Printf("Block broadcast response: %v", resp.Message)
}

func createSampleTransaction() *pb.Transaction {
	tx := &pb.Transaction{
		Version: 1,
		Inputs: []*pb.Input{
			{
				PrevTxid:     make([]byte, 32), // Genesis transaction
				PrevOutIndex: 0,
				ScriptSig:    []byte("sample signature"),
				Sequence:     0xffffffff,
			},
		},
		Outputs: []*pb.Output{
			{
				Value:        50 * 100000000, // 50 coins
				ScriptPubkey: []byte("sample public key"),
			},
		},
		Locktime: 0,
	}

	// Calculate transaction hash
	txHash := sha256.Sum256([]byte(fmt.Sprintf("%v", tx)))
	tx.Txid = txHash[:]
	return tx
}

func createSampleBlock(transactions []*pb.Transaction, prevBlock *pb.Block) *pb.Block {
	var prevBlockHash []byte
	if prevBlock != nil {
		prevBlockHash = calculateBlockHash(prevBlock)
	} else {
		prevBlockHash = make([]byte, 32) // Genesis block
	}
	merkleRoot := calculateMerkleRoot(transactions)

	return &pb.Block{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		MerkleRoot:    merkleRoot,
		Timestamp:     time.Now().Unix(),
		Bits:         545259519, // Easier difficulty target for development (2 leading zeros)
		Nonce:        0,
		Transactions: transactions,
	}
}

func mineBlock(block *pb.Block) {
	target := getDifficultyTarget(block.Bits)
	for {
		hash := calculateBlockHash(block)
		if isHashBelowTarget(hash, target) {
			break
		}
		block.Nonce++
		if block.Nonce%100000 == 0 {
			log.Printf("Mining... (nonce: %d)", block.Nonce)
		}
	}
}

func getDifficultyTarget(bits int32) []byte {
	// Convert difficulty bits to target
	exp := bits >> 24
	mant := bits & 0xffffff
	target := make([]byte, 32)
	
	// Calculate the target based on the difficulty bits
	target[32-exp] = byte(mant >> 16)
	target[32-exp+1] = byte(mant >> 8)
	target[32-exp+2] = byte(mant)
	
	return target
}

func isHashBelowTarget(hash, target []byte) bool {
	for i := 0; i < len(hash); i++ {
		if hash[i] > target[i] {
			return false
		}
		if hash[i] < target[i] {
			return true
		}
	}
	return true
}

func calculateMerkleRoot(transactions []*pb.Transaction) []byte {
	if len(transactions) == 0 {
		return make([]byte, 32)
	}

	var hashes [][]byte
	for _, tx := range transactions {
		hashes = append(hashes, tx.Txid)
	}

	for len(hashes) > 1 {
		if len(hashes)%2 == 1 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}

		var nextLevel [][]byte
		for i := 0; i < len(hashes); i += 2 {
			hash := sha256.Sum256(append(hashes[i], hashes[i+1]...))
			hash = sha256.Sum256(hash[:])
			nextLevel = append(nextLevel, hash[:])
		}
		hashes = nextLevel
	}

	return hashes[0]
}

func calculateBlockHash(block *pb.Block) []byte {
	header := append([]byte{}, block.PrevBlockHash...)
	header = append(header, block.MerkleRoot...)
	header = append(header, []byte(fmt.Sprintf("%d", block.Timestamp))...)
	header = append(header, []byte(fmt.Sprintf("%d", block.Bits))...)
	header = append(header, []byte(fmt.Sprintf("%d", block.Nonce))...)

	hash := sha256.Sum256(header)
	hash = sha256.Sum256(hash[:])
	return hash[:]
}
