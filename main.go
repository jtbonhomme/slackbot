package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"runtime"
    "os"

    "github.com/davecgh/go-spew/spew"
    "github.com/nlopes/slack"
	"github.com/kelseyhightower/envconfig"
)

type SlackToken struct {
	Token string `json:"slack-token"`
}

type DialogFlowToken struct {
	Token string `json:"dialogflow-token"`
}

type Specification struct {
	Debug bool `envconfig:"SLACKBOT_DEBUG" default:"false"`
}

// These are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.

type Message struct {
    Id      uint64 `json:"id"`
    Type    string `json:"type"`
    Channel string `json:"channel"`
    Text    string `json:"text"`
}

var (
	botKey SlackToken
	aiKey  DialogFlowToken
)

func init() {
	file, err := ioutil.ReadFile("./token.json")

	if err != nil {
		log.Fatal("File doesn't exist")
	}

	if err := json.Unmarshal(file, &botKey); err != nil {
		log.Fatal("Cannot parse token.json")
	}

	if err := json.Unmarshal(file, &aiKey); err != nil {
		log.Fatal("Cannot parse token.json")
	}
}

func main() {
	var (
        count uint64
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
    count = 1
	err := envconfig.Process("slackbot", &s)
	if err != nil {
		log.Fatal("[ERROR] Failed to process env var: %s", err.Error())
	}
	infof("DEBUG: %t\n", s.Debug)


	infof("Create new slack connexion, token: %s\n", botKey.Token)
    logger := log.New(os.Stdout, "slackbot: ", log.Lshortfile|log.LstdFlags)
    api := slack.New(botKey.Token)
    slack.OptionLog(logger)
    slack.OptionDebug(false)

    rtm := api.NewRTM()
    go rtm.ManageConnection()

    for {
        select {
        case msg := <-rtm.IncomingEvents:
            switch ev := msg.Data.(type) {
            case *slack.ConnectedEvent:
                botId := ev.Info.User.ID
                infof("I am bot ID: %s", botId)
            case *slack.TeamJoinEvent:
                // Handle new user to client
            case *slack.ConnectingEvent:
                // Handle connecting event
            case *slack.MessageEvent:
                // Handle new message to channel
                infof("[%s] %s (%s:%s)\n", ev.Msg.Timestamp, ev.Msg.Text, ev.Msg.Channel, ev.Msg.User)
                quote := ev.Msg.Text
                infof("<@%s>", rtm.GetInfo().User.ID)

                m := Message {
                    Id: count,
                    Type: ev.Msg.Type,
                    Channel: ev.Msg.Channel,
                    Text: ev.Msg.Text,
                }
                count+=1
                go func(m Message) {
                    dialogFlowResponse := GetResponse(quote, aiKey.Token)
                    m.Text = dialogFlowResponse.Fulfillment.Speech
                    infof("Message Id: %s, Type: %s, Channel: %s, Text: %s\n", m.Id, m.Type, m.Channel, m.Text)
                    rtm.SendMessage(rtm.NewOutgoingMessage(m.Text, ev.Channel))
                }(m)

//                spew.Dump(msg)
            case *slack.AckMessage:
                infof("[%s] ACK: %s\n", ev.Timestamp, ev.Text)
            case *slack.ReactionAddedEvent:
                // Handle reaction added
            case *slack.ReactionRemovedEvent:
                // Handle reaction removed
            case *slack.UserTypingEvent:
                // Handle user typing
            case *slack.LatencyReport:
                // Handle latency report
            case *slack.HelloEvent:
                // Handle Hello
            case *slack.RTMError:
                infof("Error: %s\n", ev.Error())
            case *slack.InvalidAuthEvent:
                infof("Invalid credentials")
            default:
                infof("Unknown event\n")
                spew.Dump(msg)
            }
        }
    }
}
