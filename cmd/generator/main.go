package main

import (
	"archive/tar"
	"bytes"
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
	"sort"
	"strings"
	"time"

	"github.com/blang/semver"
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
	generatorCmd.PersistentFlags().Bool("debug", false, "Whether to output debug logs.")
	generatorCmd.PersistentFlags().String("existing", "plugins.json", "An existing plugins.json to help streamline incremental updates.")

	generatorCmd.Flags().String("github-token", "", "The optional GitHub token for API requests.")
	generatorCmd.Flags().Bool("include-pre-release", false, "Whether to include pre-release versions.")
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

		existingPlugins, err := InitCommand(command)
		if err != nil {
			return err
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
			"mattermost-plugin-github",
			"mattermost-plugin-autolink",
			"mattermost-plugin-zoom",
			"mattermost-plugin-jira",
			"mattermost-plugin-welcomebot",
			"mattermost-plugin-jenkins",
			"mattermost-plugin-antivirus",
			"mattermost-plugin-custom-attributes",
			"mattermost-plugin-aws-SNS",
			"mattermost-plugin-gitlab",
			"mattermost-plugin-nps",
			"mattermost-plugin-webex",
		}

		iconPaths := map[string]string{
			"mattermost-plugin-aws-SNS": "data/icons/aws-sns.svg",
			"mattermost-plugin-github":  "data/icons/github.svg",
			"mattermost-plugin-gitlab":  "data/icons/gitlab.svg",
			"mattermost-plugin-jenkins": "data/icons/jenkins.svg",
			"mattermost-plugin-jira":    "data/icons/jira.svg",
			"mattermost-plugin-webex":   "data/icons/webex.svg",
		}

		plugins := []*model.Plugin{}

		for _, repositoryName := range repositoryNames {
			logger.Debugf("querying repository %s", repositoryName)

			var releasePlugins []*model.Plugin
			releasePlugins, err = getReleasePlugins(ctx, client, repositoryName, includePreRelease, existingPlugins)
			if err != nil {
				return errors.Wrapf(err, "failed to release plugin for repository %s", repositoryName)
			}

			for _, plugin := range releasePlugins {
				if len(plugin.IconData) == 0 {
					if iconPath, ok := iconPaths[repositoryName]; ok {
						var iconData string
						iconData, err = getIconDataFromPath(ctx, iconPath)
						if err != nil {
							return errors.Wrapf(err, "failed to fetch icon for repository %s", repositoryName)
						}
						plugin.IconData = iconData
					}
				}

				plugin.Labels = []model.Label{model.OfficialLabel}

				plugins = append(plugins, plugin)
			}
		}

		err = json.NewEncoder(os.Stdout).Encode(plugins)
		if err != nil {
			return errors.Wrap(err, "failed to encode plugins result")
		}

		return nil
	},
}

// getReleasePlugins queries GitHub for all releases of the given plugin, sorting by plugin version descending.
func getReleasePlugins(ctx context.Context, client *github.Client, repositoryName string, includePreRelease bool, existingPlugins []*model.Plugin) ([]*model.Plugin, error) {
	logger := logger.WithField("repository", repositoryName)

	repository, _, err := client.Repositories.Get(ctx, "mattermost", repositoryName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get repository")
	}

	releases, err := getReleases(ctx, client, repositoryName, includePreRelease)
	if err != nil {
		return nil, err
	}
	if len(releases) == 0 {
		logger.Warnf("no releases found for repository")
		return nil, nil
	}

	var plugins []*model.Plugin
	for _, release := range releases {
		plugin, err := getReleasePlugin(release, repository, existingPlugins)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get release plugin for %s", release.GetName())
		}

		if plugin == nil {
			logger.Warnf("no plugin found for release %s", release.GetName())
			continue
		}

		if plugin.Manifest.Version == "" {
			return nil, errors.Errorf("version is empty for manifest.Id %s", plugin.Manifest.Id)
		}

		plugins = append(plugins, plugin)
	}

	// Sort the final slice by plugin version, descending
	sort.SliceStable(
		plugins,
		func(i, j int) bool {
			return semver.MustParse(plugins[i].Manifest.Version).GT(semver.MustParse(plugins[j].Manifest.Version))
		},
	)

	return plugins, nil
}

// getReleases returns all GitHub releases for the given repository.
func getReleases(ctx context.Context, client *github.Client, repoName string, includePreRelease bool) ([]*github.RepositoryRelease, error) {
	var result []*github.RepositoryRelease
	options := &github.ListOptions{
		Page:    0,
		PerPage: 40,
	}
	for {
		releases, resp, err := client.Repositories.ListReleases(ctx, "mattermost", repoName, options)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get releases for repository %s", repoName)
		}

		for _, release := range releases {
			if release.GetDraft() {
				continue
			}

			if release.GetPrerelease() && !includePreRelease {
				continue
			}

			result = append(result, release)
		}

		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return result, nil
}

func getReleasePlugin(release *github.RepositoryRelease, repository *github.Repository, existingPlugins []*model.Plugin) (*model.Plugin, error) {
	var releaseName string
	if release.GetName() == "" {
		releaseName = release.GetTagName()
	} else {
		releaseName = fmt.Sprintf("%s (%s)", release.GetName(), release.GetTagName())
	}
	logger.Debugf("found latest release %s", releaseName)

	downloadURL := ""
	var signatureAsset github.ReleaseAsset
	var foundSignatureAsset bool
	releaseNotesURL := release.GetHTMLURL()
	var updatedAt time.Time
	for _, releaseAsset := range release.Assets {
		assetName := releaseAsset.GetName()
		if strings.Contains(assetName, "-amd64") {
			logger.Debugf("ignoring old style tar bundle %s, for release %s", assetName, releaseName)
			continue
		}

		if strings.HasSuffix(assetName, ".tar.gz") {
			downloadURL = releaseAsset.GetBrowserDownloadURL()
			timestampUpdatedAt := releaseAsset.GetUpdatedAt()
			if timestampUpdatedAt.IsZero() {
				timestampUpdatedAt = releaseAsset.GetCreatedAt()
			}

			updatedAt = timestampUpdatedAt.In(time.UTC)
		}
		if strings.HasSuffix(assetName, ".sig") || strings.HasSuffix(assetName, ".asc") {
			if foundSignatureAsset {
				return nil, errors.Errorf("found multiple signatures %s for release %s", assetName, releaseName)
			}
			signatureAsset = releaseAsset
			foundSignatureAsset = true
		}
	}

	var signature string
	if foundSignatureAsset {
		var err error
		signature, err = downloadSignature(signatureAsset.GetBrowserDownloadURL())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to download signature for release %s", releaseName)
		}
	}

	if downloadURL == "" {
		logger.Warnf("Failed to find plugin asset release %s", releaseName)
		return nil, nil
	}

	var plugin *model.Plugin
	for _, p := range existingPlugins {
		if p.DownloadURL == downloadURL {
			plugin = p
			break
		}
	}

	// If no plugin in existing database or the updated timestamp has changed, attempt to download and inspect manifest.
	if plugin == nil || updatedAt.IsZero() || plugin.UpdatedAt.Before(updatedAt) {
		switch {
		case plugin == nil:
			logger.Debug("no existing plugin")
		case updatedAt.IsZero():
			logger.Debug("no new update timestamp for plugin")
		case plugin.UpdatedAt.IsZero():
			logger.Debug("no recorded update timestamp for plugin")
		case plugin.UpdatedAt.Before(updatedAt):
			logger.Debugf("plugin release asset is newer (+%d seconds)", updatedAt.Sub(plugin.UpdatedAt)/time.Second)
		}

		logger.Debugf("fetching download url %s", downloadURL)

		plugin = &model.Plugin{}

		bundleData, err := downloadBundleData(downloadURL)
		if err != nil {
			return nil, errors.Wrapf(err, "failed download bundle data for release %s", releaseName)
		}

		manifestData, err := getFromTarFile(tar.NewReader(bytes.NewReader(bundleData)), "plugin.json")
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read manifest from plugin bundle for release %s", releaseName)
		}
		plugin.Manifest = mattermostModel.ManifestFromJson(bytes.NewReader(manifestData))
		if plugin.Manifest == nil {
			return nil, errors.Errorf("manifest nil after reading from plugin bundle for release %s", releaseName)
		}

		if plugin.Manifest.IconPath != "" {
			plugin.IconData, err = getIconDataFromTarFile(bundleData, plugin.Manifest.IconPath)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to set icon for release %s", releaseName)
			}
		}
	} else {
		logger.Debugf("skipping download since found existing plugin")
	}

	if plugin.Manifest == nil {
		return nil, fmt.Errorf("failed to find plugin manifest for release %s", releaseName)
	}

	// Reset fields, even if we found the existing plugin above.
	if plugin.Manifest.HomepageURL != "" {
		plugin.HomepageURL = plugin.Manifest.HomepageURL
	} else {
		plugin.HomepageURL = repository.GetHTMLURL()
	}
	plugin.DownloadURL = downloadURL
	plugin.ReleaseNotesURL = releaseNotesURL
	plugin.Signature = signature
	plugin.UpdatedAt = updatedAt

	return plugin, nil
}

func getFromTarFile(reader *tar.Reader, filepath string) ([]byte, error) {
	for {
		hdr, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read tar file")
		}

		// Match the filepath, assuming the tar file contains a leading folder matching the
		// plugin id.
		matched, err := path.Match(fmt.Sprintf("*/%s", filepath), hdr.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to match file %s in tar file", filepath)
		} else if !matched {
			continue
		}

		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read %s in tar file", filepath)
		}
		return data, nil
	}

	return nil, errors.Errorf("failed to find %s in tar file", filepath)
}

func downloadSignature(url string) (string, error) {
	logger.Debugf("fetching signature file from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download signature file from %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("received %d status code while downloading plugin bundle", resp.StatusCode)
	}

	sigFile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open downloaded signature file from %s", url)
	}
	return base64.StdEncoding.EncodeToString(sigFile), nil
}

func downloadBundleData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download plugin bundle")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received %d status code while downloading plugin bundle", resp.StatusCode)
	}

	gzBundleReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read gzipped plugin bundle")
	}

	bundleData, err := ioutil.ReadAll(gzBundleReader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read plugin bundle")
	}

	return bundleData, nil
}

func getIconDataFromTarFile(file []byte, path string) (string, error) {
	iconData, err := getFromTarFile(tar.NewReader(bytes.NewReader(file)), path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read icon data from plugin bundle for path %s", path)
	}

	logger.Debugf("using icon specified in manifest as %s", path)
	return fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(iconData)), nil
}

func getIconDataFromPath(ctx context.Context, path string) (string, error) {
	var icon []byte
	if strings.HasPrefix(path, "http") {
		logger.Debugf("fetching icon from url %s", path)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return "", errors.Wrapf(err, "failed to initialize request to download plugin icon at %s", path)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", errors.Wrapf(err, "failed to download plugin icon at %s", path)
		}
		defer resp.Body.Close()

		icon, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Wrapf(err, "failed to read plugin icon at %s", path)
		}
	} else {
		logger.Debugf("fetching icon from path %s", icon)

		var err error
		icon, err = ioutil.ReadFile(path)
		if err != nil {
			return "", errors.Wrapf(err, "failed to open icon at path %s", path)
		}
	}

	if svg.Is(icon) {
		return fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(icon)), nil
	}

	kind, err := filetype.Image(icon)
	if err != nil {
		return "", errors.Wrap(err, "failed to match icon to image")
	}

	return fmt.Sprintf("data:%s;base64,%s", kind.MIME, base64.StdEncoding.EncodeToString(icon)), nil
}

// InitCommand parses the persistent flags and returns a list of existing plugins, if existing flag is defined.
func InitCommand(command *cobra.Command) ([]*model.Plugin, error) {
	debug, err := command.Flags().GetBool("debug")
	if err != nil {
		return nil, err
	}

	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	var existingPlugins []*model.Plugin
	existingDatabase, err := command.Flags().GetString("existing")
	if err != nil {
		return nil, err
	}

	if existingDatabase != "" {
		file, err := os.Open(existingDatabase)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open existing database %s", existingDatabase)
		}
		defer file.Close()

		existingPlugins, err = model.PluginsFromReader(file)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read existing database %s", existingDatabase)
		}
	}

	return existingPlugins, nil
}
