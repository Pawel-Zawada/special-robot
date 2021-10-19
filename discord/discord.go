package discord

import (
	"log"
	"nextrock/borat_bot/commands"
	"os"

	"github.com/bwmarrin/discordgo"
)

var Discord *discordgo.Session

// create a bot configuration that is ready to be connected
func initBot() {
	var err error
	Discord, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	Discord.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	err = Discord.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
}

// createCommands send command creation requests for each defined command in commands.Commands
func createCommands() {
	for _, v := range commands.Commands {
		go func(v commands.Command) {
			_, err := Discord.ApplicationCommandCreate(Discord.State.User.ID, "", v.ApplicationCommand)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.ApplicationCommand.Name, err)
			}
			log.Printf("Loaded command: %v", v.ApplicationCommand.Name)
		}(v)
	}
}

// handleCommands sets a handler for each created command that is defined in commands.Command
func handleCommands() {
	Discord.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		n := interaction.ApplicationCommandData().Name
		commands.Commands[n].Handler(session, interaction)
	})
}

// Run starts up the Discord bot connection and loads in the application commands
func Run() {
	initBot()

	go createCommands()
	go handleCommands()
}

// Stop gracefully shuts down the Discord connection
func Stop() {
	log.Println("Gracefully shutdowning Discord connection")
	err := Discord.Close()
	if err != nil {
		log.Panicf("Failed to gracefully disconnect from Discord, error: '%v'", err)
	}
}
