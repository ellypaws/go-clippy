package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

const (
	deleteButton  = "delete_error_message"
	dismissButton = "dismiss_error_message"
	urlButton     = "url_button"

	readmoreDismiss = "readmore_dismiss"

	paginationButton = "pagination_button"
	okCancelButtons  = "ok_cancel_buttons"
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

	paginationButton: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Previous",
				Style:    discordgo.SecondaryButton,
				CustomID: paginationButton + "_previous",
			},
			discordgo.Button{
				Label:    "Next",
				Style:    discordgo.SecondaryButton,
				CustomID: paginationButton + "_next",
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
}

var componentHandlers = map[string]func(bot *discordgo.Session, i *discordgo.InteractionCreate){
	deleteButton: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		err := bot.ChannelMessageDelete(i.ChannelID, i.Message.ID)
		if err != nil {
			errorEphemeral(bot, i.Interaction, err)
		}
	},
}

// errorFollowup [ErrorFollowup] sends an error message as a followup message with a deletion button.
func errorFollowup(bot *discordgo.Session, i *discordgo.Interaction, errorContent ...any) {
	if errorContent == nil || len(errorContent) == 0 {
		return
	}
	var errorString string

	switch content := errorContent[0].(type) {
	case string:
		errorString = content
	case error:
		errorString = fmt.Sprint(content) // Convert the error to a string
	default:
		errorString = "An unknown error has occurred"
		errorString += "\nReceived:" + fmt.Sprint(content)
	}
	components := []discordgo.MessageComponent{components[deleteButton]}

	_, _ = bot.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content:    *sanitizeToken(&errorString),
		Components: components,
	})
}

// errorEdit [ErrorResponse] responds to the interaction with an error message and a deletion button.
func errorEdit(bot *discordgo.Session, i *discordgo.Interaction, errorContent ...any) {
	if errorContent == nil || len(errorContent) == 0 {
		return
	}
	var errorString string

	switch content := errorContent[0].(type) {
	case string:
		errorString = content
	case error:
		errorString = fmt.Sprint(content) // Convert the error to a string
	default:
		errorString = "An unknown error has occurred"
		errorString += "\nReceived:" + fmt.Sprint(content)
	}

	components := []discordgo.MessageComponent{components[deleteButton]}

	_, _ = bot.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content:    sanitizeToken(&errorString),
		Components: &components,
	})
}

// errorEphemeral [ErrorEphemeral] responds to the interaction with an ephemeral error message when the deletion button doesn't work.
func errorEphemeral(bot *discordgo.Session, i *discordgo.Interaction, errorContent ...any) {
	if errorContent == nil || len(errorContent) == 0 {
		return
	}
	var errorString string

	switch content := errorContent[0].(type) {
	case string:
		errorString = content
	case error:
		errorString = fmt.Sprint(content) // Convert the error to a string
	default:
		errorString = "An unknown error has occurred"
		errorString += "\nReceived:" + fmt.Sprint(content)
	}

	_ = bot.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Could not run the button, got an error: " + *sanitizeToken(&errorString),
		},
	})
}

func sanitizeToken(errorString *string) *string {
	if strings.Contains(*errorString, *botToken) {
		log.Println("WARNING: Bot token was found in the error message. Replacing it with \"Bot Token\"")
		log.Println("Error message:", errorString)
		sanitizedString := strings.ReplaceAll(*errorString, *botToken, "[TOKEN]")
		errorString = &sanitizedString
	}
	return errorString
}
