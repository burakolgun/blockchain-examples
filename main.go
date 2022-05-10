package main

import (
	"./block"
	"./block_chain"
	"./transaction"
	"./wallet"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("BlockChain examples")
	partSecond()
}

func partSecond() {
	fmt.Printf("%s example part two \n", strings.Repeat("=", 25))

	minerWallet := wallet.NewWallet()
	walletUserA := wallet.NewWallet()
	walletUserB := wallet.NewWallet()
	//walletUserC := wallet.NewWallet()

	chain := blockchain.New()

	t := wallet.NewTransaction(minerWallet.PrivateKey(), minerWallet.PublicKey(), minerWallet.BlockchainAddress(), walletUserB.BlockchainAddress(), 1000)
	t2 := wallet.NewTransaction(walletUserB.PrivateKey(), walletUserB.PublicKey(), walletUserB.BlockchainAddress(), walletUserA.BlockchainAddress(), 100)

	chain.BlockList = append(chain.BlockList, block.CreateGenesisBlock([]*transaction.Transaction{transaction.New("genesis", minerWallet.BlockchainAddress(), 3000)}))
	chain.Mining("0000")

	isOk := chain.AddTransaction(minerWallet.BlockchainAddress(), walletUserB.BlockchainAddress(), 1000, minerWallet.PublicKey(), t.GenerateSignature())
	fmt.Printf("is ok, %v \n", isOk)
	chain.Mining("0000")

	isOk = chain.AddTransaction(walletUserB.BlockchainAddress(), walletUserA.BlockchainAddress(), 100, walletUserB.PublicKey(), t2.GenerateSignature())
	fmt.Printf("is ok, %v \n", isOk)

	chain.Mining("0000")
	fmt.Printf("A %.1f\n", chain.CalculateTotalAmount(walletUserA.BlockchainAddress()))
	fmt.Printf("B %.1f\n", chain.CalculateTotalAmount(walletUserB.BlockchainAddress()))
	fmt.Printf("M %.1f\n", chain.CalculateTotalAmount(minerWallet.BlockchainAddress()))

	fmt.Printf("%s example part two complete\n", strings.Repeat("=", 25))
}

func partOne() {
	fmt.Printf("%s wallet example part one \n", strings.Repeat("=", 25))

	chain := blockchain.New()

	chain.BlockList = append(chain.BlockList, block.CreateGenesisBlock([]*transaction.Transaction{transaction.New("A", "B", 3)}))

	chain.TransactionPool = append(chain.TransactionPool, transaction.New("B", "B", 3))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("B", "B", 2))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("B", "B", 5))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("D", "B", 5))
	chain.Mining("0000")
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("C", "B", 3))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("C", "B", 2))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("C", "B", 5))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("FF", "B", 5))
	chain.Mining("0000")
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("E", "B", 3))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("E", "B", 2))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("E", "B", 5))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("ED", "B", 5))
	chain.Mining("0000")
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("F", "B", 3))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("F", "B", 2))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("G", "B", 5))
	chain.TransactionPool = append(chain.TransactionPool, transaction.New("H", "B", 5))
	chain.Mining("0000")

	chain.Validate()

	chain.Print()

	fmt.Printf("%s chain example one Done \n", strings.Repeat("=", 25))

	fmt.Println("hello")

	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())

	fmt.Printf("wallet address -> %s \n", w.BlockchainAddress())

	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "recipientAddress", 1.42)

	fmt.Printf("signature %s \n", t.GenerateSignature())

	fmt.Printf("%s wallet example one Done \n", strings.Repeat("=", 25))
}
