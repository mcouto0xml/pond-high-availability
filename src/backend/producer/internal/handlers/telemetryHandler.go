package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"producer/internal/config"
	"producer/internal/dto"
)

type TelemetryInstance struct{
	ctx 			*context.Context
	cloudTasks		*config.QueueImplementation
	workerURL 		string
	saEmail 		string
}

func NewTelemetryInstance(ct *config.QueueImplementation, wURL string, saEmail string) TelemetryInstance {
	return TelemetryInstance{  cloudTasks: ct, workerURL: wURL, saEmail: saEmail,  }
}

func (t *TelemetryInstance) NewData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("[internal/handlers/telemetryHandler] Requisição recebida na rota: %s %s\n", r.Method, r.URL.Path)

	var req dto.TelemetryNewDataRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := dto.TelemetryNewDataResponse{
			Message: fmt.Sprintf("Erro ao decodificar a request: %v", err),
		}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	payload := map[string]any{
		"iot_name": req.IotName,
		"temperature": req.Temperature,
		"humidity": req.Humidity,
		"presence": req.Presence,
		"vibration": req.Vibration,
		"luminosity": req.Luminosity,
		"tank_level": req.TankLevel,
	}

	msg, err := json.Marshal(payload)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		resp := dto.TelemetryNewDataResponse{  Message: fmt.Sprintf("Erro ao serializar o Json: %v", err),  }
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	err = t.cloudTasks.Context.CreateTask(msg, t.workerURL, t.saEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := dto.TelemetryNewDataResponse{  Message: fmt.Sprint(err),  }
		_ = json.NewEncoder(w).Encode(resp)
		return 
	}

	w.WriteHeader(http.StatusAccepted)
	resp := dto.TelemetryNewDataResponse{  Message: "Mensagem adicionada a fila!",  }
	_ = json.NewEncoder(w).Encode(resp)

}