package service

import (
	storage "AvitoTest/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUser(pool *pgxpool.Pool, ctx context.Context, user storage.User) (*storage.User, *storage.APIError) {
	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, storage.ErrInternal
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверяем наличие пользователя
	var cnt int64
	if err := tx.QueryRow(ctx, "SELECT count(1) FROM users WHERE user_id=$1", user.UserID).Scan(&cnt); err != nil {
		return nil, storage.ErrInternal
	}

	if cnt == 0 {
		return nil, storage.ErrUserNotFound
	}

	var ans storage.User

	rows, err := tx.Query(ctx,
		"SELECT r.pull_request_id, r.pull_request_name, r.author_id, r.status "+
			"FROM pull_requests r"+
			" JOIN (SELECT pull_request_id id FROM pull_request_reviewers WHERE reviewer_id=$1) i "+
			"ON r.author_id = i.id"+
			" WHERE team_name=$1",
		user.UserID)

	for rows.Next() {
		var tmp storage.PullRequest

		if err := rows.Scan(&tmp); err != nil {
			return nil, storage.ErrInternal
		}
		ans.PullRequests = append(ans.PullRequests, &tmp)
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return nil, storage.ErrInternal
	}

	fmt.Println("Get Reviewed: finished")
	return &ans, nil
}
