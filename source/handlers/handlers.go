package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	DB *pgxpool.Pool
}

func NewHandlers(db *pgxpool.Pool) *Handlers {
	return &Handlers{DB: db}
}

func (h *Handlers) RegisterRoutes() {
	// Teams
	http.HandleFunc("/team/add", h.AddTeamHandler)
	http.HandleFunc("/team/get", h.GetTeamHandler)

	// Users
	http.HandleFunc("/users/setIsActive", h.SetIsActiveHandler)
	http.HandleFunc("/users/getReview", h.GetReviewHandler)

	// PullRequests
	http.HandleFunc("/pullRequest/create", h.CreatePullRequestHandler)
	http.HandleFunc("/pullRequest/merge", h.MergePullRequestHandler)
	http.HandleFunc("/pullRequest/reassign", h.ReassignPullRequestHandler)
}

// Teams handlers
/*func (h *Handlers) AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	teams.AddTeamHandler(h.DB)(w, r)
}*/

/*func (h *Handlers) GetTeamHandler(w http.ResponseWriter, r *http.Request) {
	teams.GetTeamHandler(h.DB)(w, r)
}*/

// Users handlers
/*func (h *Handlers) SetIsActiveHandler(w http.ResponseWriter, r *http.Request) {
	users.SetIsActiveHandler(h.DB)(w, r)
}

func (h *Handlers) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	users.GetReviewHandler(h.DB)(w, r)
}

// PullRequests handlers
func (h *Handlers) CreatePullRequestHandler(w http.ResponseWriter, r *http.Request) {
	pullrequests.CreatePullRequestHandler(h.DB)(w, r)
}

func (h *Handlers) MergePullRequestHandler(w http.ResponseWriter, r *http.Request) {
	pullrequests.MergePullRequestHandler(h.DB)(w, r)
}

func (h *Handlers) ReassignPullRequestHandler(w http.ResponseWriter, r *http.Request) {
	pullrequests.ReassignPullRequestHandler(h.DB)(w, r)
}*/

func (h *Handlers) Close() {
	if h.DB != nil {
		h.DB.Close()
	}
}
