# GRPC email slack api service

## Setup

Copy `config.json.example to `config.json`

### Slack
- Go to `https://api.slack.com/apps?new_app=1`
- Create new app
- set up permission scope, setup requires rtm permissions
- select `Install App to Workspace`
- copy Bot User OAuth Access Token key to `config.json`
- copy the channel id from the URL
