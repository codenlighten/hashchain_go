package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"

	pb "github.com/greg/hashchain/proto"
)

// Block wraps the protobuf Block message with additional functionality
type Block struct {
	*pb.Block
	hash []byte
}

// NewBlock creates a new block with the given previous block hash and transactions
func NewBlock(prevBlockHash []byte, transactions []*pb.Transaction) *Block {
	block := &Block{
		Block: &pb.Block{
			Version:       1,
			PrevBlockHash: prevBlockHash,
			Timestamp:     time.Now().Unix(),
			Bits:         545259519, // Easier difficulty target for development (2 leading zeros)
			Nonce:        0,
			Transactions: transactions,
		},
	}
	block.CalculateMerkleRoot()
	return block
}

// CalculateMerkleRoot calculates and sets the merkle root of the block's transactions
func (b *Block) CalculateMerkleRoot() {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, calculateTxHash(tx))
	}
	b.MerkleRoot = calculateMerkleRoot(txHashes)
}

// Mine performs proof of work on the block
func (b *Block) Mine() {
	target := getDifficultyTarget(b.Bits)
	for {
		hash := b.calculateHash()
		if isHashBelowTarget(hash, target) {
			b.hash = hash
			break
		}
		b.Nonce++
	}
}

// Hash returns the block's hash, calculating it if necessary
func (b *Block) Hash() []byte {
	if b.hash == nil {
		b.hash = b.calculateHash()
	}
	return b.hash
}

func (b *Block) calculateHash() []byte {
	header := append([]byte{}, b.PrevBlockHash...)
	header = append(header, b.MerkleRoot...)
	header = append(header, []byte(fmt.Sprintf("%d", b.Timestamp))...)
	header = append(header, []byte(fmt.Sprintf("%d", b.Bits))...)
	header = append(header, []byte(fmt.Sprintf("%d", b.Nonce))...)

	hash := sha256.Sum256(header)
	hash = sha256.Sum256(hash[:]) // Double SHA256 like Bitcoin
	return hash[:]
}

func calculateTxHash(tx *pb.Transaction) []byte {
	// Simplified transaction hashing
	hash := sha256.Sum256(tx.Txid)
	return hash[:]
}

func calculateMerkleRoot(hashes [][]byte) []byte {
	if len(hashes) == 0 {
		return make([]byte, 32)
	}
	if len(hashes) == 1 {
		return hashes[0]
	}

	var newHashes [][]byte
	for i := 0; i < len(hashes)-1; i += 2 {
		hash := sha256.Sum256(append(hashes[i], hashes[i+1]...))
		newHashes = append(newHashes, hash[:])
	}
	if len(hashes)%2 == 1 {
		hash := sha256.Sum256(append(hashes[len(hashes)-1], hashes[len(hashes)-1]...))
		newHashes = append(newHashes, hash[:])
	}
	return calculateMerkleRoot(newHashes)
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
