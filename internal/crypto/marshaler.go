package crypto

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/enum"
)

// Marshaler can encode and decode an RSA key pair.
type Marshaler interface {
	Marshal(keyPair interface{}) ([]byte, []byte, error)
	Unmarshal(privateKeyBytes []byte) (interface{}, error)
}

// Stores multiple implementations of the Marshaler interface.
type MarshalerGetter struct {
	rsaMarshaler Marshaler
	eccMarshaler Marshaler
}

func NewMarshalerGetter() *MarshalerGetter {
	return &MarshalerGetter{rsaMarshaler: NewRSAMarshaler(), eccMarshaler: NewECCMarshaler()}
}

// Gets the correct marshaler if it's supported.
func (m *MarshalerGetter) GetMarshalerByAlgorithm(algorithm enum.Algorithm) (Marshaler, error) {
	switch algorithm {
	case enum.RSA:
		return m.rsaMarshaler, nil
	case enum.ECC:
		return m.eccMarshaler, nil
	default:
		return nil, errors.New("Algorithm is not supported!")
	}
}