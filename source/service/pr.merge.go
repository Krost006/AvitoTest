package service

import (
	storage "AvitoTest/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PrMerge(pool *pgxpool.Pool, ctx context.Context, pr storage.PullRequest) (*storage.PullRequest, *storage.APIError) {
	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, storage.ErrInternal
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверяем наличие pr
	var cnt int64
	if err := tx.QueryRow(ctx, "SELECT count(1) FROM pull_requests WHERE pull_request_id=$1", pr.PullRequestID).Scan(&cnt); err != nil {
		return nil, storage.ErrInternal
	}

	if cnt == 0 {
		err := storage.ErrNotFound
		err.Message = "PR " + pr.PullRequestID + " not found"
		return nil, err
	}

	if _, err := tx.Exec(ctx, "UPDATE  pull_requests SET status='MERGED' WHERE pull_request_id=$1", pr.PullRequestID); err != nil {
		return nil, storage.ErrInternal
	}

	// Собираем результат
	var ans storage.PullRequest

	rows, err := tx.Query(ctx,
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
