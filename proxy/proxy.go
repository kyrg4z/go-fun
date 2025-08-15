package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func forwardHandler(req *http.Request, clientConn net.Conn) {
	targetURL := req.URL
	if !targetURL.IsAbs() {
		host := req.Host // Example: "example.com"
		targetURL.Scheme = "https"
		targetURL.Host = host
	}
	
	newReq, err := http.NewRequest(req.Method, targetURL.String(), req.Body)
	if err != nil {
		fmt.Println("Failed to create new request:", err)
		return
	}

	// Copy headers
	for key, values := range req.Header {
		for _, value := range values {
			newReq.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(newReq)

	if err != nil {
		fmt.Println("Error forwarding request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Fprintf(clientConn, "HTTP/1.1 %s\r\n", resp.Status)

	// Write headers
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Fprintf(clientConn, "%s: %s\r\n", key, value)
		}
	}
	fmt.Fprint(clientConn, "\r\n") // End of headers

	// Stream body back to client
	io.Copy(clientConn, resp.Body)

}

func main() {

	listener, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Fatal("Failed to start listener:", err)
	}

	fmt.Println("Listening on port 8081...")

	conn, err := listener.Accept()

	if err != nil {
		log.Println("Failed to accept connection:", err)
	}
	reader := bufio.NewReader(conn)

	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Println("Failed to read HTTP request:", err)
	}

	fmt.Println("Method:", req.Method)
	fmt.Println("Path:", req.URL)
	// fmt.Println("Headers", req.Header)
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	fmt.Println("Forward/Drop")
	var input string
	_, err = fmt.Scanln(&input)

	if err != nil {
		fmt.Printf("Error input: %v", err)
	}

	action := strings.ToLower(input)


	if action == "f" {
		forwardHandler(req, conn)
	} else {
		fmt.Fprint(conn, "HTTP/1.1 403 Forbidden: Get blocked \r\nContent-Length: 0\r\n\r\n")
		conn.Close()
	}

}
