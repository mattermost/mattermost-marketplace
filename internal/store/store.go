package store

import (
	"io"
	"sort"
	"strings"

	version "github.com/hashicorp/go-version"

	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Store provides access to a store backed by the given reader.
type Store struct {
	plugins []pluginVersions
	logger  logrus.FieldLogger
}

// pluginVersions is a sorted list of all versions of a given plugin.
type pluginVersions []*model.Plugin

// New constructs a new instance of Store.
func New(reader io.Reader, logger logrus.FieldLogger) (*Store, error) {
	plugins, err := model.PluginsFromReader(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse stream")
	}

	pluginVersions, err := sortPlugins(plugins)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sort plugins")
	}

	return &Store{
		pluginVersions,
		logger,
	}, nil
}

func sortPlugins(plugins []*model.Plugin) ([]pluginVersions, error) {
	var result []pluginVersions
	versions := map[string][]*model.Plugin{}

	// Combine multiple versions of a plugin into one list
	for _, plugin := range plugins {
		if plugin.Manifest.Id == "" {
			return nil, errors.Errorf("Plugin manifest Id is empty %+v", plugin)
		}
		if _, err := version.NewVersion(plugin.Manifest.Version); err != nil {
			return nil, errors.Wrapf(err, "failed to parse manifest version for manifest.Id %s", plugin.Manifest.Id)
		}

		versions[plugin.Manifest.Id] = append(versions[plugin.Manifest.Id], plugin)
	}

	for _, v := range versions {
		sort.SliceStable(
			v,
			func(i, j int) bool {
				// Ignoring errors, already checked previously
				ver1, _ := version.NewVersion(v[i].Manifest.Version)
				ver2, _ := version.NewVersion(v[j].Manifest.Version)

				return ver1.LessThan(ver2)
			},
		)
		result = append(result, v)
	}

	// Sort alphabetically
	sort.SliceStable(
		result,
		func(i, j int) bool {
			return strings.ToLower(result[i][0].Manifest.Name) < strings.ToLower(result[j][0].Manifest.Name)
		},
	)

	return result, nil
}
