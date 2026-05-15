package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// Older plugins publish darwin bundles as "osx-amd64"; newer ones use "darwin-amd64".
const (
	OsxAmd64    = "osx-amd64"
	DarwinAmd64 = "darwin-amd64"
)

func init() {
	generatorCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Migrate existing plugins in plugins.json to the newest structure.",
	Long:    "The migrate command adds platform-specific bundles to each existing entry.",
	Example: "generator migrate",
	RunE: func(command *cobra.Command, _ []string) error {
		dbFile, err := command.Flags().GetString("database")
		if err != nil {
			return err
		}

		pluginHost, err := command.Flags().GetString("remote-plugin-store")
		if err != nil {
			return err
		}

		existingPlugins, err := pluginsFromDatabase(dbFile)
		if err != nil {
			return errors.Wrap(err, "failed to read plugins from database")
		}

		var g errgroup.Group
		toSave := []*model.Plugin{}
		for _, orig := range existingPlugins {
			orig := orig

			g.Go(func() error {
				var modified *model.Plugin
				modified, err = addPlatformSpecificBundles(orig, pluginHost)
				if err != nil {
					return errors.Wrapf(err, "failed to add platform-specific bundles for plugin %s-%s", orig.Manifest.Id, orig.Manifest.Version)
				}

				// Migrate community label to flag
				var newLabels []model.Label
				for _, l := range modified.Labels {
					switch l {
					case model.EnterpriseLabel:
						// Just drop it
					case model.CommunityLabel:
						modified.AuthorType = model.Community
					case model.BetaLabel:
						modified.ReleaseStage = model.Beta
					default:
						// Keep other labels
						newLabels = append(newLabels, l)
					}
				}
				modified.Labels = newLabels

				if modified.AuthorType == "" {
					modified.AuthorType = model.Mattermost
				}

				if modified.ReleaseStage == "" {
					modified.ReleaseStage = model.Production
				}

				toSave = append(toSave, modified)

				return nil
			})
		}

		if err = g.Wait(); err != nil {
			return errors.Wrap(err, "failed to get a migrate a plugin to new structure")
		}

		err = pluginsToDatabase(dbFile, toSave)
		if err != nil {
			return errors.Wrap(err, "failed to write plugins database")
		}
		return nil
	}}

// addPlatformSpecificBundles includes the platform-specific bundle URLs and signatures in the Marketplace entries.
func addPlatformSpecificBundles(plugin *model.Plugin, pluginHost string) (*model.Plugin, error) {
	if plugin.RepoName == "" {
		return plugin, nil
	}

	repo := plugin.RepoName
	pluginWithVersion := fmt.Sprintf("%s-v%s", repo, plugin.Manifest.Version)

	platforms, err := checkIfRemoteBundlesExist(pluginHost, pluginWithVersion)
	if err != nil {
		return nil, err
	}

	plugin.Platforms = model.PlatformBundles{}
	for _, platform := range platforms {
		fname := fmt.Sprintf("%s-%s.tar.gz", pluginWithVersion, platform)

		pluginPath := fmt.Sprintf("%s/%s", pluginHost, fname)
		sigPath := pluginPath + ".sig"

		res, err := http.Get(sigPath)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received %d status code downloading signature from %s", res.StatusCode, sigPath)
		}

		signatureBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		signatureStr := base64.StdEncoding.EncodeToString(signatureBytes)

		bundle := model.PlatformBundleMetadata{
			DownloadURL: pluginPath,
			Signature:   signatureStr,
		}

		switch platform {
		case model.LinuxAmd64:
			plugin.Platforms.LinuxAmd64 = bundle
		case OsxAmd64, DarwinAmd64:
			plugin.Platforms.DarwinAmd64 = bundle
		case model.WindowsAmd64:
			plugin.Platforms.WindowsAmd64 = bundle
		}
	}

	return plugin, nil
}

// checkIfRemoteBundlesExist checks which platform-specific bundles are available on the remote file server, as well as their signatures.
func checkIfRemoteBundlesExist(remotePluginHost, pluginWithVersion string) ([]string, error) {
	result := []string{}

	for _, platform := range []string{model.LinuxAmd64, model.WindowsAmd64} {
		exists, err := remoteBundleExists(remotePluginHost, pluginWithVersion, platform)
		if err != nil {
			return nil, err
		}
		if exists {
			result = append(result, platform)
		}
	}

	if darwin := checkDarwinBundle(remotePluginHost, pluginWithVersion); darwin != "" {
		result = append(result, darwin)
	}

	return result, nil
}

// checkDarwinBundle returns the first darwin bundle naming convention with both
// a bundle and signature available, preferring "darwin-amd64" over "osx-amd64".
// Returns an empty string if neither variant exists.
func checkDarwinBundle(remotePluginHost, pluginWithVersion string) string {
	for _, platform := range []string{DarwinAmd64, OsxAmd64} {
		exists, err := remoteBundleExists(remotePluginHost, pluginWithVersion, platform)
		if err != nil {
			logger.Debugf("Error checking darwin bundle %s: %v", platform, err)
			continue
		}
		if exists {
			return platform
		}
	}

	return ""
}

// remoteBundleExists reports whether both a platform-specific bundle and its signature
// are available on the remote file server.
func remoteBundleExists(remotePluginHost, pluginWithVersion, platform string) (bool, error) {
	bundlePath := fmt.Sprintf("%s/%s-%s.tar.gz", remotePluginHost, pluginWithVersion, platform)

	ok, err := remoteResourceExists(bundlePath)
	if err != nil {
		return false, err
	}
	if !ok {
		logger.Debugf("Platform-specific bundle not found %s %s", pluginWithVersion, bundlePath)
		return false, nil
	}

	sigPath := bundlePath + ".sig"
	ok, err = remoteResourceExists(sigPath)
	if err != nil {
		return false, err
	}
	if !ok {
		logger.Debugf("Platform-specific bundle signature not found %s %s", pluginWithVersion, sigPath)
		return false, nil
	}

	return true, nil
}

// remoteResourceExists issues a HEAD request, returning true on a 200 response.
func remoteResourceExists(url string) (bool, error) {
	res, err := http.Head(url)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}
