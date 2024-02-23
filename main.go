package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/akamensky/argparse"
	"github.com/gorilla/websocket"
	"github.com/nxadm/tail"
	"github.com/oschwald/maxminddb-golang"
)

// this struct stores the location of an IP address
type mmrecord struct {
	Location struct {
		Latitude  float64
		Longitude float64
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// parses an IP address from the beginning of a line
func parseIp(line string) (net.IP, error) {
	ipregexp := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	ipstr := ipregexp.FindString(line)

	if ipstr == "" {
		return nil, fmt.Errorf("failed to find IP address")
	}

	// parse the IP address
	ipAddr := net.ParseIP(ipstr)
	if ipAddr == nil {
		return nil, fmt.Errorf("failed to parse IP address")
	}

	return ipAddr, nil
}

// handler is the main websocket handler
func handler(w http.ResponseWriter, r *http.Request, pool *WsPool) {
	log.Println("New connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	pool.Add(conn)
}

func main() {
	// setup command line arguments
	parser := argparse.NewParser("run", "run the geolog websocket server")
	logFile := parser.String("l", "log_file", &argparse.Options{Required: true, Help: "log file to tail"})
	geoliteDb := parser.String("g", "geodb_file", &argparse.Options{Required: true, Help: "geolite db to use"})

	port := parser.String("p", "port", &argparse.Options{Required: false, Help: "port to listen on", Default: "8080"})

	// parse the command line arguments
	err := parser.Parse(os.Args)
	if err != nil {
		log.Print(parser.Usage(err))
		os.Exit(1)
	}

	// open the maxmind db
	gdb, err := maxminddb.Open(*geoliteDb)
	if err != nil {
		log.Fatalf("failed to open maxmind db: %v", err)
		os.Exit(1)
	}
	defer gdb.Close()

	// tail the log file
	tail, err := tail.TailFile(*logFile, tail.Config{
		Follow: true,
		ReOpen: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	var throttler = NewIPThrottler(5*time.Second, 60*time.Second)

	// create a new pool and start it
	pool := NewWsPool()
	go pool.Start()

	go func() {
		// read lines from the log file
		for line := range tail.Lines {

			// parse the IP address from the line
			ip, _ := parseIp(line.Text)

			if ip == nil {
				continue
			}

			// ruhig Brauner!
			if !throttler.Allow(ip.String()) {
				continue
			}

			// lookup the IP address in the maxmind db
			var record mmrecord
			err = gdb.Lookup(ip, &record)
			if err != nil {
				log.Printf("failed to lookup ip %s: %v", ip, err)
				return
			}

			// format payload
			var payload = fmt.Sprintf("[%s, %f, %f]", ip, record.Location.Latitude, record.Location.Longitude)
			log.Println(payload)
			pool.broadcast <- payload
		}
	}()

	// serve the favicon.ico file
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	// serve the index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// serve the websocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, pool)
	})

	log.Println("Listening on :" + *port)

	// start the http server
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
