package main

import (
	"flag"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	consumerKey := flag.String("consumerKey", "", "Twitter provided consumer key")
	consumerSecret := flag.String("consumerSecret", "", "Twitter provided consumer Secret")
	tokenKey := flag.String("tokenKey", "", "Twitter provided token key")
	tokenSecret := flag.String("tokenSecret", "", "Twitter provided token secret")
	screenName := flag.String("screenName", "", "Twitter provided token secret")

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*tokenKey, *tokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	user, _, _ := client.Users.Show(&twitter.UserShowParams{
		ScreenName: *screenName,
	})
	count := 0

	for ok := true; ok; {
		tweets, _, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			UserID:     user.ID,
			ScreenName: user.ScreenName,
			Count:      100,
		})

		for _, tweet := range tweets {
			fmt.Printf("Twitt of %v. ID: %v.\n", tweet.User.Name, tweet.ID)
			_, _, err := client.Statuses.Destroy(tweet.ID, nil)

			if err != nil {
				fmt.Printf("Failed to delete %v. Error: %v \n", tweet.ID, err)
			}
		}

		fmt.Printf("Deleted %v twits", len(tweets))

		count += len(tweets)
		ok = len(tweets) > 0
	}

	fmt.Printf("Total deleted %v twits", count)
}
