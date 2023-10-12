package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	deleteButton  = "delete_error_message"
	dismissButton = "dismiss_error_message"
	urlButton     = "url_button"

	readmoreDismiss = "readmore_dismiss"

	paginationButtons = "pagination_button"
	okCancelButtons   = "ok_cancel_buttons"

	awardUserSelect     = "clippy_award_user"
	awardedUserSelected = "clippy_awarded"
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
}

var componentHandlers = map[string]func(bot *discordgo.Session, i *discordgo.InteractionCreate){
	deleteButton: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		err := bot.ChannelMessageDelete(i.ChannelID, i.Message.ID)
		if err != nil {
			errorEphemeral(bot, i.Interaction, err)
		}
	},
	awardUserSelect: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		if time.Now().Sub(i.Interaction.Message.Timestamp) > 5*time.Minute {
			errorEphemeral(bot, i.Interaction, "This message is too old to award clippy points")
			return
		}
		//log.Printf("[COMPONENT] InteractionCreate: %+v\n", i)
		//log.Printf("[COMPONENT] InteractionCreate.Interaction: %+v\n", i.Interaction)
		//log.Printf("[COMPONENT] InteractionCreate.Interaction.Message: %+v\n", i.Interaction.Message)
		channels[i.Message.ID] <- true
		responses[ephemeralContentResponse].(msgResponseType)(bot, i.Interaction, "Awarding a clippy point to <@"+i.MessageComponentData().Values[0]+">")
	},
	awardedUserSelected: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[ephemeralContentResponse].(msgResponseType)(bot, i.Interaction, "Already awarded to xyz")
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

	logError(errorString, i)

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

	logError(errorString, i)

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

	logError(errorString, i)

	// decode ED4245 to int
	color, _ := strconv.ParseInt("ED4245", 16, 64)

	blankEmbed := []*discordgo.MessageEmbed{
		{
			Type: discordgo.EmbedTypeRich,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Error",
					Value:  *sanitizeToken(&errorString),
					Inline: false,
				},
			},
			Color: int(color),
		},
	}

	var action string
	var id string

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		action = "command"
		id = i.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		action = "button"
		id = i.MessageComponentData().CustomID
	}

	_ = bot.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags: discordgo.MessageFlagsEphemeral & discordgo.MessageFlagsSupressEmbeds,
			Content: fmt.Sprintf(
				"Could not run the [%v](https://github.com/ellypaws/go-clippy) `%v` on message https://discord.com/channels/%v/%v/%v",
				action,
				id,
				i.GuildID,
				i.ChannelID,
				i.Message.ID,
			),
			Embeds: blankEmbed,
		},
	})
}

func logError(errorString string, i *discordgo.Interaction) {
	log.Printf("WARNING: A command failed to execute: %v", errorString)
	if i.Type == discordgo.InteractionMessageComponent {
		log.Printf("Command: %v", i.MessageComponentData().CustomID)
	}
	log.Printf("User: %v", i.Member.User.Username)
	if i.Type == discordgo.InteractionMessageComponent {
		log.Printf("Link: https://discord.com/channels/%v/%v/%v", i.GuildID, i.ChannelID, i.Message.ID)
	}
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
