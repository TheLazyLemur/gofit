package views

import (
	"context"

	"github.com/TheLazyLemur/gofit/src/internal/db"
)

func getLoggedInUser(ctx context.Context) (db.User, bool) {
	user, ok := ctx.Value("user").(db.User)
	return user, ok
}
