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
	if r.Method == http.MethodPost {
		h.createSignatureDevice(w, r)
		return
	} else if r.Method == http.MethodGet {
		h.listSignatureDevices(w, r)
	} else {
		api.MethodNotAllowed(w, r)
		return
	}
}

func (h *deviceHandler) listSignatureDevices(w http.ResponseWriter, r *http.Request) {
	var response DeviceListResponse
	for _, device := range h.deviceService.GetAllDevices() {
		responseElem := DeviceListResponseElement{
			Id:      device.Uuid,
			Label:   device.Label,
			Counter: device.Counter,
		}
		response.Elems = append(response.Elems, responseElem)
	}
	api.WriteAPIResponse(w, 200, response)
}

type DeviceListResponseElement struct {
	Id      uuid.UUID `json:"id"`
	Label   string    `json:"label"`
	Counter int       `json:"counter"`
}

type DeviceListResponse struct {
	Elems []DeviceListResponseElement `json:"elems"`
}

func (h *deviceHandler) createSignatureDevice(w http.ResponseWriter, r *http.Request) {
	var deviceRequest CreateDeviceRequest
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

	api.WriteAPIResponse(w, 201, CreateDeviceResponse{Id: uuid})
}

type CreateDeviceRequest struct {
	Id        uuid.UUID `json:"id" validate:"required,uuid"`
	Algorithm Algorithm `json:"algorithm" validate:"required"`
	Label     string    `json:"label,omitempty"`
}

type CreateDeviceResponse struct {
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
