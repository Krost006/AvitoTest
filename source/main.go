package main

import (
	handlers "AvitoTest/handlers"
	"AvitoTest/storage"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	ctx := context.Background()

	pool, err := storage.InitDB(ctx, 5)

	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		os.Exit(1)
	}

	defer pool.Close()

	h := handlers.NewHandlers(pool)

	h.RegisterRoutes()

	// Register handlers
	// Teams
	//http.HandleFunc("/team/add", h.AddTeamHandler)
	/*http.HandleFunc("/team/get", teams.GetTeamHandler)

	// Users
	http.HandleFunc("/users/setIsActive", users.SetIsActiveHandler)
	http.HandleFunc("/users/getReview", users.GetReviewHandler)

	// PullRequests
	http.HandleFunc("/pullRequest/create", pullrequests.CreatePullRequestHandler)
	http.HandleFunc("/pullRequest/merge", pullrequests.MergePullRequestHandler)
	http.HandleFunc("/pullRequest/reassign", pullrequests.ReassignPullRequestHandler)*/

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
