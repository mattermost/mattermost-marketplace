package model

import (
	"encoding/json"
	"io"

	mattermostModel "github.com/mattermost/mattermost-server/model"
)

// PluginSignature is a public key signature of a plugin and the corresponding public key hash for use in verifying a plugin downloaded from the marketplace.
type PluginSignature struct {
	// Signature represents a signature of a plugin saved in base64 encoding.
	Signature string `json:"signature"`
	// PublicKeyHash represents first arbitrary number of symbols of the
	// public key fingerprint, hashed using SHA-1 algorithm.
	PublicKeyHash string `json:"public_key_hash"`
}

// Plugin represents a Mattermost plugin in the marketplace.
type Plugin struct {
	HomepageURL string                    `json:"homepage_url"`
	IconData    string                    `json:"icon_data"`
	DownloadURL string                    `json:"download_url"`
	Manifest    *mattermostModel.Manifest `json:"manifest"`
	Signatures  []*PluginSignature        `json:"signatures"`
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
