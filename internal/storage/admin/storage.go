package admin

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminStorage struct {
	conn *pgxpool.Pool
}

func (a *AdminStorage) PaySalary(ctx context.Context, UserID int, Sum int) error {
	//TODO implement me
	panic("implement me")
}
