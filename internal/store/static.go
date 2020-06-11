package store

import (
	"io"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// StaticStore provides access to a store backed by a static set of plugins.
type StaticStore struct {
	plugins []*model.Plugin
	logger  logrus.FieldLogger
}

// NewStatic constructs a new instance of a static store, parsing the plugins from the given reader.
func NewStaticFromReader(reader io.Reader, logger logrus.FieldLogger) (*StaticStore, error) {
	plugins, err := model.PluginsFromReader(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse stream")
	}

	return NewStatic(plugins, logger)
}

// NewStatic constructs a new instance of a static store using the given plugins.
func NewStatic(plugins []*model.Plugin, logger logrus.FieldLogger) (*StaticStore, error) {
	if err := validatePlugins(plugins); err != nil {
		return nil, errors.Wrap(err, "failed to validate plugins")
	}

	return &StaticStore{
		plugins,
		logger,
	}, nil
}

func validatePlugins(plugins []*model.Plugin) error {
	for _, plugin := range plugins {
		err := plugin.Manifest.IsValid()
		if err != nil {
			return errors.Wrapf(err, "invalid manifest for plugin %s", plugin.Manifest.Id)
		}

		if plugin.Manifest.Version == "" {
			return errors.Errorf("missing version in manifest for plugin%s", plugin.Manifest.Id)
		}
	}

	return nil
}

func pluginMatchesFilter(plugin *model.Plugin, filter string) bool {
	filter = strings.ToLower(filter)
	if strings.ToLower(plugin.Manifest.Id) == filter {
		return true
	}

	if strings.Contains(strings.ToLower(plugin.Manifest.Name), filter) {
		return true
	}

	if strings.Contains(strings.ToLower(plugin.Manifest.Description), filter) {
		return true
	}

	return false
}

// GetPlugins fetches the given page of plugins. The first page is 0.
func (store *StaticStore) GetPlugins(pluginFilter *model.PluginFilter) ([]*model.Plugin, error) {
	if pluginFilter.PerPage == 0 {
		return nil, nil
	}

	plugins, err := store.getPlugins(pluginFilter.ServerVersion)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get plugins")
	}

	filter := strings.TrimSpace(pluginFilter.Filter)
	if filter != "" {
		var filteredPlugins []*model.Plugin
		for _, plugin := range plugins {
			if pluginMatchesFilter(plugin, filter) {
				filteredPlugins = append(filteredPlugins, plugin)
			}
		}
		plugins = filteredPlugins
	}

	if len(plugins) == 0 {
		return nil, nil
	}
	if pluginFilter.PerPage == model.AllPerPage {
		return plugins, nil
	}

	start := (pluginFilter.Page) * pluginFilter.PerPage
	end := (pluginFilter.Page + 1) * pluginFilter.PerPage
	if start >= len(plugins) {
		return nil, nil
	}
	if end > len(plugins) {
		end = len(plugins)
	}

	return plugins[start:end], nil
}

// getPlugins returns all plugins compatible with the given server version, sorted by name ascending.
func (store *StaticStore) getPlugins(serverVersion string) ([]*model.Plugin, error) {
	var result []*model.Plugin
	plugins := map[string]*model.Plugin{}

	for _, storePlugin := range store.plugins {
		if serverVersion != "" && storePlugin.Manifest.MinServerVersion != "" {
			meetsMinServerVersion, err := storePlugin.Manifest.MeetMinServerVersion(serverVersion)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to check minServerVersion for manifest.Id %s", storePlugin.Manifest.Id)
			}

			if !meetsMinServerVersion {
				continue
			}
		}

		if plugins[storePlugin.Manifest.Id] == nil {
			plugins[storePlugin.Manifest.Id] = storePlugin
			continue
		}

		lastSeenPluginVersion, err := semver.Parse(plugins[storePlugin.Manifest.Id].Manifest.Version)
		if err != nil {
			return nil, errors.Errorf("failed to parse manifest.Version for manifest.Id %s", storePlugin.Manifest.Id)
		}

		storePluginVersion := semver.MustParse(storePlugin.Manifest.Version)

		// Replace the existing plugin if this version is newer, or if it's the same but
		// appears later in the list.
		if storePluginVersion.GTE(lastSeenPluginVersion) {
			plugins[storePlugin.Manifest.Id] = storePlugin
		}
	}

	for _, plugin := range plugins {
		result = append(result, plugin)
	}

	// Sort the final slice by plugin name, ascending
	sort.SliceStable(
		result,
		func(i, j int) bool {
			return strings.ToLower(result[i].Manifest.Name) < strings.ToLower(result[j].Manifest.Name)
		},
	)

	return result, nil
}
