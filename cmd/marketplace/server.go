package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-marketplace/internal/api"
	"github.com/mattermost/mattermost-marketplace/internal/store"
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var instanceID string

func init() {
	instanceID = model.NewId()

	serverCmd.PersistentFlags().String("database", "plugins.json", "The read-only JSON file backing the server.")
	serverCmd.PersistentFlags().String("listen", ":8085", "The interface and port on which to listen.")
	serverCmd.PersistentFlags().Bool("debug", false, "Whether to output debug logs.")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the provisioning server.",
	RunE: func(command *cobra.Command, args []string) error {
		command.SilenceUsage = true

		debug, _ := command.Flags().GetBool("debug")
		if debug {
			logger.SetLevel(logrus.DebugLevel)
		}

		database, _ := command.Flags().GetString("database")
		databaseFile, err := os.Open(database)
		if err != nil {
			return errors.Wrapf(err, "failed to open %s", database)
		}
		defer databaseFile.Close()

		fileStore, err := store.New(databaseFile, logger)
		if err != nil {
			return errors.Wrap(err, "failed to initialize store")
		}

		logger := logger.WithField("instance", instanceID)
		logger.Info("Starting Plugin Marketplace")

		router := mux.NewRouter()

		api.Register(router, &api.Context{
			Store:  fileStore,
			Logger: logger,
		})

		listen, _ := command.Flags().GetString("listen")
		srv := &http.Server{
			Addr:           listen,
			Handler:        router,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    time.Second * 60,
			MaxHeaderBytes: 1 << 20,
			ErrorLog:       log.New(&logrusWriter{logger}, "", 0),
		}

		go func() {
			logger.WithField("addr", srv.Addr).Info("Listening")
			listenErr := srv.ListenAndServe()
			if listenErr != nil && listenErr != http.ErrServerClosed {
				logger.WithField("err", listenErr).Error("Failed to listen and serve")
			}
		}()

		c := make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c
		logger.Info("Shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err = srv.Shutdown(ctx)
		if err != nil {
			logger.WithField("err", err).Error("Failed to shutdown")
		}

		return nil
	},
}
