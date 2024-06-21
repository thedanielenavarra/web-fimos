package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	Protocol        string `json:"protocol"`
	DestinationPort int    `json:"destinationPort"`
	Duration        int    `json:"duration"`
	Token           string `json:"token"`
}

// add flags, the config file is first and mandatory, the --regen flag is optional as well as the --host and --port flags
var (
	configFile = flag.String("config", "", "Config file")
	regen      = flag.Bool("regen", false, "Regenerate the token")
	host       = flag.String("host", "localhost", "Host to connect to")
	port       = flag.Int("port", 9091, "Port to connect to")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", "webFimos")
		flag.PrintDefaults()
	}
	flag.Parse()
	configInit(*configFile)
	if *regen {
		if *host != "" {
			config.Host = *host
		}
		if *port != 0 {
			config.Port = *port
		}
		tokenRegeneration()
		writeConfig()
		fmt.Println("Token regenerated, config file updated")
		return
	}

	// start an http server
	println("Starting an http server on " + config.Host + ":" + strconv.Itoa(config.Port) + "...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// check if the request method is POST
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		// take protocol and destination port from the request's body, than take the source IP. After that call addRule
		var request Request
		fmt.Println("Request received from " + r.RemoteAddr)

		// decode the request's body
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println("Error elaborating request: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tokenVerified := verifyToken(request.Token)
		if !tokenVerified {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		remoteIP := r.RemoteAddr
		// remove the second part after the : if present
		remoteIP = strings.Split(remoteIP, ":")[0]

		// call the addRule function from dbus.go
		err = addRule(remoteIP, request.DestinationPort, "public", request.Protocol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Rule added for " + remoteIP + " on port " + strconv.Itoa(request.DestinationPort) + " with protocol " + request.Protocol)
		// send a response to the client
		time.AfterFunc(time.Duration(request.Duration)*time.Second, func() {
			err = removeRule(remoteIP, request.DestinationPort, "public", request.Protocol)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("Rule removed for " + remoteIP + " on port " + strconv.Itoa(request.DestinationPort) + " with protocol " + request.Protocol)
		})
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Rule added for %s on port %d with protocol %s", remoteIP, request.DestinationPort, request.Protocol)
	})
	err := http.ListenAndServe(config.Host+":"+strconv.Itoa(config.Port), nil)
	if err != nil {
		fmt.Println("Error starting the http server: ", err)
	}

}
