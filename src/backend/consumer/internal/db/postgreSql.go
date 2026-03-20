package db

import (
	"function.com/consumer/function/internal/utils"
	"fmt"
	"github.com/go-pg/pg/v10"
	"function.com/consumer/function/internal/models"
	"crypto/tls"
)

type PostgreSql struct{
	DB 		*pg.DB
}


func StartDB() (*PostgreSql, error) {
    dbUrl := utils.LoadEnv("DATABASE_URL", "")

    if dbUrl == "" {
        return nil, fmt.Errorf("DATABASE_URL não definida")
    }

    dbOpts, err := pg.ParseURL(dbUrl)
    if err != nil {
        return nil, fmt.Errorf("Aconteceu um erro ao parsear a URL do Banco de Dados: %v", err)
    }

    // Configurações extras recomendadas para Supabase
    dbOpts.TLSConfig = &tls.Config{
        InsecureSkipVerify: true, // Supabase exige TLS mas o cert pode variar
    }

    db := pg.Connect(dbOpts)

    if err := db.Ping(db.Context()); err != nil {
        return nil, fmt.Errorf("Aconteceu um erro de conexão com o Banco de Dados: %v", err)
    }

    return &PostgreSql{DB: db}, nil
}

func (p *PostgreSql) Ping() error {
	if err := p.DB.Ping(p.DB.Context()); err != nil {
		return fmt.Errorf("A conexão com o banco não está saudável: %v", err)
	}
	return nil
}

func (p *PostgreSql) CreateTelemetry(m *models.Telemetry) error {
	_, err := p.DB.Model(m).Returning("*").Insert()

	if err != nil{
		return fmt.Errorf("Erro ao inserir um novo dado no PostgreSQL: %v", err)
	}

	return nil
}

func (p *PostgreSql) CreateDevice(m *models.Device) error {
	_, err := p.DB.Model(m).Returning("*").Insert()

	if err != nil{
		return fmt.Errorf("Erro ao inserir um novo dado no PostgreSQL: %v", err)
	}

	return nil
}

func (p *PostgreSql) CreateTelemetryBasedOnDeviceName(m *models.Telemetry, deviceName string) error {
	_, err := p.DB.Model(m).Value("iot_id", "(SELECT id FROM devices WHERE name = ?)", deviceName).Returning("*").Insert()
	
	if err != nil{
		return fmt.Errorf("Erro ao criar Telemetry no Banco de Dados: %v", err)
	}
	fmt.Printf("Novo registro de Telemetria criado com sucesso! %d\n", m.ID)
	return nil
}

func (p *PostgreSql) GetDeviceByName(n string) (*models.Device, error) {
	device := models.Device{}

	err := p.DB.Model(&device).Where("name = ?", n).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("Device '%s' não encontrado", n)
		}
		return nil, fmt.Errorf("Erro ao buscar Device no Banco de Dados: %v", err)
	}

	return &device, nil
}