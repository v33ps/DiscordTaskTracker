package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func stringInSlice(a string, list []string) string {
	for _, b := range list {
		if strings.Contains(a, b) {
			return b
		}
	}
	return ""
}

func main() {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)
	timezoneList := []string{"pst", "est", "zst"}
	text := "8/3/2018 3pm pst"
	timezone := stringInSlice(text, timezoneList)
	fmt.Println(timezone)
	r, err := w.Parse(text, time.Now())
	if err != nil {
		// an error has occurred
	}
	if r == nil {
		// no matches found
	}

	fmt.Println(
		"the time",
		r.Time.String(),
		"mentioned in",
		text[r.Index:r.Index+len(r.Text)],
	)
	/*
		t, err := dateparse.ParseAny("8/1/2018 1700")
		fmt.Println(t)
		// Create a new Discord session using the provided bot token.
		dg, err := discordgo.New("Bot " + Token)
		if err != nil {
			fmt.Println("error creating Discord session,", err)
			return
		}

		// Register the messageCreate func as a callback for MessageCreate events.
		dg.AddHandler(messageCreate)

		// Open a websocket connection to Discord and begin listening.
		err = dg.Open()
		if err != nil {
			fmt.Println("error opening connection,", err)
			return
		}

		// Wait here until CTRL-C or other term signal is received.
		fmt.Println("Bot is now running.  Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc

		// Cleanly close down the Discord session.
		dg.Close()
	*/
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// extract the reminder time
	if strings.Contains(m.Content, "!r ") {
		message := strings.Split("Reminder set for: "+strings.Split(m.Content, "!r ")[1], " !t")[0]
		s.ChannelMessageSend(m.ChannelID, message)
	}

	// extract the task
	if strings.Contains(m.Content, "!t ") {
		s.ChannelMessageSend(m.ChannelID, "Task is: "+strings.Split(m.Content, "!t ")[1])
	}

}
