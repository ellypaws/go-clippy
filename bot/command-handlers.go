package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sahilm/fuzzy"
	"go-clippy/database/clippy"
	"go-clippy/database/functions"
	"log"
	"regexp"
	"slices"
	"strings"
	"time"
)

// Available methods for *discordgo.Session:

var channels = map[string]chan []string{
	awardUserSelect: make(chan []string), // Not really used, only for debugging
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
		responses[thinkResponse].(newResponseType)(bot, i)
		optionMap := getOpts(i.ApplicationCommandData())

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

		var validChannel bool
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

		// check if thread author/forum author is the same as the one running the command
		if st.OwnerID != i.Member.User.ID && !moderator {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("You are not the owner of <#%v>", channel))
			return
		}

		// fetch forum
		forum, err := bot.Channel(st.ParentID)
		if err != nil {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("Encountered an error while fetching forum: %v", err))
			return
		}

		availableTags := forum.AvailableTags
		if len(availableTags) == 0 {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("<#%v> does not have any tags", forum.Name))
			return
		}
		var tagsToApply []string

		// Mark the thread as solved
		re := regexp.MustCompile(`(?i)\b(re)?solved\b`)
		for _, forumTag := range availableTags {
			if re.MatchString(forumTag.Name) {

				// check if channel already has this tag
				channelAppliedTags := st.AppliedTags

				for _, appliedTag := range channelAppliedTags {
					if appliedTag == forumTag.ID {
						errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("<#%v> is already solved", channel))
						return
					}
				}

				tagsToApply = append(tagsToApply, forumTag.ID)

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
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("<#%v> does not have a valid tag.\n Available tags: %v", forum.ID, availableTags))
		}

		//Award the user
		var username string
		var snowflakes []string
		if option, ok := optionMap[maskedUser]; ok {
			val := option.UserValue(bot)
			username = val.Username
			if val.ID != "" {
				snowflakes = append(snowflakes, val.ID)
			}
		}

		if slices.Contains(snowflakes, i.Member.User.ID) {
			errorEphemeralFollowup(bot, i.Interaction, "You can't award clippy points to yourself")
			return
		}

		if username == "" {
			awardSuggest := responses[ephemeralAwardFollowup].(msgReturnType)(bot, i.Interaction)
			if awardSuggest == nil || awardSuggest.ID == "" {
				errorEdit(bot, i.Interaction, "Encountered an error while sending followup message.")
				return
			}
			channels[awardSuggest.ID] = make(chan []string)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				select {
				case <-time.After(5 * time.Minute):
					channels[awardSuggest.ID] <- []string{}
				case <-ctx.Done():
					return
				}
			}()

			snowflakes = <-channels[awardSuggest.ID]

			cancel()

			// close channel and delete from map
			close(channels[awardSuggest.ID])
			delete(channels, awardSuggest.ID)

			if len(snowflakes) == 0 {
				return
			}

			responses[followupEdit].(editResponseType)(
				bot,
				i.Interaction,
				awardSuggest,
				fmt.Sprintf("Awarding %v and closing channel <#%v>", "<@"+strings.Join(snowflakes, ">, <@")+">", channel),
				components[awardUserSelectDisabled],
				// TODO: Edit or delete message after awarding
				//components[undoAward],
			)
		} else {
			responses[ephemeralFollowup].(msgReturnType)(bot, i.Interaction,
				fmt.Sprintf("Awarding %v and closing channel <#%v>", "<@"+strings.Join(snowflakes, ">, <@")+">", channel))
			// TODO: Edit or delete message after awarding
			recordAward(i)
		}
		responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction,
			fmt.Sprintf("<#%v> is now solved and awarded to users %v", channel,
				"<@"+strings.Join(snowflakes, ">, <@")+">"),
		)
	},

	// awardMessage is similar to solvedCommand, but we don't need to mark the thread as solved
	// this is a right-click context menu where we can run it on a message directly
	// we grab the channel ID from the message and check the type as before to make sure it's valid
	// we check the author of the message to make sure it's the same as the one running the command
	// then we grab the author of the message and check if it's not the same as the one running the command
	// then we award the user
	awardMessage: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[thinkResponse].(newResponseType)(bot, i)

		resolved := i.ApplicationCommandData().Resolved

		var channel string
		var messageID string
		var authorSnowflake string
		var authorUsername string
		for key, message := range resolved.Messages {
			messageID = key
			channel = message.ChannelID
			authorSnowflake = message.Author.ID
			authorUsername = message.Author.Username
		}

		st, err := bot.Channel(channel)
		if err != nil {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("Encountered an error while fetching channel: %v", err))
			err := bot.InteractionResponseDelete(i.Interaction)
			if err != nil {
				log.Printf("Encountered an error while deleting interaction response: %v", err)
			}
			return
		}

		// grab the channel type and check if it's valid
		var validChannel bool
		validChannel = slices.Contains([]discordgo.ChannelType{
			discordgo.ChannelTypeGuildForum,
			discordgo.ChannelTypeGuildNewsThread,
			discordgo.ChannelTypeGuildPublicThread,
			discordgo.ChannelTypeGuildPrivateThread,
		}, st.Type)

		if !validChannel {
			errorEphemeralFollowup(bot, i.Interaction, fmt.Sprintf("<#%v> is not a valid thread", channel))
			err := bot.InteractionResponseDelete(i.Interaction)
			if err != nil {
				log.Printf("Encountered an error while deleting interaction response: %v", err)
			}
			return
		}

		// check if user is the same as the one running the command
		if authorSnowflake == i.Member.User.ID {
			responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction,
				"You can't award clippy points to yourself",
				components[deleteButton],
			)
			return
		}

		responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction, fmt.Sprintf("Awarding clippy points to <@%v>",
			authorSnowflake))

		// TODO: We can probably extract this
		if _, ok := clippy.GetCache().GetConfig(authorSnowflake); !ok {
			clippy.User{
				Username:  authorUsername,
				Snowflake: authorSnowflake,
				OptOut:    false,
				Private:   false,
				Points:    1,
			}.Record()
		}

		guild, err := bot.Guild(i.GuildID)
		if err != nil {
			errorEdit(bot, i.Interaction, err)
			return
		}

		// TODO: We can probably extract this
		// check if award already exists
		a := clippy.Awards.FindByKeys(authorSnowflake + "/" + st.ID)
		if len(a) > 0 {
			responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction,
				fmt.Sprintf("<@%v> already has clippy points in <#%v>", authorSnowflake, st.ID),
				components[deleteButton],
			)
			return
		}

		clippy.Award{
			Username:        authorUsername,
			Snowflake:       authorSnowflake,
			GuildName:       guild.Name,
			GuildID:         guild.ID,
			Channel:         st.Name,
			ChannelID:       st.ID,
			MessageID:       messageID,
			OriginUsername:  i.Member.User.Username,
			OriginSnowflake: i.Member.User.ID,
			InteractionID:   i.ID,
		}.Record()

		responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction, fmt.Sprintf("Awarded clippy points to <@%v>",
			authorSnowflake))
	},

	functionCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		optionMap := getOpts(i.ApplicationCommandData())
		var input string
		var platform = "excel"
		var defaultURL = "https://support.microsoft.com/en-us/office/excel-functions-alphabetical-b3944572-255d-4efb-bb96-c6d90033e188"
		if option, ok := optionMap[requiredFunctionAutocomplete]; ok {
			input = option.StringValue()
		}
		if option, ok := optionMap[optionalPlatformSelection]; ok {
			platform = option.StringValue()
		}
		if platform == "sheets" {
			defaultURL = "https://support.google.com/docs/table/25273?hl=en"
		}
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			responses[thinkResponse].(newResponseType)(bot, i)
			f, err := functions.GetCollection(platform).FindByKey(strings.ToUpper(input))
			if err != nil {
				responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction,
					fmt.Sprintf("[Function](%v) `%v` is not found", defaultURL, input),
					components[deleteButton])
				return
			}
			b, _ := json.MarshalIndent(f, "", "  ")
			// ensure we're under 2000 char limit
			if len(b) > 2000 {
				b = b[:2000]
			}
			if f.URL == "" {
				f.URL = defaultURL
			}
			responses[editInteractionResponse].(msgReturnType)(bot, i.Interaction, string(b),
				discordgo.ActionsRow{[]discordgo.MessageComponent{
					discordgo.Button{
						Label: "Read more",
						Style: discordgo.LinkButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "📜",
						},
						URL: f.URL,
					},
					discordgo.Button{
						Label:    "Delete",
						Style:    discordgo.DangerButton,
						CustomID: deleteButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "🗑️",
						},
					},
					// Autocomplete options introduce a new interaction type (8) for returning custom autocomplete results.
				}})
		case discordgo.InteractionApplicationCommandAutocomplete:
			var choices []*discordgo.ApplicationCommandOptionChoice
			if input != "" {
				log.Printf("Querying for [%v](%v) %v", platform, defaultURL, input)

				cache := functions.Cached(platform)
				results := fuzzy.FindFrom(input, cache)

				for i, r := range results {
					if i > 25 {
						break
					}
					var nameToUse string
					switch {
					case cache[r.Index].Syntax.Layout != "":
						nameToUse = cache[r.Index].Syntax.Layout
					default:
						nameToUse = cache[r.Index].Name
					}

					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  nameToUse,
						Value: cache[r.Index].Name,
					})
				}
			} else {
				choices = []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Type a function name",
						Value: "placeholder",
					},
				}
			}

			if input != "" {
				choices = append(choices[:min(24, len(choices))], &discordgo.ApplicationCommandOptionChoice{
					Name:  input,
					Value: input,
				})
			}

			// make sure we're under 100 char limit and under 25 choices
			for i, choice := range choices {
				if len(choice.Name) > 100 {
					choices[i].Name = choice.Name[:100]
				}
			}
			_ = bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices[:min(25, len(choices))], // This is basically the whole purpose of autocomplete interaction - return custom options to the user.
				},
			})
		}
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
