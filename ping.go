package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func querySites(urls []string) {
	for _, url := range urls {
		querySite(url)
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

func main() {
	app := &cli.App{
		Name: "ping",
		Usage: "get latencies of web sites",
		Action: func(c *cli.Context) error {
			urls := c.Args().Slice()
			querySites(urls)
			return nil
		},
	}

	app.Run(os.Args)
}
