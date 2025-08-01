package metrics

import (
	"context"
	"encoding/xml"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	cachePkg "github.com/chris-ramon/golang-scaffolding/cache"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
	"github.com/chris-ramon/golang-scaffolding/drawio/gantt"
)

func TestGenerateGanttDrawIOFromPullRequests(t *testing.T) {
	// Create a test service
	cache := cachePkg.New()
	httpClient := &http.Client{}
	service, err := NewService(cache, httpClient)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

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
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test_gantt.drawio")
	err = os.WriteFile(testFilePath, drawioContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

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
		Limit:         5,
	}

	ctx := context.Background()
	result, err := service.GeneratePullRequestsGantt(ctx, params)
	if err != nil {
		t.Fatalf("Failed to generate Gantt: %v", err)
	}

	// Verify result
	if len(result.Parts) == 0 {
		t.Error("Expected at least one part to be generated")
	}

	firstPart := result.Parts[0]
	if firstPart.Limit <= 0 || firstPart.Limit > 5 {
		t.Errorf("Expected Limit to be between 1 and 5, got %d", firstPart.Limit)
	}

	if firstPart.UUID == "" {
		t.Error("Expected UUID to be set")
	}

	if firstPart.FilePath == "" {
		t.Error("Expected FilePath to be set")
	}

	// Verify file exists
	if _, err := os.Stat(firstPart.FilePath); os.IsNotExist(err) {
		t.Errorf("Generated file does not exist: %s", firstPart.FilePath)
	}

	// Verify file content
	content, err := os.ReadFile(firstPart.FilePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	var mxFile gantt.MxFile
	err = xml.Unmarshal(content, &mxFile)
	if err != nil {
		t.Fatalf("Generated file is not valid XML: %v", err)
	}

	// Clean up all generated files
	for _, part := range result.Parts {
		defer os.Remove(part.FilePath)
	}

	t.Logf("Successfully generated %d Gantt file(s): first file %s (UUID: %s)", len(result.Parts), firstPart.FilePath, firstPart.UUID)
}
