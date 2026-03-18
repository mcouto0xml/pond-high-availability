package endpoints

import (
	"net/http"
	"producer/internal/config"
	"producer/internal/handlers"
)


type Router struct{
	mux 					*http.ServeMux
	healthInstance  		*handlers.HealthInstance
	telemetryInstance 		*handlers.TelemetryInstance
	cloudTasks 				*config.QueueImplementation
}

func NewRouter(mux *http.ServeMux, ct *config.QueueImplementation, wURL string) *Router {
	// Aqui inicializa as structs de Healthz e Telemtry
	health := handlers.NewHealthHandler()
	telemetry := handlers.NewTelemetryInstance(ct, wURL)

	r := &Router{	mux: mux, healthInstance: &health, telemetryInstance: &telemetry,  }
	r.registerRoutes()
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) registerRoutes() {
	r.mux.HandleFunc("GET /healthz", r.healthInstance.Healthz)
	r.mux.HandleFunc("POST /telemetry", r.telemetryInstance.NewData)
}