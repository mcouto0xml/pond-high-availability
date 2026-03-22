package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"producer/internal/config"
	"producer/internal/endpoints"
	"time"
	"os"
)

type Server struct {
	mux 		*http.ServeMux
	httpServer 	*http.Server
	ctx 		*context.Context
	Addr 		string
}

func NewServer(addr string, ctx *context.Context) Server {
	return Server{  Addr: addr,  ctx: ctx,  }
}

func loadEnv(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func (sa *Server) Start() {
	
	projectID := loadEnv("PROJECT_ID", "ponderada")
	location := loadEnv("QUEUE_LOCATION", "us-central1")
	queueID := loadEnv("QUEUE_ID", "fila-bacana")
	workerURL := loadEnv("WORKER_URL", "google.com")
	saEmail := loadEnv("SERVICE_ACCOUNT_EMAIL", "sa@gmail.com")
	
	cloudTasks, err := config.NewTaskEnqueuer(sa.ctx, projectID, location, queueID)
	queueImplementation := &config.QueueImplementation{  Context: &cloudTasks,  }

	sa.mux = http.NewServeMux()

	r := endpoints.NewRouter(sa.mux, queueImplementation, workerURL, saEmail) // <- passar o client aqui

	sa.httpServer = &http.Server{
		Addr: sa.Addr,
		Handler: r,
	}

	fmt.Printf("Servidor escutando na porta %s\n", sa.Addr)
	err = sa.httpServer.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Fatalf("O servidor HTTP fechou por um erro inesperado!")
		sa.shutdown()
	} else {
		log.Printf("Terminando o Server HTTP com honra!")
	}


}


func (sa *Server) shutdown() {
	if sa.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := sa.httpServer.Shutdown(ctx)
		if err != nil {
			log.Panicf("O honorado shutdown falhou!: %v", err)
		} else {
			sa.httpServer = nil
		}
	}
}