package handlers

import (
	"AvitoTest/service"
	"AvitoTest/storage"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetReviewHandler handles GET /users/getReview
func (h *Handlers) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("Team Get")

		var res storage.User
		res.TeamName = r.URL.Query().Get("user_id")

		w.Header().Add("Content-Type", "application/json")
		ans, err := service.GetUser(h.DB, r.Context(), res)

		if err != nil {
			storage.SendJSONError(w, err)
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(ans)
		}
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
