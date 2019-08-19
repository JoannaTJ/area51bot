package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	// NOTE: godotenv, might screw deployment?
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

const (
	// Meme API
	memeAPIEndpoint = "https://meme-api.herokuapp.com/gimme"

	// Bot text
	helpText = "Commands:\n/meme\n/dadjoke\n/roast\n"
)

var (
	slackClient *slack.Client
)

// MemeRes : JSON response received from meme endpoint
type MemeRes struct {
	PostLink  string
	Subreddit string
	Title     string
	URL       string
}

func main() {
	// Set-up slack listener
	err := godotenv.Load("environment.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slackClient = slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))

	// Defines "<host>/receive" as the API endpoint
	http.HandleFunc("/receive", slashCommandHandler)

	log.Println("[INFO] Server listening on port", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprint(":", os.Getenv("PORT")), nil)
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VARIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Switch statement to handle custom commands
	switch s.Command {
	case "/roast":
		// Match @username handle
		// TODO: Consiuder nivubg this to a function
		matched, err := regexp.MatchString("^@([a-z0-9][a-z0-9._-]+)$", s.Text)
		if !matched || err != nil {
			w.Write([]byte("Invalid Command"))
			return
		}
		// s.Text: should be the target
		roastText := fmt.Sprintf("<@%s> said to <%s>: %s", s.UserID, s.Text, getRoastString())
		slackClient.PostMessage(s.ChannelID, slack.MsgOptionText(roastText, false))
		w.WriteHeader(http.StatusOK)

	case "/meme":
		go func() {
			attachment, err := getMeme()
			if err != nil {
				log.Println("[ERROR] getMeme")
				slackClient.PostMessage(s.ChannelID, slack.MsgOptionText("[ERROR] failed to get meme", false))
				return
			}
			slackClient.PostMessage(s.ChannelID, slack.MsgOptionText(fmt.Sprintf("<@%s> requested a meme", s.UserID), false), slack.MsgOptionAttachments(attachment))
		}()
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getRoastString() string {
	// Generate a roast string from file
	file, err := os.Open("./static/roast.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	roastText, err := ioutil.ReadAll(file)
	if err != nil {
		// STRETCH: Improve error handling
		log.Fatal(err)
	}
	roastSlice := strings.Split(string(roastText), "\n")
	return roastSlice[rand.Intn(len(roastSlice))]
}

func getMeme() (slack.Attachment, error) {
	var resBody MemeRes
	var attachment slack.Attachment

	res, err := http.Get(memeAPIEndpoint)
	if err != nil {
		fmt.Println("[ERROR] Meme API Endpoint")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), &resBody)

		// Debugging
		log.Println("[INFO] Response Body:", resBody)

		attachment = slack.Attachment{
			Title:     string(resBody.Title),
			TitleLink: string(resBody.PostLink),
			Text:      string(resBody.Subreddit),
			ImageURL:  string(resBody.URL),
		}

		return attachment, nil
	}

	return attachment, fmt.Errorf("[ERROR] StatusCode: %d", res.StatusCode)
}

func getComplimentString() string {
	// Generate a roast string from file
	file, err := os.Open("./static/compliments.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	complimentText, err := ioutil.ReadAll(file)
	if err != nil {
		// STRETCH: Improve error handling
		log.Fatal(err)
	}
	complimentSlice := strings.Split(string(complimentText), "\n")
	return complimentSlice[rand.Intn(len(complimentSlice))]
}
