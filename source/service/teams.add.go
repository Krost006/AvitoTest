package service

import (
	storage "AvitoTest/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddTeam(pool *pgxpool.Pool, ctx context.Context, team storage.Team) (*storage.Team, *storage.APIError) {
	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, storage.ErrInternal
	}
	// Всегда откатываем при ошибке (коммит отменит откат)
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверяем наличие команды
	var cnt int64
	if err := tx.QueryRow(ctx, "SELECT count(1) FROM teams WHERE team_name=$1", team.TeamName).Scan(&cnt); err != nil {
		return nil, storage.ErrInternal
	}

	if cnt == 0 {
		if _, err := tx.Exec(ctx, "INSERT INTO teams (team_name, created_at, updated_at) VALUES ($1, now(), now())", team.TeamName); err != nil {
			return nil, storage.ErrInternal
		}
	} else {
		err := storage.ErrTeamAlreadyExists
		err.Message = team.TeamName + "already exists"
		return nil, err

	}

	// Вставляем/обновляем участников
	for _, v := range team.Members {
		if v == nil || v.UserID == "" || v.Username == "" {
			return nil, storage.ErrInternal
		}

		var ucnt int64
		if err := tx.QueryRow(ctx, "SELECT count(1) FROM users WHERE user_id=$1", v.UserID).Scan(&ucnt); err != nil {
			return nil, storage.ErrInternal
		}

		if ucnt == 0 {
			if _, err := tx.Exec(ctx,
				"INSERT INTO users (user_id, username, team_name, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, now(), now())",
				v.UserID, v.Username, team.TeamName, v.IsActive); err != nil {
				return nil, storage.ErrInternal
			}
		} else {
			if _, err := tx.Exec(ctx,
				"UPDATE users SET username=$2, team_name=$3, is_active=$4, updated_at=now() WHERE user_id=$1",
				v.UserID, v.Username, team.TeamName, v.IsActive); err != nil {
				return nil, storage.ErrInternal
			}
		}
	}

	// Собираем результат в рамках той же транзакции
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

	fmt.Println("AddTeam: finished")
	return &ans, nil
}
