package blockchain

import (
	"errors"
	"fmt"
	"sync"

	pb "github.com/greg/hashchain/proto"
)

type Blockchain struct {
	blocks     []*Block
	mempool    []*pb.Transaction
	mu         sync.RWMutex
}

func NewBlockchain() *Blockchain {
	chain := &Blockchain{
		blocks:  make([]*Block, 0),
		mempool: make([]*pb.Transaction, 0),
	}
	chain.addGenesisBlock()
	return chain
}

func (bc *Blockchain) addGenesisBlock() {
	genesisBlock := NewBlock(make([]byte, 32), nil)
	genesisBlock.Mine()
	bc.blocks = append(bc.blocks, genesisBlock)
}

func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// For non-genesis blocks, validate against previous block
	if len(bc.blocks) > 0 {
		lastBlock := bc.blocks[len(bc.blocks)-1]
		if err := bc.validateNextBlock(lastBlock, block); err != nil {
			return fmt.Errorf("invalid block: %v", err)
		}
	} else {
		// For genesis block, just verify proof of work
		target := getDifficultyTarget(block.Bits)
		if !isHashBelowTarget(block.Hash(), target) {
			return errors.New("invalid block: proof of work verification failed")
		}
	}

	bc.blocks = append(bc.blocks, block)
	bc.removeTransactionsFromMempool(block.Transactions)
	return nil
}

func (bc *Blockchain) AddTransaction(tx *pb.Transaction) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Basic transaction validation would go here
	bc.mempool = append(bc.mempool, tx)
	return nil
}

func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.blocks) == 0 {
		return nil
	}
	return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) GetBlock(hash []byte) *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	for _, block := range bc.blocks {
		if string(block.Hash()) == string(hash) {
			return block
		}
	}
	return nil
}

func (bc *Blockchain) GetPendingTransactions() []*pb.Transaction {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return append([]*pb.Transaction{}, bc.mempool...)
}

func (bc *Blockchain) removeTransactionsFromMempool(transactions []*pb.Transaction) {
	txMap := make(map[string]struct{})
	for _, tx := range transactions {
		txMap[string(tx.Txid)] = struct{}{}
	}

	newMempool := make([]*pb.Transaction, 0)
	for _, tx := range bc.mempool {
		if _, exists := txMap[string(tx.Txid)]; !exists {
			newMempool = append(newMempool, tx)
		}
	}
	bc.mempool = newMempool
}

func (bc *Blockchain) validateNextBlock(prevBlock, newBlock *Block) error {
	// Check if the new block points to the previous block
	if string(newBlock.PrevBlockHash) != string(prevBlock.Hash()) {
		return errors.New("previous block hash mismatch")
	}

	// Check if the timestamp is greater than the previous block
	if newBlock.Timestamp <= prevBlock.Timestamp {
		return errors.New("timestamp must be greater than previous block")
	}

	// Verify the proof of work
	target := getDifficultyTarget(newBlock.Bits)
	if !isHashBelowTarget(newBlock.Hash(), target) {
		return errors.New("proof of work verification failed")
	}

	return nil
}
