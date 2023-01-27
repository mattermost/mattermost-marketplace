package model

// Label represents a label shown in the Plugin Marketplace UI.
type Label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       string `json:"color"`
}

var AllLabels = []Label{
	PartnerLabel,
	CommunityLabel,
	BetaLabel,
	ExperimentalLabel,
	EnterpriseLabel,
}

var PartnerLabel = Label{
	Name:        "Partner",
	Description: "This plugin is maintained by a Mattermost Partner.",
	URL:         "https://mattermost.com/pl/default-partner-plugins",
}
var CommunityLabel = Label{
	Name:        "Community",
	Description: "This plugin is maintained by the Open Source Community.",
	URL:         "https://mattermost.com/pl/default-community-plugins",
}

var BetaLabel = Label{
	Name:        "Beta",
	Description: "This plugin is currently in Beta and is not recommended for use in production.",
	URL:         "https://mattermost.com/pl/default-beta-plugins",
}

var ExperimentalLabel = Label{
	Name:        "Experimental",
	Description: "This plugin is marked as experimental and not meant for production use. Please use with caution.",
	URL:         "https://mattermost.com/pl/default-experimental-plugins",
}

var EnterpriseLabel = Label{
	Name:        "Professional/Enterprise",
	Description: "This plugin requires a Professional or Enterprise subscription.",
	URL:         "https://mattermost.com/pricing/",
}
