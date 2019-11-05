package store

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

var ErrNotFound = errors.New("Plugin not found.")

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
		n := 0
		for _, plugin := range plugins {
			if pluginMatchesFilter(plugin, filter) {
				plugins[n] = plugin
				n++
			}
		}
		plugins = plugins[:n]
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

func (store *Store) getPlugins(serverVersion string) ([]*model.Plugin, error) {
	var result []*model.Plugin

	for _, plugin := range store.plugins {
		p, err := getPlugin(serverVersion, plugin)
		if err != nil && err != ErrNotFound {
			return nil, errors.Wrapf(err, "failed to get plugin")
		} else if err != ErrNotFound {
			result = append(result, p)
		}
	}

	return result, nil
}

// getPlugin gets the first plugin from the sorted pluginVersions slice that satisfies serverVersion.
func getPlugin(serverVersion string, pluginVersions pluginVersions) (*model.Plugin, error) {
	if len(pluginVersions) == 0 {
		return nil, errors.New("plugins should not be empty.")
	}

	// Get the latest plugin if no server version is provided
	if serverVersion == "" {
		return pluginVersions[len(pluginVersions)-1], nil
	}

	for i := len(pluginVersions) - 1; i >= 0; i-- {
		// Missing MinServerVersion means it's compatible with all servers
		if pluginVersions[i].Manifest.MinServerVersion == "" {
			return pluginVersions[i], nil
		}

		meetMinServerVersion, err := pluginVersions[i].Manifest.MeetMinServerVersion(serverVersion)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to MeetMinServerVersion version for manifest.Id) %s", pluginVersions[i].Manifest.Id)
		}

		if meetMinServerVersion {
			return pluginVersions[i], nil
		}
	}

	return nil, ErrNotFound
}
