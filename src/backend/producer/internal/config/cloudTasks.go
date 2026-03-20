package config

import (
	"context"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

type TaskEnqueuer struct{
	client 		*cloudtasks.Client
	projectID 	string
	location  	string
	queueID 	string
	ctx 		*context.Context
	QueuePath 	string
}

func NewTaskEnqueuer(ctx *context.Context, pID string, loc string, qID string) (TaskEnqueuer, error) {

	t := TaskEnqueuer{
		projectID: pID,
		location: loc,
		queueID: qID,
		ctx: ctx,
	}

	err := t.newClient()
	if err != nil {
		return TaskEnqueuer{}, err
	}

	return t, nil
}

func (t *TaskEnqueuer) newClient() error {
	ctClient, err := cloudtasks.NewClient(*t.ctx)
	if err != nil {
		return fmt.Errorf("Aconteceu um erro ao criar o Cloud Tasks Client: %v", err)
	} 
	t.client = ctClient
	t.QueuePath = fmt.Sprintf("projects/%s/locations/%s/queues/%s", t.projectID, t.location, t.queueID)

	return nil
}

func (t *TaskEnqueuer) CreateTask(body []byte, workerUrl string, saEmail string) error {
	req := &taskspb.CreateTaskRequest{
		Parent: t.QueuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url: workerUrl,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: body,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: saEmail,
							Audience: workerUrl,
						},
					},
				},
			},
		},
	}

	task, err := t.client.CreateTask(*t.ctx, req)
	if err != nil {
		return fmt.Errorf("Aconteceu um erro ao enviar a Task: %v", err)
	}

	fmt.Printf("Task criada com sucesso!: %s\n", task.Name)

	return nil
}