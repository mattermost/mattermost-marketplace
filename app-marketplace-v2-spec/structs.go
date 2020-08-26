package cloudappmanifest

type App struct {
	DisplayName string
	Description string
	SiteURL     string

	CommandHandlers []Function

	SlashCommands []SlashCommand
	ModalDialogs  []ModalDialog
	BotDialogs    []BotDialog
	Widgets       []Widget

	// Subscriptions
}

type Function struct {
	// Remote function to execute
	Name string
	Args []Arg
}

type Callable interface {
	// whatever this means
}

type FunctionCall struct {
	Name string
	// #TODO instructions on what to include/expand in the context
	// #TODO other invocation metadata
}

type ModalCall struct {
	Name string
	// #TODO other invocation metadata
}

// A command launches a function, a dialog, or has sub-commands
type SlashCommand struct {
	Callable
	Item

	// Role of the user who should be able to see the autocomplete info of this command
	RoleID      string
	Subcommands []*SlashCommand
}

type ModalDialog struct {
	Callable
	Name    string
	Title   string
	Header  string
	Footer  string
	IconURL string
}

type Widget struct {
	Callable
	// TODO pre-configured arg values
}

type BotDialog struct {
	Callable
	// TODO initial values
	// TODO pre-configured arg values
}

type ItemType string

const (
	TypeCommand       = ItemType("command")
	TypeText          = ItemType("text")
	TypeStaticSelect  = ItemType("static_select")
	TypeDynamicSelect = ItemType("dynamic_select")
	TypeUserSelect    = ItemType("user_select")
	TypeChannelSelect = ItemType("channel_select")
	TypeTime          = ItemType("time")
	TypeBool          = ItemType("bool")
)

type Item struct {
	Name        string
	Type        ItemType
	DisplayName string // for modals
	Hint        string
	Help        string
	IsHidden    bool
	IsRequired  bool
}

type Arg struct {
	Item
}

type TextArg struct {
	Arg
	Subtype   string
	MinLength int
	MaxLength int
	// options - encoding, regexp, etc.
}

type SelectArg struct {
	Arg
	IsMulti bool
}

type UserArg struct {
	SelectArg
}

type ChannelArg struct {
	SelectArg
}

type DynamicSelectArg struct {
	SelectArg
	URLPath string
}

type StaticSelectArg struct {
	SelectArg
	Items []Item
}
