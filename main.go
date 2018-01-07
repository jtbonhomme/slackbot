package main

import (
	//    "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	//    "net/http"
	//    "os"
	"path"
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type Token struct {
	Token string `json:"slack-token"`
}

type Specification struct {
	Debug bool `envconfig:"SLACKBOT_DEBUG" default:"false"`
}

type User struct {
	Info   slack.User
	Rating int
}

type ActiveUsers []User

var (
	api    *slack.Client
	botKey Token

	activeUsers ActiveUsers
	//    userMessages Messages
	botId string

//    botCommandChannel chan *BotCentral
//    botReplyChannel chan AttachmentChannel
)

func init() {
	file, err := ioutil.ReadFile("./token.json")

	if err != nil {
		log.Fatal("File doesn't exist")
	}

	if err := json.Unmarshal(file, &botKey); err != nil {
		log.Fatal("Cannot parse token.json")
	}
}

func main() {
	var (
		s     Specification
		infof = func(format string, a ...interface{}) {
			msg := fmt.Sprintf(format, a...)
			function, fileName, fileLine, ok := runtime.Caller(1)
			if ok {
				_, file := path.Split(fileName)
				fmt.Printf("INFO::%s(%d)/%s %s\n", file, fileLine, runtime.FuncForPC(function).Name(), msg)
			} else {
				fmt.Printf("INFO::%s\n", msg)
			}
		}
	)

	err := envconfig.Process("slackbot", &s)
	if err != nil {
		log.Fatal("[ERROR] Failed to process env var: %s", err.Error())
	}
	infof("DEBUG: %t\n", s.Debug)

	infof("Create new slack connexion, token: %s\n", botKey.Token)
	api = slack.New(botKey.Token)
	infof("Configure logger and set debug option\n")
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	infof("Get PostMessageParameters\n")
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
		// Uncomment the following part to send a field too

		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "a",
				Value: "no",
			},
		},
	}
	infof("Send Message\n")
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage("C8NTUSS9Z", "I am connected !", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)

	infof("Get Channels\n")
	channels, err := api.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, channel := range channels {
		fmt.Printf("channel: %s (%s)\n", channel.Name, channel.ID)
		// channel is of type conversation & groupConversation
		// see all available methods in `conversation.go`
	}

	/*
		infof("Get Groups\n")
		groups, err := api.GetGroups(false)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		for _, group := range groups {
			fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
		}

	*/
	infof("Create new Real Time Messaging connexion\n")
	rtm := api.NewRTM()
	infof("Manage connexion in a go routine\n")
	go rtm.ManageConnection()
Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello
				fmt.Printf("Hello !\n")

			case *slack.ConnectedEvent:
				botId = ev.Info.User.ID
				fmt.Printf("bot: %s (%s)\n", ev.Info.User.Name, botId)

				for _, u := range ev.Info.Users {
					if u.RealName != "" {
						user := User{
							Info: u,
						}
						activeUsers = append(activeUsers, user)
						fmt.Printf("User: %s (%s)\n", user.Info.Name, user.Info.ID)
						//rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C8NAP5WEL"))
					}
				}

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)
				fmt.Printf("prefix: %s\n", prefix)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					rtm.SendMessage(rtm.NewOutgoingMessage("What's up buddy!?!?", ev.Channel))
				}

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// Ignore other events..
				// fmt.Printf("Event: %v\nData: %v\n", ev, msg.Data)
			}
		}
	}
}
