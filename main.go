package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	user2 "minibank-backend/internal/service/user"
	"minibank-backend/internal/storage/user"
	user3 "minibank-backend/internal/transport/user"
	"minibank-backend/pkg/auth"
	"minibank-backend/pkg/db"
	"minibank-backend/pkg/server"
)

func main() {
	err := Init()
	if err != nil {
		return
	}

	rtr := httprouter.New()
	PGConfig, err := db.InitPGConfig()
	if err != nil {
		return
	}
	PGClient, err := db.GetPGClient(context.Background(), PGConfig)
	if err != nil {
		return
	}

	options := cors.Options{
		AllowedOrigins:         []string{"http://localhost:3000", "http://185.185.68.187", "https://keys-store.online", "*"},
		AllowOriginFunc:        nil,
		AllowOriginRequestFunc: nil,
		AllowedMethods:         []string{"POST", "PATCH", "GET", "DELETE"},
		AllowedHeaders:         []string{"Access-Control-Allow-Origin", "Authorization", "Content-Type"},
		ExposedHeaders:         nil,
		MaxAge:                 0,
		AllowCredentials:       true,
		AllowPrivateNetwork:    false,
		OptionsPassthrough:     false,
		OptionsSuccessStatus:   0,
		Debug:                  false,
	}
	c := cors.New(options)
	handler := c.Handler(rtr)

	TokenManager := auth.NewTokenManager()
	MiddleWare := auth.NewMiddleWare(TokenManager)

	UserStorage := user.NewUserStorage(PGClient)

	UserService := user2.NewUserService(UserStorage, TokenManager)

	user3.NewUserHandler(MiddleWare, UserService).Register(rtr)

	server := server.NewHTTPServer(handler)
	err = server.Run()
	if err != nil {
		return
	}
}

func Init() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
