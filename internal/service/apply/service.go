package apply

import (
	"context"
	"minibank-backend/internal/entity"
)

type IApplyStorage interface {
	GetAppliesOfEmployee(ctx context.Context, UserID int) ([]entity.EApply, error)
	GetAppliesOfUser(ctx context.Context, UserID int) ([]entity.EApply, error)
	UpdateStatus(ctx context.Context, UserID int, ApplyID int, Status string) error
	CreateApply(ctx context.Context, UserID int, applyData map[string]interface{}) (int, error)
	GetApplyByID(ctx context.Context, UserID int, applyID int) (*entity.EApply, error)
}

type ApplyService struct {
	applyStorage IApplyStorage
}

func (a *ApplyService) GetAllApplies(ctx context.Context, UserID int, isUser bool) ([]entity.EApply, error) {
	if isUser {
		return a.applyStorage.GetAppliesOfUser(ctx, UserID)
	} else {
		return a.applyStorage.GetAppliesOfEmployee(ctx, UserID)
	}
}

func (a *ApplyService) GetApplyByID(ctx context.Context, UserID int, applyID int) (*entity.EApply, error) {
	return a.applyStorage.GetApplyByID(ctx, UserID, applyID)
}

func (a *ApplyService) CreateApply(ctx context.Context, UserID int, applyData map[string]interface{}) (int, error) {
	return a.applyStorage.CreateApply(ctx, UserID, applyData)
}

func (a *ApplyService) AcceptApply(ctx context.Context, UserID int, applyID int) error {
	return a.applyStorage.UpdateStatus(ctx, UserID, applyID, "accepted")
}

func (a *ApplyService) DenyApply(ctx context.Context, UserID int, applyID int) error {
	return a.applyStorage.UpdateStatus(ctx, UserID, applyID, "denied")
}
