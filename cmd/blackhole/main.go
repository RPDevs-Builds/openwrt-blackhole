package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	ipAddr     = flag.String("ip", "0.0.0.0", "IP address to bind to")
	port       = flag.String("port", "8080", "Port to listen on")
	logFile    = flag.String("log", "blackhole.log", "Log file path")
	rootDir    = flag.String("root", "blackhole_root", "Root directory for mirroring")
	contentDir = flag.String("content", "blackhole_content", "Directory containing real content to serve")
)

const version = "0.1.0"

var trackingPixel = []byte{
	0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00, 0x80, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0x21, 0xf9, 0x04, 0x01, 0x00,
	0x00, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
	0x00, 0x02, 0x02, 0x44, 0x01, 0x00, 0x3b,
}

type RequestLog struct {
	Timestamp  string              `json:"timestamp"`
	RemoteAddr string              `json:"remote_addr"`
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Headers    map[string][]string `json:"headers"`
}

func main() {
	flag.Parse()

	// Ensure root directory exists
	if err := os.MkdirAll(*rootDir, 0755); err != nil {
		log.Fatalf("Failed to create root directory: %v", err)
	}
	// Ensure content directory exists
	if err := os.MkdirAll(*contentDir, 0755); err != nil {
		log.Fatalf("Failed to create content directory: %v", err)
	}

	http.HandleFunc("/", handleRequest)

	addr := fmt.Sprintf("%s:%s", *ipAddr, *port)
	fmt.Printf("Blackhole server v%s listening on %s\n", version, addr)
	fmt.Printf("Logging to: %s\n", *logFile)
	fmt.Printf("Mirroring to: %s/\n", *rootDir)
	fmt.Printf("Serving content from: %s/\n", *contentDir)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 1. Log the request
	go logRequest(r)

	// 2. Check if real content exists and serve it
	cleanPath := path.Clean("/" + r.URL.Path)
	contentPath := filepath.Join(*contentDir, cleanPath)

	info, err := os.Stat(contentPath)
	if err == nil {
		// If it's a file and has content
		if !info.IsDir() && info.Size() > 0 {
			http.ServeFile(w, r, contentPath)
			return
		}
		// If it's a directory, check for index.html
		if info.IsDir() {
			idxPath := filepath.Join(contentPath, "index.html")
			idxInfo, idxErr := os.Stat(idxPath)
			if idxErr == nil && !idxInfo.IsDir() && idxInfo.Size() > 0 {
				http.ServeFile(w, r, contentPath)
				return
			}
		}
	}

	// 3. Mirror the path to filesystem
	go mirrorPath(r.URL.Path)

	// 4. Respond with tracking pixel
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)
	w.Write(trackingPixel)
}

func logRequest(r *http.Request) {
	entry := RequestLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		RemoteAddr: r.RemoteAddr,
		Method:     r.Method,
		URL:        r.URL.String(),
		Headers:    r.Header,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

func mirrorPath(rawPath string) {
	// Sanitize and clean path
	cleanPath := path.Clean("/" + rawPath)
	if cleanPath == "/" {
		return
	}

	targetPath := filepath.Join(*rootDir, cleanPath)

	// Logic: If it has an extension, treat as file. Else treat as directory.
	ext := path.Ext(cleanPath)
	if ext != "" {
		// It's a file. Create the parent directories first.
		dir := filepath.Dir(targetPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Error creating directory %s: %v", dir, err)
			return
		}

		// Create the empty file
		f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error creating file %s: %v", targetPath, err)
			return
		}
		f.Close()
	} else {
		// It's a directory. Create it (and parents).
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			log.Printf("Error creating directory %s: %v", targetPath, err)
		}
	}
}
