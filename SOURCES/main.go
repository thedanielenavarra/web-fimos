package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func main() {
	// execute the init() function from auth.go
	println("Executing the init() function from auth.go")
	configInit(os.Args[1])
	// if the program is executed with the --regen flag, call the tokenRegeneration function from auth.go and exit
	if len(os.Args) > 2 && os.Args[2] == "--regen" {
		println("Calling the tokenRegeneration function from auth.go")
		// check if the program is executed with the --host <hostname> --port <port> flags
		if len(os.Args) > 5 && os.Args[3] == "--host" && os.Args[5] == "--port" {
			port, err := strconv.Atoi(os.Args[6])
			if err != nil {
				fmt.Println("Invalid port number")
				return
			}
			host := os.Args[4]
			tokenRegeneration(host, port)
			return
		}

		tokenRegeneration(config.Host, config.Port)
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
	http.ListenAndServe(config.Host+":"+strconv.Itoa(config.Port), nil)
}
