package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Algorithm int

const (
	RSA Algorithm = iota
	ECC
	UNKNOWN
)

type DeviceRequest struct {
	Id            uuid.UUID `json:"id"`
	Algorithm     Algorithm `json:"algorithm"`
	Label         string    `json:"label,omitempty"`
	Counter       int
	LastSignature time.Time
}

func (algorithm *Algorithm) UnmarshalJSON(b []byte) error {
	var algorithmString string
	if err := json.Unmarshal(b, &algorithmString); err != nil {
		return err
	}
	switch strings.ToUpper(algorithmString) {
	case "ECC":
		*algorithm = ECC
	case "RSA":
		*algorithm = RSA
	default:
		*algorithm = UNKNOWN
	}

	return nil
}
