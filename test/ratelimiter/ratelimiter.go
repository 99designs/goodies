package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var target = flag.String("url", "", "The URL to test rate-limiting against")
var total = flag.Int("count", 100, "The number of requests to test")

var ratelimited_status = flag.Int("ratelimited_status", 429, "The response code recieved when you have been rate limited")
var success_status = flag.Int("success_status", 200, "The response code recieved when you have not been ratelimited")

var min_successes = flag.Int("min_successes", 10, "The minimum number of requests that should have a success response code")
var min_limited_responses = flag.Int("min_limited_responses", 85, "The minimum number of requests that should receive a ratelimiting code")

var max_rate = flag.Int("max_rate", 1, "The expected long-term maximum requests per second")

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	errors := make(chan error, *total)
	results := make(chan *http.Response, *total)

	start := time.Now()
	for i := 0; i < *total; i++ {
		wg.Add(1)
		go func() {
			resp, err := http.Get(*target)
			if err != nil {
				errors <- err
			} else {
				results <- resp
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(errors)
	close(results)

	allowanceForTimeTaken := int(time.Since(start).Seconds()) * (*max_rate)

	anyErrors := false
	statusCodeHist := make(map[int]int)
	for err := range errors {
		if err != nil {
			fmt.Println("Error: " + err.Error())
			anyErrors = true
		}
	}

	for res := range results {
		if res != nil {
			statusCodeHist[res.StatusCode] += 1
		}
	}

	fmt.Printf("Response Histogram: %+v\n", statusCodeHist)
	if (*min_successes) > statusCodeHist[*success_status] {
		fmt.Printf("Needed %d requests to be successful, got %d\n", *min_successes, statusCodeHist[*success_status])
		os.Exit(1)
	}
	if (*(min_limited_responses) - allowanceForTimeTaken) > statusCodeHist[*ratelimited_status] {
		fmt.Printf("Needed %d requests to be ratelimited, got %d allowing %d for time taken\n", *min_limited_responses, statusCodeHist[*ratelimited_status], allowanceForTimeTaken)
		os.Exit(1)
	}
	if anyErrors {
		os.Exit(1)
	}
}
