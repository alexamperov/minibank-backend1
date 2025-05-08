package apply

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"minibank-backend/internal/entity"
)

type ApplyStorage struct {
	conn *pgxpool.Pool
}

func (a *ApplyStorage) GetAppliesOfEmployee(ctx context.Context, UserID int) ([]entity.EApply, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ApplyStorage) GetAppliesOfUser(ctx context.Context, UserID int) ([]entity.EApply, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ApplyStorage) UpdateStatus(ctx context.Context, UserID int, ApplyID int, Status string) error {
	//TODO implement me
	panic("implement me")
}

func (a *ApplyStorage) CreateApply(ctx context.Context, UserID int, applyData map[string]interface{}) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ApplyStorage) GetApplyByID(ctx context.Context, UserID int, applyID int) (*entity.EApply, error) {
	//TODO implement me
	panic("implement me")
}
