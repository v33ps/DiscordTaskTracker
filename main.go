package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
// channels used to send and get tasks and reminders
var (
	Token string
	wg    sync.WaitGroup
	c     chan string
	u     chan string
)

/**
* @brief initalize the command line flag parser
 */
func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	c = make(chan string, 1)
	u = make(chan string, 1)

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
	// wait for output from the channel and send it to the user
	for {
		remindedUser := <-u
		remindedTask := <-c
		dg.ChannelMessageSend(remindedUser, "Hi! Have you completed \n"+remindedTask+"\nyet?")
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
/**
* @brief gets all incoming messages to the bot
*
* Gets the messages sent to the bot, and determines if we should start a reminder. This is
* determined by a check to see if we have a reminder time and task. It'll send the message
* to be processed and then send a confirmation message back to the user
* This message is called due to the AddHandler in @main()
*
* @param s		discordgo session variable
* @param m		discordgo message handler
* @return		none
 */
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// if we have a reminder time and a task, process the message and send a confirmation back
	if strings.Contains(m.Content, "!r ") && strings.Contains(m.Content, "!t ") {
		message := processMessage(m.Content, m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, message)
	}

	// if the message doesn't conform to the proper format, send a help message to the user
	if strings.Contains(m.Content, "!h") {
		helpMsg := "Help:\n!r 3 hour/s/minute/s/second/s"
		s.ChannelMessageSend(m.ChannelID, helpMsg)
	}
}

/**
* @brief process the message from the user
*
* Determine how long to wait for the task and form the confirmation message.
* If the time format is bad, send an error message back
* Also starts the goroutine for the reminder.
*
* @param content	the message from the user
* @param username	the discord user id of the user who sent the message
* @return string	the confirmation message formatted, or an error message
 */
func processMessage(content string, username string) string {
	fmt.Println("[-] in processMessage")
	var message string
	var task string
	var reminder string
	timeTypes := []string{"hour", "hours", "minute", "minutes", "second", "seconds"}

	// gets the amount of time to remind on
	timeType := findStringInList(strings.Split(content, "!r ")[1], timeTypes)

	if len(timeType) > 0 {
		reminder = strings.Split(strings.Split(content, "!r ")[1], " !t")[0]
		message = "Reminder set for: " + reminder + "\n"
	} else {
		message = "[!] Incorrect reminder format. Enter !h for help"
		return message
	}
	task = strings.Split(content, "!t ")[1]
	message = message + "Task is: " + task

	// add to the waitgroup so we don't end before we have everything wrapped up
	wg.Add(1)
	// start the goroutine that acts as the reminder countdown for the task
	go sendReminder(task, reminder, c, username, u)

	return message
}

/**
* @brief wait the task reminder time
*
* Process the reminder time and sleep for it's duration. Once it's done sleeping,
* put the task into a channel, and the discord user id into another channel
*
* @param task		the task to be reminded on
* @param reminder	the amount of time to wait on the reminder
* @param c			the channel to send the task over
* @param username	the user to be reminded
* @param u			the channel to send the user id over
* @return			none
 */
func sendReminder(task string, reminder string, c chan string, username string, u chan string) {
	// this has the wait group finish once the function returns, removing it from the
	// list of waitgroup members
	defer wg.Done()
	// get the amount of time to sleep for in a proper format (int)
	duration, err := strconv.Atoi(strings.Split(reminder, " ")[0])
	if err != nil {
		log.Fatal(err)
	}

	format := strings.Split(reminder, " ")[1]
	if format == "second" || format == "seconds" {
		time.Sleep(time.Duration(duration) * time.Second)
	} else if format == "minute" || format == "minutes" {
		time.Sleep(time.Duration(duration) * time.Minute)
	} else if format == "hour" || format == "hours" {
		time.Sleep(time.Duration(duration) * time.Hour)
	}

	// send the task and user id over channels to be read in @main()
	c <- task
	u <- username
}

/**
* @brief find a string in a list
*
* @param a			the string to find in the list
* @param list		the list of strings to search through for @a
* @return string	the string that was found, or empty string if not found
 */
func findStringInList(a string, list []string) string {
	for _, b := range list {
		if strings.Contains(a, b) {
			return b
		}
	}
	return ""
}
