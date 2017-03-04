package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		os.Exit(1)
	}
	apiServer := os.Args[1] + "/api/v1/pods"

	for {
		resp, err := http.Get("http://" + apiServer)
		if err != nil {
			fmt.Println("Can't find the apiServer")
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(strconv.Itoa(resp.StatusCode) + " : " + string(body))
		}
	}

	time.Sleep(1000 * time.Millisecond)
}
