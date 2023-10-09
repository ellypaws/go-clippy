package discord

import "github.com/bwmarrin/discordgo"

func deleteButton() discordgo.MessageComponent {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Delete this message",
				Style:    discordgo.DangerButton,
				CustomID: "delete_error_message",
			},
		},
	}
}
