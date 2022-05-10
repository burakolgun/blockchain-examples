package main

import (
	"../block"
	"../block_chain"
	"../transaction"
	"../utils"
	"../wallet"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"strings"
)

var chainStore = make(map[string]*blockchain.BlockChain)
var walletStore = make(map[string]*wallet.Wallet)

type TransactionResponse struct {
	SenderAddress    string  `json:"senderAddress"`
	RecipientAddress string  `json:"recipientAddress"`
	Value            float32 `json:"value"`
}

func init() {
	log.SetPrefix("Blockchain Server: ")

	fmt.Printf("%s migration part started \n", strings.Repeat("=", 25))

	minerWallet := wallet.NewWallet()
	walletUserA := wallet.NewWallet()
	walletUserB := wallet.NewWallet()
	//walletUserC := wallet.NewWallet()

	chain := blockchain.New()

	t := wallet.NewTransaction(minerWallet.PrivateKey(), minerWallet.PublicKey(), minerWallet.BlockchainAddress(), walletUserB.BlockchainAddress(), 1000)
	t2 := wallet.NewTransaction(walletUserB.PrivateKey(), walletUserB.PublicKey(), walletUserB.BlockchainAddress(), walletUserA.BlockchainAddress(), 100)

	chain.BlockList = append(chain.BlockList, block.CreateGenesisBlock([]*transaction.Transaction{transaction.New("genesis", minerWallet.BlockchainAddress(), 3000)}))
	chain.Mining("00000")

	isOk := chain.AddTransaction(minerWallet.BlockchainAddress(), walletUserB.BlockchainAddress(), 1000, minerWallet.PublicKey(), t.GenerateSignature())
	fmt.Printf("is ok, %v \n", isOk)
	chain.Mining("00000")

	isOk = chain.AddTransaction(walletUserB.BlockchainAddress(), walletUserA.BlockchainAddress(), 100, walletUserB.PublicKey(), t2.GenerateSignature())
	fmt.Printf("is ok, %v \n", isOk)

	chain.Mining("00000")
	fmt.Printf("A %.1f\n", chain.CalculateTotalAmount(walletUserA.BlockchainAddress()))
	fmt.Printf("B %.1f\n", chain.CalculateTotalAmount(walletUserB.BlockchainAddress()))
	fmt.Printf("M %.1f\n", chain.CalculateTotalAmount(minerWallet.BlockchainAddress()))

	fmt.Printf("%s migration part completed\n", strings.Repeat("=", 25))

	go chain.RecursiveMiner()
	chainStore["blockchain"] = chain
	walletStore["minerWallet"] = minerWallet
	walletStore["walletUserA"] = walletUserA
	walletStore["walletUserB"] = walletUserB
}

func main() {
	fmt.Println("server run")

	server := echo.New()

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	server.GET("/random-wallet", func(c echo.Context) error {
		w := walletStore["minerWallet"]

		walletJson, err := w.MarshallJson()

		if err != nil {
			panic(err)
		}

		return c.JSONBlob(http.StatusOK, walletJson)
	})

	server.GET("/blocks", func(c echo.Context) error {
		bc := chainStore["blockchain"]

		blockListJson, err := bc.MarshalJSON()

		if err != nil {
			panic(err)
		}

		return c.JSONBlob(http.StatusOK, blockListJson)
	})

	server.GET("/wallet-statuses", func(c echo.Context) error {
		addresses := make(map[string]float32)
		bc := chainStore["blockchain"]

		for _, bc := range bc.BlockList {
			for _, t := range bc.Data {
				if _, ok := addresses[t.GetSenderAddress()]; !ok {
					addresses[t.GetSenderAddress()] = 0.0
				}

				if _, ok := addresses[t.GetRecipientAddress()]; !ok {
					addresses[t.GetRecipientAddress()] = 0.0
				}
			}
		}

		for address, _ := range addresses {
			addresses[address] = bc.CalculateTotalAmount(address)
		}

		walletJson, err := json.Marshal(addresses)

		if err != nil {
			panic(err)
		}

		return c.JSONBlob(http.StatusOK, walletJson)
	})

	server.GET("/wallet-balance/:walletAddress", func(c echo.Context) error {
		bc := chainStore["blockchain"]

		walletAddress := c.Param("walletAddress")

		balance := bc.CalculateTotalAmount(walletAddress)

		walletBalanceJSON, err := json.Marshal(struct {
			Balance float32 `json:"balance"`
		}{
			Balance: balance,
		})

		if err != nil {
			panic(err)
		}

		return c.JSONBlob(http.StatusOK, walletBalanceJSON)
	})

	server.GET("/transaction-pool", func(c echo.Context) error {
		bc := chainStore["blockchain"]

		response := struct {
			TransactionPool []TransactionResponse `json:"transactionPool"`
			Length          int                   `json:"length"`
		}{
			Length: len(bc.TransactionPool),
		}

		for _, t := range bc.TransactionPool {
			response.TransactionPool = append(response.TransactionPool, TransactionResponse{
				SenderAddress:    t.GetSenderAddress(),
				RecipientAddress: t.GetRecipientAddress(),
				Value:            t.GetValue(),
			})
		}

		m, err := json.Marshal(response)

		if err != nil {
			panic(err)
		}

		return c.JSONBlob(http.StatusOK, m)
	})

	server.POST("/transactions", HandleTransaction)
	server.Logger.Fatal(server.Start(":5001"))
}

func HandleTransaction(c echo.Context) error {
	req := blockchain.TransactionRequest{}
	err := c.Bind(&req)

	if err != nil {
		panic(err)
	}

	if !req.Validate() {
		panic("bad request")
	}

	publicKey := utils.PublicKeyFromString(*req.SenderPublicKey)
	signature := utils.SignatureFromString(*req.Signature)

	bc := chainStore["blockchain"]

	isCreated := bc.CreateTransaction(*req.SenderBlockchainAddress, *req.RecipientBlockchainAddress, *req.Value, publicKey, signature)

	if isCreated == true {
		err := c.NoContent(http.StatusCreated)

		if err != nil {
			panic(err)
		}
		return nil
	}

	err = c.JSON(http.StatusBadRequest, struct {
		Message string
	}{
		Message: "transaction didnt completed",
	})

	if err != nil {
		panic(err)
	}
	return nil
}
