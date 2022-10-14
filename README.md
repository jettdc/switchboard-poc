# switchboard

### TODO
To find todo items, search the project for `// TODO:` comments

## Switchboard Core
### Setup
In the `/switchboard` directory
- `go get`

### Running
In the `/switchboard` directory
- `go run .`

You will need to ensure that any plugins that you've specified actually exist.

### Manual Testing
- [Download insomnia](https://insomnia.rest/download) (or postman, but insomnia is a bit more lightweight and simple)
- On the left panel, click the +, then select "Websocket Connection"
- For the route, enter `ws://localhost:8080/{your route you want to connect to}`

### Feature Backlog
- Different PubSub providers
- Message backlog resending
  - If someone subscribes to a channel and there have been messages sent before the connection, send those messages before sending any new messages. 
- Environment variables for pubsub connection
- Secure websockets
- Health check

## Backend Services
### Setup
### Running

## Frontend
### Setup
### Running

## Repo Organization

**High level plan: do all dev in switchboard repo, migrate Pizza app to separate repo at the end**
```
switchboard/
	src/
    	config/
        plugins/
        main.go
	Dockerfile
app-backend/
	src/
    	services-publishing-to-redis/
	Dockerfile
app-frontend/
	src/
    	index.html
.env
docker-compose
```
