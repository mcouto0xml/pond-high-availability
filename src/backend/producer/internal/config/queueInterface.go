package config

type QueueInterface interface{
	CreateTask(body []byte, workerUrl string, saEmail string) error
	Ping() error
}

type QueueImplementation struct{
	Context 	QueueInterface
}