package handlers

import (
	"AvitoTest/service"
	"AvitoTest/storage"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetTeamHandler handles GET /team/get
func (h *Handlers) GetTeamHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	if r.Method == "GET" {
		fmt.Println("Team Get")

		var res storage.Team
		res.TeamName = r.URL.Query().Get("team_name")

		w.Header().Add("Content-Type", "application/json")
		ans, err := service.GetTeam(h.DB, r.Context(), res)

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
