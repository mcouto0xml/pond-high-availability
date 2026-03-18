package dto


type HealthzResponse struct{
	Status 							string		`json: "status"` 	
	CloudTasksConnectionHealth		string		`json: "cloud_tasks_connection_health`
}