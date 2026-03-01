package crypto

import (
	"crypto/rand"
	"golang.org/x/crypto/curve25519"
)

type KeyPair struct {
	PublicKey  [32]byte
	PrivateKey [32]byte
}

func GenerateKeyPair() (*KeyPair, error) {
	var priv [32]byte

	if _, err := rand.Read(priv[:]); err != nil {
		return nil, err
	}

	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64

	pub, err := curve25519.X25519(priv[:], curve25519.Basepoint)
	if err != nil {
		return nil, err
	}

	kp := &KeyPair{}
	copy(kp.PrivateKey[:], priv[:])
	copy(kp.PublicKey[:], pub[:])

	return kp, nil
}

func ECDH(privateKey, peerPublicKey [32]byte) ([32]byte, error) {
	shared, err := curve25519.X25519(privateKey[:], peerPublicKey[:])
	if err != nil {
		return [32]byte{}, err
	}

	var out [32]byte
	copy(out[:], shared)
	return out, nil
}
