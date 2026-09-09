package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"sync"

	"carbonite/admin/internal/api"
	"carbonite/admin/internal/data"
	"carbonite/admin/internal/handler"
	"carbonite/admin/internal/redis"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	redisstore "github.com/rbcervilla/redisstore/v9"
)

type (
	application struct {
		requester   handler.Requester
		store       handler.Storer
		cacher      data.Cacher
		env         string
		sessionName string
		wg          sync.WaitGroup
	}

	options struct {
		port *int
		addr *string
	}

	Option func(options *options) error
)

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("port should be positive")
		}
		options.port = &port
		return nil
	}
}

func WithAddr(addr string) Option {
	return func(options *options) error {
		options.addr = &addr
		return nil
	}
}

const (
	version         = "1.0.0"
	defaultHTTPPort = 3000
	defaultAddr     = "127.0.0.1"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Client API
	requester, err := api.New(
		os.Getenv("API_CLIENT_ID"),
		os.Getenv("API_CLIENT_SECRET"),
		os.Getenv("API_URL"),
		os.Getenv("ENV"),
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Client Redis
	cacher, err := redis.New(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	store, err := redisstore.NewRedisStore(context.Background(), cacher)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}

	store.KeyPrefix("session_")
	store.Options(sessions.Options{
		Path: "/",
		// Domain: "planet55.fr",
		Secure:   true,
		MaxAge:   600,
		HttpOnly: true,
	})

	app := &application{
		requester:   requester,
		store:       store,
		cacher:      cacher,
		env:         os.Getenv("ENV"),
		sessionName: os.Getenv("SESSION_NAME"),
	}

	// Launch the Web Server
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	addr := os.Getenv("ADDR")
	if err = app.serve(WithPort(port), WithAddr(addr)); err != nil {
		log.Fatal(err)
	}
}
