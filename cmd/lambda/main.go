package main

import (
	"bytes"
	_ "embed"

	"github.com/akrylysov/algnhsa"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/mattermost/mattermost-marketplace/internal/api"
	"github.com/mattermost/mattermost-marketplace/internal/store"
)

var (
	// upstreamURL may be compiled into the binary by defining $BUILD_UPSTREAM_URL
	upstreamURL = ""

	//go:embed plugins.json
	database []byte
)

var logger *logrus.Logger

func main() {
	err := listenAndServe()
	if err != nil {
		panic("failed to listen and serve: " + err.Error())
	}
}

func newStaticStore(logger logrus.FieldLogger) (*store.StaticStore, error) {
	staticStore, err := store.NewStaticFromReader(bytes.NewReader(database), logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize store")
	}

	return staticStore, nil
}

func listenAndServe() error {
	logger = logrus.New()

	var apiStore store.Store
	var err error
	apiStore, err = newStaticStore(logger)
	if err != nil {
		return err
	}

	if upstreamURL != "" {
		upstreamStore, err := store.NewProxy(upstreamURL, logger)
		if err != nil {
			return errors.Wrap(err, "failed to initialize upstream store")
		}

		apiStore = store.NewMerged(logger, apiStore, upstreamStore)
	}

	router := mux.NewRouter()
	api.Register(router, &api.Context{
		Store:  apiStore,
		Logger: logger,
	})

	algnhsa.ListenAndServe(router, &algnhsa.Options{
		UseProxyPath: true,
	})

	return nil
}
