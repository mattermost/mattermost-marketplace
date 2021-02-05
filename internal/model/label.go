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
	EnterpriseLabel,
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
	Name:        "Enterprise (E20 and Cloud)",
	Description: "This plugin only works on self-managed deployments (E20) and Mattermost Cloud workspaces.",
	URL:         "https://mattermost.com/pricing/",
}
