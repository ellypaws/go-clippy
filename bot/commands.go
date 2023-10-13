package discord

import "github.com/bwmarrin/discordgo"

const (
	requiredFunction                = "function"
	optionalFunctionAutocomplete    = "function-autocomplete"
	optionalDescriptionAutocomplete = "description-autocomplete"
	optionalPlatformSelection       = "platform-selection"

	maskedUser    = "user"
	maskedChannel = "channel"
	maskedForum   = "threads"
	maskedRole    = "role"
)

const (
	helloCommand    = "hello"
	solvedCommand   = "solved"
	functionCommand = "function"
	searchCommand   = "search"
	addModerator    = "moderator"
)

var commandOptions = map[string]*discordgo.ApplicationCommandOption{
	requiredFunction: {
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        requiredFunction,
		Description: "The function to look up",
		Required:    true,
	},
	optionalFunctionAutocomplete: {
		Type:         discordgo.ApplicationCommandOptionString,
		Name:         optionalFunctionAutocomplete,
		Description:  "The function to look up",
		Required:     false,
		Autocomplete: true,
	},
	optionalDescriptionAutocomplete: {
		Type:         discordgo.ApplicationCommandOptionString,
		Name:         optionalDescriptionAutocomplete,
		Description:  "The description to look up",
		Required:     false,
		Autocomplete: true,
	},
	optionalPlatformSelection: {
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        optionalPlatformSelection,
		Description: "Excel or Google Sheets",
		Required:    false,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Excel",
				Value: "excel",
			},
			{
				Name:  "Google Sheets",
				Value: "sheets",
			},
		},
	},
}

var maskedOptions = map[string]*discordgo.ApplicationCommandOption{
	maskedUser: {
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        maskedUser,
		Description: "Choose a user",
		Required:    false,
	},
	maskedChannel: {
		Type:        discordgo.ApplicationCommandOptionChannel,
		Name:        maskedChannel,
		Description: "Choose a channel to close",
		// Channel type mask
		ChannelTypes: []discordgo.ChannelType{
			discordgo.ChannelTypeGuildText,
			discordgo.ChannelTypeGuildVoice,
		},
		Required: false,
	},
	maskedForum: {
		Type:        discordgo.ApplicationCommandOptionChannel,
		Name:        maskedForum,
		Description: "Choose a thread to mark as solved",
		ChannelTypes: []discordgo.ChannelType{
			discordgo.ChannelTypeGuildForum,
			discordgo.ChannelTypeGuildNewsThread,
			discordgo.ChannelTypeGuildPublicThread,
			discordgo.ChannelTypeGuildPrivateThread,
		},
	},
	maskedRole: {
		Type:        discordgo.ApplicationCommandOptionRole,
		Name:        maskedRole,
		Description: "Choose a role to add",
		Required:    false,
	},
}

var commands = map[string]*discordgo.ApplicationCommand{
	helloCommand: {
		Name: helloCommand,
		// All commands and options must have a description
		// Commands/options without description will fail the registration
		// of the command.
		Description: "Say hello to the bot",
		Type:        discordgo.ChatApplicationCommand,
	},
	solvedCommand: {
		Name:        solvedCommand,
		Description: "Mark a question as solved and optionally close the thread",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			maskedOptions[maskedUser],
			maskedOptions[maskedForum],
		},
	},
	functionCommand: {
		Name:        functionCommand,
		Description: "Look up a function in Excel or Google Sheets",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			commandOptions[requiredFunction],
			commandOptions[optionalPlatformSelection],
		},
	},
	searchCommand: {
		Name:        searchCommand,
		Description: "Search for a function in Excel or Google Sheets",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			commandOptions[optionalFunctionAutocomplete],
			commandOptions[optionalDescriptionAutocomplete],
			commandOptions[optionalPlatformSelection],
		},
	},
	addModerator: {
		Name:        addModerator,
		Description: "Add a moderator role to the server",
		Type:        discordgo.ChatApplicationCommand,
	},
}

// ----- UNUSED COMMANDS -----

//var unusedCommands = map[string]*discordgo.ApplicationCommand{
//	"Basic response testing": {
//		Name:        "responses",
//		Description: "Interaction responses testing initiative",
//		Options: []*discordgo.ApplicationCommandOption{
//			{
//				Name:        "resp-type",
//				Description: "Response type",
//				Type:        discordgo.ApplicationCommandOptionInteger,
//				Choices: []*discordgo.ApplicationCommandOptionChoice{
//					{
//						Name:  "Channel message with source",
//						Value: 4,
//					},
//					{
//						Name:  "Deferred response With Source",
//						Value: 5,
//					},
//				},
//				Required: true,
//			},
//		},
//	},
//
//	"Followup testing": {
//		Name:        "followups",
//		Description: "Followup messages",
//	},
//
//	"permission-overview": {
//		Name:                     "permission-overview",
//		Description:              "Command for demonstration of default command permissions",
//		DefaultMemberPermissions: &defaultMemberPermissions,
//		DMPermission:             &dmPermission,
//	},
//
//	"basic-command-with-files": {
//		Name:        "basic-command-with-files",
//		Description: "Basic command with files",
//	},
//
//	"localization": {
//		Name:        "localized-command",
//		Description: "Localized command. Description and name may vary depending on the Language setting",
//		NameLocalizations: &map[discordgo.Locale]string{
//			discordgo.ChineseCN: "本地化的命令",
//		},
//		DescriptionLocalizations: &map[discordgo.Locale]string{
//			discordgo.ChineseCN: "这是一个本地化的命令",
//		},
//		Options: []*discordgo.ApplicationCommandOption{
//			{
//				Name:        "localized-option",
//				Description: "Localized option. Description and name may vary depending on the Language setting",
//				NameLocalizations: map[discordgo.Locale]string{
//					discordgo.ChineseCN: "一个本地化的选项",
//				},
//				DescriptionLocalizations: map[discordgo.Locale]string{
//					discordgo.ChineseCN: "这是一个本地化的选项",
//				},
//				Type: discordgo.ApplicationCommandOptionInteger,
//				Choices: []*discordgo.ApplicationCommandOptionChoice{
//					{
//						Name: "First",
//						NameLocalizations: map[discordgo.Locale]string{
//							discordgo.ChineseCN: "一的",
//						},
//						Value: 1,
//					},
//					{
//						Name: "Second",
//						NameLocalizations: map[discordgo.Locale]string{
//							discordgo.ChineseCN: "二的",
//						},
//						Value: 2,
//					},
//				},
//			},
//		},
//	},
//
//	"subcommands": {
//		Name:        "subcommands",
//		Description: "Subcommands and command groups example",
//		Options: []*discordgo.ApplicationCommandOption{
//			// When a command has subcommands/subcommand groups
//			// It must not have top-level options, they aren't accesible in the UI
//			// in this case (at least not yet), so if a command has
//			// subcommands/subcommand any groups registering top-level options
//			// will cause the registration of the command to fail
//
//			{
//				Name:        "subcommand-group",
//				Description: "Subcommands group",
//				Options: []*discordgo.ApplicationCommandOption{
//					// Also, subcommand groups aren't capable of
//					// containing options, by the name of them, you can see
//					// they can only contain subcommands
//					{
//						Name:        "nested-subcommand",
//						Description: "Nested subcommand",
//						Type:        discordgo.ApplicationCommandOptionSubCommand,
//					},
//				},
//				Type: discordgo.ApplicationCommandOptionSubCommandGroup,
//			},
//			// Also, you can create both subcommand groups and subcommands
//			// in the command at the same time. But, there's some limits to
//			// nesting, count of subcommands (top level and nested) and options.
//			// Read the intro of slash-commands docs on Discord dev portal
//			// to get more information
//			{
//				Name:        "subcommand",
//				Description: "Top-level subcommand",
//				Type:        discordgo.ApplicationCommandOptionSubCommand,
//			},
//		},
//	},
//}

// ----- UNUSED OPTIONS -----

//var unusedOptions = []*discordgo.ApplicationCommandOption{
//	{
//		Type:        discordgo.ApplicationCommandOptionInteger,
//		Name:        "integer-option",
//		Description: "Integer option",
//		MinValue:    &integerOptionMinValue,
//		MaxValue:    10,
//		Required:    true,
//	},
//	{
//		Type:        discordgo.ApplicationCommandOptionNumber,
//		Name:        "number-option",
//		Description: "Float option",
//		MaxValue:    10.1,
//		Required:    true,
//	},
//	{
//		Type:        discordgo.ApplicationCommandOptionBoolean,
//		Name:        "bool-option",
//		Description: "Boolean option",
//		Required:    true,
//	},
//
//	//Required options must be listed first since optional parameters
//	//always come after when they're used.
//	//The same concept applies to Discord's Slash-commands API
//
//	{
//		Type:        discordgo.ApplicationCommandOptionChannel,
//		Name:        "channel-option",
//		Description: "Channel option",
//		// Channel type mask
//		ChannelTypes: []discordgo.ChannelType{
//			discordgo.ChannelTypeGuildText,
//			discordgo.ChannelTypeGuildVoice,
//		},
//		Required: false,
//	},
//	{
//		Type:        discordgo.ApplicationCommandOptionUser,
//		Name:        "user-option",
//		Description: "User option",
//		Required:    false,
//	},
//	{
//		Type:        discordgo.ApplicationCommandOptionRole,
//		Name:        "role-option",
//		Description: "Role option",
//		Required:    false,
//	},
//}
