package store

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// Merged is a store that merges the results of multiple stores together.
//
// If a plugin is present in multiple stores, the later version is preferred. If a plugin with
// the same version is present in multiple stores, the one from the later store (as initialized)
// is preferred.
type Merged struct {
	stores []Store
	logger logrus.FieldLogger
}

// NewMerged creates a new instance of the merged store.
func NewMerged(logger logrus.FieldLogger, stores ...Store) *Merged {
	return &Merged{
		stores: stores,
		logger: logger,
	}
}

// GetPlugins fetches the given page of plugins. The first page is 0.
func (store *Merged) GetPlugins(pluginFilter *model.PluginFilter) ([]*model.Plugin, error) {
	// Short-circuit if only one store is configured.
	if len(store.stores) == 1 {
		return store.stores[0].GetPlugins(pluginFilter)
	}

	plugins := []*model.Plugin{}
	for i, store := range store.stores {
		storePlugins, err := store.GetPlugins(&model.PluginFilter{
			Page:              0,
			PerPage:           model.AllPerPage,
			Filter:            pluginFilter.Filter,
			ServerVersion:     pluginFilter.ServerVersion,
			EnterprisePlugins: pluginFilter.EnterprisePlugins,
			Platform:          pluginFilter.Platform,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to query store %d", i)
		}

		plugins = append(plugins, storePlugins...)
	}

	staticStore, err := NewStatic(plugins, store.logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize static store")
	}

	return staticStore.GetPlugins(pluginFilter)
}
