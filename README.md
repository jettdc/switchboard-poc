# switchboard

## Switchboard Core
### Setup
In the `/switchboard` directory
- `go get`
- `go run .`

### Manual Testing
- [Download insomnia](https://insomnia.rest/download) (or postman, but insomnia is a bit more lightweight and simple)
- On the left panel, click the +, then select "Websocket Connection"
- For the route, enter `ws://localhost:8080/{your route you want to connect to}`

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
