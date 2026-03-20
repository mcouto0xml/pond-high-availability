package main

import (
	"example.com/consumer/internal/db"
	"example.com/consumer/internal/dbContext"
	"example.com/consumer/internal/dto"
	"example.com/consumer/internal/models"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"github.com/joho/godotenv"
)

var once sync.Once
var dbRepository dbContext.DbImplementation

func init() {
	functions.HTTP("PostTelemetry", PostTelemetry)

	once.Do(func ()  {
		var err error
		err = godotenv.Load()
		if err != nil {
			fmt.Printf("Erro ao fazer o carregamento do .env: %v", err)
		}
		postgreSql, err := db.StartDB()
		if err != nil {
			log.Fatalf(err.Error())
		}
		dbRepository = dbContext.DbImplementation{  Context: postgreSql, }

		fmt.Print("Function iniciada com sucesso!")
	})
}

func PostTelemetry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("[internal/function/posTelemetry] A rota %s %s foi chamada!\n", r.Method, r.URL.Path)

	var req dto.ConsumerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("BadRequest :C : %v", err)
		return
	}

	telemetryInstance := models.Telemetry{
		Temperature: req.Temperature,
		Humidity: req.Humidity,
		Presence: req.Presence,
		Vibration: req.Vibration,
		Luminosity: req.Luminosity,
		TankLevel: req.TankLevel,
	}

	err = dbRepository.Context.CreateTelemetryBasedOnDeviceName(&telemetryInstance, req.IotName)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}