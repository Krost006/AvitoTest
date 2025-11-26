package service

import (
	storage "AvitoTest/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetUser(pool *pgxpool.Pool, ctx context.Context, user storage.User) (*storage.User, *storage.APIError) {
	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		storage.ErrInternal.Message = err.Error()
		return nil, storage.ErrInternal
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверяем наличие пользователя
	var cnt int64
	if err := tx.QueryRow(ctx, "SELECT count(1) FROM users WHERE user_id=$1", user.UserID).Scan(&cnt); err != nil {
		storage.ErrNotFound.Message = err.Error()
		return nil, storage.ErrNotFound
	}

	if cnt == 0 {
		return nil, storage.ErrUserNotFound
	}

	if _, err := tx.Exec(ctx,
		"UPDATE users SET is_active=$2, updated_at=now() WHERE user_id=$1",
		user.UserID, user.IsActive); err != nil {
		storage.ErrInternal.Message = err.Error()
		return nil, storage.ErrInternal
	}

	var ans storage.User
	err = tx.QueryRow(ctx,
		"SELECT user_id, username, team_name, is_active FROM users WHERE user_id=$1",
		user.UserID).Scan(&ans.UserID, &ans.Username, &ans.TeamName, &ans.IsActive)
	if err != nil {
		storage.ErrInternal.Message = err.Error()

		return nil, storage.ErrInternal
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		storage.ErrInternal.Message = err.Error()
		return nil, storage.ErrInternal
	}

	fmt.Println("Set Active: finished")
	return &ans, nil
}
