package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const OnboardingTaskQueue = "ONBOARDING_TASK_QUEUE"

type Server struct {
	client client.Client
}

func NewServer(client client.Client) *Server {
	return &Server{
		client: client,
	}
}

type OnboardingRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

type OnboardingApprovalRequest struct {
	Approved bool
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, OnboardingTaskQueue, worker.Options{})
	w.RegisterWorkflow(OnboardingWorkflow)
	w.RegisterActivity(CreateUser)

	s := NewServer(c)

	http.HandleFunc("POST /onboarding", s.handleCreateOnboarding)
	http.HandleFunc("POST /onboardings/{id}/approve", s.ApproveOnboarding)
	http.HandleFunc("POST /onboardings/{id}/deny", s.DenyOnboarding)

	go func() {
		fmt.Println("Starting web server on http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Starting temporal worker")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func (s *Server) ApproveOnboarding(w http.ResponseWriter, r *http.Request) {
	runId := r.PathValue("id")

	req := OnboardingApprovalRequest{
		Approved: true,
	}

	err := s.client.SignalWorkflow(context.Background(), "onboarding-workflow", runId, "onboarding-approval", req)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("Approved")
}

func (s *Server) DenyOnboarding(w http.ResponseWriter, r *http.Request) {
	runId := r.PathValue("id")

	req := OnboardingApprovalRequest{
		Approved: false,
	}

	err := s.client.SignalWorkflow(context.Background(), "onboarding-workflow", runId, "onboarding-approval", req)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("Denied")
}

func (s *Server) handleCreateOnboarding(w http.ResponseWriter, r *http.Request) {
	log.Printf("onboarding request")
	w.Header().Set("Content-Type", "application/json")

	var request OnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		json.NewEncoder(w).Encode("Unable to parse body")
		return
	}

	options := client.StartWorkflowOptions{
		ID:        "onboarding-workflow",
		TaskQueue: OnboardingTaskQueue,
	}

	we, err := s.client.ExecuteWorkflow(context.Background(), options, OnboardingWorkflow, request)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var fullname string
	err = we.Get(context.Background(), &fullname)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	json.NewEncoder(w).Encode(fullname)
}

func OnboardingWorkflow(ctx workflow.Context, request OnboardingRequest) (*string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// Manual interaction
	var signal OnboardingApprovalRequest
	workflow.GetSignalChannel(ctx, "onboarding-approval").Receive(ctx, &signal)

	if !signal.Approved {
		return nil, errors.New("was not approved")
	}

	// Resume run

	// Create user
	var fullname string
	err := workflow.ExecuteActivity(ctx, CreateUser, request).Get(ctx, &fullname)

	return &fullname, err
}

func CreateUser(ctx context.Context, request OnboardingRequest) (string, error) {
	fullname := fmt.Sprintf("%s %s", request.Firstname, request.Lastname)

	return fullname, nil
}
