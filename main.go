package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

var msgPlaceholder = []string{"IT'S BURNNNTTTT", "YOU'RE AN IDIOT SANDWICH", "THIS IS RAWWWW"}

func getToken() string {
	return "xoxb-2279330878-710758657270-bAamNY5INpDxicxD6VfYJdeC"
}

func main() {
	// Set-up slack listener
	token := getToken()
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for msg := range rtm.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)

		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// userID := rtm.GetInfo().User.ID
			// botTagString := fmt.Sprintf("<@%s>", userID)
			// if !strings.Contains(ev.Msg.Text, botTagString) {
			// 	continue
			// }

			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("%s", roast(ev.User)), ev.Channel))

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			break Loop

		default:
			// Do Nothing
		}
	}
}
