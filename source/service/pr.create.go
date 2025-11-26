package service

import (
	storage "AvitoTest/storage"
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PrCreate(pool *pgxpool.Pool, ctx context.Context, pr storage.PullRequest) (*storage.PullRequest, *storage.APIError) {
	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, storage.ErrInternal
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверяем наличие pr
	var cnt int64
	if err := tx.QueryRow(ctx, "SELECT count(1) FROM users WHERE user_id=$1", pr.AuthorID).Scan(&cnt); err != nil {
		return nil, storage.ErrInternal
	}

	if cnt == 0 {
		return nil, storage.ErrUserNotFound
	}

	// Проверяем наличие pr
	if err := tx.QueryRow(ctx, "SELECT count(1) FROM pull_requests WHERE pull_request_id=$1", pr.PullRequestID).Scan(&cnt); err != nil {
		return nil, storage.ErrInternal
	}

	if cnt == 0 {
		if _, err := tx.Exec(ctx,
			"INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, created_at) VALUES ($1, $2, $3, now())",
			pr.PullRequestID, pr.PullRequestName, pr.AuthorID); err != nil {
			return nil, storage.ErrInternal
		}
	} else {
		err := storage.ErrTeamAlreadyExists
		err.Message = "PR " + pr.PullRequestID + " already exists"
		return nil, err

	}
	//Назначаем ревюеров
	var potentional_reviewers []string
	var reviewers []string

	rows, err := tx.Query(ctx,
		"SELECT user_id FROM users"+
			"WHERE team_name=(SELECT team_name FROM users WHERE user_id=$1) "+
			"AND is_active=true",
		pr.AuthorID)
	if err != nil {
		return nil, storage.ErrUserNotFound
	}
	defer rows.Close()
	for rows.Next() {
		var tmp string
		rows.Scan(&tmp)
		potentional_reviewers = append(potentional_reviewers, tmp)
	}

	if len(potentional_reviewers) < 2 {
		reviewers = append(reviewers, potentional_reviewers[0])
	} else {
		v1 := rand.Intn(len(potentional_reviewers))
		v2 := rand.Intn(len(potentional_reviewers))
		reviewers = append(reviewers, potentional_reviewers[v1])
		if v1 == v2 {
			v2 = (v1 + 1) % len(potentional_reviewers)
		}
		reviewers = append(reviewers, potentional_reviewers[v2])
	}

	// Вставляем ревьюеров
	for _, v := range reviewers {
		if _, err := tx.Exec(ctx, "INSERT INTO pull_request_reviewers VALUES (default, $1, $2, now() )", pr.PullRequestID, v); err != nil {
			return nil, storage.ErrInternal
		}
	}

	// Собираем результат
	var ans storage.PullRequest

	rows, err = tx.Query(ctx,
		"SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id=$1",
		pr.PullRequestID)
	if err != nil {
		return nil, storage.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var tmp string
		if err := rows.Scan(&tmp); err != nil {
			return nil, storage.ErrInternal
		}

		ans.AssignedReviewers = append(ans.AssignedReviewers, tmp)
	}

	if err := tx.QueryRow(ctx,
		"SELECT pull_request_id, pull_request_name, author_id, status"+
			" FROM pull_requests WHERE pull_request_id=$1",
		pr.PullRequestID).Scan(&ans.PullRequestID, &ans.PullRequestName, &ans.AuthorID, &ans.Status); err != nil {
		return nil, storage.ErrInternal
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return nil, storage.ErrInternal
	}

	fmt.Println("CreatePR: finished")
	return &ans, nil
}
