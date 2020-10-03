package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

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

		bundle := &model.PlatformBundleMetadata{
			DownloadURL: pluginPath,
			Signature:   signatureStr,
		}

		switch platform {
		case model.LinuxAmd64:
			plugin.Platforms.LinuxAmd64 = bundle
		case model.OsxAmd64:
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

	platforms := []string{model.LinuxAmd64, model.OsxAmd64, model.WindowsAmd64}
	for _, platform := range platforms {
		path := fmt.Sprintf("%s/%s-%s.tar.gz", remotePluginHost, pluginWithVersion, platform)

		// Check if plugin bundle exists on remote file server
		res, err := http.Head(path)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			logger.Infof("Platform-specific bundle not found %s", path)
			continue
		}

		// Check if signature exists on remote file server
		sigPath := path + ".sig"
		res, err = http.Head(sigPath)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			logger.Infof("Platform-specific bundle signature not found %s", sigPath)
			continue
		}

		result = append(result, platform)
	}

	return result, nil
}
