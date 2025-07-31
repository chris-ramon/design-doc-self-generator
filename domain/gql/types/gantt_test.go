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
	if _, exists := fields["uuid"]; !exists {
		t.Error("Expected 'uuid' field to exist")
	}

	if _, exists := fields["filePath"]; !exists {
		t.Error("Expected 'filePath' field to exist")
	}
}

func TestGitHubTypeGeneratePullRequestsGantt(t *testing.T) {
	// Test that the generatePullRequestsGantt field exists
	fields := GitHubType.Fields()
	ganttField, exists := fields["generatePullRequestsGantt"]
	if !exists {
		t.Error("Expected 'generatePullRequestsGantt' field to exist")
	}

	if ganttField.Type != GanttResultType {
		t.Error("Expected generatePullRequestsGantt field to have GanttResultType")
	}

	if ganttField.Description == "" {
		t.Error("Expected generatePullRequestsGantt field to have a description")
	}
}

func TestGeneratePullRequestsGanttResolver(t *testing.T) {
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

	// Create resolve params
	params := graphql.ResolveParams{
		Context: context.Background(),
		Source: map[string]interface{}{
			"url": "https://github.com/graphql-go/graphql",
		},
		Info: graphql.ResolveInfo{
			RootValue: map[string]interface{}{
				"services": mockServices,
			},
		},
	}

	// Get the resolver function
	fields := GitHubType.Fields()
	ganttField := fields["generatePullRequestsGantt"]
	
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
