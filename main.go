package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"profile-golang/config"
	"profile-golang/reader"
	"profile-golang/writer"
)

func main() {
	mode := flag.String("mode", "reader", "Server mode: reader or writer")
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	switch *mode {
	case "reader":
		server := reader.NewServer(cfg, *port)
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start reader server: %v", err)
		}
	case "writer":
		server := writer.NewServer(cfg, *port)
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start writer server: %v", err)
		}
	default:
		fmt.Printf("Unknown mode: %s\n", *mode)
		os.Exit(1)
	}
}
