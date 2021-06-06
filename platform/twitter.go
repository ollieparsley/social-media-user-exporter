package platform

import (
	"errors"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Twitter platform
type Twitter struct {
	metricPrefix   string
	client         *twitter.Client
	logger         *log.Logger
	factory        promauto.Factory
	screenName     string
	likesGauge     prometheus.Gauge
	statusesGauge  prometheus.Gauge
	followersGauge prometheus.Gauge
	friendsGauge   prometheus.Gauge
}

// NewTwitter create new instance of the Twitter user fetcher
func NewTwitter(logger *log.Logger, metricPrefix string, factory *promauto.Factory, screenName string, clientID string, clientSecret string, accessToken string, accessTokenSecret string) (Twitter, error) {

	// Check params
	if screenName == "" {
		return Twitter{}, errors.New("a screen name is required")
	}
	if clientID == "" {
		return Twitter{}, errors.New("a client ID is required")
	}
	if clientSecret == "" {
		return Twitter{}, errors.New("a client secret is required")
	}
	if accessToken == "" {
		return Twitter{}, errors.New("an access token is required")
	}
	if accessTokenSecret == "" {
		return Twitter{}, errors.New("an access token secret is required")
	}

	// oauth1 configures a client that uses app credentials to keep a fresh token
	twitterOAuthConfig := oauth1.NewConfig(clientID, clientSecret)
	twitterToken := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := twitterOAuthConfig.Client(oauth1.NoContext, twitterToken)

	// Verify credentials
	client := twitter.NewClient(httpClient)

	// Fetch rate limits
	_, _, err := client.Users.Show(&twitter.UserShowParams{ScreenName: screenName})
	if err != nil {
		return Twitter{}, fmt.Errorf("Error when checking that the screen name is valid %s: %s", screenName, err.Error())
	}

	return Twitter{
		metricPrefix: metricPrefix,
		logger:       logger,
		client:       client,
		screenName:   screenName,
		likesGauge: factory.NewGauge(prometheus.GaugeOpts{
			Name: metricPrefix + "twitter_likes",
			Help: "The number of likes the user has",
			ConstLabels: prometheus.Labels{
				"screen_name": screenName,
			},
		}),
		statusesGauge: factory.NewGauge(prometheus.GaugeOpts{
			Name: metricPrefix + "twitter_statuses",
			Help: "The number of statuses the user has",
			ConstLabels: prometheus.Labels{
				"screen_name": screenName,
			},
		}),
		followersGauge: factory.NewGauge(prometheus.GaugeOpts{
			Name: metricPrefix + "twitter_followers",
			Help: "The number of followers the user has",
			ConstLabels: prometheus.Labels{
				"screen_name": screenName,
			},
		}),
		friendsGauge: factory.NewGauge(prometheus.GaugeOpts{
			Name: metricPrefix + "twitter_friends",
			Help: "The number of friends the user has",
			ConstLabels: prometheus.Labels{
				"screen_name": screenName,
			},
		}),
	}, nil
}

// Fetch Format key with lavels
func (p Twitter) Fetch() error {

	// Fetch rate limits
	twitterUser, _, err := p.client.Users.Show(&twitter.UserShowParams{ScreenName: p.screenName})
	if err != nil {
		return fmt.Errorf("Error when requesting user details for %s: %s", p.screenName, err.Error())
	}

	p.likesGauge.Set(float64(twitterUser.FavouritesCount))
	p.statusesGauge.Set(float64(twitterUser.StatusesCount))
	p.followersGauge.Set(float64(twitterUser.FollowersCount))
	p.friendsGauge.Set(float64(twitterUser.FriendsCount))

	return nil
}
