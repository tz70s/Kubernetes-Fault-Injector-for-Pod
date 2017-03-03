package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
)

var count int = 0
var policy string
var statusCode int

type Backend struct {
	host string
	port string
}

// receive query string and change policy

func injectPolicy(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		io.WriteString(w, "Policy not change ! Do you type the right query string ? policy=xxx \nCurrent Inject Policy: "+policy+"\n")
		return
	}
	if req.Form.Get("policy") != "" {
		policy = req.Form.Get("policy")
		io.WriteString(w, "Change the policy to "+policy+"\n")
		return
	}
	if req.Form.Get("status") != "" {
		policy = "statusResponse"
		tmpCode, err := strconv.Atoi(req.Form.Get("status"))

		if err != nil {
			io.WriteString(w, "Wrong status code \n")
			return
		}
		statusCode = tmpCode
		io.WriteString(w, "Status Code : "+req.Form.Get("status")+"\n")
		return
	}

	io.WriteString(w, "Policy not change ! Do you type the right query string ? policy=xxx \nCurrent Inject Policy: "+policy+"\n")
}

// select fault injection actions

func injectSelect(w http.ResponseWriter, req *http.Request) {
	// check current policy
	switch policy {
	case "simpleResponse":
		injectResponse(w, req)
		break
	case "delayResponse":
		delayResponse(w, req)
		break
	case "statusResponse":
		statusResponse(w, req)
		break
	default:
		injectResponse(w, req)
	}
	count++
}

// Simple response

func injectResponse(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, this is the simple async response from injector\n")
}

// delay response

func delayResponse(w http.ResponseWriter, req *http.Request) {
	time.Sleep(3000 * time.Millisecond)
	io.WriteString(w, "Hello, this is the simple async delay response from injector\n")
}

// status code

func statusResponse(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(statusCode)
	io.WriteString(w, "Injector response with the status code : "+strconv.Itoa(statusCode)+"\n")
}

// Indirect Get to the backend
func (node *Backend) injectGetRedirect(w http.ResponseWriter, req *http.Request) {
	remote, err := url.Parse("http://" + node.host + ":" + node.port)
	if err != nil {
		log.Fatal("Parse : ", err)
		return
	}

	if count%5 == 0 {
		injectSelect(w, req)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, req)
	count++
}

func main() {
	// input port from args
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No enough arguments")
		os.Exit(1)
	}
	node_server := &Backend{host: "127.0.0.1", port: args[1]}
	// default reverse proxy
	http.Handle("/", http.HandlerFunc(node_server.injectGetRedirect))

	// injector with no query string
	http.Handle("/injector", http.HandlerFunc(injectPolicy))

	fmt.Println("Start Listen Injector and serve at http://localhost:12345")
	err := http.ListenAndServe("localhost:12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
