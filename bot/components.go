package discord

import (
	"github.com/bwmarrin/discordgo"
)

const (
	deleteButton  = "delete_error_message"
	dismissButton = "dismiss_error_message"
	urlButton     = "url_button"

	readmoreDismiss = "readmore_dismiss"

	paginationButtons = "pagination_button"
	okCancelButtons   = "ok_cancel_buttons"

	awardUserSelect         = "clippy_award_user"
	awardUserSelectDisabled = "clippy_award_user_disabled"
	awardedUserSelected     = "clippy_awarded"
	undoAward               = "clippy_undo"

	roleSelect = "role_select"
)

var components = map[string]discordgo.MessageComponent{
	deleteButton: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Delete this message",
				Style:    discordgo.DangerButton,
				CustomID: deleteButton,
			},
		},
	},
	urlButton: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Read more",
				Style:    discordgo.LinkButton,
				CustomID: urlButton,
			},
		},
	},
	dismissButton: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Dismiss",
				Style:    discordgo.SecondaryButton,
				CustomID: deleteButton,
			},
		},
	},

	readmoreDismiss: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Read more",
				Style:    discordgo.LinkButton,
				CustomID: urlButton,
			},
			discordgo.Button{
				Label:    "Dismiss",
				Style:    discordgo.SecondaryButton,
				CustomID: deleteButton,
			},
		},
	},

	paginationButtons: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Previous",
				Style:    discordgo.SecondaryButton,
				CustomID: paginationButtons + "_previous",
			},
			discordgo.Button{
				Label:    "Next",
				Style:    discordgo.SecondaryButton,
				CustomID: paginationButtons + "_next",
			},
		},
	},
	okCancelButtons: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "OK",
				Style:    discordgo.SuccessButton,
				CustomID: okCancelButtons + "_ok",
			},
			discordgo.Button{
				Label:    "Cancel",
				Style:    discordgo.DangerButton,
				CustomID: okCancelButtons + "_cancel",
			},
		},
	},

	awardUserSelect: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				MenuType:     discordgo.UserSelectMenu,
				CustomID:     awardUserSelect,
				Placeholder:  "Pick a user to award a clippy point to",
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			},
		},
	},

	awardUserSelectDisabled: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				MenuType:     discordgo.UserSelectMenu,
				CustomID:     awardUserSelect,
				Placeholder:  "Pick a user to award a clippy point to",
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				Disabled:     true,
			},
		},
	},

	undoAward: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Undo",
				Style:    discordgo.DangerButton,
				CustomID: undoAward,
			},
		},
	},

	roleSelect: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				MenuType:    discordgo.RoleSelectMenu,
				CustomID:    roleSelect,
				Placeholder: "Pick a role",
			},
		},
	},
}
