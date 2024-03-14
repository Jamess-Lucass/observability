package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const GreetingTaskQueue = "GREETING_TASK_QUEUE"

type Server struct {
	client client.Client
}

func NewServer(client client.Client) *Server {
	return &Server{
		client: client,
	}
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, GreetingTaskQueue, worker.Options{})
	w.RegisterWorkflow(GreetingWorkflow)
	w.RegisterActivity(ComposeGreeting)

	s := NewServer(c)

	http.HandleFunc("GET /greeting", s.handleGreeting)

	go func() {
		fmt.Println("Starting web server on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	fmt.Println("Starting temporal worker")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func (s *Server) handleGreeting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("making request /api/v1/test")

	options := client.StartWorkflowOptions{
		ID:        "greeting-workflow",
		TaskQueue: GreetingTaskQueue,
	}

	name := "World"
	we, err := s.client.ExecuteWorkflow(context.Background(), options, GreetingWorkflow, name)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var greeting string
	err = we.Get(context.Background(), &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	json.NewEncoder(w).Encode("test")
}

func GreetingWorkflow(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	err := workflow.ExecuteActivity(ctx, ComposeGreeting, name).Get(ctx, &result)

	return result, err
}

func ComposeGreeting(ctx context.Context, name string) (string, error) {
	greeting := fmt.Sprintf("Hello %s!", name)
	return greeting, nil
}
