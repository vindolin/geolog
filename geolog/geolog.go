package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"

	"github.com/akamensky/argparse"
	"github.com/hpcloud/tail"
	"github.com/oschwald/maxminddb-golang"
)

type mmrecord struct {
	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
	} `maxminddb:"location"`
}

var lastIpstr string

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

	ip := net.ParseIP(ipstr)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse IP address")
	}

	return ip, nil
}

func handleLine(line string, geoDb *maxminddb.Reader) (*mmrecord, error) {
	ip, err := parseIp(line)
	if err != nil {
		log.Printf("failed to parse IP address: %v", err)
		return nil, err
	}

	if ip == nil {
		return nil, nil
	}

	var record mmrecord
	err = geoDb.Lookup(ip, &record)
	if err != nil {
		log.Printf("failed to lookup ip %s: %v", ip, err)
		return nil, err
	}
	fmt.Printf("ip: %s, lat: %f, long: %f\n", ip, record.Location.Latitude, record.Location.Longitude)

	return &record, nil
}

func main() {
	parser := argparse.NewParser("run", "run the geolog websocket server")
	logFile := parser.String("l", "log_file", &argparse.Options{Required: true, Help: "log file to tail"})
	geoliteDb := parser.String("g", "geodb_file", &argparse.Options{Required: true, Help: "geolite db to use"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	gdb, err := maxminddb.Open(*geoliteDb)
	if err != nil {
		log.Fatalf("failed to open maxmind db: %v", err)
		os.Exit(1)
	}
	defer gdb.Close()

	tailHandle, err := tail.TailFile(*logFile, tail.Config{
		Follow:   true,
		Location: &tail.SeekInfo{Whence: io.SeekEnd},
	})
	if err != nil {
		log.Fatalf("failed to tail %s: %v", "ips.log", err)
	}
	for line := range tailHandle.Lines {
		_, err := handleLine(line.Text, gdb)
		if err != nil {
			log.Printf("failed to handle line: %v", err)
		}
	}
}
