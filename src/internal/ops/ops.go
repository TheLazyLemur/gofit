package ops

import (
	"context"
	"database/sql"

	"github.com/TheLazyLemur/gofit/src/internal/db"
	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, dbc *sql.DB, querier db.Querier, username, password, email string) (token string, err error) {
	tx, err := dbc.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	userID, err := querier.CreateUser(ctx, dbc, db.CreateUserParams{
		Name:         username,
		Email:        email,
		PasswordHash: password,
	})
	if err != nil {
		return "", err
	}

	return createSession(ctx, querier, tx, userID)
}

func LoginUser(ctx context.Context, dbc *sql.DB, querier db.Querier, email, password string) (string, error) {
	user, err := querier.GetUserByEmailAndPassword(ctx, dbc, db.GetUserByEmailAndPasswordParams{
		Email:        email,
		PasswordHash: password,
	})
	if err != nil {
		return "", err
	}

	return createSession(ctx, querier, dbc, user.ID)
}

func createSession(ctx context.Context, q db.Querier, dbtx db.DBTX, userID int64) (string, error) {
	uuid := uuid.New().String()

	_, err := q.CreateSession(ctx, dbtx, db.CreateSessionParams{
		UserID: userID,
		Token:  uuid,
	})
	if err != nil {
		return "", err
	}

	return uuid, nil
}
