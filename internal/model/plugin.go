package model

import (
	"encoding/json"
	"io"

	mattermostModel "github.com/mattermost/mattermost-server/model"
)

// Plugin represents a Mattermost plugin in the marketplace.
type Plugin struct {
	HomepageURL       string
	DownloadURL       string
	DownloadSignature []byte
	Manifest          *mattermostModel.Manifest
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
