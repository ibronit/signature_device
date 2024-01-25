package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/device"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/enum"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type deviceHandler struct {
	deviceService device.DeviceService
	logger        *slog.Logger
}

func NewDeviceHandler(deviceService device.DeviceService, logger *slog.Logger) *deviceHandler {
	return &deviceHandler{deviceService: deviceService, logger: logger}
}

func (h *deviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.createSignatureDevice(w, r)
		return
	} else if r.Method == http.MethodGet {
		h.listSignatureDevices(w, r)
	} else {
		MethodNotAllowed(w, r)
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
	WriteAPIResponse(w, http.StatusOK, response)
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
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Couldn't serialize the payload"})
		return
	}

	err := validator.New().Struct(deviceRequest)
	if err != nil {
		h.logger.Error("Validation failed", "error", err)
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Provided json payload is not valid", err.Error()})
		return
	}

	uuid, err := h.deviceService.CreateSignatureDevice(deviceRequest.Id, deviceRequest.Algorithm, deviceRequest.Label)
	if err != nil {
		h.logger.Error("Couldn't create the signature device", "error", err)
		WriteInternalError(w)
		return
	}

	WriteAPIResponse(w, http.StatusCreated, CreateDeviceResponse{Id: uuid})
}

type CreateDeviceRequest struct {
	Id        uuid.UUID      `json:"id" validate:"required,uuid"`
	Algorithm enum.Algorithm `json:"algorithm" validate:"required"`
	Label     string         `json:"label,omitempty"`
}

type CreateDeviceResponse struct {
	Id uuid.UUID `json:"id"`
}
