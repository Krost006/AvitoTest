package storage

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

//func (e *APIError) Error() string { return e.Message }

var (
	// General / transport
	ErrInvalidRequest   = &APIError{Code: "INVALID_REQUEST", Message: "Invalid request", HTTPStatus: http.StatusBadRequest}
	ErrMethodNotAllowed = &APIError{Code: "METHOD_NOT_ALLOWED", Message: "Method not allowed", HTTPStatus: http.StatusMethodNotAllowed}
	ErrInternal         = &APIError{Code: "INTERNAL_ERROR", Message: "Internal server error", HTTPStatus: http.StatusInternalServerError}
	ErrConflict         = &APIError{Code: "CONFLICT", Message: "Conflict", HTTPStatus: http.StatusConflict}
	ErrNotFound         = &APIError{Code: "NOT_FOUND", Message: "Resource not found", HTTPStatus: http.StatusNotFound}
	ErrDB               = &APIError{Code: "DB_ERROR", Message: "Database error", HTTPStatus: http.StatusInternalServerError}

	// Teams
	ErrTeamNotFound      = &APIError{Code: "TEAM_NOT_FOUND", Message: "Team not found", HTTPStatus: http.StatusNotFound}
	ErrTeamAlreadyExists = &APIError{Code: "TEAM_ALREADY_EXISTS", Message: "Team already exists", HTTPStatus: http.StatusConflict}

	// Users
	ErrUserNotFound      = &APIError{Code: "USER_NOT_FOUND", Message: "User not found", HTTPStatus: http.StatusNotFound}
	ErrUserAlreadyExists = &APIError{Code: "USER_ALREADY_EXISTS", Message: "User already exists", HTTPStatus: http.StatusConflict}
	ErrUserInactive      = &APIError{Code: "USER_INACTIVE", Message: "User is not active", HTTPStatus: http.StatusUnprocessableEntity}
	ErrInvalidUserID     = &APIError{Code: "INVALID_USER_ID", Message: "Invalid user_id", HTTPStatus: http.StatusBadRequest}

	// Pull Requests
	ErrPRNotFound           = &APIError{Code: "PR_NOT_FOUND", Message: "Pull request not found", HTTPStatus: http.StatusNotFound}
	ErrPRAlreadyMerged      = &APIError{Code: "PR_ALREADY_MERGED", Message: "Pull request already merged", HTTPStatus: http.StatusUnprocessableEntity}
	ErrPRAlreadyExists      = &APIError{Code: "PR_ALREADY_EXISTS", Message: "Pull request already exists", HTTPStatus: http.StatusConflict}
	ErrCannotModifyMergedPR = &APIError{Code: "CANNOT_MODIFY_MERGED_PR", Message: "Cannot modify reviewers of a merged pull request", HTTPStatus: http.StatusUnprocessableEntity}

	// Reviewers / assignment
	ErrNoAvailableReviewers = &APIError{Code: "NO_AVAILABLE_REVIEWERS", Message: "No available active reviewers in the team", HTTPStatus: http.StatusUnprocessableEntity}
	ErrReviewerNotFound     = &APIError{Code: "REVIEWER_NOT_FOUND", Message: "Reviewer not found", HTTPStatus: http.StatusNotFound}
	ErrReviewerInactive     = &APIError{Code: "REVIEWER_INACTIVE", Message: "Reviewer is not active", HTTPStatus: http.StatusUnprocessableEntity}
	ErrReviewerNotAssigned  = &APIError{Code: "REVIEWER_NOT_ASSIGNED", Message: "Reviewer is not assigned to this pull request", HTTPStatus: http.StatusUnprocessableEntity}
	ErrMaxReviewersReached  = &APIError{Code: "MAX_REVIEWERS_REACHED", Message: "Maximum number of reviewers reached", HTTPStatus: http.StatusConflict}
)

// SendJSONError writes APIError to http.ResponseWriter as JSON.
func SendJSONError(w http.ResponseWriter, apiErr *APIError) {
	if apiErr == nil {
		apiErr = ErrInternal
	}
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.HTTPStatus)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": apiErr,
	})
}
