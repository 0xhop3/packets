package message

import "time"

type EncryptedMessage struct {
	ID               string
	SenderID         string
	RecipientID      string
	Ciphertext       []byte
	Nonce            []byte
	SenderCipherText []byte
	SenderNonce      []byte
	CreatedAt        time.Time
}
