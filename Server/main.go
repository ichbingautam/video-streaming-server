package main

import (
	"bufio"
	"github.com/urfave/cli"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// VERSION indicates which version of the binary is running.
var VERSION string

var (
	port     string
	protocol = "tcp" // following protocol can be used: tcp, tcp4, tcp6, unix, or unix packet.
	timeout  = time.Microsecond * 2000
	wg       sync.WaitGroup
)

const (
	DEFAULT_PORT = "8080"
)

func main() {
	VERSION = "0.0.1"

	a := cli.NewApp()
	a.Name = "Streaming Server (Video or Audio)"
	a.Usage = "Streaming Server let different sources stream video/audio files"
	a.Author = "Valentyn Ponomarenko"
	a.Email = "ValentynPonomarenko@gmail.com"
	a.Version = VERSION

	a.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Value: DEFAULT_PORT,
			Usage: "Streaming Server end-point port, e.g. --port 8080",
		},
	}

	a.Action = func(c *cli.Context) error {
		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		}

		if c.IsSet("port") {
			port = c.String("port")
			log.Printf("set paremeter: 'port' to %s\n", port)
		}

		if !c.IsSet("help") {
			// start listening webSocket connections
			startSocket()
		}

		return nil
	}

	err := a.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

// processing webSocket connection
func startSocket() {
	log.Printf("Port:%s/n", port)
	lister, err := net.Listen(protocol, ":"+port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lister.Accept()
		if err != nil {
			panic(err)
		}
		go handleSocketConnection(conn)
	}
}

func handleSocketConnection(conn net.Conn) {
	name := conn.RemoteAddr()
	log.Printf("connected: %+v", name)

	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		//TODO: next to change to Bytes, Text using just for testing
		bytes := scanner.Bytes()
		log.Printf("receiverd from %+v package size:%d", name, len(bytes))
	}
}
