package wallet

import (
	"../utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

const mainNetPrefix = 0x00

func NewWallet() *Wallet {
	// 1 Creating ECDSA private  and public keys
	wallet := Wallet{}
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		panic(err)
	}

	wallet.privateKey = privateKey
	wallet.publicKey = &privateKey.PublicKey

	// 2 Perform SHA-256 hashing on the public key

	hash2 := sha256.New()
	hash2.Write(wallet.publicKey.X.Bytes())
	hash2.Write(wallet.publicKey.Y.Bytes())
	digest2 := hash2.Sum(nil)

	// 3 Perform RIPEMD-160 hashing on the result of SHA-256
	hash3 := ripemd160.New()
	hash3.Write(digest2)
	digest3 := hash3.Sum(nil)

	// 4 Add version byte in front of RIPEMD-160 hash (0x00 for main net)
	vd4 := make([]byte, 21)
	vd4[0] = mainNetPrefix

	copy(vd4[1:], digest3[:])

	// 5 Perform SHA-256 hash on the extended RIPEMD-160 result
	hash5 := sha256.New()
	hash5.Write(vd4)
	digest5 := hash5.Sum(nil)

	// 6 Perform SHA-256 hash on the result of the prev SHA-256 hash
	hash6 := sha256.New()
	hash6.Write(digest5)
	digest6 := hash6.Sum(nil)

	// 7 Take the first 4 bytes of the second SHA-256 hash for checksum.
	chsum := digest6[:4]

	// 8 Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4)
	copy(dc8[21:], chsum)

	// 9 Convert the result from a byte string into base58
	address := base58.Encode(dc8)
	wallet.blockchainAddress = address

	return &wallet
}

func (wallet *Wallet) BlockchainAddress() string {
	return wallet.blockchainAddress
}

func (wallet *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return wallet.privateKey
}

func (wallet *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", wallet.privateKey.D.Bytes())
}

func (wallet *Wallet) PublicKey() *ecdsa.PublicKey {
	return wallet.publicKey
}

func (wallet *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", wallet.publicKey.X.Bytes(), wallet.publicKey.Y.Bytes())
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, senderAddress string, recipientAddress string, value float32) *Transaction {
	return &Transaction{
		senderPrivateKey:           privateKey,
		senderPublicKey:            publicKey,
		senderBlockchainAddress:    senderAddress,
		recipientBlockchainAddress: recipientAddress,
		value:                      value,
	}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, err := json.Marshal(t)

	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(m)
	r, s, err := ecdsa.Sign(rand.Reader, t.senderPrivateKey, hash[:])

	if err != nil {
		panic(err)
	}

	return &utils.Signature{
		R: r,
		S: s,
	}

}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"senderBlockchainAddress"`
		RecipientAddress string  `json:"recipientBlockchainAddress"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.senderBlockchainAddress,
		RecipientAddress: t.recipientBlockchainAddress,
		Value:            t.value,
	})
}

func (wallet *Wallet) MarshallJson() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"privateKey"`
		PublicKey         string `json:"publicKey"`
		BlockchainAddress string `json:"blockchainAddress"`
	}{
		PrivateKey:        wallet.PrivateKeyStr(),
		PublicKey:         wallet.PublicKeyStr(),
		BlockchainAddress: wallet.blockchainAddress,
	})
}
