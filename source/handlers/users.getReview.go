package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetReviewHandler handles GET /users/getReview
func (h *Handlers) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("Get Reviewd")

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]map[string]string{
			"error": {
				"code":    "METHOD_NOT_ALLOWED",
				"message": "Method not allowed",
			},
		})
	}
}
