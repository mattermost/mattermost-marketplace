package store

import (
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

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
func (store *Store) GetPlugins(pluginFilter *model.PluginFilter) ([]*model.Plugin, error) {
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
func (store *Store) getPlugins(serverVersion string) ([]*model.Plugin, error) {
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
		if storePluginVersion.GT(lastSeenPluginVersion) {
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
