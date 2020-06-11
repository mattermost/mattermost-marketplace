package store

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/mattermost/mattermost-marketplace/internal/api"
	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// Proxy is a store that fetches its result from some remote marketplace server.
type Proxy struct {
	marketplaceURL string
	logger         logrus.FieldLogger
}

// NewProxy creates a new instance of a proxy store.
func NewProxy(marketplaceURL string, logger logrus.FieldLogger) (*Proxy, error) {
	return &Proxy{
		marketplaceURL: marketplaceURL,
		logger:         logger.WithField("marketplace_url", marketplaceURL),
	}, nil
}

// GetPlugins fetches the given page of plugins. The first page is 0.
func (store *Proxy) GetPlugins(pluginFilter *model.PluginFilter) ([]*model.Plugin, error) {
	client := api.NewClient(store.marketplaceURL)

	plugins, err := client.GetPlugins(&api.GetPluginsRequest{
		Page:          pluginFilter.Page,
		PerPage:       pluginFilter.PerPage,
		Filter:        pluginFilter.Filter,
		ServerVersion: pluginFilter.ServerVersion,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to reach upstream store")
	}

	return plugins, nil
}
