package user

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"minibank-backend/internal/entity"
)

//id SERIAL PRIMARY KEY,
//username VARCHAR(255),
//first_name VARCHAR(255),
//last_name VARCHAR(255),
//role VARCHAR(50),
//password_hash VARCHAR(255) NOT NULL,
//property_type VARCHAR(100),
//address TEXT,
//phone VARCHAR(50),
//created_at TIMESTAMPTZ,

type UserStorage struct {
	conn *pgxpool.Pool
}

func NewUserStorage(conn *pgxpool.Pool) *UserStorage {
	return &UserStorage{conn: conn}
}

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (u *UserStorage) GetUserByID(ctx context.Context, UserID int) (entity.EUser, error) {

	var out = entity.EUser{}
	query, args, err := sq.Select(
		"username", "first_name", "last_name",
		"user_role", "property_type",
		"address", "phone", "created_at").From("users").Where(squirrel.Eq{"id": UserID}).ToSql()
	if err != nil {
		return entity.EUser{}, err
	}

	err = u.conn.QueryRow(ctx, query, args...).Scan(
		&out.Username, &out.FirstName, &out.LastName,
		&out.Role, &out.PropertyType,
		&out.Address, &out.Phone, &out.CreatedAt)
	if err != nil {
		log.Println(err)
		return entity.EUser{}, err
	}
	return out, nil
}

func (u *UserStorage) GetUsers(ctx context.Context) ([]entity.EUser, error) {
	query, args, err := sq.Select(
		"username", "first_name", "last_name",
		"user_role", "property_type",
		"address", "phone", "created_at").From("users").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := u.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []entity.EUser{}
	for rows.Next() {
		var user entity.EUser
		err = rows.Scan(&user.Username, &user.FirstName, &user.LastName, &user.Role, &user.PropertyType, &user.Address, &user.CreatedAt, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserStorage) CreateUser(ctx context.Context, data map[string]interface{}) (int, error) {
	var ID int

	query, args, err := sq.Insert("users").Suffix("RETURNING id").SetMap(data).ToSql()
	if err != nil {
		return 0, err
	}
	err = u.conn.QueryRow(ctx, query, args...).Scan(&ID)
	if err != nil {
		return 0, err
	}

	return ID, err
}

func (u *UserStorage) GetUserIDByCredentials(ctx context.Context, username string, passwordHash string) (int, string, error) {
	var ID int
	var Role string
	query, args, err := sq.Select("id", "user_role").From("users").Where(squirrel.Eq{"password_hash": passwordHash, "username": username}).ToSql()
	if err != nil {
		return 0, "", err
	}

	err = u.conn.QueryRow(ctx, query, args...).Scan(&ID, &Role)
	if err != nil {
		return 0, "", err
	}
	return ID, Role, nil
}
