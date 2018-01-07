# slackbot

## Installation

* Install dep manager: https://github.com/golang/dep
```
$ dep ensure
```

* Create a new bot user integration on your Slack (https://my.slack.com/services/new/bot) and note the token API (something like xoxb-XXXXXX ) that will be generated
* Create a new DialogFlow agent (https://console.dialogflow.com/api-client/) and note it's Client access token, select Small Talk pre computed model to start quickly, and enter some answers to classical quotes
* Create a file token.json which follows the format of the token_sample.json file provided with the Slack bot token and the DialogFlow Client access token
* Then run the follownig command

```
$ go build
$ ./slackbot
```
Type CRTL+C to stop the bot.

## Debug

```
$ export SLACKBOT_DEBUG=true
$ ./slackbot
```

## Sources

* https://www.opsdash.com/blog/slack-bot-in-golang.html

## Dependencies

* https://github.com/kelseyhightower/envconfig
* https://github.com/mlabouardy/dialogflow-go-client/