package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/model"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"
)

const OnboardingTaskQueue = "ONBOARDING_TASK_QUEUE"

type Server struct {
	client *executor.WorkflowExecutor
}

func NewServer(client *executor.WorkflowExecutor) *Server {
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
	c := client.NewAPIClient(
		nil,
		settings.NewHttpSettings(
			"http://localhost:8080/api",
		),
	)
	taskRunner := worker.NewTaskRunnerWithApiClient(c)
	taskRunner.StartWorker("create_user", CreateUser, 5, 100*time.Millisecond)

	workflowExecutor := executor.NewWorkflowExecutor(c)
	workflowExecutor.RegisterWorkflow(true, CreateOnboardingWorkflow(workflowExecutor).ToWorkflowDef())

	s := NewServer(workflowExecutor)

	http.HandleFunc("POST /onboarding", s.handleCreateOnboarding)
	http.HandleFunc("POST /onboardings/{id}/approve", s.ApproveOnboarding)
	http.HandleFunc("POST /onboardings/{id}/deny", s.DenyOnboarding)

	fmt.Println("Starting web server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *Server) ApproveOnboarding(w http.ResponseWriter, r *http.Request) {
	runId := r.PathValue("id")

	signal := OnboardingApprovalRequest{
		Approved: true,
	}

	s.client.

	err := s.client.SignalWorkflow(context.Background(), "onboarding-workflow", runId, "onboarding-approval", signal)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("Approved")
}

func (s *Server) DenyOnboarding(w http.ResponseWriter, r *http.Request) {
	runId := r.PathValue("id")

	signal := OnboardingApprovalRequest{
		Approved: false,
	}

	err := s.client.SignalWorkflow(context.Background(), "onboarding-workflow", runId, "onboarding-approval", signal)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("Denied")
}

func (s *Server) handleCreateOnboarding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("onboarding request")

	var request OnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		json.NewEncoder(w).Encode("Unable to parse body")
		return
	}

	we := model.NewStartWorkflowRequest(
		"onboarding-workflow",
		1,
		"",
		request,
	)

	workflowId, err := s.client.ExecuteWorkflow(
		we,
		"",
	)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	// var fullname string
	// err = we.Get(context.Background(), &fullname)
	// if err != nil {
	// 	log.Fatalln("unable to get Workflow result", err)
	// }

	json.NewEncoder(w).Encode(workflowId.Output)
}

func CreateOnboardingWorkflow(workflowExecutor *executor.WorkflowExecutor) *workflow.ConductorWorkflow {
	return workflow.NewConductorWorkflow(workflowExecutor).
		Name("onboarding-workflow").
		Version(1).
		InputParameters("firstname", "lastname", "email").
		Add(
			workflow.NewWaitTask(),
			workflow.NewSimpleTask("create_user", "create_user").
				Input("firstname", "${workflow.input.firstname}").
				Input("lastname", "${workflow.input.lastname}"),
		)
}

// func OnboardingWorkflow(ctx workflow.Context, request OnboardingRequest) error {
// 	options := workflow.ActivityOptions{
// 		StartToCloseTimeout: 10 * time.Second,
// 	}

// 	ctx = workflow.WithActivityOptions(ctx, options)

// 	// Manual interaction
// 	var signal OnboardingApprovalRequest
// 	signalChan := workflow.GetSignalChannel(ctx, "onboarding-approval")
// 	signalChan.Receive(ctx, &signal)

// 	if !signal.Approved {
// 		return errors.New("was not approved")
// 	}

// 	// Resume run

// 	// Create user
// 	var fullname string
// 	err := workflow.ExecuteActivity(ctx, CreateUser, request).Get(ctx, &fullname)

// 	return err
// }

func CreateUser(task *model.Task) (interface{}, error) {
	log.Printf("DATA: %v", task.TaskDefinition.InputTemplate)
	// fullname := fmt.Sprintf("%s %s", request.Firstname, request.Lastname)

	// return fullname, nil
	return "asd", nil
}
