package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// OSX-specific bundle URLs are stored in the plugin store as `osx` rather than `darwin`
const (
	OsxAmd64 = "osx-amd64"
	OsxArm64 = "osx-arm64"
)

func init() {
	generatorCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Migrate existing plugins in plugins.json to the newest structure.",
	Long:    "The migrate command adds platform-specific bundles to each existing entry.",
	Example: "generator migrate",
	RunE: func(command *cobra.Command, args []string) error {
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

		signatureBytes, err := ioutil.ReadAll(res.Body)
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
		case model.LinuxArm64:
			plugin.Platforms.LinuxArm64 = bundle
		case OsxAmd64:
			plugin.Platforms.DarwinAmd64 = bundle
		case OsxArm64:
			plugin.Platforms.DarwinArm64 = bundle
		case model.WindowsAmd64:
			plugin.Platforms.WindowsAmd64 = bundle
		}
	}

	return plugin, nil
}

// checkIfRemoteBundlesExist checks which platform-specific bundles are available on the remote file server, as well as their signatures.
func checkIfRemoteBundlesExist(remotePluginHost, pluginWithVersion string) ([]string, error) {
	result := []string{}

	platforms := []string{model.LinuxAmd64, model.LinuxArm64, OsxAmd64, OsxArm64, model.WindowsAmd64}
	for _, platform := range platforms {
		path := fmt.Sprintf("%s/%s-%s.tar.gz", remotePluginHost, pluginWithVersion, platform)

		// Check if plugin bundle exists on remote file server
		res, err := http.Head(path)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			logger.Debugf("Platform-specific bundle not found %s %s", pluginWithVersion, path)
			continue
		}

		// Check if signature exists on remote file server
		sigPath := path + ".sig"
		res, err = http.Head(sigPath)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			logger.Debugf("Platform-specific bundle signature not found %s %s", pluginWithVersion, sigPath)
			continue
		}

		result = append(result, platform)
	}

	return result, nil
}
