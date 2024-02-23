package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/akamensky/argparse"
	"github.com/gorilla/websocket"
	"github.com/nxadm/tail"
	"github.com/oschwald/maxminddb-golang"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	lastIpstr string
)

type mmrecord struct {
	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
	} `maxminddb:"location"`
}

func parseIp(line string) (net.IP, error) {
	ipregexp := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	ipstr := ipregexp.FindString(line)

	if ipstr == lastIpstr {
		return nil, nil
	}
	lastIpstr = ipstr

	if ipstr == "" {
		return nil, fmt.Errorf("failed to find IP address")
	}

	ipAddr := net.ParseIP(ipstr)
	if ipAddr == nil {
		return nil, fmt.Errorf("failed to parse IP address")
	}

	return ipAddr, nil
}

func handler(w http.ResponseWriter, r *http.Request, tail *tail.Tail, gdb *maxminddb.Reader) {
	log.Println("New connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for line := range tail.Lines {
		ipAddr, _ := parseIp(line.Text)

		if ipAddr == nil {
			continue
		}

		var record mmrecord
		err = gdb.Lookup(ipAddr, &record)
		if err != nil {
			log.Printf("failed to lookup ip %s: %v", ipAddr, err)
			return
		}

		log.Printf("ip: %s, lat: %f, long: %f\n", ipAddr, record.Location.Latitude, record.Location.Longitude)
		var payload = fmt.Sprintf("[%s, %f, %f]", ipAddr, record.Location.Latitude, record.Location.Longitude)

		err = conn.WriteMessage(websocket.TextMessage, []byte(payload))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	parser := argparse.NewParser("run", "run the geolog websocket server")
	logFile := parser.String("l", "log_file", &argparse.Options{Required: true, Help: "log file to tail"})
	geoliteDb := parser.String("g", "geodb_file", &argparse.Options{Required: true, Help: "geolite db to use"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Print(parser.Usage(err))
		os.Exit(1)
	}

	gdb, err := maxminddb.Open(*geoliteDb)
	if err != nil {
		log.Fatalf("failed to open maxmind db: %v", err)
		os.Exit(1)
	}
	defer gdb.Close()

	tail, err := tail.TailFile(*logFile, tail.Config{
		Follow: true,
		ReOpen: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, tail, gdb)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
