package user

import (
	"context"
	"log"
	"minibank-backend/internal/entity"
	"minibank-backend/pkg/auth"
)

type IUserStorage interface {
	GetUserByID(ctx context.Context, UserID int) (entity.EUser, error)
	GetUsers(ctx context.Context) ([]entity.EUser, error)
	CreateUser(ctx context.Context, data map[string]interface{}) (int, error)
	GetUserIDByCredentials(ctx context.Context, username string, passwordHash string) (int, string, error)
}

type UserService struct {
	userStorage IUserStorage
	tm          auth.TokenManager
}

func NewUserService(userStorage IUserStorage, tm auth.TokenManager) *UserService {
	return &UserService{userStorage: userStorage, tm: tm}
}

func (u *UserService) GetUserByID(ctx context.Context, UserID int) (entity.EUser, error) {

	return u.userStorage.GetUserByID(ctx, UserID)
}

func (u *UserService) GetUsers(ctx context.Context) ([]entity.EUser, error) {
	return u.userStorage.GetUsers(ctx)
}

// TODO NOT IMPLEMENTED
func (u *UserService) UpdateUserByID(ctx context.Context, UserID int) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) SignUp(ctx context.Context, data map[string]interface{}) (int, string, error) {
	UserID, err := u.userStorage.CreateUser(ctx, data)
	if err != nil {
		return 0, "", err
	}
	token, err := u.tm.GenerateToken(UserID, "user")
	if err != nil {
		return 0, "", err
	}

	return UserID, token, nil

}

func (u *UserService) SignIn(ctx context.Context, data map[string]interface{}) (int, string, error) {
	UserID, Role, err := u.userStorage.GetUserIDByCredentials(ctx, data["username"].(string), data["password_hash"].(string))
	if err != nil {
		return 0, "", err
	}
	log.Println("UserID, Role = ", UserID, Role)

	token, err := u.tm.GenerateToken(UserID, Role)
	if err != nil {
		return 0, "", err
	}
	return UserID, token, nil
}
