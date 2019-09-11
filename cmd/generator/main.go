package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v28/github"
	"github.com/h2non/filetype"
	svg "github.com/h2non/go-is-svg"
	mattermostModel "github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

func init() {
	generatorCmd.PersistentFlags().String("github-token", "", "The optional GitHub token for API requests.")
	generatorCmd.PersistentFlags().Bool("debug", false, "Whether to output debug logs.")
	generatorCmd.PersistentFlags().Bool("include-pre-release", true, "Whether to include pre-release versions.")
}

func main() {
	if err := generatorCmd.Execute(); err != nil {
		logger.WithError(err).Error("command failed")
		os.Exit(1)
	}
}

var generatorCmd = &cobra.Command{
	Use:   "generator",
	Short: "Generator is a tool to generate the plugins.json database",
	// SilenceErrors allows us to explicitly log the error returned from generatorCmd below.
	SilenceErrors: true,
	RunE: func(command *cobra.Command, args []string) error {
		command.SilenceUsage = true

		debug, _ := command.Flags().GetBool("debug")
		if debug {
			logger.SetLevel(logrus.DebugLevel)
		}

		includePreRelease, _ := command.Flags().GetBool("include-pre-release")

		githubToken, _ := command.Flags().GetString("github-token")

		var client *github.Client

		if githubToken != "" {
			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: githubToken},
			)
			tc := oauth2.NewClient(ctx, ts)

			client = github.NewClient(tc)
		} else {
			client = github.NewClient(nil)
		}

		ctx := context.Background()

		repositoryNames := []string{
			"mattermost-plugin-demo",
			"mattermost-plugin-github",
			"mattermost-plugin-autolink",
			"mattermost-plugin-zoom",
			"mattermost-plugin-jira",
			"mattermost-plugin-autotranslate",
			"mattermost-plugin-profanity-filter",
			"mattermost-plugin-welcomebot",
			"mattermost-plugin-jenkins",
			"mattermost-plugin-antivirus",
			"mattermost-plugin-walltime",
			"mattermost-plugin-custom-attributes",
			"mattermost-plugin-skype4business",
			"mattermost-plugin-aws-SNS",
			"mattermost-plugin-gitlab",
			"mattermost-plugin-nps",
		}

		type iconURL struct {
			URL string
		}

		iconURLs := map[string]iconURL{
			"mattermost-plugin-aws-SNS":        iconURL{"https://unpkg.com/simple-icons@latest/icons/amazonaws.svg"},
			"mattermost-plugin-github":         iconURL{"https://unpkg.com/simple-icons@latest/icons/github.svg"},
			"mattermost-plugin-gitlab":         iconURL{"https://unpkg.com/simple-icons@latest/icons/gitlab.svg"},
			"mattermost-plugin-jenkins":        iconURL{"https://unpkg.com/simple-icons@latest/icons/jenkins.svg"},
			"mattermost-plugin-jira":           iconURL{"https://unpkg.com/simple-icons@latest/icons/jira.svg"},
			"mattermost-plugin-skype4business": iconURL{"https://unpkg.com/simple-icons@latest/icons/skype.svg"},
		}

		plugins := []*model.Plugin{}

		for _, repositoryName := range repositoryNames {
			logger.Debugf("querying repository %s", repositoryName)

			plugin, err := getReleasePlugin(ctx, client, repositoryName, includePreRelease)
			if err != nil {
				return errors.Wrapf(err, "failed to release plugin for repository %s", repositoryName)
			}

			if iconURL, ok := iconURLs[repositoryName]; ok {
				icon, err := getIcon(ctx, iconURL.URL)
				if err != nil {
					return errors.Wrapf(err, "failed to fetch icon for repository %s", repositoryName)
				}
				if svg.Is(icon) {
					plugin.IconData = fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(icon))
				} else {
					kind, err := filetype.Image(icon)
					if err != nil {
						return errors.Wrapf(err, "failed to match icon at %s to image", iconURL.URL)
					}

					plugin.IconData = fmt.Sprintf("data:%s;base64,%s", kind.MIME, base64.StdEncoding.EncodeToString(icon))
				}
			}

			plugins = append(plugins, plugin)
		}

		encoder := json.NewEncoder(os.Stdout)
		err := encoder.Encode(plugins)
		if err != nil {
			return errors.Wrap(err, "failed to encode plugins result")
		}

		return nil
	},
}

func getReleasePlugin(ctx context.Context, client *github.Client, repositoryName string, includePreRelease bool) (*model.Plugin, error) {
	logger := logger.WithField("repository", repositoryName)

	repository, _, err := client.Repositories.Get(ctx, "mattermost", repositoryName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get repository")
	}

	latestRelease, err := getLatestRelease(ctx, client, repositoryName, includePreRelease)
	if err != nil {
		return nil, err
	}
	if latestRelease == nil {
		logger.Warnf("no latest release found for repository")
		return nil, nil
	}

	var releaseName string
	if latestRelease.GetName() == "" {
		releaseName = latestRelease.GetTagName()
	} else {
		releaseName = fmt.Sprintf("%s (%s)", latestRelease.GetName(), latestRelease.GetTagName())
	}
	logger.Debugf("found latest release %s", releaseName)

	downloadURL := ""
	for _, releaseAsset := range latestRelease.Assets {
		if strings.HasSuffix(releaseAsset.GetName(), ".tar.gz") {
			downloadURL = releaseAsset.GetBrowserDownloadURL()
		}
	}

	if downloadURL == "" {
		logger.Warnf("Failed to find plugin asset release %s", releaseName)
		return nil, nil
	}
	logger.Debugf("fetching download url %s", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download plugin bundle for release %s", releaseName)
	}
	defer resp.Body.Close()

	plugin := model.Plugin{
		HomepageURL:       repository.GetHTMLURL(),
		IconData:          "",
		DownloadURL:       downloadURL,
		DownloadSignature: []byte{},
	}

	gzBundleReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read gzipped plugin bundle for release %s", releaseName)
	}

	bundleReader := tar.NewReader(gzBundleReader)
	for {
		hdr, err := bundleReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read plugin bundle for release %s", releaseName)
		}

		if path.Base(hdr.Name) != "plugin.json" {
			continue
		}
		manifest := mattermostModel.ManifestFromJson(bundleReader)

		plugin.Manifest = manifest
		break
	}

	if plugin.Manifest == nil {
		return nil, fmt.Errorf("failed to find plugin manifest for release %s", releaseName)
	}

	return &plugin, nil
}

func getLatestRelease(ctx context.Context, client *github.Client, repoName string, includePreRelease bool) (*github.RepositoryRelease, error) {
	releases, _, err := client.Repositories.ListReleases(ctx, "mattermost", repoName, &github.ListOptions{
		Page:    0,
		PerPage: 10,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get releases for repository %s", repoName)
	}

	var latestRelease *github.RepositoryRelease
	for _, release := range releases {
		if release.GetDraft() {
			continue
		}

		if release.GetPrerelease() && !includePreRelease {
			continue
		}

		if latestRelease == nil || release.GetPublishedAt().After(latestRelease.GetPublishedAt().Time) {
			latestRelease = release
		}
	}

	return latestRelease, nil
}

func getIcon(ctx context.Context, url string) ([]byte, error) {
	logger.Debugf("fetching icon url %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download plugin icon at %s", url)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
