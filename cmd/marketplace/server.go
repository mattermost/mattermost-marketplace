package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mattermost/mattermost-marketplace/internal/api"
	"github.com/mattermost/mattermost-marketplace/internal/store"
)

var (
	instanceID string

	// upstreamURL may be compiled into the binary by defining $BUILD_UPSTREAM_URL
	upstreamURL string
)

func init() {
	instanceID = model.NewId()

	serverCmd.PersistentFlags().String("database", "plugins.json", "The read-only JSON file backing the server.")
	serverCmd.PersistentFlags().String("database-url", "", "A remote URL to fetch plugins.json from (e.g. a raw GitLab/GitHub URL). Overrides --database.")
	serverCmd.PersistentFlags().Duration("database-refresh-interval", 5*time.Minute, "How often to re-fetch plugins from --database-url.")
	serverCmd.PersistentFlags().String("listen", ":8085", "The interface and port on which to listen.")
	serverCmd.PersistentFlags().String("upstream", upstreamURL, "An upstream marketplace server with which to merge results.")
	serverCmd.PersistentFlags().Bool("debug", false, "Whether to output debug logs.")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the provisioning server.",
	RunE: func(command *cobra.Command, _ []string) error {
		command.SilenceUsage = true

		debug, _ := command.Flags().GetBool("debug")
		if debug {
			logger.SetLevel(logrus.DebugLevel)
		}

		var apiStore store.Store

		databaseURL, _ := command.Flags().GetString("database-url")
		if databaseURL != "" && strings.HasPrefix(databaseURL, "http") {
			refreshInterval, _ := command.Flags().GetDuration("database-refresh-interval")
			logger.WithFields(logrus.Fields{
				"url":      databaseURL,
				"interval": refreshInterval,
			}).Info("Using remote database URL")

			remoteStore, remoteErr := store.NewRemote(databaseURL, refreshInterval, logger)
			if remoteErr != nil {
				return errors.Wrap(remoteErr, "failed to initialize remote store")
			}
			defer remoteStore.Stop()
			apiStore = remoteStore
		} else {
			database, _ := command.Flags().GetString("database")
			databaseFile, fileErr := os.Open(database)
			if fileErr != nil {
				return errors.Wrapf(fileErr, "failed to open %s", database)
			}
			defer databaseFile.Close()

			var staticErr error
			apiStore, staticErr = store.NewStaticFromReader(databaseFile, logger)
			if staticErr != nil {
				return errors.Wrap(staticErr, "failed to initialize store")
			}
		}

		upstreamURL, _ := command.Flags().GetString("upstream")
		if upstreamURL != "" {
			upstreamStore, err := store.NewProxy(upstreamURL, logger)
			if err != nil {
				return errors.Wrap(err, "failed to initialize upstream store")
			}

			logger.WithField("upstream", upstreamURL).Info("Proxying to upstream marketplace")

			apiStore = store.NewMerged(logger, apiStore, upstreamStore)
		}

		logger := logger.WithField("instance", instanceID)
		logger.Info("Starting Plugin Marketplace")

		router := mux.NewRouter()

		api.Register(router, &api.Context{
			Store:  apiStore,
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
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithField("err", err).Error("Failed to shutdown")
		}

		return nil
	},
}
