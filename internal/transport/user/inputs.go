package user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

type CreateUserInput struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`

	PropertyType string `json:"property_type,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
}

//CREATE TABLE users (
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
//);

func (i *CreateUserInput) ToMap() (map[string]interface{}, error) {
	//TODO Collect all errors and append to array and return
	m := make(map[string]interface{})
	if i.Username != "" {
		m["username"] = i.Username
	} else {
		return nil, errors.New("username is required")
	}

	if i.FirstName != "" {
		m["first_name"] = i.FirstName
	} else {
		return nil, errors.New("first_name is required")
	}

	if i.LastName != "" {
		m["last_name"] = i.LastName
	} else {
		return nil, errors.New("last_name is required")
	}

	if i.Email != "" {
		m["email"] = i.Email
	} else {
		return nil, errors.New("email is required")
	}
	if i.Password != "" {
		PasswordHash := HashPassword(i.Password)
		m["password_hash"] = PasswordHash
	} else {
		return nil, errors.New("password is required")
	}

	if i.PropertyType != "" {
		m["property_type"] = i.PropertyType
	}

	if i.Address != "" {
		m["address"] = i.Address
	}

	if i.Phone != "" {
		m["phone"] = i.Phone
	}

	m["user_role"] = "user"
	m["created_at"] = time.Now()
	return m, nil
}

type LoginUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (i *LoginUserInput) ToMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if i.Username != "" {
		m["username"] = i.Username
	} else {
		return nil, errors.New("username is required")
	}

	if i.Password != "" {
		Hash := HashPassword(i.Password)
		m["password_hash"] = Hash
	}
	return m, nil
}

func HashPassword(password string) string {
	// Создаём хеш SHA256 от пароля
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashBytes := hasher.Sum(nil)

	// Преобразуем хеш в строку в hex-формате
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

type UpdateUserInput struct {
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PropertyType string `json:"property_type,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
}

func (i *UpdateUserInput) ToMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if i.Username != "" {
		m["username"] = i.Username
	}
	if i.FirstName != "" {
		m["first_name"] = i.FirstName
	}
	if i.LastName != "" {
		m["last_name"] = i.LastName
	}
	if i.Email != "" {
		m["email"] = i.Email
	}
	if i.Phone != "" {
		m["phone"] = i.Phone
	}
	if i.PropertyType != "" {
		m["property_type"] = i.PropertyType
	}
	if i.Address != "" {
		m["address"] = i.Address
	}

	return m, nil
}
