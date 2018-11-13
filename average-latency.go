package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func querySites(sites []string) {
	for _, site := range sites {
		querySite(site)
	}
}

func querySite(site string) {
	fmt.Println("Querying ", site)
	conn, err := net.Dial("tcp", site)
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
		Name: "average-latency",
		Usage: "get average latency of web sites",
		Action: func(c *cli.Context) error {
			sites := c.Args().Slice()
			querySites(sites)
			return nil
		},
	}

	app.Run(os.Args)
}
