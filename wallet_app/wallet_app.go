package main

import (
	"../block_chain"
	"../utils"
	"../wallet"
	"./requests/transaction_request"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	log.SetPrefix("Wallet App: ")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	fmt.Println("server run")

	t := &Template{
		templates: template.Must(template.ParseGlob("/Users/burak.olgun/projects/blockchain-examples/wallet_app/templates/*.html")),
	}

	server := echo.New()

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	server.Renderer = t

	server.GET("/", WalletApp)

	server.GET("/balance/:walletAddress", func(c echo.Context) error {
		walletAddress := c.Param("walletAddress")

		res, err := http.Get("http://localhost:5001/wallet-balance/" + walletAddress)

		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		return c.JSONBlob(http.StatusOK, body)
	})

	server.POST("/send-transaction", HandleTransaction)
	server.Logger.Fatal(server.Start(":5000"))
}

func WalletApp(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", "test")
}

func HandleTransaction(c echo.Context) error {
	req := transaction_request.TransactionRequest{}
	err := c.Bind(&req)

	if err != nil {
		panic(err)
	}

	publicKey := utils.PublicKeyFromString(*req.SenderPublicKey)
	privateKey := utils.PrivateKeyFromString(*req.SenderPrivateKey, publicKey)

	fmt.Println(publicKey)
	fmt.Println(privateKey)

	t := wallet.NewTransaction(privateKey, publicKey, *req.SenderBlockchainAddress, *req.RecipientBlockchainAddress, *req.Value)

	sign := t.GenerateSignature()
	signStr := sign.String()

	bt := blockchain.TransactionRequest{
		SenderBlockchainAddress:    req.SenderBlockchainAddress,
		RecipientBlockchainAddress: req.RecipientBlockchainAddress,
		SenderPublicKey:            req.SenderPublicKey,
		Value:                      req.Value,
		Signature:                  &signStr,
	}

	m, err := json.Marshal(bt)

	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(m)

	res, err := http.Post("http://localhost:5001/transactions", "application/json", buf)

	if err != nil {
		panic(err)
	}

	fmt.Println(res.StatusCode)

	return c.JSON(http.StatusOK, req)
}
