package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/google/go-github/github"


	//"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type GithubMessage struct {
	message string
}

var (
	ch chan string
)

func init() {
 ch = make(chan string, 10)
	
}


func main() {
	ch := make(chan string, 10)

	go startServer(ch)

	go makeBot(ch)


	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	
	

}



func handleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte("your-secret-key"))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		// fmt.Println(e.GetPusher().GetName())
		// fmt.Println(e.GetSender().GetName())
		// fmt.Println(e.GetRepo().GetFullName())
		// fmt.Println(e.GetRef())
		// for _, commit := range e.Commits {
		// 	fmt.Println(commit.GetMessage())
		// 	for _, modified := range commit.Modified {
		// 		fmt.Println(modified)
		// 	}
		// } 

		ch <- e.GetPusher().GetName()

		



		// this is a commit push, do something with it
	case *github.PullRequestEvent:
		// this is a pull request, do something with it
	case *github.WatchEvent:
		// https://developer.github.com/v3/activity/events/types/#watchevent
		// someone starred our repository
		if e.Action != nil && *e.Action == "starred" {
			fmt.Printf("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

func startServer(channel chan string)  {
	log.Println("server started")
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8081", nil))

}


func makeBot(channel chan string)  {
	Token := "sample-token"
	// Create a new Discord session using the provided bot token.
	fmt.Println("Hello Discord")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.ChannelMessageSend("your-channel-id", "I don show!")
	dg.ChannelMessageSend("your-channel-id", <-ch)

	fmt.Println(<-ch)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	

	// Cleanly close down the Discord session.
	dg.Close()
}



// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

