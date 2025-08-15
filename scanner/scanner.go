package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func worker(ports, results chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	const numWorkers = 1000
	const maxPorts = 1024
	
	ports := make(chan int, numWorkers)
	results := make(chan int, maxPorts)
	var openports []int
	var wg sync.WaitGroup
	
	// Start workers
	wg.Add(numWorkers)
	for range numWorkers {
		go worker(ports, results, &wg)
	}
	
	// Send ports to scan
	go func() {
		defer close(ports)
		for i := 1; i <= maxPorts; i++ {
			ports <- i
		}
	}()
	
	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Collect results
	for port := range results {
		if port != 0 {
			openports = append(openports, port)
		}
	}
	
	// Sort open ports & print
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}