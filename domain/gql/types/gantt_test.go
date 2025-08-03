package types

import (
	"context"
	"net/http"
	"testing"

	"github.com/graphql-go/graphql"

	cachePkg "github.com/chris-ramon/golang-scaffolding/cache"
	"github.com/chris-ramon/golang-scaffolding/domain/internal/services"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics"
)

func TestGanttResultType(t *testing.T) {
	// Test that the GanttResultType is properly defined
	if GanttResultType.Name() != "GanttResultType" {
		t.Errorf("Expected type name 'GanttResultType', got '%s'", GanttResultType.Name())
	}

	fields := GanttResultType.Fields()
	if _, exists := fields["limit"]; !exists {
		t.Error("Expected 'limit' field to exist")
	}

	if _, exists := fields["uuid"]; !exists {
		t.Error("Expected 'uuid' field to exist")
	}

	if _, exists := fields["filePath"]; !exists {
		t.Error("Expected 'filePath' field to exist")
	}
}

func TestPullRequestTextType(t *testing.T) {
	// Test that the PullRequestTextType is properly defined
	if PullRequestTextType.Name() != "PullRequestTextType" {
		t.Errorf("Expected type name 'PullRequestTextType', got '%s'", PullRequestTextType.Name())
	}

	fields := PullRequestTextType.Fields()
	if _, exists := fields["uuid"]; !exists {
		t.Error("Expected 'uuid' field to exist")
	}

	if _, exists := fields["filePath"]; !exists {
		t.Error("Expected 'filePath' field to exist")
	}

	// Test that fields have resolvers that return empty strings
	params := graphql.ResolveParams{}
	
	uuidField := fields["uuid"]
	result, err := uuidField.Resolve(params)
	if err != nil {
		t.Errorf("Expected no error from uuid resolver, got: %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string from uuid resolver, got: %v", result)
	}

	filePathField := fields["filePath"]
	result, err = filePathField.Resolve(params)
	if err != nil {
		t.Errorf("Expected no error from filePath resolver, got: %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string from filePath resolver, got: %v", result)
	}
}

func TestGitHubPullRequestType(t *testing.T) {
	// Test that the GitHubPullRequestType is properly defined
	if GitHubPullRequestType.Name() != "GitHubPullRequestType" {
		t.Errorf("Expected type name 'GitHubPullRequestType', got '%s'", GitHubPullRequestType.Name())
	}

	fields := GitHubPullRequestType.Fields()
	textField, exists := fields["text"]
	if !exists {
		t.Error("Expected 'text' field to exist")
	}

	// Check that the text field has the correct type
	if textField.Type != PullRequestTextType {
		t.Error("Expected 'text' field to have PullRequestTextType")
	}

	// Test that the text field resolver returns empty object
	params := graphql.ResolveParams{}
	result, err := textField.Resolve(params)
	if err != nil {
		t.Errorf("Expected no error from text resolver, got: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result from text resolver")
	}
	if resultMap, ok := result.(map[string]interface{}); !ok {
		t.Error("Expected result from text resolver to be map[string]interface{}")
	} else if len(resultMap) != 0 {
		t.Errorf("Expected empty map from text resolver, got: %v", resultMap)
	}

	// Test that the export field exists
	exportField, exists := fields["export"]
	if !exists {
		t.Error("Expected 'export' field to exist")
	}

	// Check that the export field has the correct type
	if exportField.Type != graphql.String {
		t.Error("Expected 'export' field to have String type")
	}
}

func TestGitHubTypeGantt(t *testing.T) {
	// Test that the gantt field exists
	gitHubFields := GitHubType.Fields()
	ganttField, exists := gitHubFields["gantt"]
	if !exists {
		t.Error("Expected 'gantt' field to exist")
	}

	// Check if the type is a list
	if listType, ok := ganttField.Type.(*graphql.List); !ok {
		t.Error("Expected gantt field to have List type")
	} else if listType.OfType != GanttResultType {
		t.Error("Expected gantt field to have List of GanttResultType")
	}

	if ganttField.Description == "" {
		t.Error("Expected gantt field to have a description")
	}

	// Test that the limit argument exists
	if ganttField.Args == nil {
		t.Error("Expected gantt field to have arguments")
	} else if len(ganttField.Args) == 0 {
		t.Error("Expected gantt field to have at least one argument")
	}
}

func TestGitHubTypePullRequests(t *testing.T) {
	// Test that the pullRequests field exists
	gitHubFields := GitHubType.Fields()
	pullRequestsField, exists := gitHubFields["pullRequests"]
	if !exists {
		t.Error("Expected 'pullRequests' field to exist")
	}

	// Check if the type is a list
	if listType, ok := pullRequestsField.Type.(*graphql.List); !ok {
		t.Error("Expected pullRequests field to have List type")
	} else if listType.OfType != GitHubPullRequestType {
		t.Error("Expected pullRequests field to have List of GitHubPullRequestType")
	}

	if pullRequestsField.Description == "" {
		t.Error("Expected pullRequests field to have a description")
	}

	// Test that the resolver returns empty array
	params := graphql.ResolveParams{}
	result, err := pullRequestsField.Resolve(params)
	if err != nil {
		t.Errorf("Expected no error from pullRequests resolver, got: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result from pullRequests resolver")
	}
	if resultArray, ok := result.([]interface{}); !ok {
		t.Error("Expected result from pullRequests resolver to be []interface{}")
	} else if len(resultArray) != 0 {
		t.Errorf("Expected empty array from pullRequests resolver, got: %v", resultArray)
	}
}

func TestGanttResolver(t *testing.T) {
	// Create mock services
	cache := cachePkg.New()
	httpClient := &http.Client{}
	metricsService, err := metrics.NewService(cache, httpClient)
	if err != nil {
		t.Fatalf("Failed to create metrics service: %v", err)
	}

	mockServices := &services.Services{
		MetricsService: metricsService,
	}

	// Create resolve params with limit argument
	params := graphql.ResolveParams{
		Context: context.Background(),
		Source: map[string]interface{}{
			"url": "https://github.com/graphql-go/graphql",
		},
		Args: map[string]interface{}{
			"limit": 10,
		},
		Info: graphql.ResolveInfo{
			RootValue: map[string]interface{}{
				"services": mockServices,
			},
		},
	}

	// Get the resolver function
	gitHubFields := GitHubType.Fields()
	ganttField := gitHubFields["gantt"]
	
	// Test that the resolver doesn't panic (we can't test the full functionality without GitHub token)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Resolver panicked: %v", r)
		}
	}()

	// Call the resolver - it may fail due to missing GitHub token, but shouldn't panic
	_, err = ganttField.Resolve(params)
	// We expect this to potentially fail due to missing GitHub token or network issues
	// The important thing is that it doesn't panic and the structure is correct
	
	t.Log("Resolver executed without panicking")
}
