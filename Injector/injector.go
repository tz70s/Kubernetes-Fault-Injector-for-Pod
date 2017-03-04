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

var FaultInject struct {
	policy         string
	statusCode     int
	boundedRetries int
	boundCount     int
	timeout        time.Duration
}

type Backend struct {
	host string
	port string
}

// receive query string and change policy

func injectPolicy(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		io.WriteString(w, "Policy not change ! Do you type the right query string ? policy=xxx \nCurrent Inject Policy: "+FaultInject.policy+"\n")
		return
	}
	if req.Form.Get("policy") != "" {
		FaultInject.policy = req.Form.Get("policy")
		io.WriteString(w, "Change the policy to "+FaultInject.policy+"\n")
		return
	}
	if req.Form.Get("status") != "" {
		FaultInject.policy = "statusResponse"
		tmpCode, err := strconv.Atoi(req.Form.Get("status"))

		if err != nil {
			io.WriteString(w, "Wrong status code \n")
			return
		}
		FaultInject.statusCode = tmpCode
		io.WriteString(w, "Status Code : "+req.Form.Get("status")+"\n")
		return
	}

	if req.Form.Get("boundedRetries") != "" {
		FaultInject.policy = "boundedRetries"
		tmpBound, err := strconv.Atoi(req.Form.Get("boundedRetries"))

		if err != nil {
			io.WriteString(w, "Wrong bound \n")
			return
		}
		FaultInject.boundedRetries = tmpBound
		FaultInject.boundCount = 0
		io.WriteString(w, "Set bound : "+req.Form.Get("boundedRetries")+"\n")
		return
	}

	if req.Form.Get("timeout") != "" {
		FaultInject.policy = "timeout"
		tmpTimeout, err := strconv.Atoi(req.Form.Get("timeout"))

		if err != nil {
			io.WriteString(w, "Wrong timeout \n")
			return
		}
		FaultInject.timeout = time.Duration(tmpTimeout)
		io.WriteString(w, "Set timeout : "+req.Form.Get("timeout")+"\n")
		return
	}

	io.WriteString(w, "Policy not change ! Do you type the right query string ? policy=xxx \nCurrent Inject Policy: "+FaultInject.policy+"\n")
}

// select fault injection actions

func injectSelect(w http.ResponseWriter, req *http.Request) {

	count++
	// check current policy
	switch FaultInject.policy {
	case "simpleResponse":
		injectResponse(w, req)
		break
	case "timeout":
		delayResponse(w, req)
		break
	case "abortResponse":
		abortResponse(w, req)
		break
	case "boundedRetries":
		boundedRetries(w, req)
		break
	case "statusResponse":
		statusResponse(w, req)
		break
	default:
		injectResponse(w, req)
	}
}

// Simple response

func injectResponse(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, this is the simple async response from injector\n")
}

// delay response

func delayResponse(w http.ResponseWriter, req *http.Request) {
	time.Sleep(FaultInject.timeout * time.Second)
	io.WriteString(w, "Hello, this is the simple async delay response from injector\n")
}

// status code

func statusResponse(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(FaultInject.statusCode)
	io.WriteString(w, "Injector response with the status code : "+strconv.Itoa(FaultInject.statusCode)+"\n")
}

// abort request
func abortResponse(w http.ResponseWriter, req *http.Request) {
}

// bounded retries testing
func boundedRetries(w http.ResponseWriter, req *http.Request) {
	FaultInject.boundCount++
	if FaultInject.boundCount >= FaultInject.boundedRetries {
		io.WriteString(w, "Bounded Retries testing failed!\n")
	}
	abortResponse(w, req)
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
	count++
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, req)
}

func main() {
	// input port from args
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No enough arguments")
		os.Exit(1)
	}
	node_server := &Backend{host: "127.0.0.1", port: args[1]}

	// initialized fault injection params
	FaultInject.boundCount = 0
	FaultInject.boundedRetries = 10
	FaultInject.policy = "simpleResponse"
	FaultInject.statusCode = 200
	FaultInject.timeout = 10

	fmt.Println("Initializing fault injection params : ")
	fmt.Println("Policy : " + FaultInject.policy)
	fmt.Printf("BoundedRetries : %d\n", FaultInject.boundedRetries)
	fmt.Printf("Timeout : %d\n", FaultInject.timeout)

	// default reverse proxy
	http.Handle("/", http.HandlerFunc(node_server.injectGetRedirect))

	// injector with no query string
	http.Handle("/injector", http.HandlerFunc(injectPolicy))

	fmt.Println("Start Listen Injector and serve at :8282")
	err := http.ListenAndServe(":8282", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
