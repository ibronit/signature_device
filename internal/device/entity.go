package device

import (
	"time"

	"github.com/google/uuid"
)

type Algorithm int

const (
	RSA Algorithm = iota
	ECC
	UNKNOWN
)

type DeviceEntity struct {
	Uuid          uuid.UUID
	Algorithm     Algorithm
	Label         string
	Counter       int
	LastSignature time.Time
	PrivateKey    []byte
	PublicKey     []byte
}
