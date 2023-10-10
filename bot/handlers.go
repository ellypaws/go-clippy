package discord

import (
	"github.com/bwmarrin/discordgo"
)

// Available methods for *discordgo.Session:

var commandHandlers = map[string]func(bot *discordgo.Session, i *discordgo.InteractionCreate){
	helloCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[helloResponse].(regularResponseType)(bot, i)
	},
	solvedCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		responses[pendingResponse].(regularResponseType)(bot, i)
		responses[errorResponse].(errorResponseType)(bot, i.Interaction, "This command is not implemented yet")
		responses[errorFollowup].(errorResponseType)(bot, i.Interaction, "Testing followup error message")
	},
	functionCommand: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {

		responses[ephemeralResponse].(regularResponseType)(bot, i)
		responses[thinkResponse].(regularResponseType)(bot, i)
	},
}

const (
	thinkResponse = iota
	pendingResponse
	ephemeralResponse
	helloResponse
	errorResponse
	errorFollowup
)

type regularResponseType func(bot *discordgo.Session, i *discordgo.InteractionCreate)
type editResponseType func(bot *discordgo.Session, i *discordgo.Interaction)
type errorResponseType func(bot *discordgo.Session, i *discordgo.Interaction, errorContent ...any)
type followupResponseType editResponseType

var responses = map[int]any{
	thinkResponse: func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		})
		ErrorHandler(bot, i.Interaction, err)
	},
	pendingResponse: regularResponseType(func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// Note: this isn't documented, but you can use that if you want to.
				// This flag just allows you to create messages visible only for the caller of the command
				// (user who triggered the command)
				//Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Bot is responding...",
			},
		})
		ErrorHandler(bot, i.Interaction, err)
	}),
	ephemeralResponse: regularResponseType(func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// Note: this isn't documented, but you can use that if you want to.
				// This flag just allows you to create messages visible only for the caller of the command
				// (user who triggered the command)
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Bot is responding...",
			},
		})
		ErrorHandler(bot, i.Interaction, err)
	}),
	helloResponse: regularResponseType(func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there! Congratulations, you just executed your first slash command",
			},
		})
		ErrorHandler(bot, i.Interaction, err)
	}),
	errorResponse: errorResponseType(ErrorHandler),
	errorFollowup: errorResponseType(ErrorFollowup),
}

// ----- UNUSED COMMAND HANDLERS -----

//var unusedCommandHandlers = map[string]func(bot *discordgo.Session, i *discordgo.InteractionCreate){
//	"basic-command-with-files": func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
//		bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: "Hey there! Congratulations, you just executed your first slash command with a file in the response",
//				Files: []*discordgo.File{
//					{
//						ContentType: "text/plain",
//						Name:        "test.txt",
//						Reader:      strings.NewReader("Hello Discord!!"),
//					},
//				},
//			},
//		})
//	},
//	"localized-command": func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
//		responses := map[discordgo.Locale]string{
//			discordgo.ChineseCN: "你好！ 这是一个本地化的命令",
//		}
//		response := "Hi! This is a localized message"
//		if r, ok := responses[i.Locale]; ok {
//			response = r
//		}
//		err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: response,
//			},
//		})
//		if err != nil {
//			panic(err)
//		}
//	},
//	"options": func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
//		// Access options in the order provided by the user.
//		options := i.ApplicationCommandData().Options
//
//		// Or convert the slice into a map
//		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
//		for _, opt := range options {
//			optionMap[opt.Name] = opt
//		}
//
//		// This example stores the provided arguments in an []interface{}
//		// which will be used to format the bot's response
//		margs := make([]interface{}, 0, len(options))
//		msgformat := "You learned how to use command options! " +
//			"Take a look at the value(s) you entered:\n"
//
//		// Get the value from the option map.
//		// When the option exists, ok = true
//		if option, ok := optionMap["function"]; ok {
//			// Option values must be type asserted from interface{}.
//			// Discordgo provides utility functions to make this simple.
//			margs = append(margs, option.StringValue())
//			msgformat += "> string-option: %s\n"
//		}
//		bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			// Ignore type for now, they will be discussed in "responses"
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: fmt.Sprintf(
//					msgformat,
//					margs...,
//				),
//			},
//		})
//	},
//	"permission-overview": func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
//		perms, err := bot.ApplicationCommandPermissions(bot.State.User.ID, i.GuildID, i.ApplicationCommandData().ID)
//
//		var restError *discordgo.RESTError
//		if errors.As(err, &restError) && restError.Message != nil && restError.Message.Code == discordgo.ErrCodeUnknownApplicationCommandPermissions {
//			bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: ":x: No permission overwrites",
//				},
//			})
//			return
//		} else if err != nil {
//			panic(err)
//		}
//
//		if err != nil {
//			panic(err)
//		}
//		format := "- %s %s\n"
//
//		channels := ""
//		users := ""
//		roles := ""
//
//		for _, o := range perms.Permissions {
//			emoji := "❌"
//			if o.Permission {
//				emoji = "☑"
//			}
//
//			switch o.Type {
//			case discordgo.ApplicationCommandPermissionTypeUser:
//				users += fmt.Sprintf(format, emoji, "<@!"+o.ID+">")
//			case discordgo.ApplicationCommandPermissionTypeChannel:
//				allChannels, _ := discordgo.GuildAllChannelsID(i.GuildID)
//
//				if o.ID == allChannels {
//					channels += fmt.Sprintf(format, emoji, "All channels")
//				} else {
//					channels += fmt.Sprintf(format, emoji, "<#"+o.ID+">")
//				}
//			case discordgo.ApplicationCommandPermissionTypeRole:
//				if o.ID == i.GuildID {
//					roles += fmt.Sprintf(format, emoji, "@everyone")
//				} else {
//					roles += fmt.Sprintf(format, emoji, "<@&"+o.ID+">")
//				}
//			}
//		}
//
//		bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Embeds: []*discordgo.MessageEmbed{
//					{
//						Title:       "Permissions overview",
//						Description: "Overview of permissions for this command",
//						Fields: []*discordgo.MessageEmbedField{
//							{
//								Name:  "Users",
//								Value: users,
//							},
//							{
//								Name:  "Channels",
//								Value: channels,
//							},
//							{
//								Name:  "Roles",
//								Value: roles,
//							},
//						},
//					},
//				},
//				AllowedMentions: &discordgo.MessageAllowedMentions{},
//			},
//		})
//	},
//	"subcommands": func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
//		options := i.ApplicationCommandData().Options
//		content := ""
//
//		// As you can see, names of subcommands (nested, top-level)
//		// and subcommand groups are provided through the arguments.
//		switch options[0].Name {
//		case "subcommand":
//			content = "The top-level subcommand is executed. Now try to execute the nested one."
//		case "subcommand-group":
//			options = options[0].Options
//			switch options[0].Name {
//			case "nested-subcommand":
//				content = "Nice, now you know how to execute nested commands too"
//			default:
//				content = "Oops, something went wrong.\n" +
//					"Hol' up, you aren't supposed to see this message."
//			}
//		}
//
//		bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: content,
//			},
//		})
//	},
//	"responses": func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
//		// Responses to a command are very important.
//		// First of all, because you need to react to the interaction
//		// by sending the response in 3 seconds after receiving, otherwise
//		// interaction will be considered invalid and you can no longer
//		// use the interaction token and ID for responding to the user's request
//
//		content := ""
//		// As you can see, the response type names used here are pretty self-explanatory,
//		// but for those who want more information see the official documentation
//		switch i.ApplicationCommandData().Options[0].IntValue() {
//		case int64(discordgo.InteractionResponseChannelMessageWithSource):
//			content =
//				"You just responded to an interaction, sent a message and showed the original one. " +
//					"Congratulations!"
//			content +=
//				"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
//		case int64(discordgo.InteractionResponseDeferredChannelMessageWithSource):
//			err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
//			})
//			if err != nil {
//				bot.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
//					Content: "Something went wrong 01",
//				})
//			}
//			time.AfterFunc(time.Second*5, func() {
//				content = "Now we're responding after 5 seconds of waiting. "
//
//				_, err = bot.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
//					Content: &content,
//				})
//				if err != nil {
//					bot.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
//						Content: "Something went wrong 02",
//					})
//					return
//				}
//			})
//			return
//		}
//
//		err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
//			Data: &discordgo.InteractionResponseData{
//				Content: content,
//			},
//		})
//		if err != nil {
//			bot.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
//				Content: "Something went wrong",
//			})
//			return
//		}
//		time.AfterFunc(time.Second*5, func() {
//			content := content + "\n\nWell, now you know how to create and edit responses. " +
//				"But you still don't know how to delete them... so... wait 10 seconds and this " +
//				"message will be deleted."
//			_, err = bot.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
//				Content: &content,
//			})
//			if err != nil {
//				bot.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
//					Content: "Something went wrong",
//				})
//				return
//			}
//			time.Sleep(time.Second * 10)
//			bot.InteractionResponseDelete(i.Interaction)
//		})
//	},
//	"followups": func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
//		// Followup messages are basically regular messages (you can create as many of them as you wish)
//		// but work as they are created by webhooks and their functionality
//		// is for handling additional messages after sending a response.
//
//		bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				// Note: this isn't documented, but you can use that if you want to.
//				// This flag just allows you to create messages visible only for the caller of the command
//				// (user who triggered the command)
//				Flags:   discordgo.MessageFlagsEphemeral,
//				Content: "Surprise!",
//			},
//		})
//		msg, err := bot.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
//			Content: "Followup message has been created, after 5 seconds it will be edited",
//		})
//		if err != nil {
//			bot.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
//				Content: "Something went wrong",
//			})
//			return
//		}
//		time.Sleep(time.Second * 5)
//
//		content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
//		bot.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
//			Content: &content,
//		})
//
//		time.Sleep(time.Second * 10)
//
//		bot.FollowupMessageDelete(i.Interaction, msg.ID)
//
//		bot.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
//			Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
//				"take a unicorn :unicorn: as reward!\n" +
//				"Also, as bonus... look at the original interaction response :D",
//		})
//	},
//}
