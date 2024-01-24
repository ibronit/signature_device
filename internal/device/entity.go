package device

import (
	"time"

	"github.com/google/uuid"
)

type Algorithm uint

const (
	RSA Algorithm = iota + 1
	ECC
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
