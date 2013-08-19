package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var target = flag.String("url", "", "The URL to test rate-limiting against")
var total = flag.Int("count", 100, "The number of requests to test")

var ratelimited_status = flag.Int("ratelimited_status", 429, "The response code recieved when you have been rate limited")
var success_status = flag.Int("success_status", 200, "The response code recieved when you have not been ratelimited")

var min_successes = flag.Int("min_successes", 10, "The minimum number of requests that should have a success response code")
var min_limited_responses = flag.Int("min_limited_responses", 85, "The minimum number of requests that should receive a ratelimiting code")

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
	fmt.Printf("Response Histogram: %+v\n", statusCodeHist)
	if *min_successes > statusCodeHist[*success_status] {
		fmt.Printf("Needed %d requests to be successful, got %d\n", *min_successes, statusCodeHist[*success_status])
		os.Exit(1)
	}
	if *min_limited_responses > statusCodeHist[*ratelimited_status] {
		fmt.Printf("Needed %d requests to be ratelimited, got %d\n", *min_limited_responses, statusCodeHist[*ratelimited_status])
		os.Exit(1)
	}
}
