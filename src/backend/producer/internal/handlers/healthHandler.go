package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"producer/internal/config"
	"producer/internal/dto"
)

type HealthInstance struct {
	ct *config.QueueImplementation
}

func NewHealthHandler(ct *config.QueueImplementation) HealthInstance {
	return HealthInstance{ct: ct}
}

func (hh *HealthInstance) Healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Conferir a saúde da conexão com o Cloud Tasks
	if err := hh.ct.Context.Ping(); err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		resp := dto.HealthzResponse{  Status: "unhealthy", CloudTasksConnectionHealth: fmt.Sprint("unhealthy: %v", err)}
		json.NewEncoder(w).Encode(resp)
		return
	}


	resp := dto.HealthzResponse{  Status: "okay", CloudTasksConnectionHealth: "okay", }
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}