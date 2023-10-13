package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"slices"
	"time"
)

// Available methods for *discordgo.Session:

var channels = map[string]chan string{
	awardUserSelect: make(chan string), // Not really used, only for debugging
}

func getOpts(data discordgo.ApplicationCommandInteractionData) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := data.Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

var commandHandlers = map[string]func(bot *discordgo.Session, i *discordgo.InteractionCreate){
	helloCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[helloResponse].(newResponseType)(bot, i)
	},
	solvedCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		optionMap := getOpts(i.ApplicationCommandData())

		responses[thinkResponse].(newResponseType)(bot, i)

		var channel string
		if option, ok := optionMap[maskedChannel]; ok {
			channel = option.ChannelValue(bot).ID
		} else {
			channel = i.ChannelID
		}

		st, err := bot.Channel(channel)
		if err != nil {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("Encountered an error while fetching channel: %v", err))
			return
		}

		// TODO: Collection
		validRoles := []string{
			"",
		}

		var moderator bool

		roles := i.Member.Roles
		for _, role := range roles {
			if slices.Contains(validRoles, role) {
				moderator = true
				break
			}
		}

		// check if thread author/forum author is the same as the one running the command
		if st.OwnerID != i.Member.User.ID && !moderator {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("You are not the owner of <#%v>", channel))
			return
		}

		// check if thread is already solved
		re := regexp.MustCompile(`(?i)\b(re)?solved\b`)
		appliedTags := st.AppliedTags
		for _, tag := range appliedTags {
			if re.MatchString(tag) {
				errorFollowup(bot, i.Interaction, fmt.Sprintf("<#%v> is already solved", channel))
				return
			}
		}

		var validChannel bool
		forum, err := bot.Channel(st.ParentID)
		if err != nil {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("Encountered an error while fetching forum: %v", err))
			return
		}

		validChannel = slices.Contains([]discordgo.ChannelType{
			discordgo.ChannelTypeGuildForum,
			discordgo.ChannelTypeGuildNewsThread,
			discordgo.ChannelTypeGuildPublicThread,
			discordgo.ChannelTypeGuildPrivateThread,
		}, st.Type)

		if !validChannel {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("<#%v> is not a valid thread", channel))
			return
		}

		tags := forum.AvailableTags
		if len(tags) == 0 {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("<#%v> does not have any tags", forum.Name))
			return
		}
		var tagsToApply []string

		// Mark the thread as solved
		for _, tag := range tags {
			if re.MatchString(tag.Name) {
				tagsToApply = append(tagsToApply, tag.ID)
				_, err := bot.ChannelEdit(channel, &discordgo.ChannelEdit{
					AppliedTags: &tagsToApply,
				})
				responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction,
					fmt.Sprintf("<#%v> is now solved", channel),
				)
				if err != nil {
					errorEphemeral(bot, i.Interaction, fmt.Sprintf("Encountered an error while editing channel: %v", err))
					return
				}
				break
			}
		}
		if len(tagsToApply) == 0 {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("<#%v> does not have a valid tag.\n Available tags: %v", forum.ID, tags))
		}

		//Award the user
		var username string
		var snowflake string
		if option, ok := optionMap[maskedUser]; ok {
			val := option.UserValue(bot)
			username = val.Username
			if val.ID != "" {
				snowflake = val.ID
			} else {
				snowflake = val.Username
			}
		}

		if username == "" {
			awardSuggest := responses[ephemeralAwardFollowup].(msgReturnType)(bot, i.Interaction)
			channels[awardSuggest.ID] = make(chan string)
			username = <-channels[awardSuggest.ID]
			responses[followupEdit].(editResponseType)(
				bot,
				i.Interaction,
				awardSuggest,
				fmt.Sprintf("Awarding <@%v> and closing channel <#%v>", username, channel),
				components[awardUserSelectDisabled],
				//components[undoAward],
			)
		} else {
			responses[followupResponse].(msgReturnType)(bot, i.Interaction, fmt.Sprintf("Awarding <@%v> and closing channel <#%v>", snowflake, channel))
		}
		responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction,
			fmt.Sprintf("<#%v> "+
				"is now solved and awarded to <@%v>", channel, username),
		)
	},
	functionCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {

		responses[thinkResponse].(newResponseType)(bot, i)
		responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction, "Testing editing response")
		time.Sleep(time.Second * 2)
		responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction, "Now with a button", components[okCancelButtons])

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
		// TODO: Check if options are populated
		responses[messageResponse].(msgResponseType)(bot, i.Interaction, "Searching for "+i.ApplicationCommandData().Options[0].StringValue()+"...")
	},

	addModerator: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[ephemeralContent].(msgResponseType)(bot, i.Interaction, "Adding moderator role to server")
	},
}
