package device

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type deviceHandler struct {
	deviceService DeviceService
	logger        *slog.Logger
}

func NewDeviceHandler(deviceService DeviceService, logger *slog.Logger) *deviceHandler {
	return &deviceHandler{deviceService: deviceService, logger: logger}
}

func (h *deviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}

	var deviceRequest DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&deviceRequest); err != nil {
		h.logger.Error("Error during deserializing json body", "error", err)
		api.WriteErrorResponse(w, 400, []string{"Couldn't serialize the payload"})
		return
	}

	err := validator.New().Struct(deviceRequest)
	if err != nil {
		h.logger.Error("Validation failed", "error", err)
		api.WriteErrorResponse(w, 400, []string{"Provided json payload is not valid", err.Error()})
		return
	}

	uuid, err := h.deviceService.CreateSignatureDevice(deviceRequest.Id, deviceRequest.Algorithm, deviceRequest.Label)
	if err != nil {
		h.logger.Error("Couldn't create the signature device", "error", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteAPIResponse(w, 201, DeviceResponse{Id: uuid})
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 Method not allowed"))
}

type DeviceRequest struct {
	Id        uuid.UUID `json:"id" validate:"required,uuid"`
	Algorithm Algorithm `json:"algorithm" validate:"required"`
	Label     string    `json:"label,omitempty"`
}

type DeviceResponse struct {
	Id uuid.UUID `json:"id"`
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
		return errors.New("Algorithm is not supported!")
	}

	return nil
}
