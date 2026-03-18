package handlers

import (
	"encoding/json"
	"net/http"
	"producer/internal/dto"
)

type HealthInstance struct {
}

func NewHealthHandler() HealthInstance {
	return HealthInstance{}
}

func (hh *HealthInstance) Healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Conferir a saúde da conexão com o Cloud Tasks

	resp := dto.HealthzResponse{  Status: "okay", CloudTasksConnectionHealth: "okay", }
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}