package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	ipAddress := "10.0.0.11"
	requests := 200

	for i := 0; i < requests; i++ {
		start := time.Now()
		resp, err := http.Get("http://" + ipAddress)
		if err != nil {
			fmt.Printf("Error occurred: %s\n", err)
			continue
		}
		elapsed := time.Since(start)
		fmt.Printf("Request %d: Status - %s, Time taken - %s\n", i+1, resp.Status, elapsed)
		resp.Body.Close()
	}
}
