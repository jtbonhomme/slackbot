# slackbot

## Installation

* Install dep manager: https://github.com/golang/dep
```
$ dep ensure
```

* Create a new bot user integration on your Slack
* Create a file token.json which follows the format of the token_sample.json file provided with the Slack Bot Token
* Then run the follownig command

```
$ go run main.go
```

## Debug

```
$ export SLACKBOT_DEBUG=true
$ go run main.go
```

## Sources

* https://github.com/pricelinelabs/leaderboard

## Dependencies

* https://github.com/nlopes/slack
* https://github.com/kelseyhightower/envconfig
