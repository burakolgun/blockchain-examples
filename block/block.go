package block

import (
	"../hash"
	"../transaction"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Index      int64
	TimeStamp  int64
	PrevHash   string
	Hash       string
	Data       []*transaction.Transaction
	Nonce      int64
	Difficulty int
}

func CreateGenesisBlock(tx []*transaction.Transaction) *Block {
	now := time.Now().UnixMilli()
	data := "Genesis Block"
	return &Block{
		Index:      0,
		TimeStamp:  now,
		PrevHash:   "GENESIS BLOCK",
		Hash:       hash.CalculateHash(fmt.Sprintf("%d-%s-%s-%s-%d-%d", 0, now, "", data, 0, 0)),
		Data:       tx,
		Nonce:      0,
		Difficulty: 1,
	}
}

func (b *Block) Print() {
	fmt.Printf("index           %d\n", b.Index)
	fmt.Printf("timestamp       %d\n", b.TimeStamp)
	fmt.Printf("previous_hash   %x\n", b.PrevHash)
	fmt.Printf("nonce           %d\n", b.Nonce)
	fmt.Printf("difficulty      %s\n", b.Difficulty)
	for _, t := range b.Data {
		t.Print()
	}
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64                      `json:"timestamp"`
		Nonce        int64                      `json:"nonce"`
		PreviousHash string                     `json:"previous_hash"`
		Transactions []*transaction.Transaction `json:"transactions"`
		Difficulty   int                        `json:"difficulty"`
		Hash         string                     `json:"hash"`
	}{
		Timestamp:    b.TimeStamp,
		Nonce:        b.Nonce,
		PreviousHash: b.PrevHash,
		Transactions: b.Data,
		Difficulty:   b.Difficulty,
		Hash:         b.Hash,
	})
}
