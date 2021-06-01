package main

import (
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/docopt/docopt.go"
)

func main() {
	usage := `Twitter clean

Usage:
  trm tweets (--consumer-key=<ck> | -k <ck>) (--consumer-secret=<cs> | -s <cs>) (--token-key=<tk> | -y <tk>) (--token-secret=<ts> | -t <ts>) <user-name>
  trm likes (--consumer-key=<ck> | -k <ck>) (--consumer-secret=<cs> | -s <cs>) (--token-key=<tk> | -y <tk>) (--token-secret=<ts> | -t <ts>) <user-name>
  trm all (--consumer-key=<ck> | -k <ck>) (--consumer-secret=<cs> | -s <cs>) (--token-key=<tk> | -y <tk>) (--token-secret=<ts> | -t <ts>) <user-name>
  trm -h | --help

Options:
  --consumer-key=<ck>,-k <ck>     Twitter provided consumer key.
  --consumer-secret=<cs>,-s <cs>  Twitter provided consumer Secret.
  --token-key=<tk>,-y <tk>        Twitter provided token key.
  --token-secret=<ts>,-t <ts>     Twitter provided token secret.
  --help -h                       Print help.`

	doc, err := docopt.ParseDoc(usage)
	if err != nil {
		parsingError(err)
	}
	consumerKey, err := doc.String("--consumer-key")
	if err != nil {
		parsingError(err)
	}
	consumerSecret, err := doc.String("--consumer-secret")
	if err != nil {
		parsingError(err)
	}
	tokenKey, err := doc.String("--token-key")
	if err != nil {
		parsingError(err)
	}
	tokenSecret, err := doc.String("--token-secret")
	if err != nil {
		parsingError(err)
	}
	screenName, err := doc.String("<user-name>")
	if err != nil {
		parsingError(err)
	}
	tweets, err := doc.Bool("tweets")
	if err != nil {
		parsingError(err)
	}
	likes, err := doc.Bool("likes")
	if err != nil {
		parsingError(err)
	}
	all, err := doc.Bool("all")
	if err != nil {
		parsingError(err)
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(tokenKey, tokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: screenName,
	})
	if err != nil {
		panic(fmt.Errorf("retriving user: %s", err.Error()))
	}

	if tweets {
		deleteTweets(client, user)
	}
	if likes {
		deleteLikes(client, user)
	}
	if all {
		deleteTweets(client, user)
		deleteLikes(client, user)
	}
}

func deleteTweets(client *twitter.Client, user *twitter.User) {
	_, err := os.Stdout.WriteString(fmt.Sprintln(""))
	if err != nil {
		panic(err)
	}

	count := 0
	for ok := true; ok; {
		tweets, _, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			UserID:     user.ID,
			ScreenName: user.ScreenName,
			Count:      100,
		})

		for _, tweet := range tweets {
			_, _, err := client.Statuses.Destroy(tweet.ID, nil)
			if err != nil {
				_, err = os.Stderr.WriteString(fmt.Sprintf("\nFailed to delete %v. Error: %v", tweet.ID, err))
				if err != nil {
					panic(err)
				}
			}
		}

		count += len(tweets)
		ok = len(tweets) > 0

		_, err := os.Stdout.WriteString(fmt.Sprintf("Deleted %d tweets\r", count))
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("\nTotal deleted %d tweets", count)
}

func deleteLikes(client *twitter.Client, user *twitter.User) {
	_, err := os.Stdout.WriteString(fmt.Sprintln(""))
	if err != nil {
		panic(err)
	}

	count := 0
	for ok := true; ok; {
		likes, _, err := client.Favorites.List(&twitter.FavoriteListParams{
			UserID:     user.ID,
			ScreenName: user.ScreenName,
			Count:      1000,
		})
		if err != nil {
			panic(err)
		}

		for _, like := range likes {
			_, _, err := client.Favorites.Destroy(&twitter.FavoriteDestroyParams{
				ID: like.ID,
			})
			if err != nil {
				_, err = os.Stderr.WriteString(fmt.Sprintf("\nFailed to delete %v. Error: %v", like.ID, err))
				if err != nil {
					panic(err)
				}
			}
		}

		count = count + len(likes)
		ok = len(likes) > 0

		_, err = os.Stdout.WriteString(fmt.Sprintf("Deleted %d likes\r", count))
		if err != nil {
			panic(err)
		}
	}

	_, err = os.Stdout.WriteString(fmt.Sprintf("\nTotal deleted %d likes", count))
	if err != nil {
		panic(err)
	}
}

func parsingError(err error) {
	_, err = os.Stderr.WriteString(err.Error())
	if err != nil {
		panic(err)
	}
	os.Exit(2)
}
