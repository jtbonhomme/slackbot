package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"strings"

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
	ws, id := slackConnect(botKey.Token)
	var msg Message
	msg.Text = "Je suis connecté"
	msg.Channel = "C8NTUSS9Z"
	msg.Id = 0
	postMessage(ws, msg)

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// m.Text
			quote := m.Text
			infof("Incoming message %s for me", quote)
			postMessage(ws, m)
			go func(m Message) {
				// fetch answer to the incoming quote
				dialogFlowResponse := GetResponse(quote, aiKey.Token)
				m.Text = dialogFlowResponse.Fulfillment.Speech
				infof("Message Id:%d, Type: %s, Channel: %s, Text: %s\n", m.Id, m.Type, m.Channel, m.Text)
				postMessage(ws, m)
			}(m)
		}
	}
}
