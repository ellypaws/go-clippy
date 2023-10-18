package discord

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "go-clippy/database"
	"log"
	"os"
	"os/signal"
	"strings"
)

type BOT struct {
	guildID            *string
	token              *string
	removeCommands     *bool
	cleanCommands      *bool
	printFlags         *bool
	registeredCommands map[string]*discordgo.ApplicationCommand
	session            *discordgo.Session
}

var bot *BOT

func GetBot() *BOT {
	if bot != nil {
		return bot
	}
	token := flag.String("token", "", "Bot access token")
	if token == nil || *token == "" {
		tokenEnv := os.Getenv("BOT_TOKEN")
		if tokenEnv != "" {
			if tokenEnv == "YOUR_BOT_TOKEN_HERE" {
				log.Fatalf("Invalid bot token: %v\n"+
					"Did you edit the .env or run the program with -token ?", tokenEnv)
			}
			token = &tokenEnv
		}
	}
	session, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	return &BOT{
		guildID:            flag.String("guild", "", "Guild ID. If not passed - bot registers commands globally"),
		token:              token,
		removeCommands:     flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not"),
		cleanCommands:      flag.Bool("clcmd", false, "Remove all commands before starting"),
		printFlags:         flag.Bool("printflags", false, "Print all flags"),
		registeredCommands: make(map[string]*discordgo.ApplicationCommand),
		session:            session,
	}
}

func init() { flag.Parse() }

func init() {
	_ = godotenv.Load()
	bot = GetBot()

	if *bot.printFlags {
		log.Printf("Guild ID: %v", *bot.guildID)
		log.Printf("Bot Token: %v", *bot.token)
		log.Printf("Remove Commands: %v", *bot.removeCommands)
		log.Printf("Clean Commands: %v", *bot.cleanCommands)
	}
	log.Println("Bot initialized:", bot)
}

func (BOT) Run() {
	registerHandlers(bot.session)
	err := bot.session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// remove all commands at startup when -clcmd flag is passed
	if *bot.cleanCommands {
		log.Println("Removing all commands...")

		commandsToRemove, _ := bot.session.ApplicationCommands(bot.session.State.User.ID, *bot.guildID)
		if len(commandsToRemove) == 0 {
			log.Println("No commands to remove")
		}
		for _, command := range commandsToRemove {
			log.Printf("Attempting to remove '%v' command...", command.Name)
			err := bot.session.ApplicationCommandDelete(bot.session.State.User.ID, *bot.guildID, command.ID)
			if err != nil {
				log.Println("Cannot delete '%v' command: %v", command.Name, err)
				continue
			}
			log.Print("... success! ", command.ID)
		}
	}

	log.Println("Adding commands...")
	registerCommands(bot.session)

	defer func(session *discordgo.Session) {
		err := session.Close()
		if err != nil {
			log.Fatalf("Cannot close the session: %v", err)
		}
	}(bot.session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *bot.removeCommands {
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, command := range bot.registeredCommands {
			log.Println("Removing commands:", command.Name)
			err := bot.session.ApplicationCommandDelete(bot.session.State.User.ID, *bot.guildID, command.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", command.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}

func registerHandlers(session *discordgo.Session) {
	session.AddHandler(func(bot *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// commands
		case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(bot, i)
			}
		//buttons
		case discordgo.InteractionMessageComponent:
			log.Printf("Component with customID `%v` was pressed, attempting to respond\n", i.MessageComponentData().CustomID)
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				log.Printf(
					"Handler found, executing on message `%v`\nRan by: <@%v>\nUsername: %v",
					i.Message.ID,
					i.Member.User.ID,
					i.Member.User.Username,
				)
				log.Printf("https://discord.com/channels/%v/%v/%v", i.GuildID, i.ChannelID, i.Message.ID)
				h(bot, i)
			}
		}
	})
	session.AddHandler(func(bot *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", bot.State.User.Username, bot.State.User.Discriminator)
	})
}

func registerCommands(session *discordgo.Session) {
	bot.registeredCommands = make(map[string]*discordgo.ApplicationCommand, len(commands))
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
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, *bot.guildID, command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", command.Name, err)
		}
		bot.registeredCommands[key] = cmd
		log.Println("Registered command:", cmd.Name, cmd.ID)
	}
}
