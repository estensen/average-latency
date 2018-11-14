package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"
)

type queryResult struct {
	url string
	firstLatency time.Duration
	lastLatency time.Duration
}

func querySites(urls []string, avgFlag bool) {
	var wg sync.WaitGroup
	latencies := make(chan queryResult, 100)
	urlsLen := len(urls)

	for _, url := range urls {
		wg.Add(1)
		firstLatency, lastLatency := getLatency(url)
		latencies <- queryResult{url, firstLatency, lastLatency}
		wg.Done()
	}

	wg.Wait()
	close(latencies)

	var sumFirstLatency, sumLastLatency time.Duration

	for latency := range latencies {
		url, firstLatency, lastLatency := latency.url, latency.firstLatency, latency.lastLatency
		sumFirstLatency, sumLastLatency = sumFirstLatency + firstLatency, sumLastLatency + lastLatency

		if !avgFlag {
			println(url)
			println("First latency:", firstLatency)
			println("Last latency:", lastLatency)
		}
	}

	if avgFlag {
		avgFirstLatency := sumFirstLatency.Seconds() / float64(urlsLen) * 1000
		avgAllLatency := sumLastLatency.Seconds() / float64(urlsLen) * 1000
		fmt.Println("Aggregated stats")
		fmt.Printf("Avg first latency: %fms\n", avgFirstLatency)
		fmt.Printf("Avg last latency: %fms\n", avgAllLatency)
	}
}

func getLatency(url string) (time.Duration, time.Duration) {
	fmt.Println("Querying ", url)
	conn, err := net.Dial("tcp", url+":80")
	defer conn.Close()

	conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))

	start := time.Now()
	oneByte := make([]byte, 1)

	_, err = conn.Read(oneByte)
	if err != nil {
		panic(err)
	}
	firstLatency := time.Since(start)

	_, err = ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	lastLatency := time.Since(start)

	return firstLatency, lastLatency
}

func main() {
	app := &cli.App{
		Name: "ping",
		Usage: "get latencies of web sites",

		Flags: []cli.Flag {
			&cli.BoolFlag {
				Name: "average",
				Aliases: []string{"a"},
				Usage: "Display average latency",
			},
		},
		Action: func(c *cli.Context) error {
			urls := c.Args().Slice()
			avgFlag := c.Bool("average")

			querySites(urls, avgFlag)
			return nil
		},
	}

	app.Run(os.Args)
}
