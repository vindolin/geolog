package main

import (
	"fmt"
	"html/template"
	"io"
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
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
	} `maxminddb:"location"`
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
	conn, err := upgrader.Upgrade(w, r, nil)
	log.Println("New connection from:", conn.RemoteAddr())
	if err != nil {
		log.Println(err)
		return
	}

	pool.Add(conn)
}

func main() {
	// setup command line arguments
	parser := argparse.NewParser("run", "run the geolog websocket server")
	logFile := parser.String("l", "log_file",
		&argparse.Options{Required: true, Help: "log file to tail"})
	geoliteDb := parser.String("g", "geodb_file",
		&argparse.Options{Required: true, Help: "geolite db to use"})

	port := parser.String("p", "port",
		&argparse.Options{Required: false, Help: "port to listen on", Default: "8080"})

	darkMode := parser.Flag("d", "dark", &argparse.Options{Help: "dark mode"})

	const PING_INTERVAL = 10

	// parse the command line arguments
	err := parser.Parse(os.Args)
	if err != nil {
		log.Print(parser.Usage(err))
		os.Exit(1)
	}

	// this struct holds the data that will be passed to the index.html template
	type IndexTemplateData struct {
		DarkMode bool
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
		Follow:   true,
		ReOpen:   true,
		Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
	})
	if err != nil {
		log.Fatal(err)
	}

	// create a new pool and start it
	pool := NewWsPool()
	go pool.Start()

	go func() {
		var throttler = NewIPThrottler(5*time.Second, 60*time.Second)

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
			var payload = fmt.Sprintf(
				`["%s", %f, %f]`,
				ip, record.Location.Latitude,
				record.Location.Longitude)

			log.Println(payload)

			// broadcast the payload to all connected clients
			pool.broadcast <- payload
		}
	}()

	// send a ping every n seconds
	// this is used on the javascript side to keep the websocket alive
	go func() {
		for {
			time.Sleep(PING_INTERVAL * time.Second)

			var payload = fmt.Sprintf(
				"ping %d", time.Now().Unix())
			pool.broadcast <- payload
		}
	}()

	// serve the favicon.ico file
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	// serve the index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.ParseFiles("index.html")).Execute(
			w, IndexTemplateData{DarkMode: *darkMode})
	})

	// serve the websocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, pool)
	})

	log.Println("Listening on :" + *port)

	// start the http server
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
