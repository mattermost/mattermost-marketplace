package store

import (
	"strings"

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
	var plugins []*model.Plugin

	filter := strings.TrimSpace(pluginFilter.Filter)
	if filter == "" {
		plugins = store.plugins
	} else {
		for _, plugin := range store.plugins {
			if !pluginMatchesFilter(plugin, filter) {
				continue
			}

			plugins = append(plugins, plugin)
		}
	}

	if len(plugins) == 0 || pluginFilter.PerPage == 0 {
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
