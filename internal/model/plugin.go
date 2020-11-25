package model

import (
	"encoding/json"
	"io"
	"time"

	mattermostModel "github.com/mattermost/mattermost-server/v5/model"
)

type HostingType string

const (
	OnPrem HostingType = "on-prem"
	Cloud  HostingType = "cloud"
)

type AuthorType string

const (
	Mattermost AuthorType = "mattermost"
	Community  AuthorType = "community"
)

type ReleaseStage string

const (
	Production ReleaseStage = "production"
	Beta       ReleaseStage = "beta"
)

// Plugin represents a Mattermost plugin in the Plugin Marketplace.
type Plugin struct {
	HomepageURL     string                    `json:"homepage_url"`
	IconData        string                    `json:"icon_data"` // The base64 encoding of an svg image
	DownloadURL     string                    `json:"download_url"`
	ReleaseNotesURL string                    `json:"release_notes_url"`
	Labels          []Label                   `json:"labels,omitempty"`
	Hosting         HostingType               `json:"hosting"`       // Indicated if the plugin is limited to a certain hosting type
	AuthorType      AuthorType                `json:"author_type"`   // The maintainer of the plugin
	ReleaseStage    ReleaseStage              `json:"release_stage"` // The stage in the software release cycle that the plugin is in
	Enterprise      bool                      `json:"enterprise"`    // Indicated if the plugin is an enterprise plugin
	Signature       string                    `json:"signature"`     // A signature of a plugin saved in base64 encoding.
	RepoName        string                    `json:"repo_name"`
	Manifest        *mattermostModel.Manifest `json:"manifest"`
	Platforms       PlatformBundles           `json:"platforms"`
	UpdatedAt       time.Time                 `json:"updated_at"` // The point in time this release of the plugin was added to the Plugin Marketplace
}

// PlatformBundleMetadata holds the necessary data to fetch and verify a plugin built for a specific platform
type PlatformBundleMetadata struct {
	DownloadURL string `json:"download_url,omitempty"`
	Signature   string `json:"signature,omitempty"`
}

type PlatformBundles struct {
	LinuxAmd64   PlatformBundleMetadata `json:"linux-amd64,omitempty" yaml:"linux-amd64,omitempty"`
	DarwinAmd64  PlatformBundleMetadata `json:"darwin-amd64,omitempty" yaml:"darwin-amd64,omitempty"`
	WindowsAmd64 PlatformBundleMetadata `json:"windows-amd64,omitempty" yaml:"windows-amd64,omitempty"`
}

const (
	LinuxAmd64   = "linux-amd64"
	DarwinAmd64  = "darwin-amd64"
	WindowsAmd64 = "windows-amd64"
)

// PluginFromReader decodes a json-encoded cluster from the given io.Reader.
func PluginFromReader(reader io.Reader) (*Plugin, error) {
	cluster := Plugin{}
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&cluster)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &cluster, nil
}

// PluginsFromReader decodes a json-encoded list of plugins from the given io.Reader.
func PluginsFromReader(reader io.Reader) ([]*Plugin, error) {
	plugins := []*Plugin{}
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(&plugins)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return plugins, nil
}

// PluginsToWriter encodes a json-encoded list of plugins to the given io.Writer.
func PluginsToWriter(w io.Writer, plugins []*Plugin) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(plugins)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (p *Plugin) AddLabels() {
	if p.AuthorType == Community {
		p.Labels = append(p.Labels, CommunityLabel)
	}

	if p.ReleaseStage == Beta {
		p.Labels = append(p.Labels, BetaLabel)
	}

	if p.Enterprise {
		p.Labels = append(p.Labels, EnterpriseLabel)
	}
}

// PluginFilter describes the parameters used to constrain a set of plugins.
type PluginFilter struct {
	Page              int
	PerPage           int
	Filter            string
	ServerVersion     string
	EnterprisePlugins bool
	Cloud             bool
	Platform          string
	ReturnAllVersions bool
	PluginId          string
}
