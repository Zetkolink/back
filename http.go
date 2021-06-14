package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Zetkolink/back/http/helpers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (b *back) setupHTTPServer(config serverConfig) error {
	config.ReadTimeout *= time.Second
	config.ReadHeaderTimeout *= time.Second
	config.WriteTimeout *= time.Second
	config.IdleTimeout *= time.Second

	apiVersion := "v1"

	r := chi.NewRouter()
	r.Use(middleware.WithValue(helpers.APIVersionContextKey, apiVersion))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	r.Route(
		fmt.Sprintf("%s/%s", helpers.APIPathSuffix, apiVersion),

		func(r chi.Router) {
			r.Group(
				func(r chi.Router) {
					// catCtr := categories.NewController(
					// 	categories.ModelSet{
					// 		Categories: s.models.Categories,
					// 	},
					// 	s.rdb,
					// )
					//
					// r.Mount("/categories", catCtr.NewRouter())
				},
			)
		},
	)

	b.httpServer = &http.Server{
		Addr:              config.Bind,
		Handler:           r,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		MaxHeaderBytes:    config.MaxHeaderBytes,
	}

	return nil
}
