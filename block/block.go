package block

import "time"

type Block struct {
	Index     int64
	TimeStamp string
	PrevHash  string
	Hash      string
	Data      string
	Nonce     int64
}

func newBlock(string data) *Block {
	return &Block{
		Index:     0,
		TimeStamp: time.Now().String(),
		PrevHash: "",
		Hash: "",
		Data: "",
	
	}
}

func CalculateHash(h string) string {

}