package main

import (
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

// mock log feed
var logs = []string{
	`time="2021-01-20T21:11:43Z" level=info msg="Vela Artifactory Plugin" code="https://github.com/go-vela/vela-artifactory" docs="https://go-vela.github.io/docs/plugins/registry/artifactory" registry="https://hub.docker.com/r/target/vela-artifactory"
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

// connections
//  map that holds the data for all persistent EventSource clients
var connections map[string](chan string)

func main() {

	// initialize connections map
	// connections = make(map[string](chan string))

	// create gin server
	gin.DefaultWriter = colorable.NewColorableStderr()
	r := gin.Default()

	// enable CORS
	r.Use(CORSMiddleware())

	// handler for connecting to an EventSource and streaming data
	r.GET("/:org/:repo/builds/:build/steps/:step/logs/events", func(c *gin.Context) {

		// 1. upon GET, initialize Go channel for handling messages
		stream := make(chan string, 10)

		// 2. mock background event of log feed
		//       send event to Go channel
		go func() {

			// close the channel when we are done with it
			defer close(stream)
			log := ""
			for i := 0; i < 5; i++ {
				log += ((strconv.Itoa(i + 1)) + "   " + logs[i])
				stream <- log
				if i < 2 {
					time.Sleep(time.Second * 2)
				} else {
					time.Sleep(time.Second * 1)
				}
			}
		}()

		// 3. create a stream for this connection and send events from server to client when
		//      they are received through the channel
		c.Stream(func(w io.Writer) bool {

			//
			// range check over connected EventSource clients
			//
			// when events are sent from container logs -> Go channel they end up here
			for msg := range stream {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	r.Use(static.Serve("/", static.LocalFile("./public", true)))
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
