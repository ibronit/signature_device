package signature

import "github.com/google/uuid"

type SignatureRequest struct {
	DeviceId       uuid.UUID `json:"device_id"`
	DataToBeSigned string    `json:"data_to_be_signed"`
}
