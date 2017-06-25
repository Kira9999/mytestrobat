package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)

}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(stackoverflow(message.Text))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

//Items:
type jsonobject struct {
	Items []Item
}

//Item
type Item struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

func stackoverflow(input string) string {

	root := "http://api.stackexchange.com/2.2/similar"
	para := "?page=1&pagesize=1&order=desc&sort=relevance&site=stackoverflow&title=" + url.QueryEscape(input)

	stackoverflowEndPoint := root + para

	resp, err := http.Get(stackoverflowEndPoint)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var i jsonobject
	err = json.Unmarshal(body, &i)
	if err != nil {
		log.Println(err)
	}

	var ret string

	if len(i.Items) == 0 {
		ret = "Sorry, I can't find relevant solutions, please specify your question."
	} else {
		ret = html.UnescapeString(i.Items[0].Title) + " " + i.Items[0].Link
	}

	if len(ret) == 0 {
		ret = "Sorry, I can't find relevant solutions, please specify your question."
	}

	if strings.ToLower(input) == "hello" {
		ret = input + " +1"
	}

	return ret
}
