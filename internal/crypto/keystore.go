package crypto

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

var KeystoreParams = struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen int
}{
	Time:    3,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
	SaltLen: 16,
}

func EncryptPrivateKey(privKey [32]byte, password string) (encrypted, nonce, salt []byte, err error) {
	salt = make([]byte, KeystoreParams.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, nil, err
	}

	aesKey := argon2.IDKey(
		[]byte(password), salt,
		KeystoreParams.Time, KeystoreParams.Memory, KeystoreParams.Threads, KeystoreParams.KeyLen,
	)

	encrypted, nonce, err = Encrypt(aesKey, privKey[:])
	return encrypted, nonce, salt, err
}

func DecryptPrivateKey(encrypted, nonce, salt []byte, password string) ([32]byte, error) {
	aesKey := argon2.IDKey(
		[]byte(password), salt,
		KeystoreParams.Time, KeystoreParams.Memory, KeystoreParams.Threads, KeystoreParams.KeyLen,
	)

	plain, err := Decrypt(aesKey, encrypted, nonce)
	if err != nil {
		return [32]byte{}, err
	}

	var key [32]byte
	copy(key[:], plain)
	return key, nil
}
