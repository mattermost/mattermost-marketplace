package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// addArchSpecificBundles includes the arch-specific bundle URLs and signatures in the Marketplace entries.
func addArchSpecificBundles(plugin *model.Plugin) (*model.Plugin, error) {
	if plugin.HomepageURL == "" {
		return plugin, nil
	}

	slashIndex := strings.LastIndex(plugin.HomepageURL, "/")
	repo := plugin.HomepageURL[slashIndex+1:]
	pluginWithVersion := fmt.Sprintf("%s-v%s", repo, plugin.Manifest.Version)

	remotePluginHost := os.Getenv("REMOTE_PLUGIN_HOST")
	if remotePluginHost == "" {
		remotePluginHost = defaultRemotePluginHost
	}

	archs, err := checkIfRemoteBundlesExist(remotePluginHost, pluginWithVersion)
	if err != nil {
		return nil, err
	}

	plugin.ArchBundles = model.ArchBundles{}
	for _, arch := range archs {
		fname := fmt.Sprintf("%s-%s.tar.gz", pluginWithVersion, arch)

		pluginPath := fmt.Sprintf("%s/%s", remotePluginHost, fname)
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

		meta := &model.ArchBundleMetadata{
			DownloadURL: pluginPath,
			Signature:   signatureStr,
		}

		switch arch {
		case model.LinuxAmd64:
			plugin.ArchBundles.LinuxAmd64 = meta
		case model.DarwinAmd64:
			plugin.ArchBundles.DarwinAmd64 = meta
		case model.WindowsAmd64:
			plugin.ArchBundles.WindowsAmd64 = meta
		}
	}

	return plugin, nil
}

// checkIfRemoteBundlesExist checks which arch-specific bundles are available on the remote file server, as well as their signatures.
func checkIfRemoteBundlesExist(remotePluginHost, pluginWithVersion string) ([]string, error) {
	result := []string{}

	archs := []string{model.LinuxAmd64, model.DarwinAmd64, model.WindowsAmd64}
	for _, arch := range archs {
		path := fmt.Sprintf("%s/%s-%s.tar.gz", remotePluginHost, pluginWithVersion, arch)

		// Check if plugin bundle exists on remote file server
		res, err := http.Head(path)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			continue
		}

		// Check if signature exists on remote file server
		res, err = http.Head(path + ".sig")
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			continue
		}

		result = append(result, arch)
	}

	return result, nil
}
