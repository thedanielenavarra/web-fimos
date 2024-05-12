package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Request struct {
	Protocol        string `json:"protocol"`
	DestinationPort int    `json:"destinationPort"`
	Duration        int    `json:"duration"`
}

func main() {
	// execute the init() function from auth.go
	println("Executing the init() function from auth.go")
	configInit()
	// start an http server
	println("Starting an http server on " + config.Host + ":" + strconv.Itoa(config.Port) + "...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

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
		remoteIP := r.RemoteAddr
		// remove the second part after the : if present
		remoteIP = strings.Split(remoteIP, ":")[0]

		// call the addRule function from dbus.go
		addRule(remoteIP, request.DestinationPort, "public", request.Protocol)

	})
	// get the port from the environment variable
	http.ListenAndServe(config.Host+":"+strconv.Itoa(config.Port), nil)
}
