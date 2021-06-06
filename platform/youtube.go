package platform

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// YouTube platform
type YouTube struct {
	metricPrefix     string
	client           *youtube.Service
	logger           *log.Logger
	factory          promauto.Factory
	channelID        string
	subscribersGauge prometheus.Gauge
	viewsGauge       prometheus.Gauge
	videosGauge      prometheus.Gauge
}

// NewYouTube create new instance of the YouTube user fetcher
func NewYouTube(logger *log.Logger, metricPrefix string, factory *promauto.Factory, channelID string, clientID string, clientSecret string, accessToken string, refreshToken string) (YouTube, error) {

	// Check params
	if channelID == "" {
		return YouTube{}, errors.New("a channel ID is required")
	}
	if accessToken == "" {
		return YouTube{}, errors.New("an access token is required")
	}
	if refreshToken == "" {
		return YouTube{}, errors.New("a refresh token is required")
	}

	// OAuth token
	oauthConfig := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/youtube"},
	}
	oauthToken := oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		Expiry:       time.Now().Add(-(time.Hour * 24)), // Time in the past to force a refresh
	}
	httpClient := oauthConfig.Client(oauth2.NoContext, &oauthToken)

	// Verify credentials
	client, err := youtube.New(httpClient)
	if err != nil {
		return YouTube{}, fmt.Errorf("Error when creating youtube client %s: %s", channelID, err.Error())
	}

	// Fetch the channel
	youtube := YouTube{
		metricPrefix: metricPrefix,
		logger:       logger,
		client:       client,
		channelID:    channelID,
	}

	// Get the channel and see if the channel ID exists
	foundChannel := youtube.getChannel()
	if foundChannel == nil {
		return YouTube{}, fmt.Errorf("The authenticated user doesn't have access to the channel with ID %s", channelID)
	}

	// Add the gauges
	youtube.subscribersGauge = factory.NewGauge(prometheus.GaugeOpts{
		Name: metricPrefix + "youtube_subscribers",
		Help: "The number of subscribers the channel has",
		ConstLabels: prometheus.Labels{
			"channel_id":   channelID,
			"channel_name": foundChannel.Snippet.Title,
		},
	})
	youtube.viewsGauge = factory.NewGauge(prometheus.GaugeOpts{
		Name: metricPrefix + "youtube_views",
		Help: "The number of views the channel has",
		ConstLabels: prometheus.Labels{
			"channel_id":   channelID,
			"channel_name": foundChannel.Snippet.Title,
		},
	})
	youtube.videosGauge = factory.NewGauge(prometheus.GaugeOpts{
		Name: metricPrefix + "youtube_videos",
		Help: "The number of videos the channel has",
		ConstLabels: prometheus.Labels{
			"channel_id":   channelID,
			"channel_name": foundChannel.Snippet.Title,
		},
	})

	return youtube, nil
}

// Get channel
func (p YouTube) getChannel() *youtube.Channel {
	channelListCall := p.client.Channels.List([]string{"snippet", "contentDetails", "statistics"})
	channelListResponse, err := channelListCall.Id(p.channelID).Do()
	if err != nil {
		p.logger.Printf("Error getting channel list: %s", err.Error())
		return nil
	}

	// Get the channel and see if the channel ID exists
	var foundChannel *youtube.Channel
	for _, channel := range channelListResponse.Items {
		if channel.Id == p.channelID {
			foundChannel = channel
		}
	}
	return foundChannel
}

// Fetch Format key with lavels
func (p YouTube) Fetch() error {

	channel := p.getChannel()
	if channel == nil {
		return fmt.Errorf("Channel %s is now owned by the authenticated user", p.channelID)
	}

	p.subscribersGauge.Set(float64(channel.Statistics.SubscriberCount))
	p.viewsGauge.Set(float64(channel.Statistics.ViewCount))
	p.videosGauge.Set(float64(channel.Statistics.VideoCount))

	return nil
}
