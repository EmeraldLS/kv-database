package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", "3031", "port server is binding to")
	flag.Parse()

	conn, err := net.Dial("tcp", "127.0.0.1:"+*addr)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())
	fmt.Printf("Enter Request: ")

	for scanner.Scan() {

		req := scanner.Text()
		if req == "!q" {
			break
		}

		_, err := conn.Write([]byte(req))
		if err != nil {
			slog.Error("unable to write send request to the server", "err", err)
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err == io.EOF {
			continue
		}
		if err != nil {
			slog.Error("unable to read response from the server", "err", err)
		}

		handle_response(buf, n)

		fmt.Print("Enter Request: ")
		defer conn.Close()
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-signals
		}()
	}
}

func handle_response(buf []byte, n int) {
	var content interface{}
	err := json.Unmarshal(buf[:n], &content)
	if err != nil {
		slog.Error("unable to unmarshal json data", "err", err)
	}

	data, err := json.MarshalIndent(content, "", " ")
	if err != nil {
		slog.Error("unable to marshal indent content", "err", err)
	}
	fmt.Println(string(data))
}
