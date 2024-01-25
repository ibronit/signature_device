package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/signature"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type signatureHandler struct {
	signatureService signature.SignatureService
	logger           *slog.Logger
}

func NewSignatureHandler(signatureService signature.SignatureService, logger *slog.Logger) *signatureHandler {
	return &signatureHandler{signatureService: signatureService, logger: logger}
}

func (h *signatureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}

	var signatureRequest SignatureRequest
	if err := json.NewDecoder(r.Body).Decode(&signatureRequest); err != nil {
		h.logger.Error("Error during deserializing json body", "error", err)
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Couldn't serialize the payload"})
		return
	}

	err := validator.New().Struct(signatureRequest)
	if err != nil {
		h.logger.Error("Validation failed", "error", err)
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Provided json payload is not valid", err.Error()})
		return
	}

	device, lastSignature, signature, err := h.signatureService.SignData(signatureRequest.DeviceId, signatureRequest.DataToBeSigned)
	if err != nil {
		h.logger.Error("Couldn't create signature", "error", err)
		WriteInternalError(w)
		return
	}

	signatureData := fmt.Sprintf("%v_%v_%v", device.Counter, signatureRequest.DataToBeSigned, lastSignature)
	WriteAPIResponse(w, http.StatusCreated, SignatureResponse{Signature: signature, SignatureData: signatureData})
}

type SignatureRequest struct {
	DeviceId       uuid.UUID `json:"device_id" validate:"required,uuid"`
	DataToBeSigned string    `json:"data_to_be_signed" validate:"required"`
}

type SignatureResponse struct {
	Signature     string
	SignatureData string
}
