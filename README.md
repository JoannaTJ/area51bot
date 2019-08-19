# area51bot

Slack bot to practice Golang

## Setup

Follow the Golang environment guide on: https://golang.org/doc/install

```shell
go get github.com/rldiao/area51bot
docker build -t registry.heroku.com/area51bot/web
docker run -ti --rm -p 8080:3000 registry.heroku.com/area51bot/web
```

expose on port 3000

## Deployment

```shell
docker build -t registry.heroku.com/area51bot/web
heroku container:push web -a area51bot
heroku container:release web -a area51bot
```
