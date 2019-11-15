package store

import (
	"io"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// Store provides access to a store backed by the given reader.
type Store struct {
	plugins []*model.Plugin
	logger  logrus.FieldLogger
}

// New constructs a new instance of Store.
func New(reader io.Reader, logger logrus.FieldLogger) (*Store, error) {
	plugins, err := model.PluginsFromReader(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse stream")
	}

	if err := validatePlugins(plugins); err != nil {
		return nil, errors.Wrap(err, "failed to validate plugins")
	}

	return &Store{
		plugins,
		logger,
	}, nil
}

func validatePlugins(plugins []*model.Plugin) error {
	for _, plugin := range plugins {
		if plugin.Manifest.Id == "" {
			return errors.Errorf("plugin manifest Id is empty %+v", plugin)
		}
		if _, err := semver.Parse(plugin.Manifest.Version); err != nil {
			return errors.Wrapf(err, "failed to parse manifest version for manifest.Id %s", plugin.Manifest.Id)
		}
	}
	return nil
}
