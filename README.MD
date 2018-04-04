# Word Of The Day Retrieval

`wordDeJour` will retrieve the [word of the day](https://www.merriam-webster.com/word-of-the-day) and post it to a Slack channel.

This tool scrapes the `<title>...</title>` for the Word of the day web page.


## How to build


* `cd wordDeJour/src`
* `go get golang.org/x/net/html`
* ./buid.sh


## How to run

```
-slackChannel string
    	Slack channel (default "bot-dev")
  -slackKey string
    	Slack API Key

```

`wordDeJourLinux -slackChannel 6words -slackKey "YOUR_SLACK_API_KEY_HERE"`