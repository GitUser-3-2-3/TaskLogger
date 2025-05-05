package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (bknd *backend) serve() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	srv := &http.Server{
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Addr:         fmt.Sprintf(":%d", bknd.config.port),
		Handler:      bknd.routes(),
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  1 * time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	bknd.logger.Info().Msgf("server started at: %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	bknd.logger.Info().Msgf("server running at: %s stopped", srv.Addr)
	return nil
}
