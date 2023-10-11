package discord

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

// Available methods for *discordgo.Session:

var commandHandlers = map[string]func(bot *discordgo.Session, i *discordgo.InteractionCreate){
	helloCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[helloResponse].(newResponseType)(bot, i)
	},
	solvedCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[pendingResponse].(newResponseType)(bot, i)
		responses[ErrorResponse].(errorResponseType)(bot, i.Interaction, "This command is not implemented yet")
		responses[ErrorFollowup].(errorResponseType)(bot, i.Interaction, "Testing followup error message")

		responses[awardClippy].(newResponseType)(bot, i)
		responses[ephemeralContentResponse].(msgResponseType)(bot, i.Interaction, "Awarding clippy to <@"+i.Member.User.ID+">")

		responses[ephemeralAwardUser].(newResponseType)(bot, i)
	},
	functionCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {

		responses[thinkResponse].(newResponseType)(bot, i)
		responses[editResponse].(msgReturnType)(bot, i.Interaction, "Testing editing response")
		time.Sleep(time.Second * 2)
		responses[editResponse].(msgReturnType)(bot, i.Interaction, "Now with a button", components[okCancelButtons])

		msg := responses[followupResponse].(msgReturnType)(bot, i.Interaction, "Testing followup message", components[paginationButtons])
		time.Sleep(time.Second * 2)
		responses[followupEdit].(editResponseType)(bot, i.Interaction, msg, "Editing followup message")

		time.Sleep(time.Second * 2)
		a := responses[followupResponse].(msgReturnType)(bot, i.Interaction, "Now let's see if we're deleting the correct message")
		time.Sleep(time.Second * 2)
		responses[followupEdit].(editResponseType)(bot, i.Interaction, a, components[deleteButton])
	},

	searchCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[thinkResponse].(newResponseType)(bot, i)
		responses[messageResponse].(msgResponseType)(bot, i.Interaction, "Searching for "+i.ApplicationCommandData().Options[0].StringValue()+"...")
	},
}
