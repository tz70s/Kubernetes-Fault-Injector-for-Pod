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
	domain := args[1]

	for {
		resp, err := http.Get("http://" + domain)
		if err != nil {
			fmt.Println("Can't find the service/dns")
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(strconv.Itoa(resp.StatusCode) + " : " + string(body))
		}

		time.Sleep(1000 * time.Millisecond)

	}
}
