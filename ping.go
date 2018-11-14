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

func querySites(urls []string, avgFlag bool) {
	var wg sync.WaitGroup
	latencies := make(chan []time.Duration, 100)
	urlsLen := len(urls)

	for _, url := range urls {
		wg.Add(1)
		switch avgFlag {
		case true:
			firstLatency, lastLatency := getLatency(url)
			latencies <- []time.Duration{firstLatency, lastLatency}
		default:
			querySite(url)
		}
		wg.Done()
	}

	wg.Wait()
	close(latencies)

	if avgFlag {
		var sumFirstLatency, sumLastLatency time.Duration

		for latency := range latencies {
			firstLatency, lastLatency := latency[0], latency[1]
			sumFirstLatency, sumLastLatency = sumFirstLatency + firstLatency, sumLastLatency + lastLatency
		}

		avgFirstLatency := sumFirstLatency.Seconds() / float64(urlsLen) * 1000
		avgAllLatency := sumLastLatency.Seconds() / float64(urlsLen) * 1000
		fmt.Printf("Avg first latency: %fms\n", avgFirstLatency)
		fmt.Printf("Avg last latency: %fms\n", avgAllLatency)
	}
}

func querySite(url string) {
	fmt.Println("Querying ", url)
	conn, err := net.Dial("tcp", url + ":80")
	defer conn.Close()

	conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))

	start := time.Now()
	oneByte := make([]byte, 1)
	_, err = conn.Read(oneByte)
	if err != nil {
		panic(err)
	}
	fmt.Println("First byte:", time.Since(start))

	_, err = ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("All bytes:", time.Since(start))
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
