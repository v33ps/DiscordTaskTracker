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
var (
	Token        string
	wg           sync.WaitGroup
	c            chan string
	channelCount int
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	channelCount = 0
	channelrecv := 0
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// wait for output from the channel, print it if we get any.
	for {
		fmt.Println("lol")
		select {
		case res := <-c:
			fmt.Println(res)
		}
		// if we get the same number of outputs as inputs, stop waiting
		if channelrecv == channelCount {
			fmt.Println("[+] we got em all")
			break
		}
	}
	wg.Wait()
	close(c)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// if we have a reminder time and a task, process the message and send a response back
	if strings.Contains(m.Content, "!r ") && strings.Contains(m.Content, "!t ") {
		message := processMessage(m.Content)
		s.ChannelMessageSend(m.ChannelID, message)
	}

	if strings.Contains(m.Content, "!h") {
		helpMsg := "Help:\n!r 3 hour/s/minute/s/second/s"
		s.ChannelMessageSend(m.ChannelID, helpMsg)
	}

}

// process the message from the user

// pull out the time to wait (@reminder)
// pull out the task (@task)
// figure out how long the time is (hours/minutes/seconds)
// if the reminder time is a bad format, return an error message to give the user
// with a good reminder time and task, start a goroutine sendReminder()

// @content string		the message from the user to process
// @return string		the message to give back to the user
func processMessage(content string) string {
	fmt.Println("[-] in processMessage")
	var message string
	var task string
	var reminder string
	timeTypes := []string{"hour", "hours", "minute", "minutes", "second", "seconds"}

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

	wg.Add(1)
	go sendReminder(task, reminder, c)
	channelCount++
	return message
}

// wait for the @reminder time

// processes the @reminder time to figure out how long we have to wait
// once determined, sleep till then
// once done sleeping, put the task back out the channel to be read off
// @task string		the task that is being slept
// @reminder string	the amount of time to sleep
// @c chan string	the channel to send tasks back along
func sendReminder(task string, reminder string, c chan string) {
	defer wg.Done()
	duration, err := strconv.Atoi(strings.Split(reminder, " ")[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Split(reminder, " "), task)

	format := strings.Split(reminder, " ")[1]
	if format == "second" || format == "seconds" {
		time.Sleep(time.Duration(duration) * time.Second)
	} else if format == "minute" || format == "minutes" {
		time.Sleep(time.Duration(duration) * time.Minute)
	} else if format == "hour" || format == "hours" {
		time.Sleep(time.Duration(duration) * time.Hour)
	}
	fmt.Println("finished")
	c <- task
}

// find a string in a list

// @a 				the string to look for in the list
// @list 			the list of strings to search through
// @return string	the string that was found, or nothing if not found
func findStringInList(a string, list []string) string {
	for _, b := range list {
		if strings.Contains(a, b) {
			return b
		}
	}
	return ""
}
