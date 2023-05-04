package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func main() {
	slackKey := flag.String("slackKey", "", "Slack API Key")
	slackChannel := flag.String("slackChannel", "bot-dev", "Slack channel")
	flag.Parse()

	resp, err := http.Get("https://www.merriam-webster.com/word-of-the-day")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if title, ok := getHTMLTitle(resp.Body); ok {
		tail := strings.SplitAfter(title, ":")
		head := strings.Split(tail[1], "|")
		word := strings.TrimSpace(head[0])

		message := fmt.Sprintf("Today's word of the day, \"*%s*\"", word)

		//message := fmt.Sprintf("I can %s by six words", word)
		println(message)

		if *slackKey == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}

		if *slackChannel == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}

		postMessage(*slackKey, *slackChannel, message)

	} else {
		println("Fail to get HTML title")
	}

}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func getHTMLTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		panic("Fail to parse html")
	}

	return traverse(doc)
}

func postMessage(apiKey string, channel string, message string) []byte {
	slackMessage := &SlackMessage{
		channel,
		message,
	}
	json, err := json.Marshal(slackMessage)
	//fmt.Printf("POST: %s\n\n", json)

	url := "https://slack.com/api/chat.postMessage"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-type", "application/json; charset=utf-8")
	//req.Header.Add("Content-Type", "text/html; charset=utf-8")
	response, err := client.Do(req)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		} else {
			fmt.Println(string(contents))
		}

		return contents
	}

	return nil
}

// func postMessage(apiKey string, channel string, message string) []byte {
// 	message = url.QueryEscape(message)
// 	channel = url.QueryEscape(channel)

// 	url := "https://slack.com/api/chat.postMessage?token=" + apiKey + "&channel=" + channel + "&text=" + message

// 	response, err := http.Get(url)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	} else {
// 		defer response.Body.Close()
// 		contents, err := ioutil.ReadAll(response.Body)
// 		if err != nil {
// 			fmt.Printf("%s", err)
// 			os.Exit(1)
// 		}

// 		return contents
// 	}

// 	return nil
// }
