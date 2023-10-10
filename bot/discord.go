package discord

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
)

var (
	guildID        = flag.String("guild", "", "Guild ID. If not passed - bot registers commands globally")
	botToken       = flag.String("token", "", "Bot access token")
	removeCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
	cleanCommands  = flag.Bool("clcmd", false, "Remove all commands before starting")
)

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
)

var (
	bot                *discordgo.Session
	registeredCommands map[string]*discordgo.ApplicationCommand
)

func init() { flag.Parse() }

func init() {
	if botToken == nil || *botToken == "" {
		tokenEnv := os.Getenv("token")
		if tokenEnv != "" {
			botToken = &tokenEnv
		}
	}
}

func init() {
	var err error
	bot, err = GetBot()
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func GetBot() (*discordgo.Session, error) {
	if bot != nil {
		return bot, nil
	}
	return discordgo.New("Bot " + *botToken)
}

func Run() {
	registerHandlers(bot)
	err := bot.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// remove all commands at startup when -clcmd flag is passed
	if *cleanCommands {
		commandsToRemove, _ := bot.ApplicationCommands(*botToken, *guildID)
		for _, command := range commandsToRemove {
			err := bot.ApplicationCommandDelete(*botToken, *guildID, command.ID)
			if err != nil {
				log.Fatalf("Cannot delete '%v' command: %v", command.Name, err)
			}
		}
	}

	registerCommands(bot)

	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *removeCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := bot.ApplicationCommandDelete(bot.State.User.ID, *guildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}

func registerHandlers(bot *discordgo.Session) {
	bot.AddHandler(func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// commands
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(bot, i)
			}
		//buttons
		case discordgo.InteractionMessageComponent:
			if h, ok := interactionHandlers[i.MessageComponentData().CustomID]; ok {
				h(bot, i)
			}
		}
	})
	bot.AddHandler(func(bot *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", bot.State.User.Username, bot.State.User.Discriminator)
	})
}

func registerCommands(bot *discordgo.Session) {
	log.Println("Adding commands...")
	registeredCommands = make(map[string]*discordgo.ApplicationCommand, len(commands))
	for key, command := range commands {
		if command.Name == "" {
			// clean the key because it might be a description of some sort
			// only get the first word, and clean to only alphanumeric characters or -
			sanitized := strings.ReplaceAll(key, " ", "-")
			sanitized = strings.ToLower(sanitized)

			// remove all non valid characters
			for _, c := range sanitized {
				if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' {
					sanitized = strings.ReplaceAll(sanitized, string(c), "")
				}
			}
			command.Name = sanitized
		}
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, *guildID, command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", command.Name, err)
		}
		registeredCommands[key] = cmd
	}

}
