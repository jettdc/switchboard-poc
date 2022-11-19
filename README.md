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

### Dynamic Subscriptions
For certain use cases, there may be a need to subscribe to more than one endpoint. Instead of forcing the user to open
multiple websockets (which also would increase load on the server), switchboard provides an endpoint for dynamic 
subscriptions.

#### Establishing a Connection
Create a websocket connection to the endpoint `/multi`

#### Subscribing to Endpoints
To subscribe to new endpoints, the client must send a message to switchboard:
```json
{
  "action": "SUBSCRIBE",
  "endpoints": [
    {
      "endpoint": "/your/endpoint"
    }
  ]
}
```

If the operation fails, the server will respond with an error message:
```json
{
  "endpoint": "/multi",
  "messageType": "ERROR",
  "message": {
    "error": "No route configuration found for endpoint \"/ws/fake\""
  }
}
```

You can also do parameterized endpoints, like so:
```json
{
  "action": "SUBSCRIBE",
  "endpoints": [
    {
      "endpoint": "/your/endpoint/:id", 
      "params": {"id":  "45"}
    }
  ], 
  "requestId": "1"
}
```

Notice that we've included an optional `requestId` in the message. If we receive a response from the server, the
`requestId` will be the same.
```json
{
  "endpoint": "/multi",
  "messageType": "ERROR",
  "message": {
    "error": "No route configuration found for endpoint \"/your/endpoint/:id\""
  },
  "requestId": "1"
}
```

#### Unsubscribing from Endpoints
Unsubscribing works the same, but the action is `UNSUBSCRIBE`

```json
{
  "action": "UNSUBSCRIBE",
  "endpoints": [
    {
      "endpoint": "/your/endpoint"
    }
  ]
}
```

### Setup
In the `/switchboard` directory
- `go get`

### Running
In the `/switchboard` directory
- `go run .`

You will need to ensure that any plugins that you've specified actually exist.

### Plugins
Plugins must adhere to the `MiddlewarePlugin` and `EnrichmentPlugin` interfaces, defined in `switchboard/config/pluginInterfaces.go`. Plugins can be defined as follows (this example shows middleware):

`myPlugin.go`
```golang
package main
import (
    "net/http"
)

type myPluginName string
func (myPluginName) Process(r *http.Request) error {
    // My Plugin Implementation
}

var MiddlewarePlugin myPluginName
```

Compile the plugin to produce the .so file (i.e. `myPlugin.so`) that can then be referenced from `config.yaml`:
`go build -buildmode=plugin -o myPlugin.so myPlugin.go`

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
