package api

import (
	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/sirupsen/logrus"
)

// Store describes the interface to the backing store.
type Store interface {
	GetPlugins(filter model.PluginFilter) ([]*model.Plugin, error)
}

// Context provides the API with all necessary data and interfaces for responding to requests.
//
// It is cloned before each request, allowing per-request changes such as logger annotations.
type Context struct {
	Store     Store
	RequestID string
	Logger    logrus.FieldLogger
}

// Clone creates a shallow copy of context, allowing clones to apply per-request changes.
func (c *Context) Clone() *Context {
	return &Context{
		Store:  c.Store,
		Logger: c.Logger,
	}
}
