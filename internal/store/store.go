package store

import "github.com/mattermost/mattermost-marketplace/internal/model"

// Store describes the interface to the backing store.
type Store interface {
	GetPlugins(filter *model.PluginFilter) ([]*model.Plugin, error)
}
