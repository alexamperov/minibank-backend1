package admin

import (
	"context"
	"minibank-backend/internal/entity"
)

type AdminService struct {
	AdminStorage IAdminStorage
	UserStorage  IUserStorage
}

type IUserStorage interface {
	GetUserByID(ctx context.Context, UserID int) (entity.EUser, error)
	GetUsers(ctx context.Context) ([]entity.EUser, error)
}

type IAdminStorage interface {
	PaySalary(ctx context.Context, UserID int, Sum int) error
}

func (a *AdminService) GetEmployees(ctx context.Context) ([]entity.EUser, error) {
	users, err := a.UserStorage.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	var filteredUsers []entity.EUser
	for _, user := range users {
		if user.Role == "employee" {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers, nil
}

func (a *AdminService) PaySalary(ctx context.Context, UserID int, Sum int) error {
	return a.AdminStorage.PaySalary(ctx, UserID, Sum)
}

func (a *AdminService) GetEmployee(ctx context.Context, EmployeeID int) (entity.EUser, error) {
	return a.UserStorage.GetUserByID(ctx, EmployeeID)
}
