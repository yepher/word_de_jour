package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

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
		words := strings.SplitAfter(title, " ")
		word := strings.TrimSpace(words[4])

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
	message = url.QueryEscape(message)
	channel = url.QueryEscape(channel)

	url := "https://slack.com/api/chat.postMessage?token=" + apiKey + "&channel=" + channel + "&text=" + message

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		return contents
	}

	return nil
}