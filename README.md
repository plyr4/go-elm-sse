# SSE in Go/Elm

This repository is meant to examine the benefits and implementation patterns for [Server Sent Events](https://medium.com/conectric-networks/a-look-at-server-sent-events-54a77f8d6ff7).

SSE is a server-to-client communication pattern that involves data streaming over HTTP. 

### How it works

1. Client creates an event stream connection over HTTP via the [EventSource](https://developer.mozilla.org/en-US/docs/Web/API/EventSource) Javascript construct, only requires a URL and optional authentication credentials.
1. Server accepts connection
1. Client creates event listener to receive stream data via `message`
1. Server streams data to the client via `message`

### Run the example

```bash

# clone the repo
$ git clone git@github.com:davidvader/go-elm-sse.git
$ cd go-elm-sse

# run Elm app
$ elm-app start

# run Go gin server
$ cd sse-server
$ go build sse-server/main.go
$ ./main

```

1. Open the Elm app http://localhost:3000.
2. Click `connect`, Elm executes javascript that will: 
- create the `EventSource` using the url `http://localhost:8080/DavidVader/heyvela/builds/2/steps/3/logs/events`
- create event listener to receive stream data via `message`

3. Upon receiving the `GET` to `http://localhost:8080/DavidVader/heyvela/builds/2/steps/3/logs/events`, the `gin` server will:
- [create a Go channel](https://github.com/davidvader/go-elm-sse/blob/master/sse-server/main.go#L44-L47) to send real-time server updates to
- [create a gin c.Stream](https://github.com/davidvader/go-elm-sse/blob/master/sse-server/main.go#L67-L80) to relay events in the channel to the client's connection
- [mock real-time updates](https://github.com/davidvader/go-elm-sse/blob/master/sse-server/main.go#L49-L65) to send events to the client in real time

4. Watch the events stream data into the view.

### Pros

- Concept of `SSE` is implemented solely through HTTP, no need for websockets.
- `EventSource` is implemented natively through Javascript, fitting naturally into Elm ports.
- `SSE` has an existing `gin` implemention, meaning it will plug into current infrastructure well. 
- Events are sent from the server to the client, meaning less polling. One streaming connection per open log, per tab.
- No more refreshing for logs. Allows for real-time server update rendering. Improved UX.

### Cons

- There exists no native SSE library in `elm/http`, requiring more overhead in Elm through [ports](https://guide.elm-lang.org/interop/ports.html).
- Requires an extension of the server to accept and maintain incoming stream connections.
- Open connection maintenance with the load balancer could be complicated.

### Resources

- [Server Sent Events](https://medium.com/conectric-networks/a-look-at-server-sent-events-54a77f8d6ff7)
- [gin](https://github.com/gin-gonic/gin)
- [Elm ports](https://guide.elm-lang.org/interop/ports.html)
- [EventSource](https://developer.mozilla.org/en-US/docs/Web/API/EventSource)
