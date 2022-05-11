package blockchain

import "C"
import (
	"../block"
	"../hash"
	"../transaction"
	"../utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const BlockChainCore = "BLOCKCHAIN_CORE"
const MiningReward = 0.1
const MinerAddress = "MINER_ADDRESS"
const ExpectedMiningProcessIntervalMillisecond = 2000
const DifficultyAdjustmentIntervalBlockCount = 20

type BlockChain struct {
	BlockList       []*block.Block
	TransactionPool []*transaction.Transaction
	Difficulty      int
	lock            sync.Mutex
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress *string  `json:"recipientBlockchainAddress"`
	SenderPublicKey            *string  `json:"senderPublicKey"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (tr TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}

	return true
}

func New() *BlockChain {
	return &BlockChain{}
}

func (c *BlockChain) GetTransactionPool() []*transaction.Transaction {
	return c.TransactionPool
}

func (c *BlockChain) Mining() {
	c.lock.Lock()
	difficulty := c.GetDifficulty()
	fmt.Println("difficulty " + "----<> " + strings.Repeat("0", difficulty))
	defer c.lock.Unlock()

	fmt.Println("mining start")

	now := time.Now().UnixMilli()
	nonce := int64(0)
	prevBlock := c.BlockList[len(c.BlockList)-1]

	h := c.CalculateHash(prevBlock.Index+1, now, prevBlock.Hash, c.TransactionPool, nonce, difficulty)
	c.AddTransaction(BlockChainCore, MinerAddress, MiningReward, nil, nil)

	for {
		if strings.HasPrefix(h, strings.Repeat("0", difficulty)) {
			fmt.Println(fmt.Sprintf("found, %s", h))
			break
		}

		nonce++
		h = c.CalculateHash(prevBlock.Index+1, now, prevBlock.Hash, c.TransactionPool, nonce, difficulty)
	}

	fmt.Println("mining end")

	c.BlockList = append(c.BlockList, &block.Block{
		Index:      prevBlock.Index + 1,
		TimeStamp:  now,
		PrevHash:   prevBlock.Hash,
		Hash:       h,
		Data:       c.TransactionPool,
		Nonce:      nonce,
		Difficulty: difficulty,
	})

	c.TransactionPool = []*transaction.Transaction{}
}

func (c *BlockChain) RecursiveMiner() {
	c.Mining()
	c.RecursiveMiner()
}

func (c *BlockChain) GetDifficulty() int {
	if c.BlockList[len(c.BlockList)-1].Index%DifficultyAdjustmentIntervalBlockCount == 0 && c.BlockList[len(c.BlockList)-1].Index != 0 {
		fmt.Printf("%s difficulty adjustment started\n", strings.Repeat("=", 25))

		return c.AdjustDifficulty()
	}

	fmt.Printf("%s difficulty adjustment conditions did not occur yet \n", strings.Repeat("=", 25))
	return c.BlockList[len(c.BlockList)-1].Difficulty
}

func (c *BlockChain) AdjustDifficulty() int {
	var expectedProcessTime int64
	prevAdjustmentBlock := c.BlockList[(len(c.BlockList))-DifficultyAdjustmentIntervalBlockCount]
	expectedProcessTime = DifficultyAdjustmentIntervalBlockCount * ExpectedMiningProcessIntervalMillisecond
	timeTaken := c.BlockList[len(c.BlockList)-1].TimeStamp - prevAdjustmentBlock.TimeStamp

	fmt.Println("prevAdjustmentBlock" + "----<> " + strconv.FormatInt(prevAdjustmentBlock.Index, 10))
	fmt.Println("prevAdjustmentBlock - Difficulty " + "----<>" + strconv.FormatInt(int64(prevAdjustmentBlock.Difficulty), 10))
	fmt.Println("lastBLock" + "----<> " + strconv.FormatInt(c.BlockList[len(c.BlockList)-1].Index, 10))
	fmt.Println("time taken" + "----<> " + strconv.FormatInt(timeTaken, 10))

	if timeTaken < expectedProcessTime/2 {
		fmt.Printf("%s difficulty adjusted as %d \n", strings.Repeat("=", 25), prevAdjustmentBlock.Difficulty+1)
		return prevAdjustmentBlock.Difficulty + 1
	} else if timeTaken > expectedProcessTime*2 {
		fmt.Printf("%s difficulty adjusted as %d \n", strings.Repeat("=", 25), prevAdjustmentBlock.Difficulty-1)
		return prevAdjustmentBlock.Difficulty - 1
	}

	fmt.Printf("%s difficulty adjustment no necessary \n", strings.Repeat("=", 25))
	return prevAdjustmentBlock.Difficulty
}

func (c *BlockChain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, signature *utils.Signature) bool {
	return c.AddTransaction(sender, recipient, value, senderPublicKey, signature)
}
func (c *BlockChain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, signature *utils.Signature) bool {
	t := transaction.New(sender, recipient, value)

	if sender == BlockChainCore {
		c.TransactionPool = append(c.TransactionPool, t)
		return true
	}

	if c.VerifyTransactionSignature(senderPublicKey, signature, t) {
		if c.CalculateTotalAmount(sender) < value {
			fmt.Println("ERROR: Not enough balance in a wallet")
			return false
		}

		c.TransactionPool = append(c.TransactionPool, t)
		return true
	}

	fmt.Println("ERROR: Verify Transaction")
	return false
}

func (c *BlockChain) CalculateHash(index int64, timeStamp int64, prevHash string, data []*transaction.Transaction, nonce int64, diff int) string {
	return hash.CalculateHash(fmt.Sprintf("%d-%s-%s-%s-%d-%s", index, timeStamp, prevHash, data, nonce, strings.Repeat("0", diff)))
}

func (c *BlockChain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, signature *utils.Signature, transaction *transaction.Transaction) bool {
	m, err := json.Marshal(transaction)

	if err != nil {
		return false
	}

	h := sha256.Sum256(m)

	if err != nil {
		return false
	}

	return ecdsa.Verify(senderPublicKey, h[:], signature.R, signature.S)
}

func (c *BlockChain) Validate() bool {
	for i := 1; i < len(c.BlockList); i++ {
		if c.BlockList[i].PrevHash != c.BlockList[i-1].Hash {
			fmt.Println("chain is not valid")
			return false
		}

		if c.BlockList[i].Hash != c.CalculateHash(c.BlockList[i].Index, c.BlockList[i].TimeStamp, c.BlockList[i].PrevHash, c.BlockList[i].Data, c.BlockList[i].Nonce, c.BlockList[i].Difficulty) {
			fmt.Println("chain is not valid")
			return false
		}
	}

	fmt.Println("chain is valid")
	return true
}

func (c *BlockChain) Print() {
	for i, b := range c.BlockList {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		b.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (c *BlockChain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range c.BlockList {
		for _, t := range b.Data {
			value := t.GetValue()
			if blockchainAddress == t.GetRecipientAddress() {
				totalAmount += value
			}

			if blockchainAddress == t.GetSenderAddress() {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (c *BlockChain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*block.Block `json:"chains"`
	}{
		Blocks: c.BlockList,
	})
}
