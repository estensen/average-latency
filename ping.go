package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func querySites(urls []string, avgFlag bool) {
	for _, url := range urls {
		switch avgFlag {
		case true:
			getLatency(url)
		default:
			querySite(url)
		}
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

func getLatency(url string) {
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

	latency := time.Since(start)
	fmt.Println("First byte:", latency)

	_, err = ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
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