# social-media-user-exporter
Prometheus exporter for gathering user account info on social media platforms

## Usage

The service is controlled by environment variables. By specifying the access credentials for a social media platform you enable it to run.

## Platforms

Supported platforms:

- Twitter
- YouTube
- Facebook (coming soon)
- Instagram (coming soon)

### Environment variables

All environment variables have an `SMUE_` prefix which stands for Social Media User Exporter
#### Basic

| Name | Default | Description |
|------|---------|-------------|
| `SMUE_HTTP_PORT` | `9100` | The HTTP port used to host the `/metrics` endpoint |
| `SMUE_HTTP_PATH` | `metrics` | The path the metrics are hosted on |
| `SMUE_METRICS_PREFIX` | `social_media_user_` | Customise the metrics prefix, be careful not to overwrite any other metrics |
| `SMUE_INTERVAL_SECONDS` | `300` | How often to poll for more information. If you have trouble with rate limits, make this larger. Default of 5 minutes (300 seconds) should be plenty |

#### Twitter

| Name | Default | Description |
|------|---------|-------------|
| `SMUE_TWITTER_SCREEN_NAMES` | `""` | A comma separated list of Twitter user screen names. Leave this empty to not fetch Twitter user metrics |
| `SMUE_TWITTER_CLIENT_ID` | `""` | The Twitter app OAuth client ID |
| `SMUE_TWITTER_CLIENT_SECRET` | `""` | The Twitter app OAuth client secret |
| `SMUE_TWITTER_ACCESS_TOKEN` | `""` | The Twitter access token for a user |
| `SMUE_TWITTER_ACCESS_TOKEN_SECRET` | `""` | The Twitter access token secret for a user |

#### YouTube

| Name | Default | Description |
|------|---------|-------------|
| `SMUE_YOUTUBE_CHANNEL_IDS` | `""` | A comma separated list of YouTube channel ID's. Leave this empty to not fetch Twitter user metrics |
| `SMUE_YOUTUBE_CLIENT_ID` | `""` | The YouTube app OAuth client ID |
| `SMUE_YOUTUBE_CLIENT_SECRET` | `""` | The YouTube app OAuth client secret |
| `SMUE_YOUTUBE_ACCESS_TOKEN` | `""` | The YouTube access token for a user |
| `SMUE_YOUTUBE_REFRESH_TOKEN` | `""` | The YouTube refresh token for a user so the app can handle token refreshing |


## Example kubectl deployment

```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: social-media-user-exporter
  namespace: default
spec:
  replicas: 1
  revisionHistoryLimit: 1
  selector:
    matchLabels:
      app: social-media-user-exporter
  template:
    metadata:
      labels:
        app: social-media-user-exporter
      annotations:
        "prometheus.io/scrape": "true"
    spec:
      containers:
        - name: social-media-user-exporter
          image: ollieparsley/social-media-user-exporter:latest
          ports:
            - containerPort: 9100
          env:
            - name: SMUE_INTERVAL_SECONDS
              value: "600"
            - name: SMUE_TWITTER_SCREEN_NAMES
              value: ollieparsley
            - name: SMUE_YOUTUBE_CHANNEL_IDS
              value: UC7R4lEiVaathpWrwXArnKQg
            - name: SMUE_TWITTER_CLIENT_ID
              valueFrom:
                secretKeyRef:
                  name: twitter
                  key: client_id
            - name: SMUE_TWITTER_CLIENT_SECRET
              valueFrom:
                secretKeyRef:
                  name: twitter
                  key: client_secret
            - name: SMUE_TWITTER_ACCESS_TOKEN
              valueFrom:
                secretKeyRef:
                  name: twitter
                  key: access_token
            - name: SMUE_TWITTER_ACCESS_TOKEN_SECRET
              valueFrom:
                secretKeyRef:
                  name: twitter
                  key: access_token_secret
            - name: SMUE_YOUTUBE_CLIENT_ID
              valueFrom:
                secretKeyRef:
                  name: youtube
                  key: client_id
            - name: SMUE_YOUTUBE_CLIENT_SECRET
              valueFrom:
                secretKeyRef:
                  name: youtube
                  key: client_secret
            - name: SMUE_YOUTUBE_ACCESS_TOKEN
              valueFrom:
                secretKeyRef:
                  name: youtube
                  key: access_token
            - name: SMUE_YOUTUBE_REFRESH_TOKEN
              valueFrom:
                secretKeyRef:
                  name: youtube
                  key: refresh_token

```
