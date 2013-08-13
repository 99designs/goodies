package main

import (
	"flag"
	"fmt"
	"net/http"
)

var target = flag.String("url", "", "The URL to test rate-limiting against")
var total = flag.Int("count", 100, "The number of requests to test")

func main() {
	flag.Parse()
	errors := make(chan error)
	results := make(chan *http.Response)

	finished := 0
	for i := 0; i <= *total; i++ {
		go func() {
			resp, err := http.Get(*target)
			if err != nil {
				errors <- err
			} else {
				results <- resp
			}
			finished += 1
			if finished >= *total {
				close(errors)
				close(results)
			}
		}()
	}

	statusCodeHist := make(map[int]int)
	func() {
		for {
			select {
			case err := <-errors:
				if err != nil {
					fmt.Println("Error: " + err.Error())
				}
				if finished >= *total {
					return
				}
			case res := <-results:
				if res != nil {
					statusCodeHist[res.StatusCode] += 1
				}
				if finished >= *total {
					return
				}
			}
		}
	}()
	fmt.Println(statusCodeHist)
}
