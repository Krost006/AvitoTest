package storage

type User struct {
	UserID       string         `json:"user_id"`
	Username     string         `json:"username"`
	TeamName     string         `json:"team_name"`
	IsActive     bool           `json:"is_active"`
	PullRequests []*PullRequest `json:"pull_requests"`
}

type Team struct {
	TeamName string  `json:"team_name"`
	Members  []*User `json:"members"`
}

type PullRequestStatus string

const (
	PRStatusOpen   PullRequestStatus = "OPEN"
	PRStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	Status            PullRequestStatus
	AssignedReviewers []string
	CreatedAt         string
	MergedAt          *string
}
