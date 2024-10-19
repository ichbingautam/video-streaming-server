package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// VERSION indicates which version of the binary is running.
var VERSION string

const (
	PARAMETER_HOST_NAME     = "host"
	PARAMETER_PORT_NAME     = "port"
	PARAMETER_FILENAME_NAME = "fileName"
)

var (
	fileName   string
	host       string
	port       string
	protocol   = "tcp" // following protocol can be used: tcp, tcp4, tcp6, unix, or unix packet.
	timeout    = time.Microsecond * 2000
	bufferSize = 1024
	wg         sync.WaitGroup
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			// Perform necessary cleanup or logging operations
			// Return an appropriate error value
		}
	}()

	a := cli.NewApp()
	a.Name = "Video Streaming Client /Files"
	a.Usage = "Client for Streaming Server (Video or Audio) / Source: Stream  media files"
	a.Author = "Shubham Gautam"
	a.Email = "ichbingautam@gmail.com"
	a.Version = VERSION

	a.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  PARAMETER_HOST_NAME,
			Usage: fmt.Sprintf("Streaming Server '%s', e.g. --host 1270.0.0.1", PARAMETER_PORT_NAME),
		},
		cli.StringFlag{
			Name:  PARAMETER_PORT_NAME,
			Usage: fmt.Sprintf("Streaming Server '%s', e.g. --port 8080", PARAMETER_PORT_NAME),
		},
		cli.StringFlag{
			Name:  PARAMETER_FILENAME_NAME,
			Usage: "path to medial file",
		},
	}

	a.Action = func(c *cli.Context) error {
		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		}

		if c.IsSet(PARAMETER_HOST_NAME) {
			host = c.String(PARAMETER_HOST_NAME)
			log.Printf("set paremeter: '%s' to %s\n", PARAMETER_HOST_NAME, host)
		}

		if c.IsSet(PARAMETER_PORT_NAME) {
			port = c.String(PARAMETER_PORT_NAME)
			log.Printf("set paremeter: '%s' to %s\n", PARAMETER_PORT_NAME, port)
		}

		if c.IsSet(PARAMETER_FILENAME_NAME) {
			fileName = c.String(PARAMETER_FILENAME_NAME)
			log.Printf("set paremeter: '%s' to %s\n", PARAMETER_FILENAME_NAME, fileName)
		}

		if !c.IsSet("help") {
			transferDate()
		}

		return nil
	}

	err := a.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func transferDate() {
	//establish connection
	connection, err := net.Dial(protocol, host+":"+port)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	// Open the file for reading
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The file:%s is %d bytes long", file.Name(), fi.Size())

	// Create a buffered reader
	reader := bufio.NewReader(file)

	var totalTransfered int64
	for totalTransfered < fi.Size() {

		if int64(bufferSize) > (fi.Size() - totalTransfered) {
			bufferSize = int(fi.Size() - totalTransfered)
		}

		buffer := make([]byte, bufferSize)
		n, err := reader.Read(buffer)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}

		// send chunk of date
		_, err = connection.Write(buffer)
		totalTransfered += int64(len(buffer))
		log.Printf("send data to %s:%s, package size:%d, transfered: %d of %d", host, port, len(buffer), totalTransfered, fi.Size())

		// Check for EOF
		if err == nil && n < bufferSize {
			break
		}
	}
}
