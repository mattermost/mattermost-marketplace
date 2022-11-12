package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-github/v28/github"
	artifacts "github.com/mattermost/mattermost-marketplace/mocks/artifacts"
	mocks "github.com/mattermost/mattermost-marketplace/mocks/github/v3"
)

func TestGenerator(t *testing.T) {
	client := mocks.MockGitHubClient(t)
	ctx := context.Background()

	const orgName string = "mattermost"

	t.Run("Check the \"mocked\" GitHub client", func(t *testing.T) {
		t.Run("Assurance that the client is offline", func(t *testing.T) {
			onlineAPIServer := strings.ToLower("api.github.com")
			localhostDigitsServer := strings.ToLower("127.0.0.1")
			localhostNameServer := strings.ToLower("localhost")
			actual := strings.Split(strings.ToLower(client.BaseURL.Host), ":")[0]

			if strings.Compare(onlineAPIServer, actual) == 0 {
				t.Error("Client is still expecting to connect to main GitHub API")
			}
			if (len(actual) > 0) &&
				(strings.Compare(localhostDigitsServer, actual) != 0) &&
				(strings.Compare(localhostNameServer, actual) != 0) {
				// client does appear to be separate from the main GitHub API,
				// so add that to the log in case something breaks
				t.Log("GitHub Client interacting with API server at: " + client.BaseURL.String())
			}
		})
		t.Run("Get a Repository without erroring", func(t *testing.T) {
			_, _, err := client.Repositories.Get(ctx, "mattermost", "mattermost-plugin-github")
			if err != nil {
				t.Error("client.Repositories.Get() errored: " + err.Error())
			}
		})
		t.Run("List a Repository's Releases without erroring", func(t *testing.T) {
			_, _, err := client.Repositories.ListReleases(ctx, "mattermost", "mattermost-plugin-github", nil)
			if err != nil {
				t.Error("client.Repositories.ListReleases() errored: " + err.Error())
			}
		})
	})

	t.Run("getReleases() returns a list of Releases for a Repository", func(t *testing.T) {
		interestingRepo := "mattermost-plugin-jira"
		var releaseList []*github.RepositoryRelease
		var err error = nil
		t.Run("Requests a list of Releases", func(t *testing.T) {
			releaseList, err = getReleases(ctx, client, orgName, interestingRepo, true)
			if err != nil {
				t.Errorf("Unable to get repo %s/%s: %s", orgName, interestingRepo, err.Error())
			}
			if len(releaseList) == 0 {
				t.Errorf("Requested for repo %s/%s and got an empty list", orgName, interestingRepo)
			}
			var countPreRelease uint = 0
			for _, rel := range releaseList {
				if rel.GetPrerelease() {
					countPreRelease++
				}
			}
			if countPreRelease > 0 {
				t.Run("Optionally excludes pre-release Releases", func(t *testing.T) {
					excList, err := getReleases(ctx, client, orgName, interestingRepo, false)
					lenReleaseList := len(releaseList)
					if err == nil {
						lenExcList := len(excList)
						if lenExcList == lenReleaseList {
							t.Errorf("Release count the same with and \"without\" pre-release (%d)", lenReleaseList)
						} else if lenReleaseList != lenExcList+int(countPreRelease) {
							t.Errorf("Release list count without prereleases expected %d - %d = %d, got %d",
								lenReleaseList, countPreRelease,
								lenReleaseList-int(countPreRelease),
								lenExcList)
						}
					}
				})
			}
		})
	})

	t.Run("getReleasePlugin() functionality tests", func(t *testing.T) {
		t.Skip("Testing smaller work units first")
		/*
		 * inputs: github.RepositoryRelease, github.Repository, array of model.Plugin, name of plugin host
		 * 1 - process the release to get the release version tag name
		 * 2 - loop through the release's assets to identify the URL, signature, and most recent update time
		 * 3 - request the signature from the online version control host
		 * 4 - loop through the array of model.Plugin parameter to find one with the same download URL
		 * 5 - compare the matched model.Plugin and release asset
		 * 6 - if an updated plugin asset is available, download it, read and populate the manifest data for plugin, and extract icon
		 * 7 - always update the plugin URLs, signature, and updated timestamp
		 * 8 - populate the platform specific bundle information for the plugin
		 * 9 - return the model.Plugin as updated
		 */
	})

	t.Run("getFromTarFile() seeks and extracts a particular file from a tar file", func(t *testing.T) {
		gzreader, err := gzip.NewReader(bytes.NewReader(artifacts.MockGitHubPluginBundle))
		if err != nil {
			t.Skip("unable to load mock gzip file data")
		}
		tarreader := tar.NewReader(gzreader)
		manifest, err := getFromTarFile(tarreader, "plugin.json")
		if err != nil {
			t.Errorf("unsuccessful attempt to read plugin.json manifest from mock bundle")
		} else if len(manifest) == 0 {
			t.Errorf("0-length bundle manifest not expected from mock bundle")
		}
		icon, err := getFromTarFile(tarreader, "assets/icon.svg")
		if err != nil {
			t.Errorf("unsuccessful attempt to read icon.svg from mock bundle")
		} else if len(icon) == 0 {
			t.Errorf("0-length bundle icon not expected from mock bundle")
		} else if bytes.Equal(manifest, icon) {
			t.Errorf("getFromTarFile() returning same data for different files")
		}
	})

	t.Run("downloadSignature() downloads a signature file for an artifact", func(t *testing.T) {
		signature, err := downloadSignature(client.BaseURL.String() + "mattermost/mattermost-plugin-github/releases/download/v2.0.0/github-2.0.0.tar.gz.sig")
		if err != nil {
			t.Errorf("Unable to get signature file: %s", err.Error())
		} else if len(signature) == 0 {
			t.Errorf("0-length signature not expected for test request")
		}
	})

	t.Run("downloadBundleData() downloads and unpacks a bundle file", func(t *testing.T) {
		bundledata, err := downloadBundleData(client.BaseURL.String() + "mattermost/mattermost-plugin-github/releases/download/v0.0.0/github-0.0.0.tar.gz")
		if err != nil {
			t.Errorf("Unable to get bundle file: %s", err.Error())
		} else if len(bundledata) == 0 {
			t.Errorf("0-length bundle data not expected for test request")
		} else if bytes.Equal(bundledata[:2], []byte{0x1f, 0x8b}) {
			t.Errorf("compressed tarball not unpacked")
		}
	})

	t.Run("getIconDataFromTarFile() reads an icon from a bundle file", func(t *testing.T) {
		gzreader, err := gzip.NewReader(bytes.NewReader(artifacts.MockGitHubPluginBundle))
		if err != nil {
			t.Skip("unable to load mock bundle file data")
		}
		tardata, err := ioutil.ReadAll(gzreader)
		if err != nil {
			t.Skip("unable to load mock bundle file data")
		}
		icon, err := getIconDataFromTarFile(tardata, "assets/icon.svg")
		if err != nil {
			t.Errorf("unsuccessful attempt to read icon.svg from mock bundle")
		} else if len(icon) == 0 {
			t.Errorf("0-length bundle icon not expected from mock bundle")
		}
		if icon[:26] == "data:image/svg+xml;base64," {
			dec, _ := base64.RawStdEncoding.DecodeString(icon[26:])
			decstr := string(dec)
			if decstr[:15] != "<svg role=\"img\"" {
				t.Errorf("icon was not successfully unpacked")
			}
		} else {
			t.Log("icon data memory format has changed")
		}
	})

	t.Run("InitCommand() initializes the command line invocation using its arguments", func(t *testing.T) {
		t.Skip("Test not yet implemented")
	})

	t.Run("pluginsFromDatabase() reads plugin information from existing database", func(t *testing.T) {
		t.Skip("Test not yet implemented")
	})

	t.Run("pluginsToDatabase() writes plugin information to existing database", func(t *testing.T) {
		t.Skip("Test not yet implemented")
	})
}
