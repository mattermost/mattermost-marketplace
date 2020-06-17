package model

// Label represents a label shown in the Plugin Marketplace UI.
type Label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       string `json:"color"`
}

var AllLabels = []Label{
	CommunityLabel,
	BetaLabel,
}

var CommunityLabel Label = Label{
	Name:        "Community",
	Description: "This plugin is maintained by the Open Source Community.",
	URL:         "https://mattermost.com/pl/default-community-plugins",
}

var BetaLabel Label = Label{
	Name:        "Beta",
	Description: "This plugin is currently in Beta and is not recommended for use in production.",
	URL:         "https://mattermost.com/pl/default-beta-plugins",
}

var EnterpriseLabel Label = Label{
	Name:        "Enterprise",
	Description: "This plugin only works on servers with an E20 license.",
	URL:         "https://mattermost.com/pl/TODO",
}
