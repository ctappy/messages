# GRPC email slack api service

## Swagger
```
docker run -p 80:8080 -e SWAGGER_JSON=/foo/message.swagger.json -v $GOPATH/src/github.com/ctaperts/messages/messagepb/:/foo swaggerapi/swagger-ui
```
 
## Setup
Copy `config.json.example to `config.json`

### Slack bot setup
- Go to `https://api.slack.com/apps?new_app=1`
- Create new app
- set up permission scope, setup requires rtm permissions
- select `Install App to Workspace`
- copy Bot User OAuth Access Token key to `config.json`
- set channel name under `channel_name`

## Commands
the flag `--log-level` accepts `trace`, `info`, `warn`, `debug`, `error` and `fatal`, defaults to `error`

### Start slackbot
```
messages slackbot
```

### Start grpc server
```
messages grpc
```

## FAQ
- Email from address not being masked
see `https://stackoverflow.com/questions/13946581/spring-java-mail-the-from-address-is-being-ignored`, another option would be to change to postfix type solution
