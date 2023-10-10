package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const (
	deleteButton  = "delete_error_message"
	dismissButton = "dismiss_error_message"
	urlButton     = "url_button"
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
}

var componentHandlers = map[string]func(bot *discordgo.Session, i *discordgo.InteractionCreate){
	deleteButton: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		err := bot.ChannelMessageDelete(i.ChannelID, i.Message.ID)
		if err != nil {
			return
		}
	},
}

// ErrorFollowup sends an error message as a followup message with a deletion button.
func ErrorFollowup(bot *discordgo.Session, i *discordgo.Interaction, errorContent ...any) {
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
	}
	components := []discordgo.MessageComponent{components[deleteButton]}

	_, _ = bot.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content:    errorString,
		Components: components,
	})
}

// ErrorHandler responds to the interaction with an error message and a deletion button.
func ErrorHandler(bot *discordgo.Session, i *discordgo.Interaction, errorContent ...any) {
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
	}
	components := []discordgo.MessageComponent{components[deleteButton]}

	_, _ = bot.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content:    &errorString,
		Components: &components,
	})
}
