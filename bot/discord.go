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
	printFlags     = flag.Bool("printflags", false, "Print all flags")
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

	if *printFlags {
		log.Printf("Guild ID: %v", *guildID)
		log.Printf("Bot Token: %v", *botToken)
		log.Printf("Remove Commands: %v", *removeCommands)
		log.Printf("Clean Commands: %v", *cleanCommands)
	}
}

func init() {
	var err error
	log.Println("Initializing bot...")
	bot, err = GetBot()
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	log.Println("Bot initialized:", bot)
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
		log.Println("Removing all commands...")

		commandsToRemove, _ := bot.ApplicationCommands(bot.State.User.ID, *guildID)
		if len(commandsToRemove) == 0 {
			log.Println("No commands to remove")
		}
		for _, command := range commandsToRemove {
			log.Printf("Attempting to remove '%v' command...", command.Name)
			err := bot.ApplicationCommandDelete(bot.State.User.ID, *guildID, command.ID)
			if err != nil {
				log.Println("Cannot delete '%v' command: %v", command.Name, err)
				continue
			}
			log.Print("... success! ", command.ID)
		}
	}

	log.Println("Adding commands...")
	registerCommands(bot)

	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *removeCommands {
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, command := range registeredCommands {
			log.Println("Removing commands:", command.Name)
			err := bot.ApplicationCommandDelete(bot.State.User.ID, *guildID, command.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", command.Name, err)
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
			log.Printf("Button with customID `%v` was pressed, attempting to respond\n", i.MessageComponentData().CustomID)
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				log.Println("Handler found, executing")
				h(bot, i)
			}
		}
	})
	bot.AddHandler(func(bot *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", bot.State.User.Username, bot.State.User.Discriminator)
	})
}

func registerCommands(bot *discordgo.Session) {
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
		log.Println("Registered command:", cmd.Name, cmd.ID)
	}
}
