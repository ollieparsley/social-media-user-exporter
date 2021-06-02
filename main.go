package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ollieparsley/social-media-user-exporter/platform"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var l *log.Logger = log.New(os.Stdout, "[social-media-user-exporter] ", 2)

// getEnv Get an env vairable and set a default
func getEnv(key, fallback string) string {
	fullKey := "SMUE_" + key
	if value, ok := os.LookupEnv(fullKey); ok {
		return value
	}
	return fallback
}

func main() {
	l.Println("Starting service")

	// Fetch env variables with defaults
	httpPort := getEnv("HTTP_PORT", "9100")
	httpPath := getEnv("HTTP_PATH", "metrics")
	metricPrefix := getEnv("METRICS_PREFIX", "social_media_user_")

	promRegistry := prometheus.NewRegistry()
	promFactory := promauto.With(promRegistry)

	twitterClientID := getEnv("TWITTER_CLIENT_ID", "")
	twitterClientSecret := getEnv("TWITTER_CLIENT_SECRET", "")
	twitterAccessToken := getEnv("TWITTER_ACCESS_TOKEN", "")
	twitterAccessTokenSecret := getEnv("TWITTER_ACCESS_TOKEN_SECRET", "")

	youtubeClientID := getEnv("YOUTUBE_CLIENT_ID", "")
	youtubeClientSecret := getEnv("YOUTUBE_CLIENT_SECRET", "")
	youtubeAccessToken := getEnv("YOUTUBE_ACCESS_TOKEN", "")
	youtubeRefreshToken := getEnv("YOUTUBE_REFRESH_TOKEN", "")

	// TODO: Facebook and Instagram
	// https://github.com/huandu/facebook

	intervalSeconds := getEnv("INTERVAL_SECONDS", "300")
	intervalSecondsInt, err := strconv.Atoi(intervalSeconds)
	if err != nil {
		l.Printf("Error converting interval: %s", err.Error())
		os.Exit(1)
	}

	// API request counter counter
	counter := promFactory.NewCounter(prometheus.CounterOpts{
		Name:        metricPrefix + "counter",
		Help:        "Increment each time we call the platforms to get an update",
		ConstLabels: prometheus.Labels{},
	})

	// Twitter user info by screen names
	twitterScreenNamesList := getEnv("TWITTER_SCREEN_NAMES", "")
	twitterScreenNames := strings.Split(twitterScreenNamesList, ",")
	twitterScreenNamesFetchers := []platform.Twitter{}
	for _, twitterScreenName := range twitterScreenNames {
		if twitterScreenName == "" {
			continue
		}
		twitterScreenNameFetcher, err := platform.NewTwitter(l, metricPrefix, &promFactory, twitterScreenName, twitterClientID, twitterClientSecret, twitterAccessToken, twitterAccessTokenSecret)
		if err != nil {
			l.Fatalf("Problem setting up twitter: %s", err.Error())
		}
		twitterScreenNamesFetchers = append(twitterScreenNamesFetchers, twitterScreenNameFetcher)
	}

	// Youtue user info by screen names
	youtubeChannelIDsList := getEnv("YOUTUBE_CHANNEL_IDS", "")
	youtubeChannelIDs := strings.Split(youtubeChannelIDsList, ",")
	youtubeChannelFetchers := []platform.YouTube{}
	for _, youtubeChannelID := range youtubeChannelIDs {
		if youtubeChannelID == "" {
			continue
		}
		youtubeChannelFetcher, err := platform.NewYouTube(l, metricPrefix, &promFactory, youtubeChannelID, youtubeClientID, youtubeClientSecret, youtubeAccessToken, youtubeRefreshToken)
		if err != nil {
			l.Fatalf("Problem setting up youtube: %s", err.Error())
		}
		youtubeChannelFetchers = append(youtubeChannelFetchers, youtubeChannelFetcher)
	}

	// Increment counter every second
	go func() {
		for {
			// Show that we're making a request
			counter.Inc()

			l.Println("Calling platforms to fetch latest data")

			// Twitter user info
			for _, twitterScreenNamesFetcher := range twitterScreenNamesFetchers {
				fetchErr := twitterScreenNamesFetcher.Fetch()
				if fetchErr != nil {
					l.Printf("Error with Twitter: %s", fetchErr.Error())
				}
			}

			// YouTube channel info
			for _, youtubeChannelFetcher := range youtubeChannelFetchers {
				fetchErr := youtubeChannelFetcher.Fetch()
				if fetchErr != nil {
					l.Printf("Error with YouTube: %s", fetchErr.Error())
				}
			}

			l.Println("Updated metricss")
			l.Printf("Waiting %d seconds until the next fetch time", intervalSecondsInt)

			// Wait for the next time period
			time.Sleep(time.Duration(intervalSecondsInt) * time.Second)
		}
	}()

	// Set up the HTTP handler and block
	addr := ":" + httpPort
	l.Printf("Setting up http service %s/%s", addr, httpPath)
	http.Handle("/"+httpPath, promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{
		ErrorLog: l,
	}))
	http.ListenAndServe(addr, nil)
}
