package crypto

import "errors"

// Marshaler can encode and decode an RSA key pair.
type Marshaler interface {
	Marshal(keyPair interface{}) ([]byte, []byte, error)
	Unmarshal(privateKeyBytes []byte) (interface{}, error)
}

type MarshalerGetter struct {}

func (m *MarshalerGetter) GetMarshalerByAlgorithm(algorithm Algorithm) (Marshaler, error) {
	switch algorithm {
	case RSA:
		return NewRSAMarshaler(), nil
	case ECC:
		return NewECCMarshaler(), nil
	default:
		return nil, errors.New("Algorithm is not supported!")
	}
}