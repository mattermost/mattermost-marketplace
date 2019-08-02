package store

import (
	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// GetPlugins fetches the given page of plugins. The first page is 0.
func (store *Store) GetPlugins(filter *model.PluginFilter) ([]*model.Plugin, error) {
	if len(store.plugins) == 0 || filter.PerPage == 0 {
		return nil, nil
	}
	if filter.PerPage == model.AllPerPage {
		return store.plugins, nil
	}

	start := (filter.Page) * filter.PerPage
	end := (filter.Page + 1) * filter.PerPage
	if start >= len(store.plugins) {
		return nil, nil
	}
	if end > len(store.plugins) {
		end = len(store.plugins)
	}

	return store.plugins[start:end], nil
}
