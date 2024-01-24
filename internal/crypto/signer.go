package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"io"
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
	case ECC:
		return &ECCSigner{}, nil
	default:
		return nil, errors.New("Algorithm is not supported!")
	}
}

// RSASigner creates signature for msg.
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

type ECCSigner struct{}

func (ss *ECCSigner) CreateSignature(msg []byte, keyPair interface{}) ([]byte, error) {
	eccKeyPair, ok := keyPair.(*ECCKeyPair)
	if !ok {
		return nil, errors.New("Keypair type is not supported!")
	}

	h := md5.New()

	_, err := io.WriteString(h, "This is a message to be signed and verified by ECDSA!")
	if err != nil {
		return []byte(""), err
	}
	signhash := h.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, eccKeyPair.Private, signhash)
	if err != nil {
		return []byte(""), err
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	return signature, nil
}
