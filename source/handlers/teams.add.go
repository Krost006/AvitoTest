package handlers

import (
	"AvitoTest/service"
	"AvitoTest/storage"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AddTeamHandler handles POST /team/add
func (h *Handlers) AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Team added")
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)

		var res storage.Team
		json.Unmarshal(body, &res)
		fmt.Println("Readed msg:")
		fmt.Println(res)

		w.Header().Add("Content-Type", "application/json")

		ans, err := service.AddTeam(h.DB, r.Context(), res)

		if err != nil {
			storage.SendJSONError(w, err)
		} else {
			w.WriteHeader(201)
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
