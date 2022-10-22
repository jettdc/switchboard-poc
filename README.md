# switchboard
*Drop in infrastructure to forward backend pubsub messages to frontend websockets, with a robust plugin system*

Switchboard is a plug-and-play service for creating websocket endpoints that listen to backend pubsub topics.
The service allows you to easily define pipelines that route backend messages to frontend clients,
decoupling the message writers and readers, as well as removing unneeded complexity.

If a backend service needs to provide a frontend-relevant update, they don't need to worry about handling their own 
websocket server (which would additionally introduce various security concerns) or routing that message to another 
custom backend service, but rather can just shoot it off to a pubsub provider and trust that it will be delivered.

The plugin system allows the addition of functionality such as connection authentication, message enrichment, and missed
message recovery.

---

### Config
```yaml
server:
  host: <string>
  port: <int>
  pubsub:
    provider: redis # currently the only supported provider
  ssl: <optional> # if not present, serves on http
    mode: <none | auto | manual>
    cert: <string, optional> # path, required for manual mode
    key:  <string, optional> # path, required for manual mode
  env-file: <string, optional> # path to .env
routes:
  - endpoint: <string> # rest endpoint for initiating the websocket, e.g. "/api/ws/customers/:id/orders"
    topics:
      - <string> # topics to subscribe to on pubsub provider, e.g. "/customers/:id/orders"
    plugins: <optional>
      middleware: <optional>
        - <string> # path to go plugin file (.so)
      message-enrichment: <optional>
        - <string> # path to go plugin file (.so)

```

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