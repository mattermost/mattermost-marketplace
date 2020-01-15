package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"os"
	"regexp"
	"time"

	mattermostModel "github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

func init() {
	generatorCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add [repo] [tag]",
	Short:   "Add a plugin release to the plugins.json database",
	Example: `  generator add matterpoll v1.5.1`,
	Args:    cobra.ExactArgs(2),
	RunE: func(command *cobra.Command, args []string) error {
		command.SilenceUsage = true

		plugins, err := InitCommand(command)
		if err != nil {
			return err
		}

		repo := args[0]
		tag := args[1]
		if !regexp.MustCompile("v[0-9]+.[0-9]+.[0-9]+$").MatchString(tag) {
			return errors.Errorf("bad tag. Maybe a typo? You set this tag %s but we expected something like v2.3.4", tag)
		}

		bundleURL := "https://plugins-store.test.mattermost.com/release/" + repo + "-" + tag + ".tar.gz"
		signatureURL := bundleURL + ".sig"

		bundleData, err := downloadBundleData(bundleURL)
		if err != nil {
			return errors.Wrapf(err, "failed download bundle data")
		}

		manifestData, err := getFromTarFile(tar.NewReader(bytes.NewReader(bundleData)), "plugin.json")
		if err != nil {
			return errors.Wrap(err, "failed to read manifest from plugin bundle for release")
		}

		manifest := mattermostModel.ManifestFromJson(bytes.NewReader(manifestData))
		if manifest == nil {
			return errors.New("manifest nil after reading from plugin bundle for release")
		}

		var iconData string
		if manifest.IconPath != "" {
			iconData, err = getIconDataFromTarFile(bundleData, manifest.IconPath)
			if err != nil {
				return errors.Wrap(err, "failed to set icon")
			}
		}

		signature, err := downloadSignature(signatureURL)
		if err != nil {
			return errors.Wrap(err, "failed to download plugin signature")
		}

		plugin := &model.Plugin{
			HomepageURL:     manifest.HomepageURL,
			IconData:        iconData,
			DownloadURL:     bundleURL,
			ReleaseNotesURL: "", // Not jet supported
			Labels:          nil,
			Signature:       signature,
			Manifest:        manifest,
			UpdatedAt:       time.Now().In(time.UTC),
		}
		plugin.Signature = ""

		plugins = append(plugins, plugin)
		err = json.NewEncoder(os.Stdout).Encode(plugins)
		if err != nil {
			return errors.Wrap(err, "failed to encode plugins result")
		}

		return nil
	},
}
