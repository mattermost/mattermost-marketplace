package main

import (
	"archive/tar"
	"bytes"
	"time"

	"github.com/blang/semver"
	mattermostModel "github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

func init() {
	generatorCmd.AddCommand(addCmd)

	addCmd.Flags().Bool("beta", false, "Mark release as Beta")
	addCmd.Flags().Bool("official", false, "Mark this plugin as maintained by Mattermost")
	addCmd.Flags().Bool("community", false, "Mark this plugin as maintained by the Open Source Community")
	addCmd.Flags().Bool("enterprise", false, "Mark this plugin as only available to installations with an E20-only plugins license")
	addCmd.Flags().Bool("cloud", false, "Mark this plugin as only available to cloud installations")
	addCmd.Flags().Bool("on-prem", false, "Mark this plugin as only available to on-prem installations")
}

var addCmd = &cobra.Command{
	Use:   "add [repo] [tag]",
	Short: "Add a plugin release to the plugins.json database",
	Long: "The generator commands allows adding a specific plugin release to the database by using this command.\n\n" +
		"The release has to be built first using the /mb cutplugin command, which also uploads it to https://plugins-store.test.mattermost.com/release/. " +
		"This location is used to fetch the plugin release.",
	Example: `  generator add matterpoll v1.5.1`,
	Args:    cobra.ExactArgs(2),
	RunE: func(command *cobra.Command, args []string) error {
		command.SilenceUsage = true

		official, err := command.Flags().GetBool("official")
		if err != nil {
			return err
		}

		community, err := command.Flags().GetBool("community")
		if err != nil {
			return err
		}

		enterprise, err := command.Flags().GetBool("enterprise")
		if err != nil {
			return err
		}

		if official == community {
			return errors.New("you must either set the release as a official or as a community plugin")
		}

		cloud, err := command.Flags().GetBool("cloud")
		if err != nil {
			return err
		}

		onPrem, err := command.Flags().GetBool("on-prem")
		if err != nil {
			return err
		}

		if cloud && onPrem {
			return errors.New("if you want to make a plugin availed for cloud and on-prem, just drop both flags")
		}

		beta, err := command.Flags().GetBool("beta")
		if err != nil {
			return err
		}

		if err = InitCommand(command); err != nil {
			return err
		}

		dbFile, err := command.Flags().GetString("database")
		if err != nil {
			return err
		}

		plugins, err := pluginsFromDatabase(dbFile)
		if err != nil {
			return errors.Wrap(err, "failed to read plugins from database")
		}

		repo := args[0]
		tag := args[1]

		if _, err = semver.ParseTolerant(tag); err != nil {
			return errors.Wrapf(err, "%v is an invalid tag. Something like v2.3.4 is expected", tag)
		}

		bundleURL := "https://plugins-store.test.mattermost.com/release/" + repo + "-" + tag + ".tar.gz"
		signatureURL := bundleURL + ".sig"

		bundleData, err := downloadBundleData(bundleURL)
		if err != nil {
			return errors.Wrapf(err, "failed downloading bundle data")
		}

		manifestData, err := getFromTarFile(tar.NewReader(bytes.NewReader(bundleData)), "plugin.json")
		if err != nil {
			return errors.Wrap(err, "failed to read manifest from plugin bundle for release")
		}

		manifest := mattermostModel.ManifestFromJson(bytes.NewReader(manifestData))
		if manifest == nil {
			return errors.New("manifest nil after reading from plugin bundle for release")
		}

		err = manifest.IsValid()
		if err != nil {
			return errors.Wrap(err, "manifest is invalid")
		}

		var iconData string
		if manifest.IconPath != "" {
			iconData, err = getIconDataFromTarFile(bundleData, manifest.IconPath)
			if err != nil {
				return errors.Wrap(err, "failed to get icon")
			}
		}

		signature, err := downloadSignature(signatureURL)
		if err != nil {
			return errors.Wrap(err, "failed to download plugin signature")
		}

		labels := []model.Label{}
		if beta {
			labels = append(labels, model.BetaLabel)
		}

		if community {
			labels = append(labels, model.CommunityLabel)
		}

		if enterprise {
			labels = append(labels, model.EnterpriseLabel)
		}

		plugin := &model.Plugin{
			HomepageURL:     manifest.HomepageURL,
			IconData:        iconData,
			DownloadURL:     bundleURL,
			ReleaseNotesURL: manifest.ReleaseNotesURL,
			Labels:          labels,
			Signature:       signature,
			Manifest:        manifest,
			Enterprise:      enterprise,
			UpdatedAt:       time.Now().In(time.UTC),
		}

		if cloud {
			plugin.Hosting = model.Cloud
		}

		if onPrem {
			plugin.Hosting = model.OnPrem
		}

		plugins = append(plugins, plugin)

		err = pluginsToDatabase(dbFile, plugins)
		if err != nil {
			return errors.Wrap(err, "failed to write plugins database")
		}

		return nil
	},
}
