package types

type Solution struct {
	AnalysisSet AnalysisSet
}

type SolutionSet []Solution

type AnalysisSet []Analysis

type Analysis struct {
	InformationSet InformationSet
}

type InformationSet []Information

type Information struct {
	GitHub GitHub
}

type GitHub struct {
	Metrics Metrics
}

type Metrics struct {
	PullRequests PullRequests
}

type PullRequests []PullRequest

type PullRequest struct {
	Duration     int
	Contributors Contributors
}

type Contributors []Contributor

type Contributor struct {
	ProfileUrl string
}
