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

var minVersionSupportingEnterpriseFlags = semver.MustParse("5.25.0")

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

	plugins, err := store.getPlugins(pluginFilter.ServerVersion, pluginFilter.EnterprisePlugins, pluginFilter.Cloud, pluginFilter.Platform)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get plugins")
	}

	if !pluginFilter.ReturnAllVersions {
		plugins, err = filterToLatestVersion(plugins)
		if err != nil {
			return nil, errors.Wrap(err, "failed to filter to latest version")
		}
	}

	filter := strings.TrimSpace(pluginFilter.Filter)
	var filteredPlugins []*model.Plugin
	for _, plugin := range plugins {
		if pluginFilter.PluginId != "" && pluginFilter.PluginId != plugin.Manifest.Id {
			continue
		}
		if filter == "" || pluginMatchesFilter(plugin, filter) {
			filteredPlugins = append(filteredPlugins, plugin)
		}
	}
	plugins = filteredPlugins

	// Sort the final slice by plugin version decending then plugin name ascending.
	sort.SliceStable(
		plugins,
		func(i, j int) bool {
			if plugins[i].Manifest.Id == plugins[j].Manifest.Id {
				iVersion := semver.MustParse(plugins[i].Manifest.Version)
				jVersion := semver.MustParse(plugins[j].Manifest.Version)
				return iVersion.GT(jVersion)
			}
			return strings.ToLower(plugins[i].Manifest.Name) < strings.ToLower(plugins[j].Manifest.Name)
		},
	)

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

func filterToLatestVersion(plugins []*model.Plugin) ([]*model.Plugin, error) {
	latestVersionCollector := make(map[string]*model.Plugin)
	for _, plugin := range plugins {
		if latestVersionCollector[plugin.Manifest.Id] == nil {
			latestVersionCollector[plugin.Manifest.Id] = plugin
			continue
		}

		lastSeenPluginVersion, err := semver.Parse(latestVersionCollector[plugin.Manifest.Id].Manifest.Version)
		if err != nil {
			return nil, errors.Errorf("failed to parse manifest.Version for manifest.Id %s", plugin.Manifest.Id)
		}

		storePluginVersion := semver.MustParse(plugin.Manifest.Version)

		// Replace the existing plugin if this version is newer, or if it's the same but
		// appears later in the list.
		if storePluginVersion.GTE(lastSeenPluginVersion) {
			latestVersionCollector[plugin.Manifest.Id] = plugin
		}
	}

	result := make([]*model.Plugin, 0, len(plugins))
	for _, plugin := range latestVersionCollector {
		result = append(result, plugin)
	}

	return result, nil
}

// getPlugins returns all plugins compatible with the given server version, sorted by name ascending.
func (store *StaticStore) getPlugins(serverVersion string, includeEnterprisePlugins bool, isCloud bool, platform string) ([]*model.Plugin, error) {
	var result []*model.Plugin

	for _, storePlugin := range store.plugins {
		if storePlugin.Enterprise && !includeEnterprisePlugins {
			if serverVersion == "" {
				continue
			}

			sv, err := semver.Parse(serverVersion)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse serverVersion %s", serverVersion)
			}

			// Honor enterprise flag for server version >= 5.25.0 only.
			// Workaround for https://mattermost.atlassian.net/browse/MM-26507

			if sv.GE(minVersionSupportingEnterpriseFlags) {
				continue
			}
		}

		if isCloud && storePlugin.Hosting == model.OnPrem {
			continue
		}

		if !isCloud && storePlugin.Hosting == model.Cloud {
			continue
		}

		if serverVersion != "" && storePlugin.Manifest.MinServerVersion != "" {
			_, err := semver.Parse(serverVersion)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse serverVersion %s", serverVersion)
			}

			meetsMinServerVersion, err := storePlugin.Manifest.MeetMinServerVersion(serverVersion)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to check minServerVersion for manifest.Id %s", storePlugin.Manifest.Id)
			}

			if !meetsMinServerVersion {
				continue
			}
		}

		// Create a copy as we want to modify only the returned one
		newRef := *storePlugin
		storePlugin = &newRef

		storePlugin.AddLabels()

		if platform != "" {
			var bundle model.PlatformBundleMetadata
			switch platform {
			case model.LinuxAmd64:
				bundle = storePlugin.Platforms.LinuxAmd64
			case model.DarwinAmd64:
				bundle = storePlugin.Platforms.DarwinAmd64
			case model.WindowsAmd64:
				bundle = storePlugin.Platforms.WindowsAmd64
			}

			if bundle.DownloadURL != "" && bundle.Signature != "" {
				storePlugin.DownloadURL = bundle.DownloadURL
				storePlugin.Signature = bundle.Signature
			}
		}

		result = append(result, storePlugin)
	}

	return result, nil
}
