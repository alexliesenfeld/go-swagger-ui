package main

import (
	"flag"
	"fmt"
	swaggerui "github.com/alexliesenfeld/go-swagger-ui"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"runtime"
)

type programArguments struct {
	specFilePath    string
	persistAuth     bool
	enableFilterBar bool
}

func main() {
	args, err := parseFlags()
	if err != nil {
		log.Fatalf("failed to parse arguments: %v", err)
	}

	if args.specFilePath == "" {
		printUsage()
		return
	}

	handler, err := newHandler(args)
	if err != nil {
		log.Fatalf("cannot create handler: %v", err)
	}

	srv := httptest.NewServer(handler)

	log.Println("starting Swagger UI server at", srv.URL)
	log.Println("press Ctrl+C to stop")

	if err := openBrowser(srv.URL); err != nil {
		log.Fatalf("failed to open browser: %v", err)
	}

	<-make(chan struct{})
}

func parseFlags() (programArguments, error) {
	flag.Parse()

	if flag.NArg() > 1 {
		return programArguments{}, fmt.Errorf("too many parameters")
	}

	specFilePath := flag.Arg(0)
	persistAuth := flag.Bool("persist-auth", false, "Enables browser authentication persistence")
	enableFilterBar := flag.Bool("show-filter-bar", false, "Shows a filter bar in the UI that helps to find API operations")

	return programArguments{
		specFilePath:    specFilePath,
		persistAuth:     *persistAuth,
		enableFilterBar: *enableFilterBar,
	}, nil
}

func printUsage() {
	fmt.Println("Usage: swui <path-to-schema>")
	flag.PrintDefaults()
}

func newHandler(args programArguments) (http.HandlerFunc, error) {
	return swaggerui.NewHandler(
		swaggerui.WithSpecFilePath(args.specFilePath),
		swaggerui.WithPersistAuthorization(args.persistAuth),
		swaggerui.WithDisplayRequestDuration(true),
		swaggerui.WithCredentials(true),
		swaggerui.WithShowCommonExtensions(true),
		swaggerui.WithShowExtensions(true),
		swaggerui.WithShowMutatedRequest(true),
		swaggerui.WithTitle(args.specFilePath),
		swaggerui.WithFilter(args.enableFilterBar),
	), nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
