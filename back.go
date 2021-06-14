package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type back struct {
	db         *sql.DB
	httpServer *http.Server
	wg         sync.WaitGroup
}

type models struct {
}

type config struct {
	Database databaseConfig
	Server   serverConfig
	Redis    redisConfig
}

type databaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type serverConfig struct {
	Bind              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

type redisConfig struct {
	Addr string
	Port int
}

func newBack() (*back, error) {
	db, err := sql.Open("postgres", cfg.Database.GetConn())

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Addr, cfg.Redis.Port),
	})

	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}

	a := back{
		db: db,
	}

	err = a.setupHTTPServer(cfg.Server)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (b *back) Run() {
	b.runHTTPServer()

	log.Println(fmt.Sprintf("Starting server at port %s", b.httpServer.Addr))
}

func (b *back) runHTTPServer() {
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		err := b.httpServer.ListenAndServe()

		if err != http.ErrServerClosed {
			b.Stop()
		}
	}()
}

func (b *back) Stop() {
	err := b.httpServer.Shutdown(context.Background())

	if err != nil {
		log.Println(err)
	}

	b.wg.Wait()
}

func (d *databaseConfig) GetConn() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Database,
	)
}
