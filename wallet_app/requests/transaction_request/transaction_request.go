package transaction_request

type TransactionRequest struct {
	SenderPrivateKey           *string  `json:"senderPrivateKey"`
	SenderPublicKey            *string  `json:"senderPublicKey"`
	SenderBlockchainAddress    *string  `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress *string  `json:"recipientBlockchainAddress"`
	Value                      *float32 `json:"value"`
}
