package models

import "time"

type Telemetry struct{
	tableName 			struct{} 	`pg:"telemetry"`

	ID 					int 		`pg:"id,pk"`
	IotID 				string 		`pg:"iot_id,fk"`
	Temperature			float64 	`pg:"temperature"`
	Humidity 			float64		`pg:"humidity"`
	Presence 			bool 		`pg:"presence"`
	Vibration  			float64 	`pg:"vibration"`
	Luminosity 			float64 	`pg:"luminosity"`
	TankLevel 			float64 	`pg:"tank_level"`
	CreatedAt 			time.Time 	`pg:"created_at"`

	Device *Device `pg:"rel:has-one"`
}