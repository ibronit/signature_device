package device

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/enum"
	"github.com/google/uuid"
)

type DeviceEntity struct {
	Uuid          uuid.UUID
	Algorithm     enum.Algorithm
	Label         string
	Counter       int
	LastSignature string
	PrivateKey    []byte
	PublicKey     []byte
}
