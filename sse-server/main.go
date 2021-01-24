package main

import (
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// initialize connections map
	connections = make(map[string](chan string))

	// create gin server
	r := gin.Default()

	// enable CORS
	r.Use(CORSMiddleware())

	// handler for connecting to an EventSource and streaming data
	r.GET("/:org/:repo/builds/:build/steps/:step/logs/events", func(c *gin.Context) {

		// initialize channel for handling messages
		stream := make(chan string, 10)

		// mock background event of lo
		go func() {
			defer close(stream)

			// mock some growing log output for the demo
			log := ""
			for i := 0; i < 5; i++ {

				// line num
				n := ((strconv.Itoa(i + 1)) + "   "

				// mock a log line and send it to the channel
				log +=  n + logs[i])
				stream <- log

				time.Sleep(time.Second * 2)
			}
		}()

		// create a stream for this connection and send events from server to client when
		//      they are received through the channel
		c.Stream(func(w io.Writer) bool {

			// listen for events to this stream until the channel is closed
			for msg := range stream {
				c.SSEvent("message", msg)
				return true
			}

			// done with this stream connection
			return false
		})
	})

	r.Run()
}

// CORSMiddleware gin middleware for allowing cors
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// mock log feed
var logs = []string{
	`time="2021-01-20T21:11:43Z" level=info msg="Vela Artifactory Plugin" code="https://github.com/go-vela/vela-artifactory" docs="https://go-vela.github.io/docs/plugins/registry/artifactory" registry="https://hub.docker.com/r/vela-artifactory"
`,
	`[Info] [Thread 2] Uploading artifact: source/sample.jar
`,
	`[Error] [Thread 2] Artifactory response: 401 Unauthorized
`,
	`[Error] Failed uploading 1 artifacts.
`,
	`
`,
}
