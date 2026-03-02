package auth

type User struct {
	ID               string
	Username         string
	PasswordHash     string
	PublicKey        []byte
	EncryptedPrivKey []byte
	PrivKeySalt      []byte
}

type SessionIdentity struct {
	UserID     string
	Username   string
	PublicKey  [32]byte
	PrivateKey [32]byte
}
