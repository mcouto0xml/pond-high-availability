package dbContext

import(
	"function.com/consumer/function/internal/models"
)

type DbContext interface{
	CreateDevice(m *models.Device) error
	GetDeviceByName(n string) (*models.Device, error)
	CreateTelemetry(m *models.Telemetry) error
	CreateTelemetryBasedOnDeviceName(m *models.Telemetry, deviceName string) error
	Ping() error
}

type DbImplementation struct {
	Context 	DbContext
}