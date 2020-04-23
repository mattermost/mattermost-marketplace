package main

import (
	"github.com/akrylysov/algnhsa"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"

	_ "github.com/mattermost/mattermost-marketplace/data/statik"

	"github.com/mattermost/mattermost-marketplace/internal/api"
	"github.com/mattermost/mattermost-marketplace/internal/store"
)

var (
	// upstreamURL may be compiled into the binary by defining $BUILD_UPSTREAM_URL
	upstreamURL = ""
)

var logger *logrus.Logger

func main() {
	err := listenAndServe()
	if err != nil {
		panic("failed to listen and serve: " + err.Error())
	}
}

func newStatikStore(statikPath string, logger logrus.FieldLogger) (*store.StaticStore, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open statik fileystem")
	}

	database, err := statikFS.Open(statikPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %s", database)
	}
	defer database.Close()

	statikStore, err := store.NewStaticFromReader(database, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize store")
	}

	return statikStore, nil
}

func listenAndServe() error {
	logger = logrus.New()

	var stores []store.Store

	statikStore, err := newStatikStore("/plugins.json", logger)
	if err != nil {
		return err
	}
	stores = append(stores, statikStore)

	if upstreamURL != "" {
		upstreamStore, err := store.NewProxy(upstreamURL, logger)
		if err != nil {
			return errors.Wrap(err, "failed to initialize upstream store")
		}

		stores = append(stores, upstreamStore)
	}

	router := mux.NewRouter()
	api.Register(router, &api.Context{
		Store:  store.NewMerged(logger, stores...),
		Logger: logger,
	})

	algnhsa.ListenAndServe(router, &algnhsa.Options{
		UseProxyPath: true,
	})

	return nil
}
