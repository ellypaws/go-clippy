package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go-clippy/database/clippy"
	"log"
	"strconv"
	"strings"
	"time"
)

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

		reply := i.MessageComponentData()
		if i.Member.User.ID == reply.Values[0] {
			errorEphemeral(bot, i.Interaction, "You can't award clippy points to yourself")
			return
		}
		channels[i.Message.ID] <- reply.Values[0]
		responses[ephemeralContent].(msgResponseType)(bot, i.Interaction, "Awarding a clippy point to <@"+reply.Values[0]+">")

		recordAward(i)
	},
	awardedUserSelected: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[ephemeralContent].(msgResponseType)(bot, i.Interaction, "Already awarded to xyz")
	},
}

func newConfig(awarded *discordgo.User) clippy.User {
	return clippy.User{
		Username:  awarded.Username,
		Snowflake: awarded.ID,
		OptOut:    false,
		Private:   false,
		Points:    1,
	}
}

func recordAward(i *discordgo.InteractionCreate) {
	reply := i.MessageComponentData()

	awarded, err := bot.User(reply.Values[0])
	if err != nil {
		errorEdit(bot, i.Interaction, err)
		return
	}

	if _, ok := clippy.GetCache().GetConfig(awarded.ID); !ok {
		newConfig(awarded).Record()
	}

	guild, err := bot.Guild(i.GuildID)
	if err != nil {
		errorEdit(bot, i.Interaction, err)
		return
	}

	channel, err := bot.Channel(i.ChannelID)
	if err != nil {
		errorEdit(bot, i.Interaction, err)
		return
	}

	clippy.Award{
		Username:        awarded.Username,
		Snowflake:       reply.Values[0],
		GuildName:       guild.Name,
		GuildID:         guild.ID,
		Channel:         channel.Name,
		ChannelID:       channel.ID,
		MessageID:       i.Message.ID,
		OriginUsername:  i.Member.User.Username,
		OriginSnowflake: i.Member.User.ID,
		InteractionID:   i.ID,
	}.Record()
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
	blankEmbed, toPrint := errorEmbed(errorContent, i)

	_ = bot.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: toPrint,
			Embeds:  blankEmbed,
		},
	})
}

func errorEphemeralFollowup(bot *discordgo.Session, i *discordgo.Interaction, errorContent ...any) {
	if errorContent == nil || len(errorContent) == 0 {
		return
	}
	blankEmbed, toPrint := errorEmbed(errorContent, i)

	_, _ = bot.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content: *sanitizeToken(&toPrint),
		Embeds:  blankEmbed,
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}

func errorEmbed(errorContent []any, i *discordgo.Interaction) ([]*discordgo.MessageEmbed, string) {
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

	embed := []*discordgo.MessageEmbed{
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

	var toPrint string

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		toPrint = fmt.Sprintf(
			"Could not run the [command](https://github.com/ellypaws/go-clippy) `%v`",
			i.ApplicationCommandData().Name,
		)
	case discordgo.InteractionMessageComponent:
		toPrint = fmt.Sprintf(
			"Could not run the [button](https://github.com/ellypaws/go-clippy) `%v` on message https://discord.com/channels/%v/%v/%v",
			i.MessageComponentData().CustomID,
			i.GuildID,
			i.ChannelID,
			i.Message.ID,
		)
	}
	return embed, toPrint
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
