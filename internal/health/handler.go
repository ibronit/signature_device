package health

import (
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type HealthHandler struct{}

// Health evaluates the health of the service and writes a standardized response.
func (s *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	health := HealthResponse{
		Status:  "pass",
		Version: "v0",
	}

	api.WriteAPIResponse(w, http.StatusOK, health)
}
