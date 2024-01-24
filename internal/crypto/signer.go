package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

// Signer can generate keypairs.
type Signer interface {
	CreateSignature(msg []byte, privateKey interface{}) ([]byte, error)
}

type SignerGetter struct{}

func (m *SignerGetter) GetSignatureByAlgorithm(algorithm Algorithm) (Signer, error) {
	switch algorithm {
	case RSA:
		return &RSASigner{}, nil
	// case ECC:
	// 	return NewECCSigner(), nil
	default:
		return nil, errors.New("Algorithm is not supported!")
	}
}

// RSASigner generates a RSA key pair.
type RSASigner struct{}

func (ss *RSASigner) CreateSignature(msg []byte, keyPair interface{}) ([]byte, error) {
	rsaKeyPair, ok := keyPair.(*RSAKeyPair)
	if !ok {
		return nil, errors.New("Keypair type is not supported!")
	}

	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		return nil, err
	}
	msgHashSum := msgHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, rsaKeyPair.Private, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
