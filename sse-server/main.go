package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	// create gin server
	r := gin.Default()

	// enable CORS
	r.Use(CORSMiddleware())

	r.GET("/stream", func(c *gin.Context) {

		// set timeout
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// channels for stream data and done
		done := make(chan bool)

		// use timeout so connections dont last forever
		go func() {
			for {
				select {
				case <-c.Request.Context().Done():
					// client is done
					done <- true
					return
				case <-ctx.Done():
					// timeout
					switch ctx.Err() {
					case context.DeadlineExceeded:
						fmt.Println("timeout")
					}
					done <- true
					return
				}
			}
		}()

		client, err := docker.NewEnvClient()
		if err != nil {
			return
		}

		// tail logs
		logs, err := client.ContainerLogs(ctx, "c216a64c3c5d",
			types.ContainerLogsOptions{
				ShowStdout: true, ShowStderr: true, Follow: true, Tail: "3",
			},
		)
		if err != nil {
			fmt.Println("reader error: ", err.Error())
			done <- true
			return
		}

		// code for copying logs using pipe taken from pkg-runtime
		// https://github.com/go-vela/pkg-runtime/blob/4591dd61eeb5982bdb8e7a3c87ca2d05aac459bd/runtime/docker/container.go

		// create in-memory pipe for capturing logs
		rc, wc := io.Pipe()

		// capture all stdout and stderr logs
		go func() {

			// copy container stdout and stderr logs to our in-memory pipe
			//
			// https://godoc.org/github.com/docker/docker/pkg/stdcopy#StdCopy
			_, err := stdcopy.StdCopy(wc, wc, logs)
			if err != nil {
				logrus.Errorf("unable to copy logs for container: %v", err)
			}

			// close logs buffer
			logs.Close()

			// close in-memory pipe write closer
			wc.Close()
		}()

		// message count
		count := 0

		c.Stream(func(w io.Writer) bool {

			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {

				// send 'message' event to client
				c.Render(-1, sse.Event{
					Id:    strconv.Itoa(count),
					Event: "message",
					Data:  string(scanner.Bytes()),
				})

				count++

				return true
			}

			return false
		})
	})

	r.Run("localhost:9999")
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
