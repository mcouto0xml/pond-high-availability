package dto

type TelemetryNewDataRequest struct {
	IotName 			string 		`json:"iot_name"`
	Temperature			float64 	`json:"temperature"`
	Humidity 			float64		`json:"humidity"`
	Presence 			bool 		`json:"presence"`
	Vibration  			float64 	`json:"vibration"`
	Luminosity 			float64 	`json:"luminosity"`
	TankLevel 			float64 	`json:"tank_level"`
}