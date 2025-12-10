package web

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/VineLink-Lab/i18n/internal/parser"
)

func Web(directory string, port string) error {
	p, err := parser.NewParser(directory)
	if err != nil {
		return err
	}
	return startWebServer(p, port)
}

func startWebServer(p *parser.Parser, port string) error {
	routeHandler := NewRouteHandler(p)
	routeHandler.RegisterRoutes()

	fullPort := ":" + port
	url := fmt.Sprintf("http://localhost:%s", port)

	// Print server information
	log.Println("========================================")
	log.Printf("🌍 i18n Translation Manager")
	log.Printf("📂 Directory: %s", p.GetDirectoryPath())
	log.Printf("🚀 Server started at: %s", url)
	log.Println("========================================")
	log.Println("Press Ctrl+C to stop the server")

	// Open browser after a short delay
	go func() {
		time.Sleep(500 * time.Millisecond)
		if err := openBrowser(url); err != nil {
			log.Printf("⚠️  Failed to open browser: %v", err)
			log.Printf("Please manually open: %s", url)
		} else {
			log.Printf("✅ Browser opened successfully")
		}
	}()

	return routeHandler.StartServer(fullPort)
}

// openBrowser opens the default browser with the given URL
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}
