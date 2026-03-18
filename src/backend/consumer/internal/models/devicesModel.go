package models

import "time"

type Device struct{
	tableName 		struct{} 	`pg:"devices"`

	ID 				int 		`pg:"id,pk"`
	Name 			string  	`pg:"name"`
	Description		string 		`pg:"description"`
	CreatedAt 		time.Time	`pg:"created_at"`
}