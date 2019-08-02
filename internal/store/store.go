package store

import (
	"io"

	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	return &Store{
		plugins,
		logger,
	}, nil
}
