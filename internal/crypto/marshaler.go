package crypto

import "errors"

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
func (m *MarshalerGetter) GetMarshalerByAlgorithm(algorithm Algorithm) (Marshaler, error) {
	switch algorithm {
	case RSA:
		return m.rsaMarshaler, nil
	case ECC:
		return m.eccMarshaler, nil
	default:
		return nil, errors.New("Algorithm is not supported!")
	}
}