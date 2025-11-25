package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MergePullRequestHandler handles POST /pullRequest/merge
func (h *Handlers) MergePullRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Team added")
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Println("Read error")
			return
		}

		//var res requestTeamsAdd
		var res map[string]interface{}
		json.Unmarshal(body, &res)
		fmt.Println("Readed msg:")
		fmt.Println(res)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(201)

		json.NewEncoder(w).Encode(http.StatusCreated)
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
