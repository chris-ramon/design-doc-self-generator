package api

type Solution struct {
	AnalysisSet AnalysisSet `json:"analysis"`
}

type SolutionSet []Solution

type AnalysisSet []Analysis

type Analysis struct {
	InformationSet InformationSet `json:"information"`
}

type InformationSet []Information

type Information struct {
	GitHub GitHub `json:"github"`
}

type GitHub struct {
	Metrics Metrics `json:"metrics"`
}

type Metrics struct {
	PullRequests PullRequests `json:"pullRequests"`
}

type PullRequests []PullRequest

type PullRequest struct {
	Duration     int          `json:"duration"`
	Contributors Contributors `json:"contributors"`
}

type Contributors []Contributor

type Contributor struct {
	ProfileUrl string `json:"profileUrl"`
}
