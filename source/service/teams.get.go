package service

import (
	storage "AvitoTest/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTeam(pool *pgxpool.Pool, ctx context.Context, team storage.Team) (*storage.Team, *storage.APIError) {
	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, storage.ErrInternal
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверяем наличие команды
	var cnt int64
	if err := tx.QueryRow(ctx, "SELECT count(1) FROM teams WHERE team_name=$1", team.TeamName).Scan(&cnt); err != nil {
		return nil, storage.ErrNotFound
	}

	if cnt == 0 {
		return nil, storage.ErrInternal
	}

	// Получаем команду
	rows, err := tx.Query(ctx, "SELECT user_id, username, is_active FROM users WHERE team_name=$1", team.TeamName)
	if err != nil {
		return nil, storage.ErrInternal
	}
	defer rows.Close()

	var ans storage.Team
	ans.TeamName = team.TeamName

	for rows.Next() {
		var user storage.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.IsActive); err != nil {
			return nil, storage.ErrInternal
		}
		ans.Members = append(ans.Members, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, storage.ErrInternal
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return nil, storage.ErrInternal
	}

	fmt.Println("GetTeam: finished")
	return &ans, nil
}
