package ops

import (
	"context"
	"database/sql"

	"github.com/TheLazyLemur/gofit/src/internal/db"
	"github.com/google/uuid"
)

type dependencies interface {
	DBC() *sql.DB
	Querier() db.Querier
}

func CreateUser(ctx context.Context, d dependencies, username, password, email string) (token string, err error) {
	tx, err := d.DBC().Begin()
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

	userID, err := d.Querier().CreateUser(ctx, d.DBC(), db.CreateUserParams{
		Name:         username,
		Email:        email,
		PasswordHash: password,
	})
	if err != nil {
		return "", err
	}

	uuid := uuid.New().String()

	_, err = d.Querier().CreateSession(ctx, tx, db.CreateSessionParams{
		UserID: userID,
		Token:  uuid,
	})
	if err != nil {
		return "", err
	}

	return uuid, nil
}
