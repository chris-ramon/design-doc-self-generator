package metrics

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"

	githubClient "github.com/google/go-github/github"
	"github.com/google/uuid"

	cachePkg "github.com/chris-ramon/golang-scaffolding/cache"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/github"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
	"github.com/chris-ramon/golang-scaffolding/drawio/gantt"
)

type service struct {
	// cache is the internal cache component.
	cache *cachePkg.Cache

	// HTTPClient is the HTTP client used for GitHub API requests.
	HTTPClient *http.Client

	// GitHub is the github component.
	GitHub github.GitHubClient
}

type FindPullRequestsResult struct {
	PullRequests []*types.PullRequest
}

type findPullRequestsResult struct {
	PullRequest *types.PullRequest
}

type FindAllPullRequestsResult struct {
	PullRequests []*types.PullRequest
}

type FindAllPullRequestsParams struct {
	RepositoryURL string
}

type GeneratePullRequestsGanttPart struct {
	Limit    int
	UUID     string
	FilePath string
}

type GeneratePullRequestsGanttResult struct {
	Parts []GeneratePullRequestsGanttPart
}

type GeneratePullRequestsGanttParams struct {
	RepositoryURL string
	Limit         int
}

// `findPullRequestsCacheKey` returns cache key of `FindPullRequests`.
func (s *service) findPullRequestsCacheKey(params types.FindPullRequestsParams) (string, error) {
	key, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

// `getFindPullRequestsCacheValue` returns cached data of `FindPullRequests`.
func (s *service) getFindPullRequestsCacheValue(data any) (*FindPullRequestsResult, error) {
	result, ok := data.(*FindPullRequestsResult)
	if !ok {
		return nil, errors.New("unexpected type")
	}

	return result, nil
}

// `cacheFindPullRequestsValue` caches given result of `FindPullRequests`.
func (s *service) cacheFindPullRequestsValue(key string, data any) {
	s.cache.Add(key, data)
}

func (s *service) FindPullRequests(ctx context.Context, params types.FindPullRequestsParams) (*FindPullRequestsResult, error) {
	key, err := s.findPullRequestsCacheKey(params)
	if err != nil {
		return nil, err
	}

	findPullRequestsCacheVal, found := s.cache.Get(key)
	if found {
		return s.getFindPullRequestsCacheValue(findPullRequestsCacheVal)
	}

	result := &FindPullRequestsResult{}

	for _, pr := range params {
		r, err := s.findPullRequests(ctx, pr)
		if err != nil {
			return nil, err
		}

		result.PullRequests = append(result.PullRequests, r.PullRequest)
	}

	s.cacheFindPullRequestsValue(key, result)

	return result, nil
}

func (s *service) findPullRequests(ctx context.Context, param types.FindPullRequestParam) (*findPullRequestsResult, error) {
	// Create a GitHub client using the provided HTTP client.
	client := githubClient.NewClient(s.HTTPClient)

	// Fetch pull request information from GitHub.
	pullRequest, _, err := client.PullRequests.Get(ctx, param.Owner, param.Repo, param.Number)
	if err != nil {
		return nil, err
	}

	if pullRequest.MergedAt == nil {
		return nil, fmt.Errorf("unexpected merged at nil value")
	}

	if pullRequest.CreatedAt == nil {
		return nil, fmt.Errorf("unexpected created at nil value")
	}

	if pullRequest.Head == nil {
		return nil, fmt.Errorf("unexpected head nil value")
	}

	if pullRequest.Head.Ref == nil {
		return nil, fmt.Errorf("unexpected head ref nil value")
	}

	pullRequestContributorsParams := github.PullRequestContributorsParams{
		PullRequest: types.PullRequest{
			Owner:       param.Owner,
			Repo:        param.Repo,
			HeadRefName: *pullRequest.Head.Ref,
		},
	}
	r, err := s.GitHub.PullRequestContributors(pullRequestContributorsParams)
	if err != nil {
		return nil, err
	}

	contributors := types.Contributors{}

	for _, prNode := range r.Repository.PullRequests.Nodes {
		for _, participant := range prNode.Participants.Nodes {
			c := types.Contributor{
				ProfileURL: string(participant.URL),
			}
			contributors = append(contributors, c)
		}
	}

	// Extract pull request metrics.
	duration := pullRequest.MergedAt.Sub(*pullRequest.CreatedAt)

	var title, body string
	if pullRequest.Title != nil {
		title = *pullRequest.Title
	}
	if pullRequest.Body != nil {
		body = *pullRequest.Body
	}

	pr := &types.PullRequest{
		Duration:              duration,
		CreatedAt:             pullRequest.CreatedAt,
		MergedAt:              pullRequest.MergedAt,
		URL:                   param.URL,
		Title:                 title,
		Body:                  body,
		Contributors:          contributors,
		FormattedContributors: contributors.FormattedContributors(),
	}

	// Create the result.
	result := &findPullRequestsResult{
		PullRequest: pr,
	}

	return result, nil
}

// `findAllPullRequestsCacheKey` returns cache key of `FindAllPullRequests`.
func (s *service) findAllPullRequestsCacheKey(params FindAllPullRequestsParams) (string, error) {
	key, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

// `getFindAllPullRequestsCacheValue` returns cached data of `FindAllPullRequests`.
func (s *service) getFindAllPullRequestsCacheValue(data any) (*FindAllPullRequestsResult, error) {
	result, ok := data.(*FindAllPullRequestsResult)
	if !ok {
		return nil, errors.New("unexpected type")
	}

	return result, nil
}

// `cacheFindAllPullRequestsValue` caches given result of `FindAllPullRequests`.
func (s *service) cacheFindAllPullRequestsValue(key string, data any) {
	s.cache.Add(key, data)
}

func (s *service) FindAllPullRequests(ctx context.Context, params FindAllPullRequestsParams) (*FindAllPullRequestsResult, error) {
	key, err := s.findAllPullRequestsCacheKey(params)
	if err != nil {
		return nil, err
	}

	findAllPullRequestsCacheVal, found := s.cache.Get(key)
	if found {
		return s.getFindAllPullRequestsCacheValue(findAllPullRequestsCacheVal)
	}

	owner, repo, err := github.RepositoryFromURL(params.RepositoryURL)
	if err != nil {
		return nil, err
	}

	allPullRequestsParams := github.AllPullRequestsParams{
		Owner: owner,
		Repo:  repo,
	}

	r, err := s.GitHub.AllPullRequests(allPullRequestsParams)
	if err != nil {
		return nil, err
	}

	result := &FindAllPullRequestsResult{}

	for _, prNode := range r.Repository.PullRequests.Nodes {
		if prNode.MergedAt.Time.IsZero() || prNode.CreatedAt.Time.IsZero() {
			continue
		}

		contributors := types.Contributors{}
		for _, participant := range prNode.Participants.Nodes {
			c := types.Contributor{
				ProfileURL: string(participant.URL),
			}
			contributors = append(contributors, c)
		}

		duration := prNode.MergedAt.Time.Sub(prNode.CreatedAt.Time)
		pr := &types.PullRequest{
			Number:                int(prNode.Number),
			Owner:                 owner,
			Repo:                  repo,
			Title:                 string(prNode.Title),
			Body:                  string(prNode.Body),
			Duration:              duration,
			CreatedAt:             &prNode.CreatedAt.Time,
			MergedAt:              &prNode.MergedAt.Time,
			URL:                   string(prNode.URL),
			Contributors:          contributors,
			HeadRefName:           string(prNode.HeadRef.Name),
			FormattedContributors: contributors.FormattedContributors(),
		}

		result.PullRequests = append(result.PullRequests, pr)
	}

	s.cacheFindAllPullRequestsValue(key, result)

	return result, nil
}

// `generatePullRequestsGanttCacheKey` returns cache key of `GeneratePullRequestsGantt`.
func (s *service) generatePullRequestsGanttCacheKey(params GeneratePullRequestsGanttParams) (string, error) {
	key, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

// `getGeneratePullRequestsGanttCacheValue` returns cached data of `GeneratePullRequestsGantt`.
func (s *service) getGeneratePullRequestsGanttCacheValue(data any) (*GeneratePullRequestsGanttResult, error) {
	result, ok := data.(*GeneratePullRequestsGanttResult)
	if !ok {
		return nil, errors.New("unexpected type")
	}

	return result, nil
}

// `cacheGeneratePullRequestsGanttValue` caches given result of `GeneratePullRequestsGantt`.
func (s *service) cacheGeneratePullRequestsGanttValue(key string, data any) {
	s.cache.Add(key, data)
}

func (s *service) GeneratePullRequestsGantt(ctx context.Context, params GeneratePullRequestsGanttParams) (*GeneratePullRequestsGanttResult, error) {
	key, err := s.generatePullRequestsGanttCacheKey(params)
	if err != nil {
		return nil, err
	}

	generatePullRequestsGanttCacheVal, found := s.cache.Get(key)
	if found {
		return s.getGeneratePullRequestsGanttCacheValue(generatePullRequestsGanttCacheVal)
	}

	// Get all pull requests for the repository
	findAllPRParams := FindAllPullRequestsParams{
		RepositoryURL: params.RepositoryURL,
	}

	findAllPullRequestsResult, err := s.FindAllPullRequests(ctx, findAllPRParams)
	if err != nil {
		return nil, err
	}

	pullRequests := findAllPullRequestsResult.PullRequests
	pullRequests = s.sortPullRequestsAsc(pullRequests)

	// Extract repository name for directory structure
	owner, repo, err := github.RepositoryFromURL(params.RepositoryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract repository info: %w", err)
	}
	repoName := fmt.Sprintf("%s-%s", owner, repo)

	// Create the base directory path
	_, filename, _, _ := runtime.Caller(0)
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	baseDir := filepath.Join(repoRoot, "diagrams", "gantt", "generated", repoName)

	// Ensure the directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Divide pull requests into chunks based on limit
	var parts []GeneratePullRequestsGanttPart
	limit := params.Limit
	if limit <= 0 {
		limit = 25 // Default limit
	}

	for i := 0; i < len(pullRequests); i += limit {
		end := i + limit
		if end > len(pullRequests) {
			end = len(pullRequests)
		}

		chunk := pullRequests[i:end]

		// Generate UUID for this part
		fileUUID := uuid.New().String()

		// Generate the Gantt DrawIO file for this chunk
		drawioContent, err := s.generateGanttDrawIOFromPullRequests(chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to generate Gantt for chunk %d: %w", i/limit+1, err)
		}

		// Create the file path
		filePath := filepath.Join(baseDir, fileUUID+".drawio")

		// Write the file
		if err := os.WriteFile(filePath, drawioContent, 0644); err != nil {
			return nil, fmt.Errorf("failed to write DrawIO file: %w", err)
		}

		parts = append(parts, GeneratePullRequestsGanttPart{
			Limit:    len(chunk),
			UUID:     fileUUID,
			FilePath: filePath,
		})
	}

	result := &GeneratePullRequestsGanttResult{
		Parts: parts,
	}

	s.cacheGeneratePullRequestsGanttValue(key, result)

	return result, nil
}

func (s *service) generateGanttDrawIOFromPullRequests(pullRequests []*types.PullRequest) ([]byte, error) {
	// Get the repository root directory using runtime.Caller
	_, filename, _, _ := runtime.Caller(0)
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	templatePath := filepath.Join(repoRoot, "diagrams", "gantt", "default.drawio")

	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Parse the template
	var mxFile gantt.MxFile
	if err := xml.Unmarshal(templateData, &mxFile); err != nil {
		return nil, fmt.Errorf("failed to parse template XML: %w", err)
	}

	if len(mxFile.Diagrams) == 0 {
		return nil, errors.New("template has no diagrams")
	}

	// Keep header and calendar cells, remove only task rows (ID 63-241 based on template analysis)
	diagram := &mxFile.Diagrams[0]
	preservedCells := []gantt.MxCell{}
	maxPreservedID := 0

	for _, cell := range diagram.MxGraphModel.Root.Cells {
		// Parse cell ID as integer to check if it's a task row
		if cellIDInt, err := strconv.Atoi(cell.ID); err != nil {
			// Keep non-numeric IDs (like diagram ID)
			preservedCells = append(preservedCells, cell)
		} else if cellIDInt < 63 {
			// Keep only header cells (< 63), remove task rows (63-241) and timeline bars (242-278)
			preservedCells = append(preservedCells, cell)
			if cellIDInt > maxPreservedID {
				maxPreservedID = cellIDInt
			}
		}
		// Skip task rows (ID 63-241)
	}
	diagram.MxGraphModel.Root.Cells = preservedCells

	prNumberCell := gantt.MxCell{}
	taskNameCell := gantt.MxCell{}
	contributorsCell := gantt.MxCell{}
	durationCell := gantt.MxCell{}
	startAtCell := gantt.MxCell{}
	mergedAtCell := gantt.MxCell{}

	for _, cell := range preservedCells {
		if cell.MxGeometry != nil {
			if cell.Value == "PR #" {
				prNumberCell = cell
			}
			if cell.Value == "Task Name" {
				taskNameCell = cell
			}
			if cell.Value == "Contributors" {
				contributorsCell = cell
			}
			if cell.Value == "Duration" {
				durationCell = cell
			}
			if cell.Value == "Created At" {
				startAtCell = cell
			}
			if cell.Value == "Merged At" {
				mergedAtCell = cell
			}
		}
	}

	// Generate cells for pull requests
	startY := 380
	rowHeight := 20

	// Start new IDs after the highest preserved ID to avoid collisions
	nextID := maxPreservedID + 1
	if nextID < 63 {
		nextID = 63 // Ensure we start at least at 63 for task rows
	}

	for i, pr := range pullRequests {
		if pr.CreatedAt == nil || pr.MergedAt == nil {
			continue
		}

		y := startY + i*rowHeight
		yStr := strconv.Itoa(y)
		baseID := nextID // Use current nextID for this PR

		// PR number cell
		numberCell := gantt.MxCell{
			ID:     strconv.Itoa(baseID),
			Value:  fmt.Sprintf("#%d", pr.Number),
			Style:  "strokeColor=#DEEDFF;fillColor=#ADC3D9",
			Parent: "1",
			Vertex: "1",
			MxGeometry: &gantt.MxGeometry{
				X:      prNumberCell.MxGeometry.X,
				Y:      yStr,
				Width:  "50",
				Height: "20",
				As:     "geometry",
			},
		}

		// Task name cell (PR title)
		nameCell := gantt.MxCell{
			ID:     strconv.Itoa(baseID + 1),
			Value:  pr.Title,
			Style:  "align=left;strokeColor=#DEEDFF;fillColor=#ADC3D9",
			Parent: "1",
			Vertex: "1",
			MxGeometry: &gantt.MxGeometry{
				X:      taskNameCell.MxGeometry.X,
				Y:      yStr,
				Width:  taskNameCell.MxGeometry.Width,
				Height: "20",
				As:     "geometry",
			},
		}

		// Contributors cell
		contributorsCell := gantt.MxCell{
			ID:     strconv.Itoa(baseID + 2),
			Value:  pr.FormattedContributors,
			Style:  "align=left;strokeColor=#DEEDFF;fillColor=#ADC3D9",
			Parent: "1",
			Vertex: "1",
			MxGeometry: &gantt.MxGeometry{
				X:      contributorsCell.MxGeometry.X,
				Y:      yStr,
				Width:  contributorsCell.MxGeometry.Width,
				Height: "20",
				As:     "geometry",
			},
		}

		// Duration cell
		duration := pr.MergedAt.Sub(*pr.CreatedAt)
		durationDays := int(duration.Hours() / 24)
		if durationDays == 0 {
			durationDays = 1
		}
		durationText := fmt.Sprintf("%d days", durationDays)
		if durationDays == 1 {
			durationText = "1 day"
		}

		durationCell := gantt.MxCell{
			ID:     strconv.Itoa(baseID + 3),
			Value:  durationText,
			Style:  "strokeColor=#DEEDFF;fillColor=#ADC3D9",
			Parent: "1",
			Vertex: "1",
			MxGeometry: &gantt.MxGeometry{
				X:      durationCell.MxGeometry.X,
				Y:      yStr,
				Width:  durationCell.MxGeometry.Width,
				Height: "20",
				As:     "geometry",
			},
		}

		// Start date cell
		startDateCell := gantt.MxCell{
			ID:     strconv.Itoa(baseID + 4),
			Value:  pr.CreatedAt.Format("02.01.06"),
			Style:  "strokeColor=#DEEDFF;fillColor=#ADC3D9",
			Parent: "1",
			Vertex: "1",
			MxGeometry: &gantt.MxGeometry{
				X:      startAtCell.MxGeometry.X,
				Y:      yStr,
				Width:  startAtCell.MxGeometry.Width,
				Height: "20",
				As:     "geometry",
			},
		}

		// End date cell
		endDateCell := gantt.MxCell{
			ID:     strconv.Itoa(baseID + 5),
			Value:  pr.MergedAt.Format("02.01.06"),
			Style:  "strokeColor=#DEEDFF;fillColor=#ADC3D9",
			Parent: "1",
			Vertex: "1",
			MxGeometry: &gantt.MxGeometry{
				X:      mergedAtCell.MxGeometry.X,
				Y:      yStr,
				Width:  mergedAtCell.MxGeometry.Width,
				Height: "20",
				As:     "geometry",
			},
		}

		// Add all cells to the diagram
		diagram.MxGraphModel.Root.Cells = append(diagram.MxGraphModel.Root.Cells,
			numberCell, nameCell, contributorsCell, durationCell, startDateCell, endDateCell)

		// Increment nextID by 6 for the next PR
		nextID += 6
	}

	// Marshal back to XML
	output, err := xml.MarshalIndent(mxFile, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML declaration
	xmlDeclaration := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	return append(xmlDeclaration, output...), nil
}

// sortPullRequests sorts given pull requests by given
func (s *service) sortPullRequestsAsc(p []*types.PullRequest) []*types.PullRequest {
	prs := []types.PullRequest{}

	for _, pr := range p {
		prs = append(prs, *pr)
	}

	slices.SortFunc(prs, func(a, b types.PullRequest) int {
		if a.Number > b.Number {
			return 1
		}
		if a.Number < b.Number {
			return -1
		}
		return 0
	})

	result := []*types.PullRequest{}

	for _, pr := range prs {
		result = append(result, &pr)
	}

	return result
}

func NewService(cache *cachePkg.Cache, HTTPClient *http.Client) (*service, error) {
	gh := github.NewGitHub()

	srv := &service{
		cache:      cache,
		HTTPClient: HTTPClient,
		GitHub:     gh,
	}

	return srv, nil
}
