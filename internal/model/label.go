package model

// Label represents a label shown in the Plugin Marketplace UI.
type Label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       string `json:"color"`
}

var OfficialLabel Label = Label{
	Name:        "Official",
	Description: "This plugin is maintained by Mattermost",
	URL:         "https://mattermost.com/pl/default-community-plugins",
	Color:       "#166de0",
}
