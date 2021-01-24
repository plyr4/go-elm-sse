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
$ git clone git.target.com/DavidVader/go-elm-sse
$ cd go-elm-sse

# run Elm app
$ elm-app start

# run Go gin server
$ cd sse-server
$ go build sse-server/main.go
$ ./main

```

1. Open the Elm app http://localhost:3000

2. Click `connect` to 

- create the `EventSource` using the url `http://localhost:8080/DavidVader/heyvela/builds/2/steps/3/logs/events`
- create event listener to receive stream data via `message`

3. The `gin` server will

- receive a `GET` to `http://localhost:8080/DavidVader/heyvela/builds/2/steps/3/logs/events` and [create a channel](https://git.target.com/DavidVader/go-elm-sse/blob/master/sse-server/main.go#L44-L47) to send real-time server updates to
- [create a c.Stream](https://git.target.com/DavidVader/go-elm-sse/blob/master/sse-server/main.go#L67-L80) to relay events in the channel to the client's connection
- [mock some updates](https://git.target.com/DavidVader/go-elm-sse/blob/master/sse-server/main.go#L49-L65) to send events to the client in real time

Watch the events stream into the view

### Pros

`SSE` and `EventSource` are implemented solely through HTTP, no need for websockets.

Less polling, one connection per open log, per tab.

No more refreshing for logs. Allows for real-time server update rendering, which is the best user experience.

### Cons

There exists no native SSE library in `elm/http`, requiring more overhead in Elm through [ports](https://guide.elm-lang.org/interop/ports.html).

Requires an extension of the server to accept and maintain incoming stream connections.

Open connection maintenance with the load balancer could be complicated.

