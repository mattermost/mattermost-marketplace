package model

import (
	"encoding/json"
	"io"
	"time"

	mattermostModel "github.com/mattermost/mattermost-server/model"
)

// Plugin represents a Mattermost plugin in the marketplace.
type Plugin struct {
	HomepageURL     string `json:"homepage_url"`
	IconData        string `json:"icon_data"`
	DownloadURL     string `json:"download_url"`
	ReleaseNotesURL string `json:"release_notes_url"`
	// Signature represents a signature of a plugin saved in base64 encoding.
	Signature string                    `json:"signature"`
	Manifest  *mattermostModel.Manifest `json:"manifest"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

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

// PluginFilter describes the parameters used to constrain a set of plugins.
type PluginFilter struct {
	Page          int
	PerPage       int
	Filter        string
	ServerVersion string
}
