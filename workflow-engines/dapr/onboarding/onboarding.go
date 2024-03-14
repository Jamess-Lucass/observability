package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/workflow"
)

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
	w, err := workflow.NewWorker()
	if err != nil {
		log.Fatal(err)
	}

	if err := w.RegisterWorkflow(OnboardingWorkflow); err != nil {
		log.Fatal(err)
	}

	if err := w.RegisterActivity(CreateUser); err != nil {
		log.Fatal(err)
	}

	if err := w.Start(); err != nil {
		log.Fatal(err)
	}

	c, err := client.NewClient()
	if err != nil {
		log.Fatalf("failed to intialise client: %v", err)
	}

	s := NewServer(c)

	http.HandleFunc("POST /onboarding", s.handleCreateOnboarding)
	http.HandleFunc("POST /onboardings/{id}/approve", s.ApproveOnboarding)
	http.HandleFunc("POST /onboardings/{id}/deny", s.DenyOnboarding)

	fmt.Println("Starting web server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) ApproveOnboarding(w http.ResponseWriter, r *http.Request) {
	runId := r.PathValue("id")

	req := OnboardingApprovalRequest{
		Approved: true,
	}

	if err := s.client.RaiseEventWorkflowBeta1(context.Background(), &client.RaiseEventWorkflowRequest{
		InstanceID: runId,
		EventName:  "onboarding-approval",
		EventData:  req,
	}); err != nil {
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

	if err := s.client.RaiseEventWorkflowBeta1(context.Background(), &client.RaiseEventWorkflowRequest{
		InstanceID: runId,
		EventName:  "onboarding-approval",
		EventData:  req,
	}); err != nil {
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

	wfClient, err := workflow.NewClient(workflow.WithDaprClient(s.client))
	if err != nil {
		log.Fatalf("failed to initialise workflow client: %v", err)
	}

	options := &client.StartWorkflowRequest{
		WorkflowName: "OnboardingWorkflow",
		Options:      nil,
		Input:        request,
		SendRawInput: false,
	}

	we, err := s.client.StartWorkflowBeta1(context.Background(), options)
	if err != nil {
		log.Fatalln("unable to start Workflow", err)
	}

	metadata, err := wfClient.WaitForWorkflowCompletion(r.Context(), we.InstanceID)
	if err != nil {
		log.Fatalln("unable to wait for workflow completion", err)
	}

	// Get the results
	var fullname string
	if err := json.Unmarshal([]byte(metadata.SerializedOutput), &fullname); err != nil {
		log.Fatalln("unable to unmarshal state data", err)
	}

	json.NewEncoder(w).Encode(fullname)
}

func OnboardingWorkflow(ctx *workflow.WorkflowContext) (any, error) {
	var request OnboardingRequest
	if err := ctx.GetInput(&request); err != nil {
		return nil, err
	}

	// Manual interaction
	var signal OnboardingApprovalRequest
	err := ctx.WaitForExternalEvent("onboarding-approval", time.Duration(math.MaxInt64)).Await(&signal)
	if err != nil {
		return nil, err
	}

	if !signal.Approved {
		return nil, errors.New("was not approved")
	}

	// Resume run

	// Create user
	var fullname string
	err = ctx.CallActivity(CreateUser, workflow.ActivityInput(request)).Await(&fullname)

	return fullname, err
}

func CreateUser(ctx workflow.ActivityContext) (any, error) {
	var request OnboardingRequest
	if err := ctx.GetInput(&request); err != nil {
		return "", err
	}

	fullname := fmt.Sprintf("%s %s", request.Firstname, request.Lastname)

	return fullname, nil
}
