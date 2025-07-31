package metrics

import (
	"context"
	"encoding/xml"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	cachePkg "github.com/chris-ramon/golang-scaffolding/cache"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
	"github.com/chris-ramon/golang-scaffolding/drawio/gantt"
)

func TestGenerateGanttDrawIOFromPullRequests(t *testing.T) {
	// Get the repository root directory using runtime.Caller
	_, filename, _, _ := runtime.Caller(0)
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	
	// Create a test service
	cache := cachePkg.New()
	httpClient := &http.Client{}
	service, err := NewService(cache, httpClient)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	
	// Temporarily override the template path resolution for testing
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(repoRoot)

	// Create test pull requests
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	twoDaysAgo := yesterday.Add(-time.Hour * 24)
	
	testPRs := []*types.PullRequest{
		{
			Title:     "Add new feature",
			CreatedAt: &yesterday,
			MergedAt:  &now,
		},
		{
			Title:     "Fix bug in authentication",
			CreatedAt: &twoDaysAgo,
			MergedAt:  &yesterday,
		},
	}

	// Generate DrawIO content
	drawioContent, err := service.generateGanttDrawIOFromPullRequests(testPRs)
	if err != nil {
		t.Fatalf("Failed to generate DrawIO content: %v", err)
	}

	// Verify the content can be parsed as valid XML
	var mxFile gantt.MxFile
	err = xml.Unmarshal(drawioContent, &mxFile)
	if err != nil {
		t.Fatalf("Generated content is not valid XML: %v", err)
	}

	// Verify basic structure
	if len(mxFile.Diagrams) == 0 {
		t.Error("Expected at least one diagram")
	}

	diagram := mxFile.Diagrams[0]
	if len(diagram.MxGraphModel.Root.Cells) < 2 {
		t.Error("Expected at least root cells")
	}

	// Count cells with PR titles
	prTitleCells := 0
	for _, cell := range diagram.MxGraphModel.Root.Cells {
		if cell.Value == "Add new feature" || cell.Value == "Fix bug in authentication" {
			prTitleCells++
		}
	}

	if prTitleCells != 2 {
		t.Errorf("Expected 2 PR title cells, got %d", prTitleCells)
	}

	// Write to a test file to verify it works
	testFilePath := filepath.Join("test_gantt.drawio")
	err = os.WriteFile(testFilePath, drawioContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Clean up
	defer os.Remove(testFilePath)

	t.Logf("Successfully generated Gantt DrawIO file with %d cells", len(diagram.MxGraphModel.Root.Cells))
}

func TestGeneratePullRequestsGanttIntegration(t *testing.T) {
	// Skip this test if we don't have a GitHub token
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("Skipping integration test: GITHUB_TOKEN not set")
	}

	// Create a test service
	cache := cachePkg.New()
	httpClient := &http.Client{}
	service, err := NewService(cache, httpClient)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Test with a small public repository
	params := GeneratePullRequestsGanttParams{
		RepositoryURL: "https://github.com/graphql-go/graphql",
	}

	ctx := context.Background()
	result, err := service.GeneratePullRequestsGantt(ctx, params)
	if err != nil {
		t.Fatalf("Failed to generate Gantt: %v", err)
	}

	// Verify result
	if result.UUID == "" {
		t.Error("Expected UUID to be set")
	}

	if result.FilePath == "" {
		t.Error("Expected FilePath to be set")
	}

	// Verify file exists
	if _, err := os.Stat(result.FilePath); os.IsNotExist(err) {
		t.Errorf("Generated file does not exist: %s", result.FilePath)
	}

	// Verify file content
	content, err := os.ReadFile(result.FilePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	var mxFile gantt.MxFile
	err = xml.Unmarshal(content, &mxFile)
	if err != nil {
		t.Fatalf("Generated file is not valid XML: %v", err)
	}

	// Clean up
	defer os.Remove(result.FilePath)

	t.Logf("Successfully generated Gantt file: %s (UUID: %s)", result.FilePath, result.UUID)
}
