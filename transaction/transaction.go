package transaction

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

func New(senderAddress string, recipientAddress string, value float32) *Transaction {
	return &Transaction{
		senderAddress:    senderAddress,
		recipientAddress: recipientAddress,
		value:            value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address      %s\n", t.senderAddress)
	fmt.Printf(" recipient_blockchain_address   %s\n", t.recipientAddress)
	fmt.Printf(" value                          %.1f\n", t.value)
}

func (t *Transaction) GetValue() float32 {
	return t.value
}

func (t *Transaction) GetRecipientAddress() string {
	return t.recipientAddress
}

func (t *Transaction) GetSenderAddress() string {
	return t.senderAddress
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"senderBlockchainAddress"`
		RecipientAddress string  `json:"recipientBlockchainAddress"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.senderAddress,
		RecipientAddress: t.recipientAddress,
		Value:            t.value,
	})
}
