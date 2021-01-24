# SSE in Go/Elm

This repository is meant to examine the benefits and implementation patterns for [Server Sent Events](https://medium.com/conectric-networks/a-look-at-server-sent-events-54a77f8d6ff7).

SSE is a server-to-client communication pattern that involves data streaming over HTTP. 

### How it works

1. Client creates an event stream connection over HTTP via the [EventSource](https://developer.mozilla.org/en-US/docs/Web/API/EventSource) Javascript construct, only requires a URL and optional authentication credentials.
1. Server accepts connection
1. Client creates event listener to receive stream data via `message`
1. Server streams data to the client via `message`

### Run the example

1. Install [Elm](https://guide.elm-lang.org/install/elm.html).
```bash
$ elm --version
0.19.1
```
2. Install [Go](https://golang.org/doc/install).
```bash
$ go version
go version go1.15.3 darwin/amd64
```
3. Clone the repo.
```bash
$ git clone git@github.com:davidvader/go-elm-sse.git
$ cd go-elm-sse
```
4. Run the Elm app.
```bash
$ elm-app start
```
5. Run the Go gin server.
```bash
$ cd sse-server
$ go build main.go
$ ./main
```
6. Open the Elm app http://localhost:3000.
7. Click `connect`, Elm executes javascript that will: 
- create the `EventSource` using the url `http://localhost:8080/DavidVader/heyvela/builds/2/steps/3/logs/events`
- create event listener to receive stream data via `message`
8. Upon receiving the `GET` to `http://localhost:8080/DavidVader/heyvela/builds/2/steps/3/logs/events`, the `gin` server will:
- [create a Go channel](https://github.com/davidvader/go-elm-sse/blob/master/sse-server/main.go#L44-L47) to send real-time server updates to
- [create a gin c.Stream](https://github.com/davidvader/go-elm-sse/blob/master/sse-server/main.go#L67-L80) to relay events in the channel to the client's connection
- [mock real-time updates](https://github.com/davidvader/go-elm-sse/blob/master/sse-server/main.go#L49-L65) to send events to the client in real time
9. Watch the events stream data into the view.

### Discussion

#### Pros

- Concept of `SSE` is implemented solely through HTTP, no need for websockets.
- `EventSource` is implemented natively through Javascript, fitting naturally into Elm ports.
- `SSE` has an existing `gin` implemention, meaning it will plug into current infrastructure well. 
- Events are sent from the server to the client, meaning less polling. One streaming connection per open log, per tab.
- No more refreshing for logs. Allows for real-time server update rendering. Improved UX.

#### Cons

- Added complexity. API polling and websockets are far more straight forward and require less code to maintain client-server communication.
- An Elm native SSE library, nor are there helpers in `elm/http`, requiring more overhead through [ports](https://guide.elm-lang.org/interop/ports.html).
- Requires added functionality to the server to accept incoming and maintain active stream connections.
- Maintaining open connections and making it thread safe for the load balancer.
- `EventSource` appears to be tied directly to `GET` so it relies heavily on the URL and headers to provide UI connection data.

#### Vela-specific concerns 

- Should SSE be integrated into the worker? ie JUST logs
- Should SSE apply to all resource events? ie SSE on every resource endpoint
- How can we integrate SSE behind the scenes with better HTTP manipulation?

#### Conclusion

If you can deal with the code overhead and complexity added to the infrastructure, then

SSE > websockets > API polling

### Resources

- [Server Sent Events](https://medium.com/conectric-networks/a-look-at-server-sent-events-54a77f8d6ff7)
- [gin](https://github.com/gin-gonic/gin)
- [Elm ports](https://guide.elm-lang.org/interop/ports.html)
- [EventSource](https://developer.mozilla.org/en-US/docs/Web/API/EventSource)
